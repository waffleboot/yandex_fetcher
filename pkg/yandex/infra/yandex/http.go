package http

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const baseYandexURL = "https://yandex.ru/search/touch/?service=www.yandex&ui=webmobileapp.yandex&numdoc=50&lr=213&p=0&text=%s"

func NewHttpClient() func(context.Context, string) ([]byte, error) {
	var c http.Client
	return func(ctx context.Context, search string) ([]byte, error) {
		url := fmt.Sprintf(baseYandexURL, url.QueryEscape(search))
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}
		resp, err := c.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
}
