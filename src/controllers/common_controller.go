package controllers

import (
	"github.com/labstack/echo"
	"gostartup/src/adapters"
	"gostartup/src/config"
	"net/http"
)

var dbAdapter = adapters.DBAdapter{}
var cacheAdapter = adapters.CacheAdapter{}

type CommonController struct{}

func (a *CommonController) ServiceStatus(c echo.Context) error {
	c.String(http.StatusOK, "OK")
	return nil
}

func (a *CommonController) ServiceDeepStatus(c echo.Context) error {
	if err := dbAdapter.DeepStatus(); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	if err := cacheAdapter.DeepStatus(); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, "OK")
}

func (a *CommonController) ServiceAuthentication(username, password string, c echo.Context) (bool, error) {
	conf := config.GetConfig()
	if username == conf.GetString("auth.auth_id") &&
		password == conf.GetString("auth.auth_token") {
		return true, nil
	}
	return false, nil
}

func (a *CommonController) ApiResponse(status int, data interface{}, c echo.Context) error {
	apiError := NewApiError(status)
	response := MakeApiResponse(apiError.MsgType, apiError.MsgInfo.Message, data)
	return c.JSON(apiError.MsgType, response)
}