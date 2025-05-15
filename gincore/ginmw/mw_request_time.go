package ginmw

import (
	"time"

	"github.com/gin-gonic/gin"
)

func RequestTimeMW() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("requestTime", time.Now())
		c.Next()
	}
}
