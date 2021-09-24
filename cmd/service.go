package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"

	public_http "github.com/waffleboot/playstation_buy/pkg/interfaces/public/http"
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
	public_http.AddRoutes(r)
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
