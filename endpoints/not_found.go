package endpoints

import (
	"encoding/json"
	"net/http"

	"github.com/robertacosta/go-integration-example/client"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(client.Error{Message: "not found"})
}
