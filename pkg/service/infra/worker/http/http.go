package http

import (
	"bytes"
	"context"
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

func (b *BenchmarkSupplier) Benchmark(ctx context.Context, items []domain.YandexItem) ([]domain.StatsItem, error) {
	m := make([]domain.StatsItem, 0, len(items))
	for _, item := range items {
		ans, err := b.process(ctx, item)
		if ctx.Err() == context.DeadlineExceeded {
			return m, ctx.Err()
		} else if err != nil {
			return m, err
		}
		m = append(m, ans)
	}
	return m, nil
}

func (b *BenchmarkSupplier) process(ctx context.Context, item domain.YandexItem) (domain.StatsItem, error) {
	req := http.Request{
		Host: item.Host,
		Url:  item.Url,
	}
	body, err := json.Marshal(req)
	if err != nil {
		return domain.StatsItem{}, err
	}
	httpRequest, err := http2.NewRequestWithContext(ctx, http2.MethodPost, b.checkerUrl, bytes.NewReader(body))
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
