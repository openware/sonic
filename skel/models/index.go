package models

import (
	"io/ioutil"
	"log"

	"gorm.io/gorm"
)

// Models contains the list of registered models of the application
var Models = []interface{}{}

type LoaderFunc func([]byte) (interface{}, error)

// SeedAction contains informations needed to seed a table
type SeedAction struct {
	SeedFile string
	Loader   LoaderFunc
}

// SeedActions contains the list of seed actions to perform on start
var SeedActions = []SeedAction{}

// db pointer for sharing among models
var db *gorm.DB

// SetDB used to assign `db` connection
// after connection is established on start server
func SetDB(conn *gorm.DB) {
	db = conn
}

// RegisterModel register a model to the framework
func RegisterModel(model interface{}) {
	Models = append(Models, model)
}

// RegisterSeedAction register a seed action to perform on start of the application
func RegisterSeedAction(seedFile string, loader LoaderFunc) {
	SeedActions = append(SeedActions, SeedAction{
		SeedFile: seedFile,
		Loader:   loader,
	})
}

// Migrate create and modify database tables according to the models
func Migrate() error {
	for _, table := range Models {
		log.Printf("Migrating %T\n", table)
		err := db.AutoMigrate(table)
		if err != nil {
			return err
		}
	}
	return nil
}

// Seed execute all table seeding from yaml
func Seed() error {
	for _, action := range SeedActions {
		err := runSeedAction(db, action)
		if err != nil {
			return err
		}
	}
	return nil
}

func runSeedAction(db *gorm.DB, action SeedAction) error {
	raw, err := ioutil.ReadFile(action.SeedFile)
	if err != nil {
		return err
	}
	list, err := action.Loader(raw)
	if err != nil {
		return err
	}

	tx := db.CreateInBatches(list, 100)
	return tx.Error
}
