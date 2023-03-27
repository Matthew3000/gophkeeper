package client

import (
	"errors"
	"time"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"   envDefault:"http://localhost:8080"`
	OutputFolder  string `env:"OUTPUT_FOLDER"     envDefault:"C:/temp/gophkeeper"`
}

const dateTimeLayout = "02.01.2006 15:04:05"

const (
	logoPassFile    = "LogoPasses.json"
	textFile        = "TextData.json"
	creditCardFile  = "CreditCards.json"
	binaryListFile  = "BinaryList.json"
	UpdateDataTimer = 300 * time.Second
)

var (
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAlreadyExists      = errors.New("already exists")
	ErrEmpty              = errors.New("no data")
)
