package ipc

import (
	"context"

	"github.com/waffleboot/playstation_buy/pkg/common/domain"
	"github.com/waffleboot/playstation_buy/pkg/worker/interfaces/private/ipc"
)

type BenchmarkSupplier struct {
	channel chan ipc.ChannelItem
}

func NewBenchmarkSupplier(channel chan ipc.ChannelItem) *BenchmarkSupplier {
	return &BenchmarkSupplier{
		channel: channel,
	}
}

func (b *BenchmarkSupplier) Benchmark(ctx context.Context, items []domain.YandexItem) (chan domain.StatsItem, chan error) {
	datc := make(chan domain.StatsItem, len(items))
	errc := make(chan error, 1)
	entry := ipc.ChannelItem{
		Items: items,
		Datc:  datc,
		Errc:  errc,
	}
	select {
	case b.channel <- entry:
	case <-ctx.Done():
		break
	}
	return datc, errc
}
