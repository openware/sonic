package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/openware/pkg/database"
	"github.com/openware/pkg/ika"
	"github.com/openware/pkg/kli"
	"github.com/openware/sonic/skel/handlers"
	"github.com/openware/sonic/skel/models"
	"gorm.io/gorm"
)

// Version of the application displayed by the cli and the version endpoint
// TODO move to release system passing it as build param
var Version = "v1.0.0"

// Database struct {
// 	Driver string `yaml:"driver" env:"DATABASE_DRIVER" env-description:"Database driver"`
// 	Host   string `yaml:"host" env:"DATABASE_HOST" env-description:"Database host"`
// 	Port   string `yaml:"port" env:"DATABASE_PORT" env-description:"Database port"`
// 	Name   string `yaml:"name" env:"DATABASE_NAME" env-description:"Database name"`
// 	User   string `yaml:"user" env:"DATABASE_USER" env-description:"Database user"`
// 	Pass   string `env:"DATABASE_PASS" env-description:"Database user password"`
// } `yaml:"database"`

// Config is the application configuration structure
type Config struct {
	// TODO move to database package
	Database database.Config `yaml:"database"`

	// TODO Create a redis and vault package
	Redis struct {
		Host string `yaml:"host" env:"REDIS_HOST" env-description:"Redis Server host" env-default:"localhost"`
		Port string `yaml:"port" env:"REDIS_PORT" env-description:"Redis Server port" env-default:"6379"`
	} `yaml:"redis"`
	Port string `env:"APP_PORT" env-description:"Port for HTTP service" env-default:"6009"`
}

// Runtime configuration of the application
type Runtime struct {
	conf Config
	db   *gorm.DB
	srv  *gin.Engine
}

// App config for the application
var App Runtime

func serve() error {
	if err := boot(); err != nil {
		return err
	}
	App.srv = gin.Default()
	// Setup routes
	handlers.Setup(App.srv)
	// TODO handlers.Setup(conf)
	App.srv.Run(":" + App.conf.Port)
	return nil
}

func dbCreate() error {
	// Use existing connection
	db, err := database.Connect(&App.conf.Database)
	if err != nil {
		return err
	}
	tx := db.Exec(fmt.Sprintf("CREATE DATABASE `%s`;", App.conf.Database.Name))
	return tx.Error
}

func dbMigrate() error {
	if err := boot(); err != nil {
		return err
	}
	return models.Migrate()
}

func dbSeed() error {
	if err := boot(); err != nil {
		return err
	}
	return models.Seed()
}

// TODO: copy skel and replace package name
func appCreate() error {
	println("Creating app")
	return nil
}

// boot is executed before commands
func boot() error {
	// Connect to the database server with the config/app.yaml configure
	// TODO write boot()
	// 	conf := sonic.ParseConfig()
	// 	App.DB := database.Connect(conf)
	// 	models.Setup(conf)
	var err error

	App.db, err = database.Connect(&App.conf.Database)
	if err != nil {
		return err
	}
	models.Setup(App.db)
	return models.Migrate()

}

func main() {
	// Create new cli
	cnf := "config/app.yml"
	cli := kli.NewCli("sonic", "Fullstack micro application", Version)
	cli.StringFlag("config", "Application yaml configuration file", &cnf)

	// Create an init subcommand
	cli.NewSubCommand("create", "Create a sonic application").Action(appCreate)

	dbCmd := cli.NewSubCommand("db", "Database commands")
	dbCmd.NewSubCommand("create", "Create database").Action(dbCreate)
	dbCmd.NewSubCommand("migrate", "Run database migration").Action(dbMigrate)
	dbCmd.NewSubCommand("seed", "Run database seeding").Action(dbSeed)

	serveCmd := cli.NewSubCommand("serve", "Run the application")
	serveCmd.Action(serve)

	// read configuration from the file and environment variables
	if err := ika.ReadConfig(cnf, &App.conf); err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	if err := cli.Run(); err != nil {
		log.Fatalf("Run: %v\n", err)
	}
}
