package ipc

import (
	"context"

	"github.com/waffleboot/playstation_buy/pkg/common/domain"
	app "github.com/waffleboot/playstation_buy/pkg/worker/application"
)

type channelItem struct {
	items []domain.YandexItem
	datc  chan domain.StatsItem
	errc  chan error
}

type Endpoint struct {
	channel chan channelItem
}

func NewEndpoint(s *app.Service) *Endpoint {
	channel := make(chan channelItem, 1)
	go func() {
		for e := range channel {
			var i int
			for i < len(e.items) {
				ans, err := s.Benchmark(e.items[i])
				i++
				if err != nil {
					e.errc <- err
					break
				}
				e.datc <- ans
			}
			for i < len(e.items) {
				ans, err := s.Benchmark(e.items[i])
				i++
				if err != nil {
					continue
				}
				e.datc <- ans
			}
		}
	}()
	return &Endpoint{channel: channel}
}

func (e *Endpoint) Benchmark(ctx context.Context, items []domain.YandexItem) (chan domain.StatsItem, chan error) {
	datc := make(chan domain.StatsItem, len(items))
	errc := make(chan error, 1)
	entry := channelItem{
		items: items,
		datc:  datc,
		errc:  errc,
	}
	select {
	case e.channel <- entry:
	case <-ctx.Done():
		break
	}
	return datc, errc
}
