package http

import (
	"bytes"
	"encoding/json"
	net_http "net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	errReasonKey          = "err_reason"
	ResponseStatusFailed  = 0
	ResponseStatusSuccess = 1
)

// HTTPWriter type of HTTP response writer
type HTTPWriter struct{}

// NewHTTPWriter initiates and returns a new instance of HTTP response writer
func NewHTTPWriter() *HTTPWriter {
	return &HTTPWriter{}
}

type httpResponseJSON struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
	Error  string      `json:"error"`
}

func (h *HTTPWriter) GinHTTPResponseWriter(ctx *gin.Context, data interface{}, err error, httpStatus ...int) {
	var (
		resRaw httpResponseJSON
		code   int
	)
	if err != nil {
		resRaw.Status = ResponseStatusFailed
		resRaw.Data = data
		resRaw.Error = err.Error()
		code = net_http.StatusInternalServerError
	} else {
		resRaw.Status = ResponseStatusSuccess
		resRaw.Data = data
		code = net_http.StatusOK
	}

	if len(httpStatus) > 0 {
		code = httpStatus[0]
	}

	ctx.Writer.Header().Add("Content-Type", "application/json")
	ctx.Writer.Header().Add("Status-Code", strconv.Itoa(code))

	res, err := jsonMarshal(resRaw)
	if err != nil {
		resError := httpResponseJSON{
			Error: err.Error(),
		}
		res, _ = jsonMarshal(resError)
		ctx.Writer.Write(res)
		return
	}

	ctx.Writer.Write(res)
}

func jsonMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}
