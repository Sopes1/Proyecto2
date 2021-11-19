package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"net"

	"github.com/streadway/amqp"

	pb "servidor/proto"

	"google.golang.org/grpc"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"cloud.google.com/go/pubsub"
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

	switch os.Getenv("COLA") {
	case "rabbit":
		RabbitMQ(idGame, gameName, players)
		break
	case "kafka":
		saveGameKafka(idGame, gameName, players)
		break
	case "pubsub":
		publish(idGame, gameName, players)
		break
	}

	result := "200"
	fmt.Printf(">> SERVER: %s\n", result)

	// Creamos un nuevo objeto GreetResponse definido en el protofile
	res := &pb.GameResponse{
		Result: result,
	}

	return res, nil
}

func RabbitMQ(idGame string, gameName string, players string) {
	conn, err := amqp.Dial("amqp://guest:guest@" + os.Getenv("RHOST") + ":5672/")
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

func saveGameKafka(idGame string, gameName string, players string) {

	host := os.Getenv("QKHOST")
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": host})
	if err != nil {
		panic(err)
	}

	body := idGame + "-" + gameName + "-" + players

	// Produce messages to topic (asynchronously)
	topic := "topic2"
	for _, word := range []string{string(body)} {
		p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(word),
		}, nil)
	}
}

func publish(idGame string, gameName string, players string) error {
	fmt.Println("save to pubsub")

	// Definimos el ProjectID del proyecto
	// Este dato lo sacamos de Google Cloud
	projectID := "sopes1p2-330401" //goDotEnvVariable("PROJECT_ID")

	// Definimos el TopicId del proyecto
	// Este dato lo sacamos de Google Cloud
	topicID := "sopesp2" //goDotEnvVariable("TOPIC_ID")

	// Definimos el contexto en el que ejecutaremos PubSub
	ctx := context.Background()
	// Creamos un nuevo cliente
	client, err := pubsub.NewClient(ctx, projectID)
	// Si un error ocurrio creando el nuevo cliente, entonces imprimimos un error y salimos
	if err != nil {
		fmt.Println("error aqui")
		fmt.Println(err)
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}

	// Obtenemos el topico al que queremos enviar el mensaje
	t := client.Topic(topicID)

	body := idGame + "-" + gameName + "-" + players

	// Publicamos los datos del mensaje
	result := t.Publish(ctx, &pubsub.Message{Data: []byte(body)})

	// Bloquear el contexto hasta que se tenga una respuesta de parte de GooglePubSub
	id, err := result.Get(ctx)

	// Si hubo un error creando el mensaje, entonces mostrar que existio un error
	if err != nil {
		fmt.Println("error:")
		fmt.Println(err)
		return fmt.Errorf("Error: %v", err)
	}

	// El mensaje fue publicado correctamente
	fmt.Printf("Published a message; msg ID: %v\n", id)
	return nil
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
