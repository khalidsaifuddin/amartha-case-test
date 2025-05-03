package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ServiceName string `envconfig:"SERVICE_NAME" default:"Amartha Case 1 - Billing Engine"`
	HTTPPort    string `envconfig:"HTTP_PORT" default:"8080"`
	Environment string `envconfig:"ENVIRONMENT" default:"staging"`

	Host            string `envconfig:"PGSQL_HOST" default:""`
	Port            string `envconfig:"PGSQL_PORT" default:""`
	Username        string `envconfig:"PGSQL_USERNAME" default:""`
	Password        string `envconfig:"PGSQL_PASSWORD" default:""`
	DBName          string `envconfig:"PGSQL_DBNAME" default:""`
	LogMode         bool   `envconfig:"DB_LOG_MODE" default:"true"`
	MaxIdleConns    int    `envconfig:"DB_MAX_IDLE_CONNS" default:"5"`
	MaxOpenConns    int    `envconfig:"DB_MAX_OPEN_CONNS" default:"10"`
	ConnMaxLifetime int    `envconfig:"DB_CONN_MAX_LIFETIME" default:"10"`
	IsDebugMode     bool   `envconfig:"DEBUG_MODE" default:"true"`

	EnableRedisCache bool   `envconfig:"ENABLE_REDIS_CACHE" default:"false"`
	RedisHost        string `envconfig:"REDIS_HOST" default:"127.0.0.1"`
	RedisPort        string `envconfig:"REDIS_PORT" default:"6379"`
	RedisPassword    string `envconfig:"REDIS_PASSWORD" default:""`
	RedisMaxIdle     int    `envconfig:"REDIS_MAX_IDLE" default:"10"`
	DefaultTTL       int64  `envconfig:"DEFAULT_TTL" default:"10"`

	AccessControlAllowOrigin string `envconfig:"ACCESS_CONTROL_ALLOW_ORIGIN" default:"*"`
	AccessControlAllowMethod string `envconfig:"ACCESS_CONTROL_ALLOW_METHOD" default:"POST, HEAD, PATCH, OPTIONS, GET, PUT, DELETE"`

	APISecretKey string `envconfig:"API_SECRET_KEY" default:"APISecretKey"`

	// config value
	IsDelinquentLoanLimit int `envconfig:"IS_DELINQUENT_LOAN_LIMIT" default:"2"`
}

func Get() Config {
	cfg := Config{}

	envconfig.MustProcess("", &cfg)
	return cfg
}
