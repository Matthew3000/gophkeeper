package client

import (
	"encoding/json"
	"gophkeeper/internal/service"
	"io"
	"os"
)

type Storage interface {
	StoreAllData(listLogoPasses []service.LogoPass, listTexts []service.TextData,
		listCreditCards []service.CreditCard, binaryList service.UserBinaryList) error
	UpdateLogoPass(logoPass service.LogoPass) error
	UpdateText(Text service.TextData) error
	UpdateCreditCard(CreditCard service.CreditCard) error
	UpdateBinaryList(binary service.BinaryData) error
	GetLogoPasses() ([]service.LogoPass, error)
	GetTexts() ([]service.TextData, error)
	GetCreditCards() ([]service.CreditCard, error)
	GetBinaryList() (service.UserBinaryList, error)
}

type FileStorage struct {
	outputPath string
}

func NewStorage(path string) (*FileStorage, error) {

	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return nil, err
	}
	return &FileStorage{outputPath: path}, nil
}

func (storage *FileStorage) StoreAllData(listLogoPasses []service.LogoPass, listTexts []service.TextData,
	listCreditCards []service.CreditCard, binaryList service.UserBinaryList) error {

	jsonBytes, err := json.Marshal(listLogoPasses)
	if err != nil {
		return err
	}
	err = os.WriteFile(storage.outputPath+LogopassFile, jsonBytes, 0644)
	if err != nil {
		return err
	}

	jsonBytes, err = json.Marshal(listTexts)
	if err != nil {
		return err
	}
	err = os.WriteFile(storage.outputPath+TextFile, jsonBytes, 0644)
	if err != nil {
		return err
	}

	jsonBytes, err = json.Marshal(listCreditCards)
	if err != nil {
		return err
	}
	err = os.WriteFile(storage.outputPath+CreditCardFile, jsonBytes, 0644)
	if err != nil {
		return err
	}

	jsonBytes, err = json.Marshal(binaryList)
	if err != nil {
		return err
	}
	err = os.WriteFile(storage.outputPath+BinaryListFile, jsonBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (storage *FileStorage) UpdateLogoPass(logoPass service.LogoPass) error {
	file, err := os.Open(storage.outputPath + LogopassFile)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	var listLogoPasses []service.LogoPass
	err = json.Unmarshal(data, &listLogoPasses)
	if err != nil {
		return err
	}

	for _, existingLogoPass := range listLogoPasses {
		if existingLogoPass.Description == logoPass.Description {
			if logoPass.Overwrite {
				existingLogoPass = logoPass
			} else {
				return ErrAlreadyExists
			}
		} else {
			listLogoPasses = append(listLogoPasses, logoPass)
		}
	}

	jsonBytes, err := json.Marshal(listLogoPasses)
	if err != nil {
		return err
	}
	err = os.WriteFile(storage.outputPath+LogopassFile, jsonBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (storage *FileStorage) UpdateText(Text service.TextData) error {
	file, err := os.Open(storage.outputPath + TextFile)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	var listTexts []service.TextData
	err = json.Unmarshal(data, &listTexts)
	if err != nil {
		return err
	}

	for _, existingText := range listTexts {
		if existingText.Description == Text.Description {
			if Text.Overwrite {
				existingText = Text
			} else {
				return ErrAlreadyExists
			}
		} else {
			listTexts = append(listTexts, Text)
		}
	}

	jsonBytes, err := json.Marshal(listTexts)
	if err != nil {
		return err
	}
	err = os.WriteFile(storage.outputPath+TextFile, jsonBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (storage *FileStorage) UpdateCreditCard(CreditCard service.CreditCard) error {
	file, err := os.Open(storage.outputPath + CreditCardFile)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	var listCreditCards []service.CreditCard
	err = json.Unmarshal(data, &listCreditCards)
	if err != nil {
		return err
	}

	for _, existingCreditCard := range listCreditCards {
		if existingCreditCard.Number == CreditCard.Number {
			if CreditCard.Overwrite {
				existingCreditCard = CreditCard
			} else {
				return ErrAlreadyExists
			}
		} else {
			listCreditCards = append(listCreditCards, CreditCard)
		}
	}

	jsonBytes, err := json.Marshal(listCreditCards)
	if err != nil {
		return err
	}
	err = os.WriteFile(storage.outputPath+CreditCardFile, jsonBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (storage *FileStorage) UpdateBinaryList(binary service.BinaryData) error {
	file, err := os.Open(storage.outputPath + BinaryListFile)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	var binaryList service.UserBinaryList
	err = json.Unmarshal(data, &binaryList)
	if err != nil {
		return err
	}

	for _, existingBinary := range binaryList.BinaryList {
		if existingBinary == binary.Description {
			return ErrAlreadyExists
		} else {
			binaryList.BinaryList = append(binaryList.BinaryList, binary.Description)
		}
	}

	jsonBytes, err := json.Marshal(binaryList)
	if err != nil {
		return err
	}
	err = os.WriteFile(storage.outputPath+BinaryListFile, jsonBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (storage *FileStorage) GetLogoPasses() ([]service.LogoPass, error) {
	var listLogoPasses []service.LogoPass
	return listLogoPasses, nil
}
func (storage *FileStorage) GetTexts() ([]service.TextData, error) {
	var listTexts []service.TextData
	return listTexts, nil
}
func (storage *FileStorage) GetCreditCards() ([]service.CreditCard, error) {
	var listCreditCards []service.CreditCard
	return listCreditCards, nil
}
func (storage *FileStorage) GetBinaryList() (service.UserBinaryList, error) {
	var BinaryList service.UserBinaryList
	return BinaryList, nil
}
