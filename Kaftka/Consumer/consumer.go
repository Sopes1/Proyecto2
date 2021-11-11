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

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/go-redis/redis/v8"
)
var (
	URI = "mongodb://mongo_admin:SopesP2_2021@3" + os.Getenv("MHOST") + ":27017/admin"
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

	receiveFromKafka()	
}

func receiveFromKafka() {
	host := os.Getenv("QKHOST") +":29092"
	fmt.Println("Start receiving from Kafka")
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":host,
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		panic(err)
	}

	c.SubscribeTopics([]string{"topic2"}, nil)

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
			splittedBody := strings.Split(string(msg.Value), "-")
			idGame := splittedBody[0]
			gameName := splittedBody[1]
			players, _ := strconv.Atoi(splittedBody[2])

			rand.Seed(time.Now().UnixNano())
			winner := strconv.Itoa(rand.Intn(players) + 1)

			newGameMongo(Game{Game: idGame, Gamename: gameName, Players: splittedBody[2], Winner: winner})
			newGameRedis(Game{Game: idGame, Gamename: gameName, Players: splittedBody[2], Winner: winner})
			newLogMongo(Log{Request_Game: count, Game: idGame, Gamename: gameName, Players: splittedBody[2], Winner: winner, Worker: "Kaftka"})
			count++
			log.Print("Ganador: " + winner)
		} else {
			// The client will automatically try to recover from all errors.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}

	c.Close()

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
		Addr:     os.Getenv("RHOST"),
		Password: "",
		DB:       0,
	})

	dg, _ := json.Marshal(data)
	_, err := rdb.Do(ctx, "rpush", "Game", dg).Result()

	if err != nil {
		panic(err)
	}

}