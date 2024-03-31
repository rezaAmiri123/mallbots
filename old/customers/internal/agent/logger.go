package agent

import (
	"github.com/rezaAmiri123/edatV2/di"
	"github.com/rezaAmiri123/edatV2/log"
	"github.com/rezaAmiri123/edatV2/log/zerologger"
	"github.com/rezaAmiri123/mallbots/customers/internal/constants"
)

func (a *Agent) setupLogger() error {
	zlogger, err := zerologger.NewZeroLogger(edatlog.Config{
		Environment: a.config.Environment,
		LogLevel:    a.config.LogLevel,
	})
	if err != nil {
		return err
	}

	// edatlogger := zerologger.Logger(zlogger)
	edatlog.DefaultLogger = zerologger.Logger(zlogger)
	// logger.
	// logger, err := logging.NewZeroLogger(logging.Config{
	// 	Environment: a.config.Environment,
	// 	LogLevel:    a.config.LogLevel,
	// })
	// if err != nil {
	// 	return err
	// }

	// log.DefaultLogger = zerologto.Logger(logger)
	a.container.AddSingleton(constants.LoggerKey, func(c di.Container) (any, error) {
		return zlogger, nil
	})

	return nil
}

func (a *Agent) cleanupLogger() error {
	return nil
}
