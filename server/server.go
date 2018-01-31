package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Shopify/sarama"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/gorilla/mux"
	"github.com/robertacosta/go-integration-example/adaptor"
	"github.com/robertacosta/go-integration-example/endpoints"
)

type Config struct {
	KafkaBrokers []string
	KafkaTopic   string

	ServiceAddr     string
	TestServiceAddr string
}

type Dependencies struct {
	TestService  adaptor.Requestor
	SyncProducer sarama.SyncProducer
}

type Server struct {
	router     *mux.Router
	httpServer *http.Server

	addr string
}

func NewServer(config Config, deps Dependencies) (*Server, error) {
	s := &Server{
		router: mux.NewRouter(),
		addr:   config.ServiceAddr,
	}

	s.initializeRoutes(config, deps)

	return s, nil
}

func (s *Server) initializeRoutes(config Config, deps Dependencies) {

	exampleEndpoint := endpoints.Example{
		Requestor:    deps.TestService,
		SendMessager: deps.SyncProducer,
		Topic:        config.KafkaTopic,
	}

	// optional - if you want to see the stream and use the dashboard
	hystrixHandler := hystrix.NewStreamHandler()
	hystrixHandler.Start()

	s.router.Handle("/example", http.HandlerFunc(exampleEndpoint.Get)).Methods("GET")
	s.router.Handle("/hystrix", hystrixHandler)
	s.router.NotFoundHandler = http.HandlerFunc(endpoints.NotFound)
}

func (s *Server) Start() error {
	s.httpServer = &http.Server{Addr: s.addr, Handler: s.router}

	log.Printf("service listening on %s\n", s.addr)
	if err := s.httpServer.ListenAndServe(); err != nil {
		return fmt.Errorf("error occurred when starting up service: %s", err)
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
