package controllers

import (
	"github.com/labstack/echo"
	"github.com/rs/zerolog/log"
	"gostartup/src/usecases"
	"gostartup/src/worker"
	"strings"
)

type CallbackController struct{}

var common = CommonController{}

// HandleCallback handles callback responses from for a subscribed URI
func (a *CallbackController) HandleCallback(c echo.Context) error {
	// translate payload to provider callback response
	if err := c.Bind(&m); err != nil {
		log.Error().Err(err).Msgf("Unable to bind to input request")
		return common.ApiResponse(InvalidInput, nil, c)
	}

	// async callback to user
	callbackJob := usecases.CallbackJob{
		MsgId:    c.Param("msg_id"),
		Provider: provider,
		Payload:  m,
	}
	worker.Execute(callbackJob)

	// we are good with callback, so respond in kind
	return common.ApiResponse(Okay, nil, c)
}
