package main

import (
	"github.com/gin-gonic/gin"
	"github.com/nk-bm/gocore/gincore"
	"github.com/nk-bm/gocore/gincore/response"
	"github.com/nk-bm/gocore/gocore"
	"go.uber.org/zap"
)

func HelloWorldHandler(c *gin.Context) {
	response.Success(c, "Hello, World!")
}

func main() {
	service, err := gocore.NewDefaultApp("example_service")
	if err != nil {
		panic(err)
	}

	service.L.Info("Service initialized", zap.String("service_name", service.Name), zap.String("api_path", service.GinServer.APIRouter.BasePath()))

	service.GinServer.RegisterRoute(gincore.Route{
		Method:  "GET",
		Path:    "/",
		Handler: HelloWorldHandler,
	})

	service.Start()
}
