package model

import (
	"errors"

	"gorm.io/gorm"
)

type PostModel struct {
	gorm.Model
	ID        uint   `gorm:"primaryKey"`
	Title     string `gorm:"type:varchar(100);not null"`
	Content   string `gorm:"type:text;not null"`
	Author    string `gorm:"type:varchar(100);not null"`
	Published bool   `gorm:"type:boolean;default:false"`
}

func (model *PostModel) Create(database *gorm.DB) error {
	if err := validatePostInput(model); err != nil {
		return err
	}
	return database.Create(model).Error
}

func (model *PostModel) ReadAll(database *gorm.DB) ([]PostModel, error) {
	var posts []PostModel
	err := database.Find(&posts).Error
	return posts, err
}

func (model *PostModel) Read(database *gorm.DB) error {
	return database.First(model, model.ID).Error
}

func (model *PostModel) Update(database *gorm.DB) error {
	if err := validatePostInput(model); err != nil {
		return err
	}
	return database.Save(model).Error
}

func (model *PostModel) Delete(database *gorm.DB) error {
	return database.Delete(model).Error
}

// validatePostInput checks if the post input meets the requirements
func validatePostInput(post *PostModel) error {
	if post.Title == "" || len(post.Title) > 100 {
		return errors.New("invalid title")
	}
	if post.Content == "" {
		return errors.New("content cannot be empty")
	}
	if post.Author == "" || len(post.Author) > 100 {
		return errors.New("invalid author")
	}
	return nil
}
