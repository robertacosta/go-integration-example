package cmd

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/bsm/sarama-cluster"
	"github.com/robertacosta/go-integration-example/adaptor"
	"github.com/robertacosta/go-integration-example/server"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Runs http server",
	Run:   startServer,
}

func startServer(cmd *cobra.Command, args []string) {
	kafkaBrokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	kafkaTopic := os.Getenv("KAFKA_TOPIC")
	serviceAddr := os.Getenv("SERVICE_ADDRESS")
	testServiceAddr := os.Getenv("TEST_SERVICE_ADDRESS")

	config := server.Config{
		KafkaBrokers:    kafkaBrokers,
		KafkaTopic:      kafkaTopic,
		ServiceAddr:     serviceAddr,
		TestServiceAddr: testServiceAddr,
	}

	stop := make(chan os.Signal)
	defer close(stop)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	deps := serverDependencies(config)

	defer deps.SyncProducer.Close()

	srv, err := server.NewServer(config, deps)
	if err != nil {
		log.Printf("unable to initialize service: %s", err)
	}

	srv.Start()

	//// wait for sigint or sigterm
	<-stop

	log.Println("shutting down servers")
	ctx, cncl := context.WithTimeout(context.Background(), time.Duration(5000)*time.Millisecond)
	defer cncl()
	srv.Stop(ctx)
}

func serverDependencies(config server.Config) server.Dependencies {
	hystrixTestServiceConfig := hystrix.CommandConfig{
		Timeout:               5000, // Timeout request after 5 sec
		MaxConcurrentRequests: 100,  // Bulk head, max requests that can be concurrently running, all others rejected
		SleepWindow:           5000, // If circuit opens, try to close every 5 sec
		ErrorPercentThreshold: 50,   // If over 50% of the requests return an error, open ciruit
	}
	testService := adaptor.NewHystrixTestService(config.TestServiceAddr, hystrixTestServiceConfig)
	testServiceTimeDuraction := time.Duration(5000) * time.Millisecond
	testServiceHttpClient := &http.Client{
		Timeout: testServiceTimeDuraction,
		Transport: &http.Transport{
			TLSHandshakeTimeout: testServiceTimeDuraction,
		},
	}
	testService.SetHttpClient(testServiceHttpClient)

	kafkaClient, err := newKafkaClient(config.KafkaBrokers)
	if err != nil {
		panic(err)
	}

	producer, err := sarama.NewSyncProducerFromClient(kafkaClient)
	if err != nil {
		panic(err)
	}

	return server.Dependencies{
		TestService:  testService,
		SyncProducer: producer,
	}
}

func newKafkaClient(brokers []string) (*cluster.Client, error) {
	config := cluster.NewConfig()

	// required for sync producer
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	config.ClientID = hostname

	client, err := cluster.NewClient(brokers, config)
	if err != nil {
		return nil, err
	}

	return client, nil
}
