package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi"

	cache "github.com/waffleboot/playstation_buy/pkg/cache"

	worker_application "github.com/waffleboot/playstation_buy/pkg/worker/application"
	worker_interfaces_private_http "github.com/waffleboot/playstation_buy/pkg/worker/interfaces/private/http"
)

func run(args []string) int {

	checkerAddr := os.Getenv("CHECKER_ADDR")
	if checkerAddr == "" {
		return 1
	}

	log.Printf("Starting service on %s", checkerAddr)

	if err := startServer(checkerAddr); err != nil {
		log.Println(err)
		return 2
	}
	return 0
}

func startServer(checkerAddr string) error {

	r := chi.NewRouter()

	ctx := context.Background()

	cache := &cache.MemoryCache{}

	timeout := 3 * time.Second

	worker := worker_interfaces_private_http.NewEndpoint(
		worker_application.NewService(ctx, cache, checkersCount(10), timeout))

	worker.AddRoutes(r)

	server := &http.Server{Addr: checkerAddr, Handler: r}
	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func checkersCount(def int) int {
	s := os.Getenv("CHECKERS_COUNT")
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil || n <= 0 {
		return def
	}
	return n
}
