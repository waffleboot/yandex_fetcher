package ipc

import (
	"github.com/waffleboot/playstation_buy/pkg/common/domain"

	"github.com/waffleboot/playstation_buy/pkg/yandex/interfaces/private/ipc"
)

type Yandex struct {
	endpoint *ipc.Endpoint
}

func NewYandex(endpoint *ipc.Endpoint) *Yandex {
	return &Yandex{endpoint: endpoint}
}

func (y *Yandex) GetItems(search string) ([]domain.YandexItem, error) {
	return y.endpoint.AddQuery(search)
}