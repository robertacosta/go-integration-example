package adaptor

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/robertacosta/go-integration-example/client"
)

type Requestor interface {
	Request() (*client.Message, error)
	SetHttpClient(httpClient *http.Client)
}

type TestService struct {
	addr       string
	httpClient *http.Client
}

func NewTestService(addr string) *TestService {
	return &TestService{
		addr:       addr,
		httpClient: http.DefaultClient,
	}
}

func (t *TestService) SetHttpClient(httpClient *http.Client) {
	t.httpClient = httpClient
}

func (t *TestService) Request() (*client.Message, error) {
	url := fmt.Sprintf("http://%s/message", t.addr)
	resp, err := t.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	message := &client.Message{}
	err = json.NewDecoder(resp.Body).Decode(message)

	return message, err
}
