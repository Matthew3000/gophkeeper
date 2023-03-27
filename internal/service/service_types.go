package service

import (
	"gorm.io/gorm"
)

// LogoPass struct holds the secret login and password pair along with description. Overwrite flag used for api.
type LogoPass struct {
	gorm.Model
	Login       string `json:"-"`
	SecretLogin string `json:"secret_login"`
	SecretPass  string `json:"secret"`
	Description string `json:"description"`
	Overwrite   bool   `json:"overwrite" gorm:"-"`
}

// TextData struct holds the secret string data of any kind along with description. Overwrite flag used for api.
type TextData struct {
	gorm.Model
	Login       string `json:"-"`
	Text        string `json:"data"`
	Description string `json:"description"`
	Overwrite   bool   `json:"overwrite" gorm:"-"`
}

// CreditCard struct holds the credit card data along with description. Overwrite flag used for api.
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

// BinaryData struct holds arbitrary binary data along with description. Overwrite flag used for api.
type BinaryData struct {
	gorm.Model
	Login       string `json:"-"`
	Binary      string `json:"binary"`
	Description string `json:"description"`
	Overwrite   bool   `json:"overwrite" gorm:"-"`
}
