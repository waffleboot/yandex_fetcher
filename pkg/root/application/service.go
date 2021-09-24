package application

import "github.com/waffleboot/playstation_buy/pkg/root/domain"

type supplier interface {
	Supply(search string) ([]domain.SearchEngineItem, error)
}

type Service struct {
	supplier supplier
}

func NewService(supplier supplier) Service {
	return Service{supplier: supplier}
}
