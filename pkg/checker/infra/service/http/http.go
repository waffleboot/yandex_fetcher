package http

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/waffleboot/yandex_fetcher/pkg/service/inter/public/http"

	http2 "net/http"
)

func NewInitialService(serviceUrl string) func(string, int) error {
	var httpClient http2.Client
	return func(host string, count int) error {
		req := http.CacheUpdate{
			Host:  host,
			Count: count,
		}
		body, err := json.Marshal(req)
		if err != nil {
			return fmt.Errorf("cache update: %w", err)
		}
		httpRequest, err := http2.NewRequest(http2.MethodPost, serviceUrl, bytes.NewReader(body))
		if err != nil {
			return fmt.Errorf("cache update: %w", err)
		}
		resp, err := httpClient.Do(httpRequest)
		if err != nil {
			return fmt.Errorf("cache update: %w", err)
		}
		defer resp.Body.Close()
		return nil
	}
}
