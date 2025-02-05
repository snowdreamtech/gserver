package http

import (
	"time"

	"github.com/gin-gonic/gin"
)

// Response Response
type Response struct {
	Version   string `json:"version" xml:"version" yaml:"version" schema:"version"`
	Code      string `json:"code" xml:"code" yaml:"code" schema:"code"`
	Message   string `json:"message" xml:"message" yaml:"message" schema:"message"`
	Data      any    `json:"data" xml:"data" yaml:"data" schema:"data"`
	Timestamp int64  `json:"timestamp" xml:"timestamp" yaml:"timestamp" schema:"timestamp"`
	RequestID string `json:"requestid" xml:"requestid" yaml:"requestid" schema:"requestid"`
}

// NewResponse NewResponse
func NewResponse(code string, message string, data any) Response {
	if nil == data {
		return Response{
			Version:   "0.1",
			Code:      code,
			Message:   message,
			Data:      struct{}{},
			Timestamp: time.Now().Unix(),
		}
	}

	return Response{
		Version:   "0.1",
		Code:      code,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().Unix(),
	}
}

// ResponseSuccess ResponseSuccess
func ResponseSuccess(c *gin.Context) Response {
	return NewResponse("SUCCESS", "SUCCESS", nil)
}

// ResponseSuccessWithData ResponseSuccessWithData
func ResponseSuccessWithData(c *gin.Context, data any) Response {
	return NewResponse("SUCCESS", "SUCCESS", data)
}

// ResponseFailure ResponseFailure
func ResponseFailure(c *gin.Context) Response {
	return NewResponse("FAILURE", "FAILURE", nil)
}

// ResponseFailureWithData ResponseFailureWithData
func ResponseFailureWithData(c *gin.Context, data any) Response {
	return NewResponse("FAILURE", "FAILURE", data)
}
