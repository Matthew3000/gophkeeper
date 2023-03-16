package storage

import (
	"context"
	"errors"
	"gophkeeper/internal/service"
)

type UserStorage interface {
	CheckUserAuth(authDetails service.Authentication, ctx context.Context) error
	RegisterUser(user service.User, ctx context.Context) error
	PutLogoPass(logoPass service.LogoPass, ctx context.Context) error
	BatchGetLogoPasses(login string, ctx context.Context) ([]service.LogoPass, error)
	PutText(secret service.TextData, ctx context.Context) error
	BatchGetTexts(login string, ctx context.Context) ([]service.TextData, error)
	PutCreditCard(card service.CreditCard, ctx context.Context) error
	BatchGetCreditCards(login string, ctx context.Context) ([]service.CreditCard, error)
	PutBinary(binary service.BinaryData, ctx context.Context) error
	GetBinaryList(login string, ctx context.Context) ([]service.BinaryData, error)
	GetBinary(binary service.BinaryData, ctx context.Context) (service.BinaryData, error)
	DeleteAll()
}

var (
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAlreadyExists      = errors.New("already exists")
	ErrEmpty              = errors.New("no data")
	ErrOldData            = errors.New("newer data available on remote storage")
)
