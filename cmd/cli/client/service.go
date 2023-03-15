package client

import (
	"fmt"
	"gophkeeper/internal/service"
)

type Service interface {
	Register(login, password string) error
	UpdateAll() error
	ShowLogoPasses() error
	ShowTexts() error
	ShowCreditCards() error
	ShowBinaryList() error
	UploadLogoPass() error
	UploadText() error
	UploadCreditCard() error
	UploadBinary() error
}

type LocalService struct {
	config  Config
	api     Api
	storage Storage
	key     string
}

func NewService(config Config, api Api, storage Storage) *LocalService {
	return &LocalService{config: config, api: api, storage: storage}
}

func (svc *LocalService) Auth(login, password string) error {
	var user service.User
	user.Login = login
	user.Password = password
	err := svc.api.Login(user)
	if err != nil {
		return err
	}
	svc.key = password
	fmt.Print("Authorization successful, updating, please wait")
	err = svc.UpdateAll()
	if err != nil {
		return err
	}
	return nil
}

func (svc *LocalService) Register(login, password string) error {
	var user service.User
	user.Login = login
	user.Password = password
	err := svc.api.Register(user)
	if err != nil {
		return err
	}
	svc.key = password
	fmt.Print("Registration successful, please proceed")
	return nil
}

func (svc *LocalService) UpdateAll() error {
	/*	var listLogoPasses []service.LogoPass
		var listTexts []service.TextData
		var listCreditCards []service.CreditCard
		var binaryList service.UserBinaryList
		var err error*/

	listLogoPasses, err := svc.api.GetLogoPasses()
	if err != nil {
		return err
	}
	listTexts, err := svc.api.GetTexts()
	if err != nil {
		return err
	}
	listCreditCards, err := svc.api.GetCreditCards()
	if err != nil {
		return err
	}
	binaryList, err := svc.api.GetBinaryList()
	if err != nil {
		return err
	}

	err = svc.storage.StoreAllData(listLogoPasses, listTexts, listCreditCards, binaryList)
	if err != nil {
		return err
	}

	return nil
}

func (svc *LocalService) ShowLogoPasses() error {

	return nil
}

func (svc *LocalService) ShowTexts() error {
	return nil
}
func (svc *LocalService) ShowCreditCards() error {
	return nil
}
func (svc *LocalService) ShowBinaryList() error {
	return nil
}
func (svc *LocalService) UploadLogoPass() error {
	return nil
}
func (svc *LocalService) UploadText() error {
	return nil
}
func (svc *LocalService) UploadCreditCard() error {
	return nil
}
func (svc *LocalService) UploadBinary() error {
	return nil
}
