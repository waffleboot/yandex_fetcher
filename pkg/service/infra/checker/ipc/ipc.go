package ipc

import (
	"context"

	"github.com/waffleboot/playstation_buy/pkg/checker/interfaces/private/ipc"
	"github.com/waffleboot/playstation_buy/pkg/common/domain"
)

type BenchmarkSupplier struct {
	channel chan ipc.ChannelItem
}

func NewBenchmarkSupplier(channel chan ipc.ChannelItem) *BenchmarkSupplier {
	return &BenchmarkSupplier{
		channel: channel,
	}
}

func (b *BenchmarkSupplier) Benchmark(ctx context.Context, items []domain.YandexItem) ([]domain.StatsItem, error) {
	done := make(chan domain.StatsItem, len(items))
	errc := make(chan error, 1)
	entry := ipc.ChannelItem{
		Items: items,
		Done:  done,
		Errc:  errc,
	}
	m := make([]domain.StatsItem, 0, len(items))
	select {
	case b.channel <- entry:
		for {
			select {
			case item := <-done:
				m = append(m, item)
			case err := <-errc:
				return m, err
			case <-ctx.Done():
				return m, ctx.Err()
			}
		}
	case <-ctx.Done():
		return m, ctx.Err()
	}
}
