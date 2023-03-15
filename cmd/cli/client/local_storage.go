package client

import (
	"encoding/json"
	"gophkeeper/internal/service"
	"os"
)

type Storage interface {
	StoreAllData(listLogoPasses []service.LogoPass, listTexts []service.TextData,
		listCreditCards []service.CreditCard, binaryList service.UserBinaryList) error
	UpdateLogoPass(logoPass service.LogoPass) error
	UpdateText(Text service.TextData) error
	UpdateCreditCard(CreditCard service.CreditCard) error
	UpdateBinary(binary service.BinaryData) error
}

type FileStorage struct {
	outputPath string
}

const (
	LogopassFile   = "LogoPasses.json"
	TextFile       = "TextData.json"
	CreditCardFile = "CreditCards.json"
	BinaryListFile = "BinaryList.json"
)

func NewStorage(path string) *FileStorage {
	return &FileStorage{outputPath: path}
}

func (storage *FileStorage) StoreAllData(listLogoPasses []service.LogoPass, listTexts []service.TextData,
	listCreditCards []service.CreditCard, binaryList service.UserBinaryList) error {

	err := os.MkdirAll(storage.outputPath, os.ModePerm)
	if err != nil {
		return err
	}

	jsonBytes, err := json.Marshal(listLogoPasses)
	if err != nil {
		return err
	}
	err = os.WriteFile(LogopassFile, jsonBytes, 0644)
	if err != nil {
		return err
	}

	jsonBytes, err = json.Marshal(listTexts)
	if err != nil {
		return err
	}
	err = os.WriteFile(TextFile, jsonBytes, 0644)
	if err != nil {
		return err
	}

	jsonBytes, err = json.Marshal(listCreditCards)
	if err != nil {
		return err
	}
	err = os.WriteFile(CreditCardFile, jsonBytes, 0644)
	if err != nil {
		return err
	}

	jsonBytes, err = json.Marshal(binaryList)
	if err != nil {
		return err
	}
	err = os.WriteFile(BinaryListFile, jsonBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (storage *FileStorage) UpdateLogoPass(logoPass service.LogoPass) error {
	return nil
}

func (storage *FileStorage) UpdateText(Text service.TextData) error {
	return nil
}

func (storage *FileStorage) UpdateCreditCard(CreditCard service.CreditCard) error {
	return nil
}

func (storage *FileStorage) UpdateBinary(binary service.BinaryData) error {
	return nil
}
