package response

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success        bool           `json:"success"`
	Data           any            `json:"data"`
	Error          *ErrorResponse `json:"error,omitempty"`
	ResponseTimeMs *float64       `json:"response_time_ms,omitempty"`
}

type ErrorResponse struct {
	Code        int     `json:"code"`
	Type        *string `json:"type,omitempty"`
	Error       string  `json:"error"`
	Description string  `json:"description,omitempty"`
}

func NewResponse() *Response {
	return &Response{
		Success: true,
		Data:    struct{}{},
	}
}

func (r *Response) SetSuccess(data any) *Response {
	r.Success = true
	r.Data = data
	return r
}

func (r *Response) SetErrorString(error string, status int) *Response {
	r.Success = false
	r.Error = &ErrorResponse{
		Code:        status,
		Error:       error,
		Description: "",
	}
	return r
}

func (r *Response) SetErrorType(errorType string) *Response {
	if r.Error == nil {
		r.Error = &ErrorResponse{}
	}
	r.Error.Type = &errorType
	return r
}

func (r *Response) SetErrorDescription(description string) *Response {
	if r.Error == nil {
		r.Error = &ErrorResponse{}
	}
	r.Error.Description = description
	return r
}

func (r *Response) SetError(error error, statusCode int) *Response {
	return r.SetErrorString(error.Error(), statusCode)
}

func (r *Response) Respond(c *gin.Context) {
	if r.Error != nil {
		c.JSON(r.Error.Code, r.Error)
		return
	}

	requestTime := c.GetTime("requestTime")
	if !requestTime.IsZero() {
		duration := time.Since(requestTime)
		ms := float64(duration.Nanoseconds()) / float64(time.Millisecond)
		r.ResponseTimeMs = &ms
	}

	c.JSON(http.StatusOK, r)
}
