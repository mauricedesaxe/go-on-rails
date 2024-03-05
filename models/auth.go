package model

import (
	"crypto/rand"
	"errors"
	"net/mail"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserModel struct {
	gorm.Model
	ID       uint   `gorm:"primaryKey"`
	Email    string `gorm:"unique"`
	Password string
}

// Note that this hashes the password before storing the user.
func (model *UserModel) Create(database *gorm.DB) error {
	err := ValidateUserInput(model)
	if err != nil {
		return err
	}
	hashed, err := Hash(model.Password)
	if err != nil {
		return err
	}
	model.Password = hashed
	return database.Create(model).Error
}

func (model *UserModel) ReadAll(database *gorm.DB) ([]UserModel, error) {
	var users []UserModel
	err := database.Find(&users).Error
	return users, err
}

func (model *UserModel) Read(database *gorm.DB) error {
	return database.First(model, model.ID).Error
}

func (model *UserModel) ReadByEmail(database *gorm.DB) error {
	return database.Where("email = ?", model.Email).First(model).Error
}

// Note that this hashes the password before storing the user.
func (model *UserModel) Update(database *gorm.DB) error {
	err := ValidateUserInput(model)
	if err != nil {
		return err
	}
	hashed, err := Hash(model.Password)
	if err != nil {
		return err
	}
	model.Password = hashed
	return database.Save(model).Error
}

func (model *UserModel) Delete(database *gorm.DB) error {
	return database.Delete(model).Error
}

// Can help with email verification, password reset and magic link login. It's meant as a way to
// verify that the user has access to the email of the account.
type TokenModel struct {
	gorm.Model
	Email string `gorm:"primaryKey"`
	Value string `gorm:"primaryKey"`
}

// Note that this hashes the token value before storing the token. The unhashed value is returned for
// use in the email link.
func (model *TokenModel) Create(database *gorm.DB) (string, error) {
	// Validate the email
	_, err := mail.ParseAddress(model.Email)
	if err != nil {
		return "", errors.New("invalid email address")
	}

	// Generate a random 32 character string using crypto/rand for better randomness
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789@#%^&*"
	b := make([]byte, 32)
	_, err = rand.Read(b)
	if err != nil {
		return "", err
	}
	for i := range b {
		b[i] = charset[b[i]%byte(len(charset))]
	}
	val := string(b)

	// Hash the random string
	hashed, err := Hash(val)
	if err != nil {
		return "", err
	}

	// Store the hashed string and return the unhashed string for use in the email link
	model.Value = hashed
	return val, database.Create(model).Error
}

// Reads a token by email where CreatedAt is no older than 24 hours.
// You're meant to check the read value against another hashed value to verify the token.
func (model *TokenModel) Read(database *gorm.DB) error {
	return database.First(model, "email = ? AND created_at > ?", model.Email, time.Now().Add(-24*time.Hour)).Error
}

// Delete deletes a token by email and value.
func (model *TokenModel) Delete(database *gorm.DB) error {
	return database.Delete(model, "email = ? AND value = ?", model.Email, model.Value).Error
}

// ValidateUserInput checks if the user input meets the requirements.
func ValidateUserInput(user *UserModel) error {
	_, err := mail.ParseAddress(user.Email)
	if err != nil {
		return errors.New("invalid email address")
	}
	if len(user.Password) < 16 {
		return errors.New("password too short, must be at least 16 characters")
	}
	if len(user.Password) > 128 {
		return errors.New("password too long, must be at most 128 characters")
	}
	// check password to have a mix of upper/lower case, numbers and special characters
	upper := false
	lower := false
	number := false
	special := false
	for _, c := range user.Password {
		switch {
		case 'A' <= c && c <= 'Z':
			upper = true
		case 'a' <= c && c <= 'z':
			lower = true
		case '0' <= c && c <= '9':
			number = true
		case strings.ContainsRune("!@#$%^&*()-_=+[]{};:,.<>?/\\|", c):
			special = true
		}
	}
	if !(upper && lower && number && special) {
		return errors.New("password too weak, must contain upper/lower case, numbers and special characters")
	}
	return nil
}

// Hashes the input string using bcrypt, a one-way hashing algorithm.
func Hash(input string) (output string, err error) {
	if len(input) > 0 {
		hashBytes, err := bcrypt.GenerateFromPassword(
			[]byte(input),
			bcrypt.DefaultCost,
		)
		if err != nil {
			return "", err
		}
		output = string(hashBytes)
	}
	return output, nil
}
