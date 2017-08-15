package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(level zapcore.Level) (logger *zap.Logger) {
	opt := zap.NewDevelopmentConfig()
	opt.Development = false
	opt.DisableCaller = true
	opt.DisableStacktrace = true
	opt.Sampling = &zap.SamplingConfig{
		Initial:    100,
		Thereafter: 100,
	}
	opt.Level.SetLevel(level)
	opt.OutputPaths = []string{"stdout"}
	opt.ErrorOutputPaths = []string{"stderr"}
	logger, _ = opt.Build()
	return
}
