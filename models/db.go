package model

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB(dest string) *gorm.DB {
	var err error
	db, err := gorm.Open(sqlite.Open(dest), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&PostModel{})
	db.AutoMigrate(&UserModel{})
	db.AutoMigrate(&TokenModel{})

	return db
}
