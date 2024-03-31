package system

import (
	"github.com/rezaAmiri123/edatV2/log"
	"github.com/rezaAmiri123/edatV2/log/zerologger"
)

func (s *System) setupLogger() error {
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

func (a *Agent) cleanupLogger() error {
	return nil
}
