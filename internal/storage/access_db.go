package storage

import (
	"context"
	"errors"
	"fmt"
	"gophkeeper/internal/service"
	"gorm.io/gorm"
	"time"
)

func (dbStorage DBStorage) RegisterUser(user service.User, ctx context.Context) error {
	var dbUser service.User
	err := dbStorage.db.WithContext(ctx).Where("login = ?", user.Login).First(&dbUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			hashedPassword, err := service.GeneratePasswordHash(user.Password)
			if err != nil {
				return fmt.Errorf("error in password hashing: %s", err)
			}
			user.Password = hashedPassword
			err = dbStorage.db.WithContext(ctx).Create(&user).Error
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}
	return ErrUserExists
}

func (dbStorage DBStorage) CheckUserAuth(authDetails service.Authentication, ctx context.Context) error {
	var authUser service.User

	err := dbStorage.db.WithContext(ctx).Where("login  = 	?", authDetails.Login).First(&authUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrInvalidCredentials
		}
		return err
	}

	if !service.CheckPasswordHash(authDetails.Password, authUser.Password) {
		return ErrInvalidCredentials
	}
	return nil
}

func (dbStorage DBStorage) PutLogoPass(logoPass service.LogoPass, ctx context.Context) error {
	logoPass.UpdatedAt = time.Now()
	var checkEntry service.LogoPass

	err := dbStorage.db.WithContext(ctx).Where("login  = 	?  AND description = ?",
		logoPass.Login, logoPass.Description).First(&checkEntry).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}
	if !logoPass.Overwrite {
		if checkEntry.Login != "" {
			return ErrAlreadyExists
		}
	}

	if checkEntry.UpdatedAt.After(logoPass.UpdatedAt) {
		return ErrOldData
	}
	err = dbStorage.db.WithContext(ctx).Save(&logoPass).Error
	if err != nil {
		return err
	}

	return nil
}

func (dbStorage DBStorage) PutText(secret service.TextData, ctx context.Context) error {
	secret.UpdatedAt = time.Now()
	var checkEntry service.LogoPass

	err := dbStorage.db.WithContext(ctx).Where("login  = 	?  AND description = ?",
		secret.Login, secret.Description).First(&checkEntry).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}
	if !secret.Overwrite {
		if checkEntry.Login != "" {
			return ErrAlreadyExists
		}
	}

	if checkEntry.UpdatedAt.After(secret.UpdatedAt) {
		return ErrOldData
	}
	err = dbStorage.db.WithContext(ctx).Save(&secret).Error
	if err != nil {
		return err
	}

	return nil
}

func (dbStorage DBStorage) PutCreditCard(card service.CreditCard, ctx context.Context) error {
	card.UpdatedAt = time.Now()
	var checkEntry service.LogoPass

	err := dbStorage.db.WithContext(ctx).Where("login  = 	?  AND number = ?",
		card.Login, card.Number).First(&checkEntry).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}
	if !card.Overwrite {
		if checkEntry.Login != "" {
			return ErrAlreadyExists
		}
	}

	if checkEntry.UpdatedAt.After(card.UpdatedAt) {
		return ErrOldData
	}
	err = dbStorage.db.WithContext(ctx).Save(&card).Error
	if err != nil {
		return err
	}

	return nil
}

func (dbStorage DBStorage) PutBinary(binary service.BinaryData, ctx context.Context) error {
	binary.UpdatedAt = time.Now()
	var checkEntry service.LogoPass

	err := dbStorage.db.WithContext(ctx).Where("login  = 	?  AND description = ?",
		binary.Login, binary.Description).First(&checkEntry).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}
	if !binary.Overwrite {
		if checkEntry.Login != "" {
			return ErrAlreadyExists
		}
	}

	if checkEntry.UpdatedAt.After(binary.UpdatedAt) {
		return ErrOldData
	}
	err = dbStorage.db.WithContext(ctx).Save(&binary).Error
	if err != nil {
		return err
	}

	return nil
}

func (dbStorage DBStorage) DeleteAll() {
	dbStorage.db.Exec("DELETE FROM users")
	dbStorage.db.Exec("DELETE FROM logoPass")
	dbStorage.db.Exec("DELETE FROM textData")
	dbStorage.db.Exec("DELETE FROM creditCard")
	dbStorage.db.Exec("DELETE FROM userBinaryList")
	dbStorage.db.Exec("DELETE FROM binaryData")
}
