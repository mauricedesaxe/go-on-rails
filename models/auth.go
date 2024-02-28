package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique"`
	Email    string `gorm:"unique"`
	Password string
}

// Note that this hashes the password before storing the user.
func (u *User) Create() error {
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

func (u *User) ReadByUsername() error {
	return DB.Where("username = ?", u.Username).First(u).Error
}

func (u *User) ReadByEmail() error {
	return DB.Where("email = ?", u.Email).First(u).Error
}

// Note that this hashes the password before storing the user.
func (u *User) Update() error {
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

// Note that this hashes the token value before storing the token.
func (t *Token) Create() error {
	hashed, err := Hash(t.Value)
	if err != nil {
		return err
	}
	t.Value = hashed
	return DB.Create(t).Error
}

// Reads a token by email and value where CreatedAt is no older than 24 hours.
func (t *Token) Read() error {
	return DB.First(t, "email = ? AND value = ? AND created_at > ?", t.Email, t.Value, time.Now().Add(-24*time.Hour)).Error
}

// Delete deletes a token by email and value.
func (t *Token) Delete() error {
	return DB.Delete(t, "email = ? AND value = ?", t.Email, t.Value).Error
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

// Used to to keep the user logged in once authenticated.
type Session struct {
	gorm.Model
	ID     uint `gorm:"primaryKey"`
	UserID uint
}

func (s *Session) Create() error {
	return DB.Create(s).Error
}

func (s *Session) Read() error {
	return DB.First(s, s.ID).Error
}

func (s *Session) ReadByUserID() error {
	return DB.Where("user_id = ?", s.UserID).First(s).Error
}

func (s *Session) Delete() error {
	return DB.Delete(s).Error
}
