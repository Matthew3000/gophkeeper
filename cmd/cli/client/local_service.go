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
	StartCommunicate() error
	Auth(login, password string) error
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

func (svc *LocalService) StartCommunicate() error {
	reader := bufio.NewReader(os.Stdin)

auth:
	fmt.Println("Login: type 1\nRegister: type 2")
	choice, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input: ", err)
	}

	fmt.Println("Enter your login")
	login, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input: ", err)
	}

	fmt.Println("Enter your password")
	password, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input: ", err)
	}

	switch choice {
	case "1":
		err = svc.Auth(login, password)
	case "2":
		err = svc.Register(login, password)
	default:
		fmt.Println("Houston we gor problem ")
		goto auth
	}
	if err != nil {
		fmt.Println(err)
		goto auth
	}

initialActionChoice:
	fmt.Print("What are we going to do today?\n" +
		"View all stored login password pairs: type 1\n" +
		"Create a new login password pair:     type 2\n" +
		"View all texts:                       type 3\n" +
		"Create a new text secret:             type 4\n" +
		"View all credit cards' info:          type 5\n" +
		"Create a new credit card entry:       type 6\n" +
		"Review all binary data:               type 7\n" +
		"Upload a new binary:                  type 8\n")
	choice, err = reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input: ", err)
	}
	switch choice {
	case "1":
		err = svc.ShowLogoPasses()
	case "2":
		var updLogoPass service.LogoPass
		updLogoPass.Overwrite = false
		err = svc.UploadLogoPass(updLogoPass)
	case "3":
		err = svc.ShowTexts()
	case "4":
		err = svc.UploadText()
	case "5":
		err = svc.ShowCreditCards()
	case "6":
		err = svc.UploadCreditCard()
	case "7":
		err = svc.ShowBinaryList()
	case "8":
		err = svc.UploadBinary()
	default:
		fmt.Println("Houston we got a problem ")
		goto initialActionChoice
	}
	if err != nil {
		fmt.Print(err)
	}
	goto initialActionChoice
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
	fmt.Println("Authorization successful, updating, please wait")

	//TODO offline

	err = svc.UpdateAll()
	if err != nil {
		return err
	}
	fmt.Println("Update successful")
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
	fmt.Println("Registration successful, please proceed")
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

	//TODO decryption

	for _, logoPass := range listLogoPasses {
		row := []string{strconv.FormatUint(uint64(logoPass.ID), 10), logoPass.SecretLogin, logoPass.SecretPass,
			logoPass.Description, logoPass.UpdatedAt.Format(DateTimeLayout)}
		table.Append(row)
	}

	table.Render()

update:
	fmt.Print("If you want to update any pair enter it's ID\n" +
		"otherwise type exit\n")
	reader := bufio.NewReader(os.Stdin)
	choice, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input: ", err)
	}
	switch choice {
	case "exit":
		return nil
	default:
		var updLogoPass service.LogoPass
		updLogoPass.Overwrite = true
		id, err := strconv.ParseUint(choice, 10, 32)
		if err != nil {
			fmt.Println("I feel you bro, but i just dont get it, please try again")
			goto update
		}
		updLogoPass.ID = uint(id)
		updLogoPass.Description = ""
		for _, logoPass := range listLogoPasses {
			if updLogoPass.ID == logoPass.ID {
				updLogoPass.Description = logoPass.Description
			}
		}
		if updLogoPass.Description == "" {
			fmt.Println("There is no such ID, try again")
			goto update
		}
		err = svc.UploadLogoPass(updLogoPass)
		if err != nil {
			fmt.Println(err)
			goto update
		}
	}
	return nil
}

func (svc *LocalService) UploadLogoPass(logoPass service.LogoPass) error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Please, enter login")
	var err error
	logoPass.SecretLogin, err = reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input: ", err)
	}
	fmt.Println("Please, enter password")
	logoPass.SecretPass, err = reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input: ", err)
	}
	if logoPass.Description == "" {
		fmt.Println("Please, enter description for the pair")
		logoPass.Description, err = reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input: ", err)
		}
	}

	//TODO encryption

	err = svc.api.PutLogoPass(logoPass)
	if err != nil {
		return err
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

func (svc *LocalService) UploadText() error {
	return nil
}
func (svc *LocalService) UploadCreditCard() error {
	return nil
}
func (svc *LocalService) UploadBinary() error {
	return nil
}
