package http

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/waffleboot/playstation_buy/pkg/service/interfaces/public/http"

	http2 "net/http"
)

type InitialService struct {
	httpClient http2.Client
	serviceUrl string
}

func NewInitialService(serviceUrl string) *InitialService {
	return &InitialService{
		httpClient: http2.Client{},
		serviceUrl: serviceUrl,
	}
}

func (b *InitialService) Update(ctx context.Context, host string, count int) error {
	req := http.Update{
		Host:  host,
		Count: count,
	}
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	httpRequest, err := http2.NewRequestWithContext(ctx, http2.MethodPost, b.serviceUrl, bytes.NewReader(body))
	if err != nil {
		return err
	}
	resp, err := b.httpClient.Do(httpRequest)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
