package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"

	pb "../proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedGameServiceServer
}

func (*server) Game(ctx context.Context, req *pb.GameRequest) (*pb.GameResponse, error) {
	fmt.Printf(">> SERVER: Función Greet llamada con éxito. Datos: %v\n", req)

	// Todos los datos podemos obtenerlos desde req
	// Tendra la misma estructura que definimos en el protofile
	// Para ello utilizamos en este caso el GetGreeting
	idGame := req.GetGame().GetGame()
	gameName := req.GetGame().GetGamename()
	players, _ := strconv.Atoi(req.GetGame().GetPlayers())

	rand.Seed(time.Now().UnixNano())
	winner := rand.Intn(players) + 1
	result := "El ganador en el juego " + idGame + " " + gameName + " fue " + strconv.Itoa(winner)

	fmt.Printf(">> SERVER: %s\n", result)
	// Creamos un nuevo objeto GreetResponse definido en el protofile
	res := &pb.GameResponse{
		Result: result,
	}

	return res, nil
}

// Funcion principal
func main() {

	// Leer el host de las variables del ambiente
	host := "localhost:50051"
	fmt.Println(">> SERVER: Iniciando en ", host)

	// Primero abrir un puerto para poder escuchar
	// Lo abrimos en este puerto arbitrario
	lis, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatalf(">> SERVER: Error inicializando el servidor: %v", err)
	}

	fmt.Println(">> SERVER: Empezando server gRPC")

	// Ahora si podemos iniciar un server de gRPC
	s := grpc.NewServer()

	// Registrar el servicio utilizando el codigo que nos genero el protofile
	pb.RegisterGameServiceServer(s, &server{})

	fmt.Println(">> SERVER: Escuchando servicio...")
	// Iniciar a servir el servidor, si hay un error salirse
	if err := s.Serve(lis); err != nil {
		log.Fatalf(">> SERVER: Error inicializando el listener: %v", err)
	}
}
