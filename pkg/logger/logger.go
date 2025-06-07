package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerConfig struct {
	Mode  string
	Level string
}

func InitLogger(cfg LoggerConfig) (*zap.SugaredLogger, error) {
	// Set Logger Level
	logLevel, exists := loggerLevelMap[cfg.Level]
	if !exists {
		logLevel = zapcore.DebugLevel
	}

	// Set Logger Mode
	var config zap.Config
	if cfg.Mode == "development" {
		config = zap.NewDevelopmentConfig()
		config.Encoding = "console"
		config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006/01/02 15:04:05")
	} else if cfg.Mode == "production" {
		config = zap.NewProductionConfig()
		config.Encoding = "json"
	}

	config.Level = zap.NewAtomicLevelAt(logLevel)
	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	sugarLogger := logger.Sugar()

	return sugarLogger, nil
}

var loggerLevelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}
