package ipc

import (
	"github.com/waffleboot/playstation_buy/pkg/root/domain"
	"github.com/waffleboot/playstation_buy/pkg/yandex/interfaces/private/ipc"
)

type YandexSupplier struct {
	ipc.Endpoint
}

func NewYandexSupplier(endpoint ipc.Endpoint) YandexSupplier {
	return YandexSupplier{Endpoint: endpoint}
}

func (y YandexSupplier) Supply(search string) ([]domain.SearchEngineItem, error) {
	data, err := y.Endpoint.GetYandexItems(search)
	if err != nil {
		return nil, err
	}
	out := make([]domain.SearchEngineItem, 0, len(data))
	for _, v := range data {
		out = append(out, domain.SearchEngineItem{
			Host: v.Host,
			Url:  v.Url,
		})
	}
	return out, nil
}
