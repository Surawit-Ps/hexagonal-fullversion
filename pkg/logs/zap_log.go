package logs

import (
	"go.uber.org/zap"
)

type ZapLogger struct {
	logger *zap.Logger
}

func NewZapLogger() (*ZapLogger, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	return &ZapLogger{logger: logger}, nil
}

func (zl *ZapLogger) Info(msg string, fields ...zap.Field) {
	zl.logger.Info(msg, fields...)
}

func (zl *ZapLogger) Error(msg string, fields ...zap.Field) {
	zl.logger.Error(msg, fields...)
}

func (zl *ZapLogger) Debug(msg string, fields ...zap.Field) {
	zl.logger.Debug(msg, fields...)
}

func (zl *ZapLogger) Warn(msg string, fields ...zap.Field) {
	zl.logger.Warn(msg, fields...)
}

func (zl *ZapLogger) Sync() {
	zl.logger.Sync()
}




