/*
Package config contains all the configurations for the service
*/
package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Init initializes config for the service from env variables
func Init() {
	// Default level for this example is info, unless debug flag is present
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if err != nil {
		log.Fatal().Err(err).Msg("Error on parsing configuration.")
	}

	env := os.Getenv("GOSTART_ENV")
	if "" == env {
		env = "development"
	}

	godotenv.Load(".env." + env + ".local")
	if "test" != env {
		godotenv.Load(".env.local")
	}
	godotenv.Load(".env." + env)
	godotenv.Load() // The Original .env

	if os.Getenv("GOSTART_DEBUG") == "1" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	} else {
		// use unix timestamp for non-debug env
		zerolog.TimeFieldFormat = ""
	}
}
