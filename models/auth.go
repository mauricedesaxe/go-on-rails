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

type User struct {
	gorm.Model
	ID       uint   `gorm:"primaryKey"`
	Email    string `gorm:"unique"`
	Password string
}

// Note that this hashes the password before storing the user.
func (u *User) Create() error {
	err := ValidateUserInput(u)
	if err != nil {
		return err
	}
	hashed, err := Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = hashed
	return DB.Create(u).Error
}

func (u *User) Read() error {
	return DB.First(u, u.ID).Error
}

func (u *User) ReadByEmail() error {
	return DB.Where("email = ?", u.Email).First(u).Error
}

// Note that this hashes the password before storing the user.
func (u *User) Update() error {
	err := ValidateUserInput(u)
	if err != nil {
		return err
	}
	hashed, err := Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = hashed
	return DB.Save(u).Error
}

func (u *User) Delete() error {
	return DB.Delete(u).Error
}

// Can help with email verification, password reset and magic link login. It's meant as a way to
// verify that the user has access to the email of the account.
type Token struct {
	gorm.Model
	Email string `gorm:"primaryKey"`
	Value string `gorm:"primaryKey"`
}

// Note that this hashes the token value before storing the token. The unhashed value is returned for
// use in the email link.
func (t *Token) Create() (string, error) {
	// Validate the email
	_, err := mail.ParseAddress(t.Email)
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
	t.Value = hashed
	return val, DB.Create(t).Error
}

// Reads a token by email where CreatedAt is no older than 24 hours.
// You're meant to check the read value against another hashed value to verify the token.
func (t *Token) Read() error {
	return DB.First(t, "email = ? AND created_at > ?", t.Email, time.Now().Add(-24*time.Hour)).Error
}

// Delete deletes a token by email and value.
func (t *Token) Delete() error {
	return DB.Delete(t, "email = ? AND value = ?", t.Email, t.Value).Error
}

// ValidateUserInput checks if the user input meets the requirements.
func ValidateUserInput(u *User) error {
	_, err := mail.ParseAddress(u.Email)
	if err != nil {
		return errors.New("invalid email address")
	}
	if len(u.Password) < 16 {
		return errors.New("password too short, must be at least 16 characters")
	}
	if len(u.Password) > 128 {
		return errors.New("password too long, must be at most 128 characters")
	}
	// check password to have a mix of upper/lower case, numbers and special characters
	upper := false
	lower := false
	number := false
	special := false
	for _, c := range u.Password {
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
