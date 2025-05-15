package gincore

import (
	"github.com/gin-gonic/gin"
	"github.com/nk-bm/gocore/gincore/static"
	"github.com/nk-bm/gocore/gotypes"
)

func CtxTMAInitData(c *gin.Context) (*gotypes.TelegramMiniAppInitData, bool) {
	initData, ok := c.Value(static.TMA_INIT_DATA).(*gotypes.TelegramMiniAppInitData)
	return initData, ok
}
