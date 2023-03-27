package storage

import (
	"context"
	"errors"
	"fmt"
	"gophkeeper/internal/service"
	"gophkeeper/internal/tools"
	"gorm.io/gorm"
)

// RegisterUser puts user login and password to DB while checking that user.Login stays unique
func (dbStorage DBStorage) RegisterUser(user service.User, ctx context.Context) error {
	var dbUser service.User
	err := dbStorage.db.WithContext(ctx).Where("login = ?", user.Login).First(&dbUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			hashedPassword, passErr := tools.GeneratePasswordHash(user.Password)
			if passErr != nil {
				return fmt.Errorf("error in password hashing: %s", err)
			}
			user.Password = hashedPassword
			passErr = dbStorage.db.WithContext(ctx).Create(&user).Error
			if passErr != nil {
				return err
			}
			return nil
		}
		return err
	}
	return ErrUserExists
}

// CheckUserAuth checks if the login and password provided are valid
func (dbStorage DBStorage) CheckUserAuth(authDetails service.Authentication, ctx context.Context) error {
	var authUser service.User

	err := dbStorage.db.WithContext(ctx).Where("login  = 	?", authDetails.Login).First(&authUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrInvalidCredentials
		}
		return err
	}

	if !tools.CheckPasswordHash(authDetails.Password, authUser.Password) {
		return ErrInvalidCredentials
	}
	return nil
}

// PutLogoPass puts a secret logo-pass pair while checking that descriptions stays unique for a user
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

	logoPass.ID = checkEntry.ID

	if checkEntry.UpdatedAt.After(logoPass.UpdatedAt) {
		return ErrOldData
	}
	err = dbStorage.db.WithContext(ctx).Save(&logoPass).Error
	if err != nil {
		return err
	}

	return nil
}

// PutText puts a secret text data while checking that descriptions stays unique for a user
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

	secret.ID = checkEntry.ID

	if checkEntry.UpdatedAt.After(secret.UpdatedAt) {
		return ErrOldData
	}
	err = dbStorage.db.WithContext(ctx).Save(&secret).Error
	if err != nil {
		return err
	}

	return nil
}

// PutCreditCard puts a credit card data while checking that creditCard.Number stays unique for a user
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

	card.ID = checkEntry.ID

	if checkEntry.UpdatedAt.After(card.UpdatedAt) {
		return ErrOldData
	}
	err = dbStorage.db.WithContext(ctx).Save(&card).Error
	if err != nil {
		return err
	}

	return nil
}

// PutBinary puts an arbitrary binary data while checking that descriptions stays unique for a user
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

	binary.ID = checkEntry.ID

	if checkEntry.UpdatedAt.After(binary.UpdatedAt) {
		return ErrOldData
	}
	err = dbStorage.db.WithContext(ctx).Save(&binary).Error
	if err != nil {
		return err
	}

	return nil
}

// BatchGetLogoPasses is used to get the list of all user's stored logo-pass pairs
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

// BatchGetTexts is used to get the list of all user's stored secret texts
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

// BatchGetCreditCards is used to get the list of all user's stored secret texts
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

	err := dbStorage.db.WithContext(ctx).Table("binary_data").Select("id, login, description, updated_at").
		Where("login  = 	?", login).Find(&binaryList).Error
	if err != nil {
		return binaryList, err
	}
	if len(binaryList) == 0 {
		return binaryList, ErrEmpty
	}

	return binaryList, nil
}

func (dbStorage DBStorage) GetBinary(binary service.BinaryData, ctx context.Context) (service.BinaryData, error) {
	err := dbStorage.db.WithContext(ctx).Where("login  = 	? AND description = ?",
		binary.Login, binary.Description).First(&binary).Error
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
	dbStorage.db.Exec("DELETE FROM logo_passes")
	dbStorage.db.Exec("DELETE FROM text_data")
	dbStorage.db.Exec("DELETE FROM credit_cards")
	dbStorage.db.Exec("DELETE FROM binary_data")
}
