package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Error(c *gin.Context, err error, status int) {
	NewResponse().SetError(err, status).Respond(c)
}

func ErrorString(c *gin.Context, err string, status int) {
	NewResponse().SetErrorString(err, status).Respond(c)
}

func Success(c *gin.Context, data any) {
	NewResponse().SetSuccess(data).Respond(c)
}

func SuccessWithStatus(c *gin.Context, data any, status int) {
	NewResponse().SetSuccess(data).Respond(c)
}

func Unauthorized(c *gin.Context) {
	NewResponse().SetErrorString("Unauthorized", http.StatusUnauthorized).Respond(c)
}

func Forbidden(c *gin.Context) {
	NewResponse().SetErrorString("Forbidden", http.StatusForbidden).Respond(c)
}

func NotFound(c *gin.Context) {
	NewResponse().SetErrorString("Not Found", http.StatusNotFound).Respond(c)
}

func NotFoundWithMessage(c *gin.Context, message string) {
	NewResponse().SetErrorString(message, http.StatusNotFound).Respond(c)
}

func BadRequest(c *gin.Context) {
	NewResponse().SetErrorString("Bad Request", http.StatusBadRequest).Respond(c)
}

func BadRequestWithMessage(c *gin.Context, message string) {
	NewResponse().SetErrorString(message, http.StatusBadRequest).Respond(c)
}

func InternalServerError(c *gin.Context) {
	NewResponse().SetErrorString("Internal Server Error", http.StatusInternalServerError).Respond(c)
}

func Conflict(c *gin.Context) {
	NewResponse().SetErrorString("Conflict", http.StatusConflict).Respond(c)
}

func ConflictWithMessage(c *gin.Context, message string) {
	NewResponse().SetErrorString(message, http.StatusConflict).Respond(c)
}

func TooManyRequests(c *gin.Context) {
	NewResponse().SetErrorString("Too Many Requests", http.StatusTooManyRequests).Respond(c)
}

func UnprocessableEntity(c *gin.Context, message string) {
	NewResponse().SetErrorString(message, http.StatusUnprocessableEntity).Respond(c)
}
