package main
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"context"
	
	"github.com/gorilla/mux"
	// Libreria de Google PubSub
	"cloud.google.com/go/pubsub"
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

	publish(_game)

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

func publish(game DataGame) error {
	fmt.Println("save to pubsub")

	jsonString, err := json.Marshal(game)

	gameString := string(jsonString)
	fmt.Print(gameString)

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

	// Publicamos los datos del mensaje
	result := t.Publish(ctx, &pubsub.Message{Data: []byte(gameString)})

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