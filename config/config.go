package config

import (
	"github.com/spf13/viper"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Init initializes config for the service using a combination of env + config files
func Init() {
	// Default level for this example is info, unless debug flag is present
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// viper: load env vars + local toml config files
	viper.SetEnvPrefix("GOSTARTUP")
	viper.AutomaticEnv()

	// viper: load config files: base file first and then env file
	viper.AddConfigPath("./config")
	viper.SetConfigName("base")
	if err:= viper.ReadInConfig(); err != nil {
		log.Panic().Msgf("Can't read base config file! Error: %v", err)
	}

	viper.SetConfigName(viper.GetString("env"))
	viper.SetConfigType("toml")
	if err:= viper.MergeInConfig(); err != nil {
		log.Warn().Msgf("Can't read config file: %s", viper.GetString("env"))
	}

	if viper.GetBool("debug") {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	} else {
		// use unix timestamp for non-debug env
		zerolog.TimeFieldFormat = ""
	}
}
