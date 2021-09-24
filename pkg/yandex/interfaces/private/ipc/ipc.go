package ipc

import (
	"context"

	"github.com/waffleboot/playstation_buy/pkg/common/domain"
	app "github.com/waffleboot/playstation_buy/pkg/yandex/application"
)

type ChannelItem struct {
	Done   chan []domain.YandexItem
	Err    chan error
	Ctx    context.Context
	Search string
}

func NewChannelItem(ctx context.Context, search string) ChannelItem {
	return ChannelItem{
		Ctx:    ctx,
		Search: search,
		Done:   make(chan []domain.YandexItem, 1),
		Err:    make(chan error),
	}
}

func StartEndpoint(channel chan ChannelItem, s *app.Service) {
	go func() {
		for task := range channel {
			s.GetYandexItems(task.Ctx, task.Search, task.Done, task.Err)
		}
	}()
}
