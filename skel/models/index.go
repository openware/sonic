package models

import (
	"time"

	"gorm.io/gorm"
)

// Timestamps adding time at the end of models
type Timestamps struct {
	CreatedAt time.Time `yaml:"created_at"`
	UpdatedAt time.Time `yaml:"updated_at"`
}

// db pointer for sharing among models
var db *gorm.DB

// Tables export tables list
// Please order tables to able to delete tables when drop
// Used in migration database
var Tables = []interface{}{
	&Page{},
}

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
