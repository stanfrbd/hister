// SPDX-FileContributor: Adam Tauber <asciimoo@gmail.com>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package model

import (
	"crypto/rand"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrUserAlreadyExists = errors.New("user already exists")
)

type User struct {
	CommonFields
	Username string `gorm:"uniqueIndex" json:"username"`
	Password string `json:"-"`
	Token    string `json:"-"`
	IsAdmin  bool   `json:"is_admin"`
}

func CreateUser(username, password string, isAdmin bool) (*User, error) {
	var existing User
	if err := DB.Where("username = ?", username).First(&existing).Error; err == nil {
		return nil, ErrUserAlreadyExists
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	u := &User{Username: username, Password: string(hash), Token: rand.Text(), IsAdmin: isAdmin}
	return u, DB.Create(u).Error
}

func DeleteUser(username string) error {
	result := DB.Where("username = ?", username).Delete(&User{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}

func AuthenticateUser(username, password string) (*User, error) {
	var u User
	if err := DB.Where("username = ?", username).First(&u).Error; err != nil {
		return nil, ErrUserNotFound
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return nil, ErrInvalidPassword
	}
	return &u, nil
}

func GetUserByToken(token string) (*User, error) {
	var u User
	if err := DB.Where("token = ?", token).First(&u).Error; err != nil {
		return nil, ErrUserNotFound
	}
	return &u, nil
}

func RegenerateToken(userID uint) (string, error) {
	token := rand.Text()
	if err := DB.Model(&User{}).Where("id = ?", userID).Update("token", token).Error; err != nil {
		return "", err
	}
	return token, nil
}
