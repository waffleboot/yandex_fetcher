package http

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"strconv"

	"github.com/waffleboot/yandex_fetcher/pkg/checker/inter/private/http"

	http2 "net/http"
)

func NewBenchmarkSupplier(checkerUrl string) func(string, string) (int, error) {
	var httpClient http2.Client
	return func(host, url string) (int, error) {
		req := http.Request{
			Host: host,
			Url:  url,
		}
		body, err := json.Marshal(req)
		if err != nil {
			return 0, err
		}
		httpRequest, err := http2.NewRequest(http2.MethodPost, checkerUrl, bytes.NewReader(body))
		if err != nil {
			return 0, err
		}
		resp, err := httpClient.Do(httpRequest)
		if err != nil {
			return 0, err
		}
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return 0, err
		}
		if err := resp.Body.Close(); err != nil {
			return 0, err
		}
		count, err := strconv.Atoi(string(data))
		if err != nil {
			return 0, err
		}
		return count, nil

	}
}
