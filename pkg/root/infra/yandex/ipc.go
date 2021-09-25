package ipc

import (
	"context"

	"github.com/waffleboot/playstation_buy/pkg/common/domain"

	"github.com/waffleboot/playstation_buy/pkg/yandex/interfaces/private/ipc"
)

type Yandex struct {
	endpoint *ipc.Endpoint
}

func NewYandex(endpoint *ipc.Endpoint) *Yandex {
	return &Yandex{endpoint: endpoint}
}

func (y *Yandex) GetItems(ctx context.Context, search string) ([]domain.YandexItem, error) {
	done, errc := y.endpoint.AddQuery(ctx, search)
	select {
	case data := <-done:
		return data, nil
	case err := <-errc:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
