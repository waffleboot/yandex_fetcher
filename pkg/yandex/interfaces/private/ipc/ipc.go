package ipc

import (
	"github.com/waffleboot/playstation_buy/pkg/common/domain"
	app "github.com/waffleboot/playstation_buy/pkg/yandex/application"
)

type Endpoint struct {
	service *app.Service
}

func NewEndpoint(s *app.Service) *Endpoint {
	return &Endpoint{service: s}
}

func (s *Endpoint) GetYandexItems(search string) ([]domain.YandexItem, error) {
	return s.service.ParseYandex(search)
}
