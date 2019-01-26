package api

import (
	"gostartup/src/controllers"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// NewApiRouter for routing requests
func NewApiRouter() *echo.Echo {
	router := echo.New()
	router.Use(middleware.Logger())
	router.Use(middleware.Recover())

	// common
	commonController := new(controllers.CommonController)
	router.GET("/status/", commonController.ServiceStatus)
	router.GET("/deepstatus/", commonController.ServiceDeepStatus)

	// account
	acctController := new(controllers.AccountController)
	accountService := router.Group("/account")
	accountServiceV1 := accountService.Group("/v1")
	accountServiceV1.Group("/account",
		middleware.BasicAuth(commonController.ServiceAuthentication),
		middleware.JWTWithConfig())


	// internal - balance
	balanceController := new(controllers.BalanceController)
	balanceService := router.Group("/balance")
	balanceServiceVersion := balanceService.Group("/v1")
	balance := balanceServiceVersion.Group("/balance",
		middleware.BasicAuth(commonController.ServiceAuthentication))
	balance.POST("/", balanceController.EditBalance)
	balance.GET("/:acct_id/", balanceController.GetBalance)

	return router
}
