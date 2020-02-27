package controllers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/hashicorp/go-multierror"
	"github.com/krajeswaran/gostartup/internal/usecases"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// HelloController Ridiculous hello-world controller
type HelloController struct{}
var helloRepo = usecases.HelloRepo{}

// SayHello Returns SayHello <user.name> along with some API stats
func (h *HelloController) SayHello(c echo.Context) error {
	// fetch userId from JWT token
	token := c.Get("user").(jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userId := claims["userId"].(string)

	userName, err := helloRepo.FetchUserName(userId)
	if err != nil {
		log.Error().Msgf("failed to fetch username for userId: %s, error %v", userId, multierror.Flatten(err))
		return c.JSON(NewApiResponse(InternalError, multierror.Flatten(err).Error(), nil))
	}

	// log everything preferably in one layer: controller layer
	log.Debug().Msgf("fetched username: %s for userId: %s", userName, userId)

	return c.JSON(NewApiResponse(Okay, "", userName))
}

//GetStats returns rudimentary stats for hello api
func (h *HelloController) GetStats(c echo.Context) error {
	stats, err := helloRepo.GetApiStats()
	if err != nil {
		log.Error().Msgf("failed to fetch api stats for hello api: %v", multierror.Flatten(err))
		return c.JSON(NewApiResponse(InternalError, multierror.Flatten(err).Error(), nil))
	}

	return c.JSON(NewApiResponse(Okay, "", stats))
}

