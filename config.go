package sonic

// Config is the application configuration structure
type Config struct {
	// TODO move to database package
	Database struct {
		Driver string `yaml:"driver" env:"DATABASE_DRIVER" env-description:"Database driver"`
		Host   string `yaml:"host" env:"DATABASE_HOST" env-description:"Database host"`
		Port   string `yaml:"port" env:"DATABASE_PORT" env-description:"Database port"`
		Name   string `yaml:"name" env:"DATABASE_NAME" env-description:"Database name"`
		User   string `yaml:"user" env:"DATABASE_USER" env-description:"Database user"`
		Pass   string `env:"DATABASE_PASS" env-description:"Database user password"`
	} `yaml:"database"`
	// TODO Create a redis and vault package
	Redis struct {
		Host string `yaml:"host" env:"REDIS_HOST" env-description:"Redis Server host" env-default:"localhost"`
		Port string `yaml:"port" env:"REDIS_PORT" env-description:"Redis Server port" env-default:"6379"`
	} `yaml:"redis"`
	Port string `env:"APP_PORT" env-description:"Port for HTTP service" env-default:"6009"`
}
