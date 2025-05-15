package gocore

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitGlobalLogger(isProd bool) error {
	var config zap.Config
	if isProd {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	logger, err := config.Build()
	if err != nil {
		return err
	}

	zap.ReplaceGlobals(logger)
	return nil
}

func GetLogger() *zap.Logger {
	return zap.L()
}

func SetGlobalLogger(logger *zap.Logger) {
	zap.ReplaceGlobals(logger)
}
