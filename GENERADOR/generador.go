package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var reader = bufio.NewReader(os.Stdin)
var gamesNumber []string
var gamesName []string
var players, rungames, concurrence, timeout int

type Datos struct {
	Game     string
	GameName string
	Players  int
}

func main() {
	showMenu()
}

func showMenu() {
	fmt.Print("Cli: ")
	entrada, _ := reader.ReadString('\n')
	analizador(entrada)
	ejectuar()
}

func analizador(entrada string) {

	var cli = strings.Split(strings.Replace(entrada, "\"", "", -1), " ")

	if strings.ToLower(cli[0]) != "rungame" {
		fmt.Println("Sintaxis incorrecta en columna 1. " + cli[0])
		return
	}

	for i := 1; i < len(cli); i++ {

		switch cli[i] {

		case "--gamename":
			i++

			for true {
				if isNumeric(cli[i]) {
					gamesNumber = append(gamesNumber, cli[i])
					i++
				}

				if cli[i] != "|" {
					fmt.Println("Error de sintaxis en la columna " + strconv.Itoa(i) + " \"" + cli[i] + "\"")
					return
				}

				i++

				gamesName = append(gamesName, cli[i])

				if cli[i+1] != "|" {
					break
				}

				i = i + 2
			}

			break

		case "--players":
			i++
			players, _ = strconv.Atoi(cli[i])
			break

		case "--rungames":
			i++
			rungames, _ = strconv.Atoi(cli[i])
			break

		case "--concurrence":
			i++
			concurrence, _ = strconv.Atoi(cli[i])
			break

		case "--timeout":
			i++
			timeout, _ = strconv.Atoi(strings.Split(cli[i], "m")[0])
			break

		}

	}
}

func ejectuar() {
	var seconds = 0
	var minutos = 0
	var totalGames = 0
	var stop = false
	for true {
		for i := 0; i < concurrence; i++ {
			go enviarDatos(newJSON())
			totalGames++
			if totalGames >= rungames {
				stop = true
				break
			}
		}
		time.Sleep(1 * time.Second)
		seconds++
		if seconds == 60 {
			minutos++
			seconds = 0
		}
		if minutos >= timeout || stop {
			fmt.Println("")
			if stop {
				fmt.Println("Termina por cantidad alcanzada")
			} else {
				fmt.Println("Termina por tiempo")
			}
			break
		}
	}
}

func newJSON() Datos {

	rand.Seed(time.Now().UnixNano())
	var index = rand.Intn(len(gamesName))
	var game = gamesNumber[index]
	var name = gamesName[index]
	var nplayers = rand.Intn(players)

	return Datos{Game: game, GameName: name, Players: nplayers}
}

func enviarDatos(dato Datos) {
	fmt.Print(dato)
}

//Auxiliares
func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}
