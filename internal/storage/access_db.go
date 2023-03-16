package storage

import (
	"context"
	"errors"
	"fmt"
	"gophkeeper/internal/service"
	"gorm.io/gorm"
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
	var checkEntry service.TextData

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
	var checkEntry service.CreditCard

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
	var checkEntry service.BinaryData

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

func (dbStorage DBStorage) BatchGetLogoPasses(login string, ctx context.Context) ([]service.LogoPass, error) {
	var listLogoPasses []service.LogoPass

	err := dbStorage.db.WithContext(ctx).Where("login  = 	?", login).Find(&listLogoPasses).Error
	if len(listLogoPasses) == 0 {
		return nil, ErrEmpty
	}
	if err != nil {
		return nil, err
	}

	return listLogoPasses, nil
}

func (dbStorage DBStorage) BatchGetTexts(login string, ctx context.Context) ([]service.TextData, error) {
	var listTexts []service.TextData

	err := dbStorage.db.WithContext(ctx).Where("login  = 	?", login).Find(&listTexts).Error
	if len(listTexts) == 0 {
		return nil, ErrEmpty
	}
	if err != nil {
		return nil, err
	}

	return listTexts, nil
}

func (dbStorage DBStorage) BatchGetCreditCards(login string, ctx context.Context) ([]service.CreditCard, error) {
	var listCards []service.CreditCard

	err := dbStorage.db.WithContext(ctx).Where("login  = 	?", login).Find(&listCards).Error
	if len(listCards) == 0 {
		return nil, ErrEmpty
	}
	if err != nil {
		return nil, err
	}

	return listCards, nil
}

func (dbStorage DBStorage) GetBinaryList(login string, ctx context.Context) ([]service.BinaryData, error) {
	var binaryList []service.BinaryData

	err := dbStorage.db.WithContext(ctx).Table("binaryData").Select("id, description, updated_at").
		Where("login  = 	?", login).Find(&binaryList).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return binaryList, ErrEmpty
		}
		return binaryList, err
	}

	return binaryList, nil
}

func (dbStorage DBStorage) GetBinary(login string, ctx context.Context) (service.BinaryData, error) {
	var binary service.BinaryData
	err := dbStorage.db.WithContext(ctx).Where("login  = 	?", login).First(&binary).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return binary, ErrEmpty
		}
		return binary, err
	}
	return binary, nil
}

func (dbStorage DBStorage) DeleteAll() {
	dbStorage.db.Exec("DELETE FROM users")
	dbStorage.db.Exec("DELETE FROM logoPass")
	dbStorage.db.Exec("DELETE FROM textData")
	dbStorage.db.Exec("DELETE FROM creditCard")
	dbStorage.db.Exec("DELETE FROM binaryData")
}
