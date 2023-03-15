package client

import "errors"

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"   envDefault:"localhost:8080"`
	OutputFolder  string `env:"OUTPUT_FOLDER"     envDefault:"/AppData"`
}

const DateTimeLayout = "02.01.2006 15:04:05"

const (
	LogopassFile   = "LogoPasses.json"
	TextFile       = "TextData.json"
	CreditCardFile = "CreditCards.json"
	BinaryListFile = "BinaryList.json"
)

var (
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAlreadyExists      = errors.New("already exists")
	ErrEmpty              = errors.New("no data")
	ErrOldData            = errors.New("newer data available on remote storage")
)
