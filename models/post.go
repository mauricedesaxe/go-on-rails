package model

import (
	"errors"

	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	ID        uint   `gorm:"primaryKey"`
	Title     string `gorm:"type:varchar(100);not null"`
	Content   string `gorm:"type:text;not null"`
	Author    string `gorm:"type:varchar(100);not null"`
	Published bool   `gorm:"type:boolean;default:false"`
}

func (p *Post) Create(database *gorm.DB) error {
	if err := validatePostInput(p); err != nil {
		return err
	}
	return database.Create(p).Error
}

func (p *Post) ReadAll(database *gorm.DB) ([]Post, error) {
	var posts []Post
	err := database.Find(&posts).Error
	return posts, err
}

func (p *Post) Read(database *gorm.DB) error {
	return database.First(p, p.ID).Error
}

func (p *Post) ReadByTitle(database *gorm.DB) error {
	return database.Where("title = ?", p.Title).Or("title LIKE ?", "%"+p.Title+"%").First(p).Error
}

func (p *Post) ReadByAuthor(database *gorm.DB) ([]Post, error) {
	var posts []Post
	err := database.Where("author = ?", p.Author).Or("author LIKE ?", "%"+p.Author+"%").Find(&posts).Error
	return posts, err
}

func (p *Post) ReadByContent(database *gorm.DB) ([]Post, error) {
	var posts []Post
	err := database.Where("content LIKE ?", "%"+p.Content+"%").Find(&posts).Error
	return posts, err
}

func (p *Post) ReadByPublished(database *gorm.DB) ([]Post, error) {
	var posts []Post
	err := database.Where("published = ?", p.Published).Find(&posts).Error
	return posts, err
}

func (p *Post) Update(database *gorm.DB) error {
	if err := validatePostInput(p); err != nil {
		return err
	}
	return database.Save(p).Error
}

func (p *Post) Delete(database *gorm.DB) error {
	return database.Delete(p).Error
}

// validatePostInput checks if the post input meets the requirements
func validatePostInput(p *Post) error {
	if p.Title == "" || len(p.Title) > 100 {
		return errors.New("invalid title")
	}
	if p.Content == "" {
		return errors.New("content cannot be empty")
	}
	if p.Author == "" || len(p.Author) > 100 {
		return errors.New("invalid author")
	}
	return nil
}
