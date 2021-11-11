package main

import (
	"context"
	"fmt"
	"log"

	"net"

	"github.com/streadway/amqp"

	pb "../proto"
	"google.golang.org/grpc"
)

type DataGame struct {
	Id      string
	Name    string
	Players string
}

type server struct {
	pb.UnimplementedGameServiceServer
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func (*server) Game(ctx context.Context, req *pb.GameRequest) (*pb.GameResponse, error) {
	fmt.Printf(">> SERVER: Función Greet llamada con éxito. Datos: %v\n", req)

	// Todos los datos podemos obtenerlos desde req
	// Tendra la misma estructura que definimos en el protofile
	// Para ello utilizamos en este caso el GetGreeting
	idGame := req.GetGame().GetGame()
	gameName := req.GetGame().GetGamename()
	players := req.GetGame().GetPlayers()

	RabbitMQ(idGame, gameName, players)

	result := "200"
	fmt.Printf(">> SERVER: %s\n", result)

	// Creamos un nuevo objeto GreetResponse definido en el protofile
	res := &pb.GameResponse{
		Result: result,
	}

	return res, nil
}

func RabbitMQ(idGame string, gameName string, players string) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()
	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	body := idGame + "-" + gameName + "-" + players

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
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
