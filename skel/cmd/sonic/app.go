package main

import (
	"fmt"

	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
	"github.com/openware/pkg/database"
	"github.com/openware/pkg/kli"
	"github.com/openware/sonic/skel/config"
	"github.com/openware/sonic/skel/handlers"
	"github.com/openware/sonic/skel/models"
)

func serve() error {
	var cfg config.Config
	var app = gin.Default()

	config.Parse(&cfg)
	// Serve static files
	app.Static("/public", "./public")

	// Set up view engine
	app.HTMLRender = ginview.Default()

	// View routes
	handlers.SetUp(app)

	// Connect to the database server with the config/app.yaml configure
	db := database.ConnectDatabase(cfg.Database.Name)
	// defer db.Close()
	models.SetDB(db)
	if !cfg.SkipMigrate {
		models.Migrate()
	}
	handlers.SetPageRoutes(db, app)
	app.Run(":" + cfg.Port)
	return nil
}

func main() {
	// Create new cli
	var cnf string
	cli := kli.NewCli("sonic", "Fullstack micro application", "v1.0.0")

	// Create an init subcommand
	// TODO: copy skel and replace package name
	create := cli.NewSubCommand("create", "Create a sonic application")
	create.Action(func() error {
		println("Creating application")
		return nil
	})

	dbCmd := cli.NewSubCommand("db", "Database commands")
	dbCmd.NewSubCommand("migrate", "Run database migration").
		Action(func() error {
			println("Migrating")
			return nil
		})
	dbCmd.NewSubCommand("seed", "Run database seeding").
		Action(func() error {
			println("Seeding" + cnf)
			return nil
		})

	// Create a test subcommand that's hidden
	serveCmd := cli.NewSubCommand("serve", "Run the application")
	serveCmd.Action(serve)

	// TODO modify pkg/kli to have root flags without
	cli.StringFlag("config", "Application yaml configuration file", &cnf)
	cli.Action(func() error {
		//FIXME test OtherArgs, and test config file
		if cnf != "" {
			args := cli.OtherArgs()
			if len(args) > 0 {
				cli.Run(cli.OtherArgs()...)
				return nil
			}
		}
		cli.PrintHelp()
		return fmt.Errorf("You must select a command")
	})

	// Run!
	if err := cli.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
