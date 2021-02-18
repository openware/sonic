package sonic

import (
	kaigara "github.com/openware/kaigara/pkg/config"
	"github.com/openware/pkg/database"
	"github.com/openware/pkg/mngapi"
)

// Config is the application configuration structure
type Config struct {
	Database database.Config `yaml:"database"`
	// TODO Create a redis and vault package
	Redis struct {
		Host string `yaml:"host" env:"REDIS_HOST" env-description:"Redis Server host" env-default:"localhost"`
		Port string `yaml:"port" env:"REDIS_PORT" env-description:"Redis Server port" env-default:"6379"`
	} `yaml:"redis"`
	Port          string                `env:"APP_PORT" env-description:"Port for HTTP service" env-default:"6009"`
	MngAPI        mngapi.Config         `yaml:"mngapi"`
	KaigaraConfig kaigara.KaigaraConfig `yaml:"kaigara"`
	Vault         VaultConfig           `yaml:"vault"`
	DeploymentID  string                `yaml:"deploymentID" env:"DEPLOYMENT_ID"`
}

// VaultConfig contains Vault-related configuration
type VaultConfig struct {
	Addr  string `yaml:"addr" env:"VAULT_ADDR" env-default:"http://localhost:8200"`
	Token string `yaml:"token" env:"VAULT_TOKEN"`
}
