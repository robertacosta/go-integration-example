package endpoints

import (
	"encoding/json"
	"net/http"

	"github.com/robertacosta/go-integration-example/client"
)

type TestExmaple struct{}

func (e *TestExmaple) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(client.Message{Message: "Hello Go Meetup"})
}
