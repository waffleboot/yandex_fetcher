package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"

	cache "github.com/waffleboot/playstation_buy/pkg/cache"

	worker_application "github.com/waffleboot/playstation_buy/pkg/worker/application"
	worker_interfaces_private_ipc "github.com/waffleboot/playstation_buy/pkg/worker/interfaces/private/ipc"

	yandex_application "github.com/waffleboot/playstation_buy/pkg/yandex/application"
	yandex_infra_yandex "github.com/waffleboot/playstation_buy/pkg/yandex/infra/yandex"
	yandex_interfaces_private_ipc "github.com/waffleboot/playstation_buy/pkg/yandex/interfaces/private/ipc"

	root_application "github.com/waffleboot/playstation_buy/pkg/root/application"
	root_infra_worker "github.com/waffleboot/playstation_buy/pkg/root/infra/worker/ipc"
	root_infra_yandex "github.com/waffleboot/playstation_buy/pkg/root/infra/yandex"
	root_interfaces_public_http "github.com/waffleboot/playstation_buy/pkg/root/interfaces/public/http"
)

func run(args []string) int {
	log.Println("Starting service")
	if err := startServer(); err != nil {
		log.Println(err)
		return 2
	}
	return 0
}

func startServer() error {

	r := chi.NewRouter()

	ctx := context.Background()

	cache := &cache.MemoryCache{}

	timeout := 3 * time.Second

	yandexFetchers := 1

	yandex := yandex_interfaces_private_ipc.NewEndpoint(
		yandex_application.NewService(
			yandex_infra_yandex.NewHttpClient, yandexFetchers))

	siteFetchers := 1

	worker_channel := make(chan worker_interfaces_private_ipc.ChannelItem, 1)

	worker_interfaces_private_ipc.StartEndpoint(
		worker_application.NewService(ctx, cache, siteFetchers, timeout),
		worker_channel)

	service := root_interfaces_public_http.NewEndpoint(
		root_application.NewService(
			timeout,
			root_infra_yandex.NewYandex(yandex),
			root_infra_worker.NewBenchmarkSupplier(worker_channel),
			cache))

	service.AddRoutes(r)

	server := &http.Server{Addr: ":9000", Handler: r}
	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func signalContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer cancel()
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c
	}()
	return ctx
}
