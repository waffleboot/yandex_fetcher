package ipc

import (
	"github.com/waffleboot/playstation_buy/pkg/common/domain"
	"github.com/waffleboot/playstation_buy/pkg/yandex/interfaces/private/ipc"
)

type YandexSupplier struct {
	endpoint *ipc.Endpoint
}

func NewYandexSupplier(endpoint *ipc.Endpoint) *YandexSupplier {
	return &YandexSupplier{endpoint: endpoint}
}

func (y *YandexSupplier) Supply(search string) ([]domain.YandexItem, error) {
	return y.endpoint.GetYandexItems(search)
}
