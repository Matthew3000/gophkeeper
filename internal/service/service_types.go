package service

import (
	"gorm.io/gorm"
)

type LogoPass struct {
	gorm.Model
	Login       string `json:"-"`
	SecretLogin string `json:"secret_login"`
	SecretPass  string `json:"secret"`
	Description string `json:"description,omitempty"`
}

type TextData struct {
	gorm.Model
	Login       string `json:"-"`
	Text        string `json:"data"`
	Description string `json:"description,omitempty"`
}

type CreditCard struct {
	gorm.Model
	Login       string `json:"-"`
	CardNumber  string `json:"card_number"`
	CardHolder  string `json:"card_holder"`
	DewDate     string `json:"dew_date"`
	CVV         string `json:"cvv"`
	Description string `json:"description,omitempty"`
}

type UserBinaryList struct {
	Login      string   `json:"-" gorm:"unique"`
	BinaryList []string `json:"binary_list"`
}

type BinaryData struct {
	gorm.Model
	Login       string `json:"-"`
	Binary      []byte `json:"binary"`
	Description string `json:"description,omitempty"`
}
