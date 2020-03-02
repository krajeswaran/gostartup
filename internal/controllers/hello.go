package controllers

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/hashicorp/go-multierror"
	"github.com/krajeswaran/gostartup/internal/usecases"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"time"
)

// HelloController Ridiculous hello-world controller
type HelloController struct{}

var helloRepo = usecases.HelloRepo{}

// SayHello Returns SayHello <user.name> along with some API stats
func (h *HelloController) SayHello(c echo.Context) error {
	// fetch userId from JWT token
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userId := claims["id"].(string)

	userName, err := helloRepo.FetchUserName(userId)
	if err != nil {
		err = multierror.Flatten(err)
		log.Error().Msgf("failed to fetch username for userId: %s, error %v", userId, err)
		return c.JSON(NewApiResponse(InternalError, err.Error(), nil))
	}

	// log everything preferably in one layer: controller layer
	log.Debug().Msgf("fetched username: %s for userId: %s", userName, userId)

	return c.JSON(NewApiResponse(Okay, "", userName))
}

// GetStats returns rudimentary stats for hello api
func (h *HelloController) GetStats(c echo.Context) error {
	stats, err := helloRepo.GetApiStats()
	if err != nil {
		err = multierror.Flatten(err)
		log.Error().Msgf("failed to fetch api stats for hello api: %v", err)
		return c.JSON(NewApiResponse(InternalError, err.Error(), nil))
	}

	return c.JSON(NewApiResponse(Okay, "", stats))
}

// CreateUser creates a user and returns a JWT token for the user
func (h *HelloController) CreateUser(name string) (string, error) {
	user, err := helloRepo.CreateUser(name)
	if err != nil {
		err = multierror.Flatten(err)
		log.Error().Msgf("failed to create user %s, error %v", name, err)
		return "", err
	}

	// convert user to JWT
	token := jwt.New(jwt.SigningMethodHS256)
	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = fmt.Sprintf("%d", user.ID)
	claims["exp"] = time.Now().Add(time.Hour * viper.GetDuration("general.jwt_validity_period")).Unix()
	t, err := token.SignedString([]byte(viper.GetString("secret_key")))
	if err != nil {
		log.Error().Msgf("failed to create jwt for user %s, error %v", name, err)
		return "", err
	}

	return t, nil
}
