package ipc

import (
	"github.com/waffleboot/yandex_fetcher/pkg/common/domain"
	app "github.com/waffleboot/yandex_fetcher/pkg/yandex/application"
)

type Endpoint struct {
	service *app.Service
}

func NewEndpoint(s *app.Service) *Endpoint {
	return &Endpoint{
		service: s,
	}
}

func (e *Endpoint) AddQuery(search string) ([]domain.YandexItem, error) {
	return e.service.GetItems(search)
}
