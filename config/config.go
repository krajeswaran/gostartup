package config

import (
	"github.com/spf13/viper"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Init initializes config for the service using a combination of env + config files
func Init() {
	// Default level for this example is info, unless debug flag is present
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// viper: load env vars + local toml config files
	viper.SetEnvPrefix("GOSTARTUP")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// viper: load config files: base file first and then env file
	viper.AddConfigPath("./config")
	viper.SetConfigName("base")
	if err := viper.ReadInConfig(); err != nil {
		log.Panic().Msgf("Can't read base config file! Error: %v", err)
	}

	viper.SetConfigName(viper.GetString("env"))
	viper.SetConfigType("toml")
	if err := viper.MergeInConfig(); err != nil {
		log.Warn().Msgf("Can't read config file: %s", viper.GetString("env"))
	}

	if viper.GetBool("general.debug") {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
	} else {
		// use unix timestamp for non-debug env
		zerolog.TimeFieldFormat = ""
	}
}
