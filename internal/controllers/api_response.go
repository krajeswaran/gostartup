package controllers

import (
	"github.com/labstack/echo"
	"net/http"
)

//ApiError Denotes API error format
type ApiError struct {
	MsgType int
	MsgInfo MessageInfo
}

//MessageInfo  Standard messaging format for this API
type MessageInfo struct {
	Status  int
	Message string
}

var (
	statusMap = map[int]MessageInfo{
		InvalidRequest:  MessageInfo{http.StatusBadRequest, "invalid request format"},
		AccessDenied:    MessageInfo{http.StatusForbidden, "access denied"},
		InvalidInput:    MessageInfo{http.StatusBadRequest, "invalid request"},
		AppAccessDenied: MessageInfo{http.StatusForbidden, "access to this service denied"},
		InternalError:   MessageInfo{http.StatusInternalServerError, "service is unavailable at this moment, please try again later"},
		InvalidAuth:     MessageInfo{http.StatusForbidden, "authorization failed"},
		NotFound:        MessageInfo{http.StatusNotFound, "requested resource not found"},
		Accepted:        MessageInfo{http.StatusAccepted, "accepted"},
		Okay:            MessageInfo{http.StatusOK, "ok"},
	}
)

const (
	InvalidInput    = iota
	AppAccessDenied = iota
	AccessDenied    = iota
	InternalError   = iota
	InvalidRequest  = iota
	InvalidAuth     = iota
	NotFound        = iota
	Accepted        = iota
	Okay            = iota
)

// NewApiResponse creates API response in a standard format
func NewApiResponse(status int, msg string, data interface{}) (int, *echo.Map) {
	d := statusMap[status]

	r := echo.Map{
		"status": d.Status,
	}
	if msg != "" {
		r["message"] = msg
	} else {
		r["message"] = d.Message
	}
	if data != nil {
		r["data"] = data
	}
	return d.Status, &r
}
