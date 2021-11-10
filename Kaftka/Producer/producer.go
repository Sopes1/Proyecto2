package main
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gorilla/mux"
)

type DataGame struct {
	Id      string
	Name    string
	Players string
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/games", gamePostHandler).Methods("POST")

	log.Fatal(http.ListenAndServe(":9090", router))
	
}

func gamePostHandler(w http.ResponseWriter, r *http.Request) {

	//Retrieve body from http request
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		panic(err)
	}

	//Save data into game struct
	var _game DataGame
	err = json.Unmarshal(b, &_game)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	saveGameKafka(_game)

	//Convert game struct into json
	jsonString, err := json.Marshal(_game)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	//Set content-type http header
	w.Header().Set("content-type", "application/json")

	//Send back data as response
	w.Write(jsonString)

}

func saveGameKafka(game DataGame){
	fmt.Println("save to kafka")

	jsonString, err := json.Marshal(game)

	gameString := string(jsonString)
	fmt.Print(gameString)

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:29092"})
	if err != nil {
		panic(err)
	}

	//defer p.Close()

	// Delivery report handler for produced messages
	/*go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()*/

	// Produce messages to topic (asynchronously)
	topic := "topic2"
	for _, word := range []string{string(gameString)} {
		p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(word),
		}, nil)
	}

	// Wait for message deliveries before shutting down
	//p.Flush(15 * 1000)
}