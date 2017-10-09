package endpoints

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Shopify/sarama"
	"github.com/robertacosta/go-integration-example/client"
)

type requestor interface {
	Request() (*client.Message, error)
}

type sendMessager interface {
	SendMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error)
}

type Example struct {
	Requestor    requestor
	SendMessager sendMessager
	Topic        string
}

func (e *Example) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	// Make request to get message
	message, err := e.Requestor.Request()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(client.Error{Message: fmt.Sprintf("Error retriving message: %s", err)})
		return
	}

	b, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshalling message: %s", err)
	}

	// Publish message to kafka
	newMsg := &sarama.ProducerMessage{Topic: e.Topic, Value: sarama.ByteEncoder(b)}
	_, _, err = e.SendMessager.SendMessage(newMsg)
	if err != nil {
		log.Printf("Error publishing to Kafka topic: %s", err)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(message)
}
