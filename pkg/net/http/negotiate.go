package http

import (
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

// OFFEREDJSON []string{gin.MIMEJSON}
var OFFEREDJSON = []string{gin.MIMEJSON}

// OFFEREDHTML []string{gin.MIMEHTML}
var OFFEREDHTML = []string{gin.MIMEHTML}

// OFFEREDXML []string{gin.MIMEXML}
var OFFEREDXML = []string{gin.MIMEXML}

// OFFEREDYAML = []string{gin.MIMEYAML}
var OFFEREDYAML = []string{gin.MIMEYAML}

// OFFEREDBACKEND []string{gin.MIMEJSON, gin.MIMEXML, gin.MIMEYAML}
var OFFEREDBACKEND = []string{gin.MIMEJSON, gin.MIMEXML, gin.MIMEYAML}

// OFFEREDALL []string{gin.MIMEHTML, gin.MIMEJSON, gin.MIMEXML, gin.MIMEYAML}
var OFFEREDALL = []string{gin.MIMEHTML, gin.MIMEJSON, gin.MIMEXML, gin.MIMEYAML}

// Negotiate calls different Render according acceptable Accept format.
func Negotiate(c *gin.Context, statuscode int, config gin.Negotiate) {
	negotiate := gin.Negotiate{
		Offered:  config.Offered,
		HTMLName: config.HTMLName,
		HTMLData: config.HTMLData,
		JSONData: ResponseSuccessWithData(c, config.JSONData),
		XMLData:  ResponseSuccessWithData(c, config.XMLData),
		YAMLData: ResponseSuccessWithData(c, config.YAMLData),
		Data:     config.Data,
	}

	if negotiate.HTMLName == "" {
		negotiate.Data = ResponseSuccessWithData(c, config.Data)
	}

	c.Negotiate(statuscode, negotiate)
}

// NegotiateData calls different Render according acceptable Accept format.
func NegotiateData(c *gin.Context, statuscode int, code string, message string, data any) {
	offered := OFFEREDBACKEND

	response := NewResponse(code, message, data)

	response.RequestID = requestid.Get(c)

	negotiate := gin.Negotiate{Offered: offered, Data: response}

	c.Negotiate(statuscode, negotiate)
}

// NegotiateResponse calls different Render according acceptable Accept format.
func NegotiateResponse(c *gin.Context, statuscode int, response Response) {
	offered := OFFEREDBACKEND

	response.RequestID = requestid.Get(c)

	negotiate := gin.Negotiate{Offered: offered, Data: response}

	c.Negotiate(statuscode, negotiate)
}
