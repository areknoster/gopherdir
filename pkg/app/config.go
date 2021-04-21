package app

import (
	"flag"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
)

// LoadConfig loads configuration from envs and prints usage if requirements are not fulfilled
func LoadConfig(logger *zap.Logger, cfg interface{}) {
	if flag.Lookup("help") != nil {
		envconfig.Usage("", cfg)

	}
	err := envconfig.Process("", cfg)
	if err != nil {
		envconfig.Usage("", cfg)
		logger.Fatal("app: could not process config", zap.Error(err))
	}
}
