package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/openware/pkg/database"
	"github.com/openware/pkg/ika"
	"github.com/openware/pkg/kli"
	"github.com/openware/sonic"
	"github.com/openware/sonic/skel/handlers"
	"github.com/openware/sonic/skel/models"
)

// Version of the application displayed by the cli and the version endpoint
var Version = "1.0.0"

// App config for the application
var App sonic.Runtime

func serve() error {
	App.Srv = gin.Default()
	handlers.Setup(&App)
	App.Srv.Run(":" + App.Conf.Port)
	return nil
}

func dbCreate() error {
	// TODO move to pkg database and drop db command
	db := database.ConnectDatabase("")
	tx := db.Exec(fmt.Sprintf("CREATE DATABASE `%s`;", App.Conf.Database.Name))
	return tx.Error
}

// boot is executed before commands
func boot() error {
	App.DB = database.ConnectDatabase(App.Conf.Database.Name)
	models.Setup(&App)
	return models.Migrate()
}

func main() {
	// Create new cli
	cnf := "config/app.yml"
	cli := kli.NewCli("sonic", "Fullstack micro application", Version)
	cli.StringFlag("config", "Application yaml configuration file", &cnf)

	// TODO move create and drop to sonic pkg

	dbCmd := cli.NewSubCommand("db", "Database commands")
	dbCmd.NewSubCommand("create", "Create database").Action(dbCreate)
	dbCmd.NewSubCommand("migrate", "Run database migration").Action(boot)
	dbCmd.NewSubCommand("seed", "Run database seeding").Action(func() error {
		return models.Seed()
	})

	serveCmd := cli.NewSubCommand("serve", "Run the application")
	serveCmd.Action(serve)

	// read configuration from the file and environment variables
	if err := ika.ReadConfig(cnf, &App.Conf); err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	if err := boot(); err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	if err := cli.Run(); err != nil {
		log.Fatalf("Run: %v\n", err)
	}
}
