package logger

import (
	uberZap "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
)

type zap struct {
	logger           *uberZap.SugaredLogger
	prefixKeysValues []interface{}
}

func newZap(level Level, serviceName string) *zap {
	loggerConfig := uberZap.NewProductionConfig()

	loggerConfig.Level.SetLevel(zapcore.Level(level))

	loggerConfig.ErrorOutputPaths = []string{"stdout"}

	logger, err := loggerConfig.Build(uberZap.AddCallerSkip(1))
	if err != nil {
		log.Panicln(err)
	}

	return &zap{logger: logger.Sugar(), prefixKeysValues: []interface{}{"service", serviceName}}

}

func (l *zap) Infow(msg string, keysAndValues ...interface{}) {
	l.logger.Infow(msg, append(l.prefixKeysValues, keysAndValues...)...)
}
func (l *zap) Warnw(msg string, keysAndValues ...interface{}) {
	l.logger.Warnw(msg, append(l.prefixKeysValues, keysAndValues...)...)

}
func (l *zap) Errorw(msg string, keysAndValues ...interface{}) {
	l.logger.Errorw(msg, append(l.prefixKeysValues, keysAndValues...)...)
}

func (l *zap) Level() Level {
	switch lvl := l.logger.Level(); lvl {
	case uberZap.InfoLevel:
		return InfoLevel
	case uberZap.WarnLevel:
		return WarnLevel
	default:
		return ErrorLevel
	}
}
