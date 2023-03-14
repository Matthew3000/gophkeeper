package service

import (
	"gorm.io/gorm"
)

type LogoPass struct {
	gorm.Model
	Login       string `json:"-"`
	SecretLogin string `json:"secret_login"`
	SecretPass  string `json:"secret"`
	Description string `json:"description"`
	Overwrite   bool   `json:"overwrite" gorm:"-"`
}

type TextData struct {
	gorm.Model
	Login       string `json:"-"`
	Text        string `json:"data"`
	Description string `json:"description"`
	Overwrite   bool   `json:"overwrite" gorm:"-"`
}

type CreditCard struct {
	gorm.Model
	Login       string `json:"-"`
	Number      string `json:"number"`
	Holder      string `json:"holder"`
	DewDate     string `json:"dew_date"`
	CVV         string `json:"cvv"`
	Description string `json:"description,omitempty"`
	Overwrite   bool   `json:"overwrite" gorm:"-"`
}

type UserBinaryList struct {
	Login      string   `json:"-" gorm:"unique"`
	BinaryList []string `json:"binary_list"`
}

type BinaryData struct {
	gorm.Model
	Login       string `json:"-"`
	Binary      []byte `json:"binary"`
	Description string `json:"description"`
	Overwrite   bool   `json:"overwrite" gorm:"-"`
}
