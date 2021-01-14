package sonic

import (
	"github.com/openware/pkg/database"
)

// Config is the application configuration structure
type Config struct {
	Database database.Config `yaml:"database"`
	// TODO Create a redis and vault package
	Redis struct {
		Host string `yaml:"host" env:"REDIS_HOST" env-description:"Redis Server host" env-default:"localhost"`
		Port string `yaml:"port" env:"REDIS_PORT" env-description:"Redis Server port" env-default:"6379"`
	} `yaml:"redis"`
	Port string `env:"APP_PORT" env-description:"Port for HTTP service" env-default:"6009"`
}
