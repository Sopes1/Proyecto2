package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/go-redis/redis/v8"
)

var (
	URI = "mongodb://mongo_admin:SopesP2_2021@" + os.Getenv("MHOST") + ":27017/admin"
)

var count = 1
var ctx = context.Background()

type Log struct {
	Request_Game int    `json:"request_game"`
	Game         string `json:"game"`
	Gamename     string `json:"gamename"`
	Winner       string `json:"winner"`
	Players      string `json:"players"`
	Worker       string `json:"worker"`
}

type Game struct {
	Game     string `json:"game"`
	Gamename string `json:"gamename"`
	Winner   string `json:"winner"`
	Players  string `json:"players"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	workerRabbit()

}

func workerRabbit() {
	conn, err := amqp.Dial("amqp://guest:guest@" + os.Getenv("QRHOST") + ":5672/")
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

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			splittedBody := strings.Split(string(d.Body), "-")
			idGame := splittedBody[0]
			gameName := splittedBody[1]
			players, _ := strconv.Atoi(splittedBody[2])

			rand.Seed(time.Now().UnixNano())
			winner := strconv.Itoa(rand.Intn(players) + 1)

			newGameMongo(Game{Game: idGame, Gamename: gameName, Players: splittedBody[2], Winner: winner})
			newGameRedis(Game{Game: idGame, Gamename: gameName, Players: splittedBody[2], Winner: winner})
			newLogMongo(Log{Request_Game: count, Game: idGame, Gamename: gameName, Players: splittedBody[2], Winner: winner, Worker: "RabbitMQ"})
			count++
			log.Print("Ganador: " + winner)

		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func newGameMongo(data Game) bool {
	clientOptions := options.Client().ApplyURI(URI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Println("Mongo.connect() ERROR: ", err)
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	col := client.Database("Proyecto2Sopes").Collection("Games")

	_, insertErr := col.InsertOne(ctx, data)
	if insertErr != nil {
		fmt.Println("InsertONE Error:", insertErr)
		defer cancel()
		return false
	} else {
		defer cancel()
		return true
	}
}

func newLogMongo(data Log) bool {
	clientOptions := options.Client().ApplyURI(URI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Println("Mongo.connect() ERROR: ", err)
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	col := client.Database("Proyecto2Sopes").Collection("Logs")

	_, insertErr := col.InsertOne(ctx, data)
	if insertErr != nil {
		fmt.Println("InsertONE Error:", insertErr)
		defer cancel()
		return false
	} else {
		defer cancel()
		return true
	}
}

func newGameRedis(data Game) {

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("RHOST") + ":6379",
		Password: "",
		DB:       0,
	})

	dg, _ := json.Marshal(data)
	_, err := rdb.Do(ctx, "lpush", "Game", dg).Result()

	if err != nil {
		panic(err)
	}

}
