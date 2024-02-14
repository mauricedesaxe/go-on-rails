package model

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title     string `gorm:"type:varchar(100);not null"`
	Content   string `gorm:"type:text;not null"`
	Author    string `gorm:"type:varchar(100);not null"`
	Published bool   `gorm:"type:boolean;default:false"`
}
