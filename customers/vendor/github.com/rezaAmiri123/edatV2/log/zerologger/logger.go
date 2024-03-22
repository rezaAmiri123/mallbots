package zerologger

import (
	"os"

	edatlog "github.com/rezaAmiri123/edatV2/log"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

type zerologLogger struct {
	l zerolog.Logger
}

var _ edatlog.Logger = (*zerologLogger)(nil)

func Logger(logger zerolog.Logger) edatlog.Logger {
	zLog := logger.With().CallerWithSkipFrameCount(3).Logger()
	return &zerologLogger{l: zLog}
}
func NewZeroLogger(cfg edatlog.Config) (zerolog.Logger, error) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	var logger zerolog.Logger
	switch cfg.Environment {
	case "production":
		logger = zerolog.New(os.Stdout).
			Level(logLevelToZero(cfg.LogLevel)).
			With().
			Timestamp().
			Caller().
			Logger()
	default:
		stdout := zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
			w.TimeFormat = "03:04:05.000PM"
		})
		logger = zerolog.New(stdout).
			Level(logLevelToZero(cfg.LogLevel)).
			With().
			Timestamp().
			Caller().
			Logger()
	}

	
	return logger, nil
}

func (l *zerologLogger) Trace(msg string, fields ...edatlog.Field) {
	if l.l.GetLevel() > zerolog.DebugLevel {
		return
	}
	logger := l.fields(l.l.With(), fields).Logger()
	logger.Debug().Msg(msg)
}

func (l *zerologLogger) Debug(msg string, fields ...edatlog.Field) {
	if l.l.GetLevel() > zerolog.DebugLevel {
		return
	}
	logger := l.fields(l.l.With(), fields).Logger()
	logger.Debug().Msg(msg)

}

func (l *zerologLogger) Info(msg string, fields ...edatlog.Field) {
	if l.l.GetLevel() > zerolog.InfoLevel {
		return
	}
	logger := l.fields(l.l.With(), fields).Logger()
	logger.Info().Msg(msg)
}

func (l *zerologLogger) Warn(msg string, fields ...edatlog.Field) {
	if l.l.GetLevel() > zerolog.WarnLevel {
		return
	}
	logger := l.fields(l.l.With(), fields).Logger()
	logger.Warn().Msg(msg)

}

func (l *zerologLogger) Error(msg string, fields ...edatlog.Field) {
	if l.l.GetLevel() > zerolog.ErrorLevel {
		return
	}
	logger := l.fields(l.l.With(), fields).Logger()
	logger.Error().Msg(msg)

}

func (l *zerologLogger) Sub(fields ...edatlog.Field) edatlog.Logger {
	return &zerologLogger{
		l: l.fields(l.l.With(), fields).Logger(),
	}
}

func (l *zerologLogger) fields(ctx zerolog.Context, fields []edatlog.Field) zerolog.Context {
	for _, field := range fields {
		switch field.Type {
		case edatlog.StringType:
			ctx = ctx.Str(field.Key, field.String)
		case edatlog.IntType:
			ctx = ctx.Int(field.Key, field.Int)
		case edatlog.DurationType:
			ctx = ctx.Str(field.Key, field.Duration.String())
		case edatlog.ErrorType:
			ctx = ctx.Stack().Err(field.Error)
		}
	}

	return ctx
}

func logLevelToZero(level edatlog.Level) zerolog.Level {
	switch level {
	case edatlog.PANIC:
		return zerolog.PanicLevel
	case edatlog.ERROR:
		return zerolog.ErrorLevel
	case edatlog.WARN:
		return zerolog.WarnLevel
	case edatlog.INFO:
		return zerolog.InfoLevel
	case edatlog.DEBUG, edatlog.TRACE:
		return zerolog.DebugLevel
	default:
		return zerolog.InfoLevel

	}
}
