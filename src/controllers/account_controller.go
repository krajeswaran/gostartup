package controllers

import (
	"errors"
	"github.com/labstack/echo"
	"github.com/rs/zerolog/log"
	"gostartup/src/usecases"
)

// SmsController - Struct to logically bind all the CommonController functions
type AccountController struct{}

var accountUsecase = usecases.AccountUsecase{}

// user auth
func (a *AccountController) UserAuthentication(username, password string, c echo.Context) (bool, error) {
	// TODO jwt

	// validate with DB
	account, err := accountUsecase.AuthenticateUser(username, password)
	if err != nil {
		log.Error().Err(err).Msgf("INTERNAL_AUTH_ERR for id: %s", password)
		return false, errors.New("user auth failure")
	}

	// everything good to go
	c.Set("account", account)
	return true, nil
}
