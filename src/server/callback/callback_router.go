package callback

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gostartup/src/controllers"
)

// NewApiRouter for routing requests
func NewCallbackRouter() *echo.Echo {
	router := echo.New()
	router.Use(middleware.Logger())
	router.Use(middleware.Recover())

	commonController := new(controllers.CommonController)
	router.GET("/status/", commonController.ServiceStatus)
	router.GET("/deepstatus/", commonController.ServiceDeepStatus)

	callbackController := new(controllers.CallbackController)
	callbackService := router.Group("/callback")
	callbackServiceV1 := callbackService.Group("/v1")
	callbackServiceV1.POST("/:provider/message/:msg_id", callbackController.MessageCallback)
	callbackServiceV1.POST("/:provider/number/:number_id", callbackController.NumberCallback)

	return router
}
