package models

import (
	"errors"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
	"gorm.io/gorm"
)

// Page : Table name is `Pages`
type Page struct {
	ID          uint   `gorm:"primarykey"`
	Path        string `gorm:"uniqueIndex;size:64;not null" yaml:"path"`
	Lang        string `yaml:"lang"`
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Body        string `yaml:"body"`
	Timestamps
}

// SeedPages load from Page from config/Pages.yml to database
// TODO: Remote from model
func SeedPages(db *gorm.DB) error {
	raw, err := ioutil.ReadFile("config/seeds/pages.yml")
	if err != nil {
		return err
	}
	Pages := []Page{}
	err = yaml.Unmarshal(raw, &Pages)
	if err != nil {
		return err
	}

	tx := db.Create(&Pages)
	if tx.Error != nil {
		return err
	}
	return nil
}

// FindByPath find and return a page by path
func (p *Page) FindByPath(path string) *Page {
	page := Page{}
	tx := db.Where("path = ?", path).First(&page)

	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil
		}
		log.Fatalf("FindPageByPath failed: %s", tx.Error.Error())
		return nil
	}
	return &page
}

// List returns all pages
func (p *Page) List() []Page {
	pages := []Page{}
	tx := db.Find(&pages)

	if tx.Error != nil {
		log.Fatalf("FindPageByPath failed: %s", tx.Error.Error())
	}
	return pages
}
