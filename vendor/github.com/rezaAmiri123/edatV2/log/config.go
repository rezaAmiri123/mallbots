package edatlog

type Level string

const (
	TRACE Level = "TRACE"
	DEBUG Level = "DEBUG"
	INFO  Level = "INFO"
	WARN  Level = "WARN"
	ERROR Level = "ERROR"
	PANIC Level = "PANIC"
)

// type EnvironmentConfig = func(options []zap.Option) zap.Config
//
type Config struct {
	Environment string
	LogLevel    Level
	//	Options     []zap.Option
}