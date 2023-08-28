package wlog

import (
	"context"
	"os"
	"time"

	"github.com/gofrs/uuid"
	"github.com/rs/zerolog"
	"runtime/debug"
)

const (
	goVersionField     = "go_version"
	correlationIDField = "correlation_id"
)

func NewDefaultLogger() zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339Nano

	buildInfo, _ := debug.ReadBuildInfo()

	return zerolog.New(os.Stdout).
		Level(zerolog.InfoLevel).
		With().
		Stack().
		Caller().
		Timestamp().
		Str(goVersionField, buildInfo.GoVersion).
		Logger()
}

func Debug(ctx context.Context) *zerolog.Event {
	l := getLogger(ctx)

	return l.Debug()
}

func Info(ctx context.Context) *zerolog.Event {
	l := getLogger(ctx)

	return l.Info()
}

func Warn(ctx context.Context) *zerolog.Event {
	l := getLogger(ctx)

	return l.Warn()
}

func Error(ctx context.Context) *zerolog.Event {
	l := getLogger(ctx)

	return l.Error()
}

func getLogger(ctx context.Context) *zerolog.Logger {
	l := zerolog.Ctx(ctx)

	if l.GetLevel() == zerolog.Disabled {
		newLog := zerolog.New(os.Stdout).
			Level(zerolog.InfoLevel).
			With().
			Timestamp().
			Logger()

		newLog.UpdateContext(func(c zerolog.Context) zerolog.Context {
			correlationID, _ := uuid.NewV4()

			return c.Str(correlationIDField, correlationID.String())
		})

		return &newLog
	}

	return l
}

func AddContextString(ctx context.Context, key, value string) context.Context {
	l := getLogger(ctx)

	ctx = l.WithContext(ctx)

	l = getLogger(ctx)

	l.UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.Str(key, value)
	})

	return ctx
}
