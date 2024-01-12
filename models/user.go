package models

import (
	"html"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name        string `gorm:"size:255;not null;unique" json:"name"`
	Username    string `gorm:"size:255;not null;unique" json:"username"`
	Email       string `gorm:"size:255;not null;unique" json:"email"`
	Password    string `gorm:"size:255;not null;unique" json:"-"`
	Gender      string `gorm:"size:10;not null" json:"gender"`
	PackageType string `gorm:"size:10;default:'standard'" json:"package_type"`

	// relationship
	SwipeUserHistories []SwipeUserHistory `json:"-"`
}

type UserResult struct {
	Name        string `gorm:"size:255;not null;unique" json:"name"`
	Username    string `gorm:"size:255;not null;unique" json:"username"`
	Email       string `gorm:"size:255;not null;unique" json:"email"`
	Gender      string `gorm:"size:10;not null" json:"gender"`
	PackageType string `gorm:"size:10;default:'standard'" json:"package_type"`
}

func (u *User) SaveUser(db *gorm.DB) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return &User{}, err
	}
	u.Password = string(hashedPassword)
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))
	if err := db.Save(&u).Error; err != nil {
		return &User{}, err
	}
	return u, nil
}
