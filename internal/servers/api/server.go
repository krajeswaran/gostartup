package api

import (
	"context"
	"github.com/krajeswaran/gostartup/internal/controllers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"time"

	"github.com/rs/zerolog/log"
)

func apiRoutes() *echo.Echo {
	router := echo.New()

	// some default routing preferences
	router.Pre(middleware.AddTrailingSlash())
	router.Use(middleware.BodyLimit("5M"))
	router.Use(middleware.Secure())
	router.Use(middleware.Logger())
	router.Use(middleware.Recover())
	router.Use(middleware.Gzip())
	router.Use(middleware.CSRF())
	if viper.GetBool("general.enable_cors") {
		router.Use(middleware.CORS())
	}

	// common
	commonController := new(controllers.CommonController)
	// swagger:route
	router.GET("/status/", commonController.Status)
	// swagger:route
	router.GET("/status/full/", commonController.DeepStatus)

	// hello service
	apiV1 := router.Group("/api/v1")
	helloGrp := apiV1.Group("/hello", middleware.JWTWithConfig(
		middleware.JWTConfig{
			SigningKey: []byte(viper.GetString("secret_key")),
		}))
	helloController := new(controllers.HelloController)
	helloGrp.GET("/", helloController.SayHello)
	helloGrp.GET("/stats/", helloController.GetStats)

	return router
}

// Init initializes the API servers
func Init() {
	r := apiRoutes()
	go func() {
		log.Info().Msg("Initializing Api servers")
		if err := r.Start(viper.GetString("API_ADDRESS")); err != nil {
			log.Fatal().Err(err).Msg("Unable to bring api servers up")
		}

	}()

	// Wait for interrupt signal to gracefully shutdown the servers
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Info().Msg("Caught SIGINT, proceeding with graceful shutdown...")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := r.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Unable to shutdown incoming requests gracefully")
	}
}
