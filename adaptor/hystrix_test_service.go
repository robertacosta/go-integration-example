package adaptor

import (
	"errors"
	"net/http"

	"fmt"

	"encoding/json"

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

		// If the response was not OK, don't ding the circuit if this was a bad request
		if resp.StatusCode != http.StatusOK && (resp.StatusCode >= 500) {
			return errors.New(fmt.Sprintf("Received a non-200 status code, %d", resp.StatusCode))
		}

		err = json.NewDecoder(resp.Body).Decode(message)

		return nil
	}, nil)

	if hystrixErr != nil {
		return nil, hystrixErr
	}

	return message, nil
}
