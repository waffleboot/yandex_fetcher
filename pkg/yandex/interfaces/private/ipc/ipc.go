package ipc

import (
	"context"

	"github.com/waffleboot/playstation_buy/pkg/common/domain"
	app "github.com/waffleboot/playstation_buy/pkg/yandex/application"
)

type channelItem struct {
	done   chan []domain.YandexItem
	errc   chan error
	ctx    context.Context
	search string
}

type Endpoint struct {
	channel chan channelItem
}

func NewEndpoint(s *app.Service) *Endpoint {
	channel := make(chan channelItem, 1)
	go func() {
		for item := range channel {
			if item.ctx.Err() == context.DeadlineExceeded {
				continue
			}
			s.GetItems(item.ctx, item.search, item.done, item.errc)
		}
	}()
	return &Endpoint{
		channel: channel,
	}
}

func (e *Endpoint) AddQuery(ctx context.Context, search string) (chan []domain.YandexItem, chan error) {
	datc := make(chan []domain.YandexItem, 1)
	errc := make(chan error, 1)
	select {
	case e.channel <- channelItem{
		ctx:    ctx,
		done:   datc,
		errc:   errc,
		search: search}:
	case <-ctx.Done():
		break
	}
	return datc, errc
}
