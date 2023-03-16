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
	DueDate     string `json:"due_date"`
	CVV         string `json:"cvv"`
	Description string `json:"description,omitempty"`
	Overwrite   bool   `json:"overwrite" gorm:"-"`
}

type BinaryData struct {
	gorm.Model
	Login       string `json:"-"`
	Binary      string `json:"binary"`
	Description string `json:"description"`
	Overwrite   bool   `json:"overwrite" gorm:"-"`
}
