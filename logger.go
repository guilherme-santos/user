package user

import (
	"context"

	"github.com/sirupsen/logrus"
)

type contextKey string

var loggerCtx = contextKey("logger")

func SetLogger(ctx context.Context, logger logrus.FieldLogger) context.Context {
	return context.WithValue(ctx, loggerCtx, logger)
}

func Logger(ctx context.Context) logrus.FieldLogger {
	logger, _ := ctx.Value(loggerCtx).(logrus.FieldLogger)
	if logger == nil {
		logger = logrus.New()
	}
	return logger
}
