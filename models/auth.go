package model

import (
	"crypto/rand"
	"errors"
	"net/mail"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID    uint   `gorm:"primaryKey"`
	Email string `gorm:"unique"`
	Role  string `gorm:"type:varchar(100);not null;default:'user'"` // "user", "admin", ...
}

func (model *User) Create(database *gorm.DB) error {
	_, err := mail.ParseAddress(model.Email)
	if err != nil {
		return errors.New("invalid email address")
	}
	if model.Role == "" {
		model.Role = "user"
	}

	return database.Create(model).Error
}

func (model *User) ReadAll(database *gorm.DB) ([]User, error) {
	var users []User
	err := database.Find(&users).Error
	return users, err
}

func (model *User) Read(database *gorm.DB) error {
	return database.Order("id desc").First(model).Error
}

func (model *User) ReadByEmail(database *gorm.DB) error {
	return database.Where("email = ?", model.Email).First(model).Error
}

func (model *User) Delete(database *gorm.DB) error {
	return database.Delete(model).Error
}

// For magic link login, but could also help with email verification, password reset if that was needed.
// It's meant as a way to verify that the user has access to the email of the account.
type Token struct {
	gorm.Model
	Email string `gorm:"primaryKey"`
	Value string `gorm:"primaryKey"`
}

// Note that this generates a random hashed token value. The unhashed value is returned for
// use in the email link.
func (model *Token) Create(database *gorm.DB) (string, error) {
	// Validate the email
	_, err := mail.ParseAddress(model.Email)
	if err != nil {
		return "", errors.New("invalid email address")
	}

	// Generate a random 32 character string using crypto/rand for better randomness
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789@#%^*"
	b := make([]byte, 32)
	_, err = rand.Read(b)
	if err != nil {
		return "", err
	}
	for i := range b {
		b[i] = charset[b[i]%byte(len(charset))]
	}
	unhashedValue := string(b)

	// Hash the random string
	hashBytes, err := bcrypt.GenerateFromPassword(
		[]byte(unhashedValue),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}
	hashedValue := string(hashBytes)

	// Store the hashed string and return the unhashed string for use in the email link
	model.Value = hashedValue
	return unhashedValue, database.Create(model).Error
}

// Reads a token by email where CreatedAt is no older than 24 hours.
// You're meant to check the read value against another hashed value to verify the token.
func (model *Token) Read(database *gorm.DB) error {
	return database.Order("id desc").First(model, "email = ? AND created_at > ?", model.Email, time.Now().Add(-24*time.Hour)).Error
}

// Delete deletes a token by email and value.
func (model *Token) Delete(database *gorm.DB) error {
	return database.Delete(model, "email = ? AND value = ?", model.Email, model.Value).Error
}
