package ipc

import (
	app "github.com/waffleboot/yandex_fetcher/pkg/checker/application"
	"github.com/waffleboot/yandex_fetcher/pkg/common/domain"
)

type ChannelItem struct {
	Items []domain.YandexItem
	Done  chan domain.StatsItem
	Errc  chan error
}

func StartEndpoint(s *app.Service, channel chan ChannelItem) {
	go func() {
		for e := range channel {
			var i int
			for i < len(e.Items) {
				ans, err := s.Benchmark(e.Items[i].Host, e.Items[i].Url)
				i++
				if err != nil {
					e.Errc <- err
					break
				}
				e.Done <- domain.StatsItem{
					Host:  e.Items[i].Host,
					Count: ans,
				}
			}
			for i < len(e.Items) {
				ans, err := s.Benchmark(e.Items[i].Host, e.Items[i].Url)
				i++
				if err != nil {
					continue
				}
				e.Done <- domain.StatsItem{
					Host:  e.Items[i].Host,
					Count: ans,
				}
			}
		}
	}()
}
