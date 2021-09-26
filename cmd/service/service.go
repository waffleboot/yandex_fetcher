package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"go.uber.org/zap"

	cache "github.com/waffleboot/yandex_fetcher/pkg/cache"

	yandex_application "github.com/waffleboot/yandex_fetcher/pkg/yandex/application"
	yandex_infra_yandex "github.com/waffleboot/yandex_fetcher/pkg/yandex/infra/yandex"
	yandex_interfaces_private_ipc "github.com/waffleboot/yandex_fetcher/pkg/yandex/interfaces/private/ipc"

	root_application "github.com/waffleboot/yandex_fetcher/pkg/service/application"
	root_infra_worker "github.com/waffleboot/yandex_fetcher/pkg/service/infra/checker/http"
	root_infra_yandex "github.com/waffleboot/yandex_fetcher/pkg/service/infra/yandex"
	root_interfaces_public_http "github.com/waffleboot/yandex_fetcher/pkg/service/interfaces/public/http"
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

	log, _ := zap.NewProduction()

	cache := cache.NewMemoryCache(log)

	timeout := time.Duration(intConfig("TIMEOUT", 3)) * time.Second

	yandex := yandex_interfaces_private_ipc.NewEndpoint(
		yandex_application.NewService(
			yandex_infra_yandex.NewHttpClient, intConfig("YANDEX_FETCHERS", 1)))

	service := root_interfaces_public_http.NewEndpoint(
		root_application.NewService(
			timeout,
			root_infra_yandex.NewYandex(yandex),
			root_infra_worker.NewBenchmarkSupplier(checkerUrl),
			cache))

	service.AddRoutes(r)

	server := &http.Server{Addr: serviceAddr, Handler: r}

	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalChannel
		server.Close()
	}()

	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func intConfig(name string, def int) int {
	s := os.Getenv(name)
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil || n <= 0 {
		return def
	}
	return n
}
