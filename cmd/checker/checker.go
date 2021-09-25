package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"go.uber.org/zap"

	cache "github.com/waffleboot/playstation_buy/pkg/cache"

	root_infra_service "github.com/waffleboot/playstation_buy/pkg/checker/infra/service/http"

	worker_application "github.com/waffleboot/playstation_buy/pkg/checker/application"
	worker_interfaces_private_http "github.com/waffleboot/playstation_buy/pkg/checker/interfaces/private/http"
)

func run(args []string) int {

	checkerAddr := os.Getenv("CHECKER_ADDR")
	if checkerAddr == "" {
		return 1
	}

	serviceUrl := os.Getenv("SERVICE_URL")

	log.Printf("Starting service on %s", checkerAddr)

	if err := startServer(checkerAddr, serviceUrl); err != nil {
		log.Println(err)
		return 2
	}
	return 0
}

func startServer(checkerAddr, serviceUrl string) error {

	r := chi.NewRouter()

	log, _ := zap.NewProduction()

	cache := cache.NewMemoryCache(log)

	timeout := time.Duration(intConfig("TIMEOUT", 3)) * time.Second

	service := root_infra_service.NewInitialService(serviceUrl)

	worker := worker_interfaces_private_http.NewEndpoint(
		worker_application.NewService(cache, service, intConfig("CHECKERS_COUNT", 10), timeout))

	worker.AddRoutes(r)

	server := &http.Server{Addr: checkerAddr, Handler: r}
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
