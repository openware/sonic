package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/openware/pkg/ika"
)

// Config is the application configuration structure
type Config struct {
	Database struct {
		Driver string `yaml:"driver" env:"DATABASE_DRIVER" env-description:"Database driver"`
		Host   string `yaml:"host" env:"DATABASE_HOST" env-description:"Database host"`
		Port   string `yaml:"port" env:"DATABASE_PORT" env-description:"Database port"`
		Name   string `yaml:"name" env:"DATABASE_NAME" env-description:"Database name"`
		User   string `yaml:"user" env:"DATABASE_USER" env-description:"Database user"`
		Pass   string `env:"DATABASE_PASS" env-description:"Database user password"`
	} `yaml:"database"`
	Redis struct {
		Host string `yaml:"host" env:"REDIS_HOST" env-description:"Redis Server host" env-default:"localhost"`
		Port string `yaml:"port" env:"REDIS_PORT" env-description:"Redis Server port" env-default:"6379"`
	} `yaml:"redis"`
	Port string `env:"APP_PORT" env-description:"Port for HTTP service" env-default:"6009"`
	// Internal options
	ConfigPath  string
	SkipMigrate bool
}

// Parse the configuration and prefill the param cfg
func Parse(cfg *Config) {
	processArgs(cfg)

	// read configuration from the file and environment variables
	if err := ika.ReadConfig(cfg.ConfigPath, cfg); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}

// ProcessArgs processes and handles CLI arguments
// TODO: can be removed using commando
func processArgs(cfg *Config) {
	f := flag.NewFlagSet("server", 1)
	f.StringVar(&cfg.ConfigPath, "c", "config/app.yml", "Path to configuration file")
	f.BoolVar(&cfg.SkipMigrate, "skip-migration", false, "Skip database migration")

	fu := f.Usage
	f.Usage = func() {
		fu()
		envHelp, _ := ika.GetDescription(cfg, nil)
		fmt.Fprintln(f.Output())
		fmt.Fprintln(f.Output(), envHelp)
	}

	f.Parse(os.Args[1:])
}
