package application

import "github.com/waffleboot/playstation_buy/pkg/common/domain"

type supplier interface {
	Supply(search string) ([]domain.YandexItem, error)
}

type Service struct {
	supplier supplier
}

func NewService(supplier supplier) Service {
	return Service{supplier: supplier}
}
