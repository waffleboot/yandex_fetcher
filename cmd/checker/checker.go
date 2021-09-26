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

	"github.com/go-chi/chi/middleware"

	cache "github.com/waffleboot/yandex_fetcher/pkg/cache"

	checker_app "github.com/waffleboot/yandex_fetcher/pkg/checker/application"
	checker_infra_service "github.com/waffleboot/yandex_fetcher/pkg/checker/infra/service/http"
	checker_http "github.com/waffleboot/yandex_fetcher/pkg/checker/inter/private/http"
)

func run(args []string) int {

	checkerAddr := os.Getenv("CHECKER_ADDR")
	if checkerAddr == "" {
		return 1
	}

	serviceUrl := os.Getenv("SERVICE_URL")

	redisAddr := os.Getenv("REDIS_URL")

	log.Printf("Starting service on %s", checkerAddr)

	if err := startServer(checkerAddr, serviceUrl, redisAddr); err != nil {
		log.Println(err)
		return 2
	}
	return 0
}

func startServer(checkerAddr, serviceUrl, redisAddr string) error {

	r := chi.NewRouter()

	log, _ := zap.NewProduction()

	var cach cache.Cache

	if redisAddr != "" {
		cach = cache.NewRedisCache(redisAddr, log)
	} else {
		cach = cache.NewMemoryCache(log)
	}

	timeout := time.Duration(intConfig("TIMEOUT", 3)) * time.Second

	service := checker_infra_service.NewInitialService(serviceUrl)

	checker := checker_app.NewService(cach, service, intConfig("CHECKERS_COUNT", 10), timeout)

	checker_http.AddRoutes(checker, r)

	r.Mount("/debug", middleware.Profiler())

	server := &http.Server{Addr: checkerAddr, Handler: r}

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
