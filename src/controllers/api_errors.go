package controllers

import (
	"fmt"
	"net/http"
	"github.com/labstack/echo"
)

type ApiError struct {
	MsgType int
	MsgInfo MessageInfo
}

type MessageInfo struct {
	Status  int
	Message string
}

// NewApiError to create new error type
func NewApiError(msgType int) ApiError {

	msg, ok := MsgMap[msgType]
	if !ok {
		msg = MessageInfo{http.StatusInternalServerError, "Internal Server Error"}
	}

	apiError := ApiError{
		MsgType: msg.Status,
		MsgInfo: msg,
	}

	return apiError
}

var (
	MsgMap = map[int]MessageInfo{
		InvalidRequest:      MessageInfo{http.StatusBadRequest, "Invalid request format"},
		AccessDenied:        MessageInfo{http.StatusForbidden, "Access denied"},
		InvalidInput:        MessageInfo{http.StatusBadRequest, "Invalid request"},
		AppAccessDenied:     MessageInfo{http.StatusForbidden, "No external access to this private app"},
		InternalError:       MessageInfo{http.StatusInternalServerError, "Our API is unavailable at this moment, please try again later."},
		InsufficientBalance: MessageInfo{http.StatusPaymentRequired, "Insufficient credits on your account. Please recharge/subscribe for more credits."},
		InvalidAuth:         MessageInfo{http.StatusForbidden, "Authorization failed"},
		Queued:         MessageInfo{http.StatusAccepted, "Message Queued"},
		Sent:         MessageInfo{http.StatusAccepted, "Message Sent to Provider"},
		NotFound:         MessageInfo{http.StatusNotFound, "Not found"},
		Okay:         MessageInfo{http.StatusOK, "OK"},
	}
)

const (
	InvalidInput        = iota
	AppAccessDenied     = iota
	AccessDenied        = iota
	InternalError       = iota
	InvalidRequest      = iota
	InsufficientBalance = iota
	InvalidAuth         = iota
	Queued = iota
	NotFound = iota
	Okay = iota
	Sent = iota
)

// Error to format errors
func (c ApiError) Error() string {
	return fmt.Sprintf("%s", c.MsgInfo.Message)
}

func MakeApiResponse(status int, msg string, data interface{}) *echo.Map {
	response := echo.Map{
		"status": status,
	}
	if msg != "" {
		response["message"] = msg
		// TODO add msg links to documentation later
	}
	if data != nil {
		response["data"] = data
	}
	return &response
}
