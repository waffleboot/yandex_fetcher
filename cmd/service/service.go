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

	yandex_app "github.com/waffleboot/yandex_fetcher/pkg/yandex/application"
	yandex_yandex "github.com/waffleboot/yandex_fetcher/pkg/yandex/infra/yandex"
	yandex_ipc "github.com/waffleboot/yandex_fetcher/pkg/yandex/inter/private/ipc"

	service_app "github.com/waffleboot/yandex_fetcher/pkg/service/application"
	service_checker "github.com/waffleboot/yandex_fetcher/pkg/service/infra/checker/http"
	service_http "github.com/waffleboot/yandex_fetcher/pkg/service/inter/public/http"
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

	redisAddr := os.Getenv("REDIS_URL")

	log.Printf("Starting service on %s", serviceAddr)

	if err := startServer(serviceAddr, checkerUrl, redisAddr); err != nil {
		log.Println(err)
		return 2
	}
	return 0
}

func startServer(serviceAddr, checkerUrl, redisAddr string) error {

	r := chi.NewRouter()

	log, _ := zap.NewProduction()

	var cach cache.Cache

	if redisAddr != "" {
		cach = cache.NewRedisCache(redisAddr, log)
	} else {
		cach = cache.NewMemoryCache(log)
	}

	timeout := time.Duration(intConfig("TIMEOUT", 3)) * time.Second

	yandex := yandex_ipc.NewEndpoint(
		yandex_app.NewService(
			yandex_yandex.NewHttpClient, intConfig("YANDEX_FETCHERS", 1)))

	service := service_http.NewEndpoint(
		service_app.NewService(
			timeout,
			yandex.AddQuery,
			service_checker.NewBenchmarkSupplier(checkerUrl),
			cach))

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
