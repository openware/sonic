package main

import (
	"fmt"
	"os"

	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
	"github.com/openware/pkg/database"
	"github.com/openware/pkg/ika"
	"github.com/openware/pkg/kli"
	"github.com/openware/sonic/skel/handlers"
	"github.com/openware/sonic/skel/models"
)

// Version of the application displayed by the cli and the version endpoint
var Version = "v1.0.0"

// Config is the application configuration structure
type config struct {
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
}

// Config for the application
var Config config

func serve() error {
	var app = gin.Default()

	// Serve static files
	app.Static("/public", "./public")

	// Set up view engine
	app.HTMLRender = ginview.Default()

	// Setup routes
	handlers.Setup(app)

	// Connect to the database server with the config/app.yaml configure
	db := database.ConnectDatabase(Config.Database.Name)
	models.SetDB(db)
	models.Migrate()
	handlers.SetPageRoutes(app)
	app.Run(":" + Config.Port)
	return nil
}

func dbCreate() error {
	db := database.ConnectDatabase("")
	tx := db.Exec(fmt.Sprintf("CREATE DATABASE `%s`;", Config.Database.Name))
	return tx.Error
}

func dbMigrate() error {
	println("Migrating")
	db := database.ConnectDatabase(Config.Database.Name)
	models.SetDB(db)
	return models.Migrate()
}

func dbSeed() error {
	db := database.ConnectDatabase(Config.Database.Name)
	models.SetDB(db)
	return models.Seed()
}

func appCreate() error {
	println("Creating app")
	return nil
}

// Parse the configuration and prefill the param cfg
func parse(path string) {
	// read configuration from the file and environment variables
	fmt.Println(path)
	if err := ika.ReadConfig(path, &Config); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}

func main() {
	// Create new cli
	cnf := "config/app.yml"
	cli := kli.NewCli("sonic", "Fullstack micro application", Version)
	cli.StringFlag("config", "Application yaml configuration file", &cnf)

	// Create an init subcommand
	// TODO: copy skel and replace package name
	cli.NewSubCommand("create", "Create a sonic application").Action(appCreate)

	dbCmd := cli.NewSubCommand("db", "Database commands")
	dbCmd.NewSubCommand("create", "Create database").Action(dbCreate)
	dbCmd.NewSubCommand("migrate", "Run database migration").Action(dbMigrate)
	dbCmd.NewSubCommand("seed", "Run database seeding").Action(dbSeed)

	// Create a test subcommand that's hidden
	serveCmd := cli.NewSubCommand("serve", "Run the application")
	serveCmd.Action(serve)

	parse(cnf)

	// Run!
	if err := cli.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
