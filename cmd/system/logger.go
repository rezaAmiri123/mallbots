package system

import (
	"github.com/rezaAmiri123/edatV2/log"
	"github.com/rezaAmiri123/edatV2/log/zerologger"
	"github.com/rs/zerolog"
)
type Logger interface{
	Logger() zerolog.Logger
}

func (s *System) initLogger() error {
	zlogger, err := zerologger.NewZeroLogger(edatlog.Config{
		Environment: s.cfg.Environment,
		LogLevel:    edatlog.Level(s.cfg.LogLevel),
	})
	if err != nil {
		return err
	}
	edatlog.DefaultLogger = zerologger.Logger(zlogger)
	s.logger = zlogger
	return nil
}

func (s *System) Logger() zerolog.Logger {
	return s.logger
}
