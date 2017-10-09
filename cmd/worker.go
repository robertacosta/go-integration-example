package cmd

import (
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/bsm/sarama-cluster"
	"github.com/robertacosta/go-integration-example/worker"
	"github.com/spf13/cobra"
)

var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Runs worker",
	Run:   runWorker,
}

func runWorker(cmd *cobra.Command, args []string) {
	kafkaBrokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	kafkaTopic := os.Getenv("KAFKA_TOPIC")

	config := worker.Config{
		KafkaBrokers: kafkaBrokers,
		KafkaTopic:   kafkaTopic,
	}

	stop := make(chan os.Signal, 1)
	defer close(stop)
	signal.Notify(stop, os.Interrupt)

	deps := workerDependencies(config)

	log.Println("starting worker")
	rtsWorker := worker.New(deps, config)
	go rtsWorker.Run()

	// wait for sigint or sigterm
	<-stop

	log.Println("Stopping worker")
}

func workerDependencies(config worker.Config) worker.Dependencies {
	client, err := newKafkaClient(config.KafkaBrokers)
	if err != nil {
		panic(err)
	}
	consumer, err := cluster.NewConsumerFromClient(client, "meetup", []string{config.KafkaTopic})
	if err != nil {
		panic(err)
	}

	return worker.Dependencies{
		Consumer: consumer,
	}
}
