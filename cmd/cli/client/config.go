package client

import (
	"errors"
	"time"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"   envDefault:"http://localhost:8080"`
	OutputFolder  string `env:"OUTPUT_FOLDER"     envDefault:"C:/temp/gophkeeper"`
}

const DateTimeLayout = "02.01.2006 15:04:05"

const (
	LogoPassFile    = "LogoPasses.json"
	TextFile        = "TextData.json"
	CreditCardFile  = "CreditCards.json"
	BinaryListFile  = "BinaryList.json"
	UpdateDataTimer = 300 * time.Second
)

var (
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAlreadyExists      = errors.New("already exists")
	ErrEmpty              = errors.New("no data")
)
