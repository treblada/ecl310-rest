package api

import (
	"context"
	"log"
	"net/http"

	"github.com/goburrow/modbus"
	openapi "github.com/treblada/ecl310-rest/openapi/go"
)

type HealthApiService struct {
	openapi.HealthApiService
	client modbus.Client
}

func NewHealthApiService(client modbus.Client) openapi.HealthApiServicer {
	return &HealthApiService{
		client: client,
	}
}

func (s *HealthApiService) GetHealth(ctx context.Context) (openapi.ImplResponse, error) {
	status := "OK"
	_, err := s.client.ReadHoldingRegisters(278, 4)

	if err != nil {
		log.Printf("ECL host not reachable: %v", err)
		status = "FAIL"
	}

	body := openapi.GetHealthResponse{
		Status: status,
	}
	return openapi.Response(http.StatusOK, body), nil
}
