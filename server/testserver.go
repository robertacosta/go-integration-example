package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/robertacosta/go-integration-example/endpoints"
)

type TestConfig struct {
	TestServiceAddr string
}

type TestServer struct {
	router     *mux.Router
	httpServer *http.Server

	addr string
}

func NewTestServer(config TestConfig) (*TestServer, error) {
	s := &TestServer{
		router: mux.NewRouter(),
		addr:   config.TestServiceAddr,
	}

	s.initializeRoutes(config)

	return s, nil
}

func (s *TestServer) initializeRoutes(config TestConfig) {

	testExampleEndpoint := endpoints.TestExample{}

	s.router.Handle("/message", http.HandlerFunc(testExampleEndpoint.Get)).Methods("GET")
	s.router.NotFoundHandler = http.HandlerFunc(endpoints.NotFound)
}

func (s *TestServer) Start() error {
	s.httpServer = &http.Server{Addr: s.addr, Handler: s.router}

	log.Printf("service listening on %s\n", s.addr)
	if err := s.httpServer.ListenAndServe(); err != nil {
		return fmt.Errorf("error occurred when starting up service: %s", err)
	}

	return nil
}

func (s *TestServer) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
