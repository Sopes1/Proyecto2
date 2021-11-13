// Paquete principal, acá iniciará la ejecución
package main

// Importar dependencias, notar que estamos en un módulo llamado tuiterclient
import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"fmt"
	"log"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"

	greetpb "servidor/proto"
)

type Health struct {
	Status int
}

type Response struct {
	Ganador string
}

type Request struct {
	IdGame   string
	GameName string
	Players  string
}

type server struct{}

// Funcion que realiza una llamada unaria
func sendMessage(idGame string, gameName string, players int) string {
	server_host := "localhost:50051"

	fmt.Println(">> CLIENT: Iniciando cliente")
	fmt.Println(">> CLIENT: Iniciando conexion con el servidor gRPC ", server_host)

	// Crear una conexion con el servidor (que esta corriendo en el puerto 50051)
	// grpc.WithInsecure nos permite realizar una conexion sin tener que utilizar SSL
	cc, err := grpc.Dial(server_host, grpc.WithInsecure())
	if err != nil {
		log.Fatalf(">> CLIENT: Error inicializando la conexion con el server %v", err)
	}

	// Defer realiza una axion al final de la ejecucion (en este caso, desconectar la conexion)
	defer cc.Close()

	// Iniciar un servicio NewGreetServiceClient obtenido del codigo que genero el protofile
	// Esto crea un cliente con el cual podemos escuchar
	// Le enviamos como parametro el Dial de gRPC
	c := greetpb.NewGameServiceClient(cc)

	fmt.Println(">> CLIENT: Iniciando llamada a Unary RPC")

	// Crear una llamada de GameRequest
	// Este codigo lo obtenemos desde el archivo que generamos con protofile
	req := &greetpb.GameRequest{
		Game: &greetpb.Game{
			Game:     idGame,
			Gamename: gameName,
			Players:  strconv.Itoa(players),
		},
	}

	fmt.Println(">> CLIENT: Enviando datos al server")
	// Iniciar un greet, en background con la peticion que estamos realizando
	res, err := c.Game(context.Background(), req)
	if err != nil {
		log.Fatalf(">> CLIENT: Error realizando la peticion %v", err)
	}

	fmt.Println(">> CLIENT: El servidor nos respondio con el siguiente mensaje: ", res.Result)
	return res.Result
}

func root(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(Health{Status: 200})
}

func sendData(w http.ResponseWriter, req *http.Request) {
	var request Request
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&request)
	if err != nil {
		panic(err)
	}
	p, _ := strconv.Atoi(request.Players)
	response := sendMessage(request.IdGame, request.GameName, p)

	json.NewEncoder(w).Encode(Response{Ganador: response})

}

// Funcion principal
func main() {

	router := mux.NewRouter()
	//Ruta para combrobar conexion
	router.HandleFunc("/", root).Methods("GET")
	router.HandleFunc("/sendData", sendData).Methods("GET")

	log.Fatal(http.ListenAndServe(":3200", router))

	//sendMessage("1", "GAME 1", 30)
}
