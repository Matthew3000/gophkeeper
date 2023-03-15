package client

import (
	"bufio"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"gophkeeper/internal/service"
	"os"
	"strconv"
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
	svc.config.OutputFolder += fmt.Sprintf("_%s/", login)
	fmt.Print("Authorization successful, updating, please wait")

	//TODO offline

	err = svc.UpdateAll()
	if err != nil {
		return err
	}
	fmt.Print("Update successful")
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
	svc.config.OutputFolder += fmt.Sprintf("_%s/", login)
	fmt.Print("Registration successful, please proceed")
	return nil
}

func (svc *LocalService) UpdateAll() error {
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
	listLogoPasses, err := svc.storage.GetLogoPasses()
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Login", "Password", "Description", "Last updated"})

	for _, logoPass := range listLogoPasses {
		row := []string{strconv.FormatUint(uint64(logoPass.ID), 10), logoPass.SecretLogin, logoPass.SecretPass,
			logoPass.Description, logoPass.UpdatedAt.Format(DateTimeLayout)}
		table.Append(row)
	}

	table.Render()
	fmt.Print("If you want to update any pair enter it's ID\n" +
		"otherwise type exit")
	reader := bufio.NewReader(os.Stdin)
	choice, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input: ", err)
	}
	switch choice {
	case "1":

	}
	return nil
}

func (svc *LocalService) ShowTexts() error {
	listTexts, err := svc.storage.GetTexts()
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Description", "Text", "Last updated"})

	for _, text := range listTexts {
		row := []string{strconv.FormatUint(uint64(text.ID), 10), text.Description, text.Text, text.UpdatedAt.Format(DateTimeLayout)}
		table.Append(row)
	}

	table.Render()
	return nil
}
func (svc *LocalService) ShowCreditCards() error {
	listCreditCards, err := svc.storage.GetCreditCards()
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Number", "Holder", "Dew date", "CVV", "Description", "Last updated"})

	for _, card := range listCreditCards {
		row := []string{strconv.FormatUint(uint64(card.ID), 10), card.Number, card.Holder, card.DewDate, card.CVV,
			card.Description, card.UpdatedAt.Format(DateTimeLayout)}
		table.Append(row)
	}

	table.Render()
	return nil
}
func (svc *LocalService) ShowBinaryList() error {
	binaryList, err := svc.storage.GetBinaryList()
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Description", "Last updated"})

	i := 0
	for _, binary := range binaryList.BinaryList {
		row := []string{strconv.Itoa(i), binary}
		table.Append(row)
		i++
	}

	table.Render()
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
