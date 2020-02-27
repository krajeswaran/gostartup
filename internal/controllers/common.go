package controllers

import (
	"github.com/krajeswaran/gostartup/internal/adapters"
	"github.com/labstack/echo/v4"
)

var dbAdapter = adapters.DBAdapter{}
var cacheAdapter = adapters.CacheAdapter{}

// CommonController Contains common routes for service
type CommonController struct{}

//Status Returns OK for health check
func (a *CommonController) Status(c echo.Context) error {
	return c.JSON(NewApiResponse(Okay, "", nil))
}

//DeepStatus More meaningful health check. Checks all *must-have* dependencies before returning ok
func (a *CommonController) DeepStatus(c echo.Context) error {
	if err := dbAdapter.DeepStatus(); err != nil {
		return c.JSON(NewApiResponse(InternalError, err.Error(), nil))
	}
	if err := cacheAdapter.DeepStatus(); err != nil {
		return c.JSON(NewApiResponse(InternalError, err.Error(), nil))
	}
	return c.JSON(NewApiResponse(Okay, "", nil))
}
