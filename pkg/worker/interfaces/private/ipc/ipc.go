package ipc

import (
	"github.com/waffleboot/playstation_buy/pkg/common/domain"
	app "github.com/waffleboot/playstation_buy/pkg/worker/application"
)

type ChannelItem struct {
	Items []domain.YandexItem
	Datc  chan domain.StatsItem
	Errc  chan error
}

func StartEndpoint(s *app.Service, channel chan ChannelItem) {
	go func() {
		for e := range channel {
			var i int
			for i < len(e.Items) {
				ans, err := s.Benchmark(e.Items[i])
				i++
				if err != nil {
					e.Errc <- err
					break
				}
				e.Datc <- ans
			}
			for i < len(e.Items) {
				ans, err := s.Benchmark(e.Items[i])
				i++
				if err != nil {
					continue
				}
				e.Datc <- ans
			}
		}
	}()
}
