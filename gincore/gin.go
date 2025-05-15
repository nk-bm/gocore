package gincore

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nk-bm/gocore/gincore/ginmw"
	"go.uber.org/zap"
)

type Options struct {
	EnableCORS                bool `env:"GIN_ENABLE_CORS; default:true"`
	DisableRequestTime        bool `env:"GIN_DISABLE_REQUEST_TIME; default:false"`
	DisableHealthCheckHandler bool `env:"GIN_DISABLE_HEALTH_CHECK_HANDLER; default:false"`
}

type Config struct {
	APIPath string `env:"GIN_API_PATH; default:/api/v1"`
	Port    int    `env:"GIN_PORT; default:8080"`
	Host    string `env:"GIN_HOST; default:0.0.0.0"`
	Options Options
}

type Server struct {
	config    *Config
	Router    *gin.Engine
	APIRouter *gin.RouterGroup
	logger    *zap.Logger
}

type Route struct {
	Method  string
	Path    string
	Handler gin.HandlerFunc
}

func NewServer(config Config, logger *zap.Logger) *Server {
	router := gin.Default()

	if config.Options.EnableCORS {
		router.Use(ginmw.EnableCORS())
	}
	if !config.Options.DisableRequestTime {
		router.Use(ginmw.RequestTimeMW())
	}
	if !config.Options.DisableHealthCheckHandler {
		router.GET("/health", HealthCheckHandler)
	}
	return &Server{
		config:    &config,
		logger:    logger,
		Router:    router,
		APIRouter: router.Group(config.APIPath),
	}
}

func (s *Server) RegisterRoutes(routes []Route) {
	for _, route := range routes {
		s.RegisterRoute(route)
	}
}

func (s *Server) RegisterRoute(route Route) {
	s.APIRouter.Handle(route.Method, route.Path, route.Handler)
}

func (s *Server) Start() error {
	if s.config.Port == 0 {
		return fmt.Errorf("port is not set")
	}
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	s.logger.Info("Server started", zap.String("addr", addr))
	return s.Router.Run(addr)
}
