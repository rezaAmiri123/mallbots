package agent

import (
	"context"
	"time"

	edatlog "github.com/rezaAmiri123/edatV2/log"
)

type CleanupFiles func(ctx context.Context) error

var envFiles []string

type pgCfg struct {
	Conn string `envconfig:"CONN" desc:"a postgres DATABASE_URL or CONNECTION_STRING"`
}

type natsCfg struct {
	URL            string        `envconfig:"URL"`
	ClusterID      string        `envconfig:"CLUSTER_ID"`
	// Stream         string        `envconfig:"STREAM"`
	AckWaitTimeout time.Duration `envconfig:"ACK_WAIT_TIMEOUT" default:"30s"`
}

type kafkaCfg struct {
	Brokers []string `envconfig:"BROKERS"`
}

type WebServerCfg struct {
	Port              string        `envconfig:"PORT" default:":80"`
	CertPath          string        `envconfig:"CERT_PATH"`
	KeyPath           string        `envconfig:"KEY_PATH"`
	ReadTimeout       time.Duration `envconfig:"READ_TIMEOUT" default:"1s"`
	WriteTimeout      time.Duration `envconfig:"WRITE_TIMEOUT" default:"1s"`
	IdleTimeout       time.Duration `envconfig:"IDLE_TIMEOUT" default:"30s"`
	ReadHeaderTimeout time.Duration `envconfig:"READ_HEADER_TIMEOUT" default:"2s"`
	RequestTimeout    time.Duration `envconfig:"REQUEST_TIMEOUT" default:"60s"`
}

type webCfg struct {
	ApiPath     string       `envconfig:"API_PATH" default:"/api"`
	PingPath    string       `envconfig:"PING_PATH" default:"/ping"`
	MetricsPath string       `envconfig:"METRICS_PATH" default:"/metrics"`
	Http        WebServerCfg `envconfig:"HTTP"`
	Cors        WebCorsCfg   `envconfig:"CORS"`
}

type WebCorsCfg struct {
	Origins          []string `envconfig:"ORIGINS" default:"*"`
	AllowCredentials bool     `envconfig:"ALLOW_CREDENTIALS" default:"true"`
	MaxAge           int      `envconfig:"MAX_AGE" default:"300"`
}
type MonitoringCfg struct{
	LivenessAddress string `envconfig:"LIVENESS_ADDRESS" default:":8080"`
	MetricAddress string`envconfig:"METRIC_ADDRESS" default:":8080"`
	PproffAddress string `envconfig:"PPROFF_ADDRESS" default:":8080"`
}
type ServerCfg struct {
	Network  string `envconfig:"NETWORK" default:"tcp"`
	Address  string `envconfig:"ADDRESS" default:":8000"`
	CertPath string `envconfig:"CERT_PATH"`
	KeyPath  string `envconfig:"KEY_PATH"`
}

type ClientCfg struct {
	Address  string `envconfig:"ADDRESS"`
	CertPath string `envconfig:"CERT_PATH"`
	KeyPath  string `envconfig:"KEY_PATH"`
}

type PGConfig struct {
	PGDriver     string `mapstructure:"POSTGRES_DRIVER" envconfig:"DRIVER"`
	PGHost       string `mapstructure:"POSTGRES_HOST" envconfig:"HOST"`
	PGPort       string `mapstructure:"POSTGRES_PORT" envconfig:"PORT"`
	PGUser       string `mapstructure:"POSTGRES_USER" envconfig:"USER"`
	PGDBName     string `mapstructure:"POSTGRES_DB_NAME" envconfig:"DB_NAME"`
	PGPassword   string `mapstructure:"POSTGRES_PASSWORD" envconfig:"PASSWORD"`
	PGSearchPath string `mapstructure:"POSTGRES_SEARCH_PATH" envconfig:"SEARCH_PATH"`
}

type Config struct {
	Environment     string        `envconfig:"ENVIRONMENT" default:"production"`
	ServiceID       string        `envconfig:"SERVICE_ID" required:"true"`
	LogLevel        edatlog.Level `envconfig:"LOG_LEVEL" default:"WARN" desc:"options: [TRACE,DEBUG,INFO,WARN,ERROR,PANIC]"`
	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"30s" desc:"time to allow services to gracefully stop"`
	Web             webCfg        `envconfig:"WEB"`        // Web Config
	Rpc             ServerCfg     `envconfig:"RPC"`        // RPC Config
	Restaurant      ClientCfg     `envconfig:"RESTAURANT"` // RPC Client Config
	//Postgres        postgres.Config `envconfig:"PG"`                                                              // DataDriver / Postgres
	Postgres    pgCfg    `envconfig:"PG"`
	EventDriver string   `envconfig:"EVENT_DRIVER" default:"inmem" desc:"options: [inmem,nats,kafka]"` // "inmem", "nats", "kafka"
	Nats        natsCfg  `envconfig:"NATS"`                                                            // EventDriver / Nats Streaming Config
	Kafka       kafkaCfg `envconfig:"KAFKA"`                                                           // EventDriver / Kafka Config
}
