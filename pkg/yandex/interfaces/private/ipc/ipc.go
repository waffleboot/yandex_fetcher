package ipc

import app "github.com/waffleboot/playstation_buy/pkg/yandex/application"

type YandexItem struct {
	Host string
	Url  string
}

type Endpoint struct {
	service app.Service
}

func NewEndpoint(s app.Service) Endpoint {
	return Endpoint{service: s}
}

func (s Endpoint) GetYandexItems(search string) ([]YandexItem, error) {
	items, err := s.service.ParseYandex(search)
	if err != nil {
		return nil, err
	}
	out := make([]YandexItem, 0, len(items))
	for _, v := range items {
		out = append(out, YandexItem{
			Host: v.Host,
			Url:  v.Url,
		})
	}
	return out, nil
}
