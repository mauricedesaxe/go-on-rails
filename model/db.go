package model

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(dest string) *gorm.DB {
	var err error
	db, err := gorm.Open(sqlite.Open(dest), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&Post{})

	return db
}
