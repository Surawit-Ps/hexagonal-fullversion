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

func (zl *ZapLogger) Info(msg string, fields ...interface{}) {
	var zapFields []zap.Field

	for i := 0; i < len(fields); i += 2 {
		key := fields[i].(string)
		value := fields[i+1]
		zapFields = append(zapFields, zap.Any(key, value))
	}

	zl.logger.Info(msg, zapFields...)
}

func (zl *ZapLogger) Error(msg string, fields ...interface{}) {
	var zapFields []zap.Field
	for i := 0; i < len(fields); i += 2 {
		key := fields[i].(string)
		value := fields[i+1]
		zapFields = append(zapFields, zap.Any(key, value))
	}
	zl.logger.Error(msg, zapFields...)
}

func (zl *ZapLogger) Debug(msg string, fields ...interface{}) {
	var zapFields []zap.Field
	for i := 0; i < len(fields); i += 2 {
		key := fields[i].(string)
		value := fields[i+1]
		zapFields = append(zapFields, zap.Any(key, value))
	}
	zl.logger.Debug(msg, zapFields...)
}

func (zl *ZapLogger) Warn(msg string, fields ...interface{}) {
	var zapFields []zap.Field
	for i := 0; i < len(fields); i += 2 {
		key := fields[i].(string)
		value := fields[i+1]
		zapFields = append(zapFields, zap.Any(key, value))
	}
	zl.logger.Warn(msg, zapFields...)
}

func (zl *ZapLogger) Sync() {
	zl.logger.Sync()
}



