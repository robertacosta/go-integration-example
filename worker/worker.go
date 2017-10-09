package worker

import (
	"encoding/json"
	"log"

	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
	"github.com/robertacosta/go-integration-example/client"
)

type Config struct {
	KafkaBrokers []string
	KafkaTopic   string
}

type Dependencies struct {
	Consumer *cluster.Consumer
}

type Worker struct {
	deps   Dependencies
	config Config
}

func New(deps Dependencies, config Config) *Worker {
	return &Worker{
		deps:   deps,
		config: config,
	}
}

func (w *Worker) Run() {
	log.Println("Running kafka consumer")

	var msg *sarama.ConsumerMessage
	var err error
	msgc := w.deps.Consumer.Messages()
	errc := w.deps.Consumer.Errors()
	for {
		select {
		case err = <-errc:
			log.Printf("Error consuming from Kafka: %s", err)
		case msg = <-msgc:
			message := client.Message{}
			err := json.Unmarshal(msg.Value, &message)
			if err != nil {
				log.Printf("Error unmarshalling json")
				continue
			}

			log.Printf("Received message: %s\n", message.Message)

			w.deps.Consumer.MarkOffset(msg, "") // mark message as processed
		}
	}
}
