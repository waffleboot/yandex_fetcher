package http

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"strconv"

	"github.com/waffleboot/yandex_fetcher/pkg/checker/inter/private/http"

	http2 "net/http"
)

type BenchmarkSupplier struct {
	httpClient http2.Client
	checkerUrl string
}

func NewBenchmarkSupplier(checkerUrl string) *BenchmarkSupplier {
	return &BenchmarkSupplier{
		httpClient: http2.Client{},
		checkerUrl: checkerUrl,
	}
}

func (b *BenchmarkSupplier) Benchmark(host, url string) (int, error) {
	req := http.Request{
		Host: host,
		Url:  url,
	}
	body, err := json.Marshal(req)
	if err != nil {
		return 0, err
	}
	httpRequest, err := http2.NewRequest(http2.MethodPost, b.checkerUrl, bytes.NewReader(body))
	if err != nil {
		return 0, err
	}
	resp, err := b.httpClient.Do(httpRequest)
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
