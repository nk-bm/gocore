package ginmw

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nk-bm/gocore/gincore/response"
	"github.com/nk-bm/gocore/gincore/static"
	"github.com/nk-bm/gocore/gotypes"
	"github.com/nk-bm/gocore/goutils"
	initdata "github.com/telegram-mini-apps/init-data-golang"
)

func TelegramMiniAppAuthMW(telegramBotToken string, expIn time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := goutils.ExtractGinToken(c, static.TMA_TOKEN_KEY, "initdata")
		if err != nil {
			response.ErrorString(c, static.TMA_TOKEN_KEY+" required", http.StatusUnauthorized)
			c.Abort()
			return
		}

		if err = initdata.Validate(data, telegramBotToken, expIn); err != nil {
			response.ErrorString(c, "not valid initdata", http.StatusUnauthorized)
			c.Abort()
			return
		}

		initData, err := initdata.Parse(data)
		if err != nil {
			response.ErrorString(c, "err on parse initdata", http.StatusUnauthorized)
			c.Abort()
			return
		}

		c.Set(static.TMA_INIT_DATA, gotypes.TelegramMiniAppInitData(initData))
		c.Next()
	}
}
