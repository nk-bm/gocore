package gincore

import (
	"github.com/gin-gonic/gin"
	"github.com/nk-bm/gocore/gincore/response"
)

type HealthCheckResponse struct {
	OK bool `json:"ok"`
}

func HealthCheckHandler(c *gin.Context) {
	response.Success(c, HealthCheckResponse{OK: true})
}
