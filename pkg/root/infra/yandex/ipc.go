package ipc

import (
	"context"

	"github.com/waffleboot/playstation_buy/pkg/common/domain"
	"github.com/waffleboot/playstation_buy/pkg/yandex/interfaces/private/ipc"
)

type YandexSupplier struct {
	channel chan ipc.ChannelItem
}

func NewYandexSupplier(channel chan ipc.ChannelItem) *YandexSupplier {
	return &YandexSupplier{channel: channel}
}

func (y *YandexSupplier) GetYandexItems(ctx context.Context, search string) ([]domain.YandexItem, error) {
	item := ipc.NewChannelItem(ctx, search)
	y.channel <- item
	select {
	case data := <-item.Done:
		return data, nil
	case err := <-item.Err:
		return nil, err
	}
}
