package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Load() {
	db, err := gorm.Open(sqlite.Open("challenge.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	DB = db
}
