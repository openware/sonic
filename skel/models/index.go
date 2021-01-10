package models

import (
	"gorm.io/gorm"
)

// Tables export tables ordered list
var Tables = []interface{}{
	&Page{},
}

// db pointer for sharing among models
var db *gorm.DB

// SetDB used to assign `db` connection
// after connection is established on start server
func SetDB(conn *gorm.DB) {
	db = conn
}

// Migrate create and modify database tables according to the models
// TODO: Improve error and log
func Migrate() {
	for _, table := range Tables {
		db.AutoMigrate(table)
	}
}

// Seed execute all table seeding from yaml
// TODO rewrite generic maybe using Tables
func Seed() error {
	err := SeedPages(db)
	if err != nil {
		return err
	}
	return nil
}
