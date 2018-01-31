package adaptor

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/robertacosta/go-integration-example/client"
)

const TestServiceHystrixKey = "test_service"

type HystrixTestService struct {
	*TestService
	hystrixKey string
}

func NewHystrixTestService(addr string, hystrixCfg hystrix.CommandConfig) *HystrixTestService {
	hystrix.ConfigureCommand(TestServiceHystrixKey, hystrixCfg)

	testService := NewTestService(addr)

	return &HystrixTestService{
		TestService: testService,
		hystrixKey:  TestServiceHystrixKey,
	}
}

func (h *HystrixTestService) SetHttpClient(httpClient *http.Client) {
	h.httpClient = httpClient
}

func (h *HystrixTestService) Request() (*client.Message, error) {
	message := &client.Message{}

	hystrixErr := hystrix.Do(h.hystrixKey, func() error {
		url := fmt.Sprintf("http://%s/message", h.addr)
		resp, err := h.httpClient.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// If the response was not OK, only ding the circuit if it was 500+ status code
		if resp.StatusCode != http.StatusOK && (resp.StatusCode >= 500) {
			return errors.New(fmt.Sprintf("Received a non-200 status code, %d", resp.StatusCode))
		}

		err = json.NewDecoder(resp.Body).Decode(message)

		return nil
	}, func(err error) error {
		// If the circuit opens, then the fallback function is called
		log.Printf("Circuit Fallback, Error received: %s", err)

		message = &client.Message{Message: "Keep Calm and Eat Pizza"}

		return nil
	})

	if hystrixErr != nil {
		return nil, hystrixErr
	}

	return message, nil
}
