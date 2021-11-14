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
	

	// Libreria de Google PubSub
	"cloud.google.com/go/pubsub"
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
	for{
		consumerPub()
	}
	
	
}

func consumerPub() {
	projectID := "sopes1p2-330401"
    subID := "subscripcionsopesp2"
    client, err := pubsub.NewClient(ctx, projectID)
    if err != nil {
		fmt.Println("pubsub.NewClient: ", err)
            return 
    }
    defer client.Close()
    sub := client.Subscription(subID)
    sub.ReceiveSettings.Synchronous = true
    // MaxOutstandingMessages is the maximum number of unprocessed messages the
    // client will pull from the server before pausing.
    //
    // This is only guaranteed when ReceiveSettings.Synchronous is set to true.
    // When Synchronous is set to false, the StreamingPull RPC is used which
    // can pull a single large batch of messages at once that is greater than
    // MaxOustandingMessages before pausing. For more info, see
    // https://cloud.google.com/pubsub/docs/pull#streamingpull_dealing_with_large_backlogs_of_small_messages.
    sub.ReceiveSettings.MaxOutstandingMessages = 10
    // MaxOutstandingBytes is the maximum size of unprocessed messages,
    // that the client will pull from the server before pausing. Similar
    // to MaxOutstandingMessages, this may be exceeded with a large batch
    // of messages since we cannot control the size of a batch of messages
    // from the server (even with the synchronous Pull RPC).
    sub.ReceiveSettings.MaxOutstandingBytes = 1e10
    err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		fmt.Println("Got message: ", string(msg.Data))
		splittedBody := strings.Split(string(msg.Data), "-")
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
        msg.Ack()
    })
    if err != nil {
		fmt.Println("Receive: ", err)
            return 
    }
    return 
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