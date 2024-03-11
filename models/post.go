package model

import (
	"errors"

	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title     string `gorm:"type:varchar(100);not null"`
	Content   string `gorm:"type:text;not null"`
	AuthorID  uint   `gorm:"type:integer;not null"`
	Published bool   `gorm:"type:boolean;default:false"`
}

func (model *Post) Create(database *gorm.DB) error {
	if err := validatePostInput(model); err != nil {
		return err
	}
	return database.Create(model).Error
}

func (model *Post) ReadAll(database *gorm.DB) ([]Post, error) {
	var posts []Post
	err := database.Find(&posts).Error
	return posts, err
}

func (model *Post) Read(database *gorm.DB) error {
	return database.First(model, model.ID).Error
}

func (model *Post) Update(database *gorm.DB) error {
	if err := validatePostInput(model); err != nil {
		return err
	}
	return database.Save(model).Error
}

func (model *Post) Delete(database *gorm.DB) error {
	return database.Delete(model).Error
}

// validatePostInput checks if the post input meets the requirements
func validatePostInput(post *Post) error {
	if post.Title == "" || len(post.Title) > 100 {
		return errors.New("invalid title")
	}
	if post.Content == "" {
		return errors.New("content cannot be empty")
	}
	return nil
}
