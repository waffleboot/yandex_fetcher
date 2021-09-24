package application

import "github.com/waffleboot/playstation_buy/pkg/common/domain"

type supplier interface {
	Supply(search string) ([]byte, error)
}

type Service struct {
	supplier
}

func NewService(supplier supplier) *Service {
	return &Service{supplier: supplier}
}

func (s *Service) ParseYandex(search string) ([]domain.YandexItem, error) {
	data, err := s.supplier.Supply(search)
	if err != nil {
		return nil, err
	}
	result := parseYandexResponse(data)
	if result.Error != nil {
		return nil, result.Error
	}
	out := make([]domain.YandexItem, 0, len(result.Items))
	for _, v := range result.Items {
		out = append(out, domain.YandexItem{
			Host: v.Host,
			Url:  v.Url,
		})
	}
	return out, nil
}
