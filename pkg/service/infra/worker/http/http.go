package http

import (
	"bytes"
	"encoding/json"

	"github.com/waffleboot/playstation_buy/pkg/common/domain"
	"github.com/waffleboot/playstation_buy/pkg/worker/interfaces/private/http"

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

func (b *BenchmarkSupplier) Benchmark(item domain.YandexItem) (domain.StatsItem, error) {
	req := http.Request{
		Host: item.Host,
		Url:  item.Url,
	}
	body, err := json.Marshal(req)
	if err != nil {
		return domain.StatsItem{}, err
	}
	httpRequest, err := http2.NewRequest(http2.MethodPost, b.checkerUrl, bytes.NewReader(body))
	if err != nil {
		return domain.StatsItem{}, err
	}
	resp, err := b.httpClient.Do(httpRequest)
	if err != nil {
		return domain.StatsItem{}, err
	}
	var ans http.Response
	if err := json.NewDecoder(resp.Body).Decode(&ans); err != nil {
		return domain.StatsItem{}, err
	}
	if err := resp.Body.Close(); err != nil {
		return domain.StatsItem{}, err
	}
	return domain.StatsItem{
		Host:  ans.Host,
		Count: ans.Count,
	}, nil
}
