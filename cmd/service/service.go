package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-chi/chi"

	cache "github.com/waffleboot/playstation_buy/pkg/cache"

	yandex_application "github.com/waffleboot/playstation_buy/pkg/yandex/application"
	yandex_infra_yandex "github.com/waffleboot/playstation_buy/pkg/yandex/infra/yandex"
	yandex_interfaces_private_ipc "github.com/waffleboot/playstation_buy/pkg/yandex/interfaces/private/ipc"

	root_application "github.com/waffleboot/playstation_buy/pkg/root/application"
	root_infra_worker "github.com/waffleboot/playstation_buy/pkg/root/infra/worker/http"
	root_infra_yandex "github.com/waffleboot/playstation_buy/pkg/root/infra/yandex"
	root_interfaces_public_http "github.com/waffleboot/playstation_buy/pkg/root/interfaces/public/http"
)

func run(args []string) int {

	serviceAddr := os.Getenv("SERVICE_ADDR")
	if serviceAddr == "" {
		return 1
	}

	checkerUrl := os.Getenv("CHECKER_URL")
	if checkerUrl == "" {
		return 1
	}

	log.Printf("Starting service on %s", serviceAddr)

	if err := startServer(serviceAddr, checkerUrl); err != nil {
		log.Println(err)
		return 2
	}
	return 0
}

func startServer(serviceAddr, checkerUrl string) error {

	r := chi.NewRouter()

	cache := &cache.MemoryCache{}

	timeout := 3 * time.Second

	yandex := yandex_interfaces_private_ipc.NewEndpoint(
		yandex_application.NewService(
			yandex_infra_yandex.NewHttpClient, yandexFetchers(1)))

	service := root_interfaces_public_http.NewEndpoint(
		root_application.NewService(
			timeout,
			root_infra_yandex.NewYandex(yandex),
			root_infra_worker.NewBenchmarkSupplier(checkerUrl),
			cache))

	service.AddRoutes(r)

	server := &http.Server{Addr: serviceAddr, Handler: r}
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

func yandexFetchers(def int) int {
	s := os.Getenv("YANDEX_FETCHERS")
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil || n <= 0 {
		return def
	}
	return n
}
