package callback

import (
	"context"
	"github.com/rs/zerolog/log"
	"gostartup/src/worker"
	"os"
	"os/signal"
	"sync"
	"time"
)

var onceCallback sync.Once

// Init initializes callback server
func Init() {
	// Server should be initialized only once
	onceCallback.Do(func() {
		// start callback workers
		log.Info().Msg("Starting callback workers")
		worker.Init()
	})

	// Process requests
	r := NewCallbackRouter()
	go func() {
		log.Info().Msg("Initializing callback server")
		if err := r.Start(os.Getenv("GOSTART_CALLBACK_SERVER_PORT")); err != nil {
			log.Fatal().Err(err).Msg("Unable to bring callback service up")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Info().Msg("Got shutdown signal, proceeding with graceful shutdown...")
	worker.Shutdown()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := r.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Unable to shutdown gracefully")
	}
}
