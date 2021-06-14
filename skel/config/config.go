package config

import (
	"github.com/openware/pkg/mngapi"
	sonicCfg "github.com/openware/pkg/sonic/config"
)

// Config is the application configuration structure
type Config struct {
	Database DatabaseConfig `yaml:"database"`
	// TODO Create a redis and vault package
	Redis struct {
		Host string `yaml:"host" env:"REDIS_HOST" env-description:"Redis Server host" env-default:"localhost"`
		Port string `yaml:"port" env:"REDIS_PORT" env-description:"Redis Server port" env-default:"6379"`
	} `yaml:"redis"`
	Port                string        `env:"APP_PORT" env-description:"Port for HTTP service" env-default:"6009"`
	MngAPI              mngapi.Config `yaml:"mngapi"`
	Vault               sonicCfg.VaultConfig   `yaml:"vault"`
	DeploymentID        string        `yaml:"deploymentID" env:"DEPLOYMENT_ID"`
	Opendax             sonicCfg.OpendaxConfig `yaml:"opendax"`
	MarketsBlacklist    string        `yaml:"markets_blacklist" env:"MARKETS_BLACKLIST"`
	CurrenciesBlacklist string        `yaml:"currencies_blacklist" env:"CURRENCIES_BLACKLIST"`
}
