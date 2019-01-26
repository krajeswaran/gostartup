package api

import (
	"context"
	"gostartup/src/worker"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

var onceAPI sync.Once

// Init initializes the API server
func Init() {
	// Server should be initialized only once
	onceAPI.Do(func() {
		log.Info().Msg("Starting api workers")
		worker.Init()
	})

	r := NewApiRouter()
	go func() {
		log.Info().Msg("Initializing Api server")
		if err := r.Start(os.Getenv("GOSTART_API_SERVER_PORT")); err != nil {
			log.Fatal().Err(err).Msg("Unable to bring api server up")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Info().Msg("Got shutdown signal, proceeding with graceful shutdown...")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := r.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Unable to shutdown gracefully")
	}
	worker.Shutdown()
}
