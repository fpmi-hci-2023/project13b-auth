package config

import (
	"os"
	"time"

	"github.com/caarlos0/env/v9"
	"github.com/rs/zerolog"

	"github.com/fpmi-hci-2023/project13b-auth/pkg/logger"
)

var DefaultWriter = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC1123}

var LogInfo struct {
	Level zerolog.Level
}

type Config struct {
	Host         string `env:"HOST,required"`
	HTTPPort     string `env:"HTTP_PORT" envDefault:"8080"`
	GRPCPort     string `env:"GRPC_PORT" envDefault:"4000"`
	FiberPrefork bool   `env:"FIBER_PREFORK" envDefault:"false"`
	DbConnString string `env:"DB_CONN_STRING,required"`
	TTL          struct {
		Access  int64 `env:"TTL_ACCESS,required"`
		Refresh int64 `env:"TTL_REFRESH,required"`
	}
	Secret       string `env:"SECRET,required"`
	AES          string `env:"AES,required"`
	SecureCookie bool   `env:"SECURE_COOKIE" envDefault:"false"`
}

var GlobalConfig Config

func Init(log logger.Logger) {
	if err := env.Parse(&GlobalConfig); err != nil {
		log.Fatal(err, "env parse error")
	}
}
