package logger

import (
	"context"
	"github.com/rs/zerolog"
	"os"
)

type ZerologAdapter struct {
	logger zerolog.Logger
}

func NewZerologAdapter() *ZerologAdapter {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	return &ZerologAdapter{logger: logger}
}

func (z *ZerologAdapter) Debug(ctx context.Context, msg string, fields map[string]interface{}) {
	event := z.logger.Debug()
	for key, value := range fields {
		event = event.Interface(key, value)
	}
	event.Msg(msg)
}

func (z *ZerologAdapter) Info(ctx context.Context, msg string, fields map[string]interface{}) {
	event := z.logger.Info()
	for key, value := range fields {
		event = event.Interface(key, value)
	}
	event.Msg(msg)
}

func (z *ZerologAdapter) Warn(ctx context.Context, msg string, fields map[string]interface{}) {
	event := z.logger.Warn()
	for key, value := range fields {
		event = event.Interface(key, value)
	}
	event.Msg(msg)
}

func (z *ZerologAdapter) Error(ctx context.Context, msg string, fields map[string]interface{}) {
	event := z.logger.Error()
	for key, value := range fields {
		event = event.Interface(key, value)
	}
	event.Msg(msg)
}

func (z *ZerologAdapter) WithField(ctx context.Context, key string, value interface{}) Logger {
	newLogger := z.logger.With().Interface(key, value).Logger()
	return &ZerologAdapter{logger: newLogger}
}
