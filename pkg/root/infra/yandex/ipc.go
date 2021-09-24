package ipc

import (
	"github.com/waffleboot/playstation_buy/pkg/common/domain"
	"github.com/waffleboot/playstation_buy/pkg/yandex/interfaces/private/ipc"
)

type YandexSupplier struct {
	ipc.Endpoint
}

func NewYandexSupplier(endpoint ipc.Endpoint) YandexSupplier {
	return YandexSupplier{Endpoint: endpoint}
}

func (y YandexSupplier) Supply(search string) ([]domain.YandexItem, error) {
	return y.Endpoint.GetYandexItems(search)
}
