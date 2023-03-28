// Package service holds universal types used in the App for creating DB and sending requests
package service

import (
	"gorm.io/gorm"
)

// User struct holds unique App user
type User struct {
	gorm.Model
	Login    string `json:"login" gorm:"unique"`
	Password string `json:"password"`
}

// Authentication struct is same as user, but doesn't have gorm.Model. Used only for auth requests.
type Authentication struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
