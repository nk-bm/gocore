package ginmw

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nk-bm/gocore/gincore/response"
	"github.com/nk-bm/gocore/goutils"
)

func AuthMW(secretKey string, idKey, authType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := goutils.ExtractGinToken(c, "Authorization", "Bearer")
		if err != nil {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		id, err := goutils.ExtractIDFromJWT(secretKey, token, idKey, authType)
		if err != nil {
			response.Error(c, err, http.StatusUnauthorized)
			c.Abort()
			return
		}

		c.Set(idKey, id)
		c.Next()
	}
}
