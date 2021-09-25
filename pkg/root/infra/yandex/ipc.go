package ipc

import (
	"context"

	"github.com/waffleboot/playstation_buy/pkg/common/domain"
)

type Endpoint interface {
	AddQuery(context.Context, string) (chan []domain.YandexItem, chan error)
}

type Yandex struct {
	endpoint Endpoint
}

func NewYandex(endpoint Endpoint) *Yandex {
	return &Yandex{endpoint: endpoint}
}

func (y *Yandex) GetItems(ctx context.Context, search string) ([]domain.YandexItem, error) {
	datc, errc := y.endpoint.AddQuery(ctx, search)
	select {
	case data := <-datc:
		return data, nil
	case err := <-errc:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
