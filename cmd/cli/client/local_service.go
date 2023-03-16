package client

import (
	"bufio"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"gophkeeper/internal/service"
	"gophkeeper/internal/tools"
	"os"
	"strconv"
	"time"
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
	UploadLogoPass(logoPass service.LogoPass) error
	UploadText(text service.TextData) error
	UploadCreditCard(creditCard service.CreditCard) error
	UploadBinary(binary service.BinaryData) error
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
		var newLogoPass service.LogoPass
		newLogoPass.Overwrite = false
		err = svc.UploadLogoPass(newLogoPass)
	case "3":
		err = svc.ShowTexts()
	case "4":
		var newText service.TextData
		newText.Overwrite = false
		err = svc.UploadText(newText)
	case "5":
		err = svc.ShowCreditCards()
	case "6":
		var newCreditCard service.CreditCard
		newCreditCard.Overwrite = false
		err = svc.UploadCreditCard(newCreditCard)
	case "7":
		err = svc.ShowBinaryList()
	case "8":
		var newBinary service.BinaryData
		newBinary.Overwrite = false
		err = svc.UploadBinary(newBinary)
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

	for _, logoPass := range listLogoPasses {

		logoPass.SecretPass, err = tools.DecryptString(logoPass.SecretPass, svc.key)
		if err != nil {
			return err
		}
		logoPass.SecretLogin, err = tools.DecryptString(logoPass.SecretLogin, svc.key)
		if err != nil {
			return err
		}

		row := []string{strconv.FormatUint(uint64(logoPass.ID), 10), logoPass.SecretLogin, logoPass.SecretPass,
			logoPass.Description, logoPass.UpdatedAt.Format(DateTimeLayout)}
		table.Append(row)
	}

	table.Render()

updateLogoPass:
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
			goto updateLogoPass
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
			goto updateLogoPass
		}
		err = svc.UploadLogoPass(updLogoPass)
		if err != nil {
			fmt.Println(err)
			goto updateLogoPass
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

	logoPass.SecretPass, err = tools.EncryptString(logoPass.SecretPass, svc.key)
	if err != nil {
		return err
	}
	logoPass.SecretLogin, err = tools.EncryptString(logoPass.SecretLogin, svc.key)
	if err != nil {
		return err
	}
	logoPass.UpdatedAt = time.Now()

	err = svc.storage.UpdateLogoPass(logoPass)
	if err != nil {
		return err
	}
	fmt.Println("Successfully saved to local storage")
	err = svc.api.PutLogoPass(logoPass)
	if err != nil {
		return err
	}
	fmt.Println("Successfully saved to remote")
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
		text.Text, err = tools.DecryptString(text.Text, svc.key)
		if err != nil {
			return err
		}

		row := []string{strconv.FormatUint(uint64(text.ID), 10), text.Description, text.Text, text.UpdatedAt.Format(DateTimeLayout)}
		table.Append(row)
	}

	table.Render()

updateText:
	fmt.Print("If you want to update any text enter it's ID\n" +
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
		var updText service.TextData
		updText.Overwrite = true
		id, err := strconv.ParseUint(choice, 10, 32)
		if err != nil {
			fmt.Println("I feel you bro, but i just dont get it, please try again")
			goto updateText
		}
		updText.ID = uint(id)
		updText.Description = ""
		for _, text := range listTexts {
			if updText.ID == text.ID {
				updText.Description = text.Description
			}
		}
		if updText.Description == "" {
			fmt.Println("There is no such ID, try again")
			goto updateText
		}
		err = svc.UploadText(updText)
		if err != nil {
			fmt.Println(err)
			goto updateText
		}
	}

	return nil
}

func (svc *LocalService) UploadText(text service.TextData) error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Please, enter text")
	var err error
	text.Text, err = reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input: ", err)
	}
	if text.Description == "" {
		fmt.Println("Please, enter description for the text")
		text.Description, err = reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input: ", err)
		}
	}

	text.Text, err = tools.EncryptString(text.Text, svc.key)
	if err != nil {
		return err
	}
	text.UpdatedAt = time.Now()

	err = svc.storage.UpdateText(text)
	if err != nil {
		return err
	}
	fmt.Println("Successfully saved to local storage")
	err = svc.api.PutText(text)
	if err != nil {
		return err
	}
	fmt.Println("Successfully saved to remote")
	return nil
}

func (svc *LocalService) ShowCreditCards() error {
	listCreditCards, err := svc.storage.GetCreditCards()
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Number", "Holder", "Due date", "CVV", "Description", "Last updated"})

	for _, card := range listCreditCards {

		card.Holder, err = tools.DecryptString(card.Holder, svc.key)
		if err != nil {
			return err
		}
		card.DueDate, err = tools.DecryptString(card.DueDate, svc.key)
		if err != nil {
			return err
		}
		card.CVV, err = tools.DecryptString(card.CVV, svc.key)
		if err != nil {
			return err
		}

		row := []string{strconv.FormatUint(uint64(card.ID), 10), card.Number, card.Holder, card.DueDate, card.CVV,
			card.Description, card.UpdatedAt.Format(DateTimeLayout)}
		table.Append(row)
	}

	table.Render()

updateCard:
	fmt.Print("If you want to update any credit card info enter it's ID\n" +
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
		var updCreditCard service.CreditCard
		updCreditCard.Overwrite = true
		id, err := strconv.ParseUint(choice, 10, 32)
		if err != nil {
			fmt.Println("I feel you bro, but i just dont get it, please try again")
			goto updateCard
		}
		updCreditCard.ID = uint(id)
		updCreditCard.Number = ""
		for _, card := range listCreditCards {
			if updCreditCard.ID == card.ID {
				updCreditCard.Number = card.Number
			}
		}
		if updCreditCard.Number == "" {
			fmt.Println("There is no such ID, try again")
			goto updateCard
		}
		err = svc.UploadCreditCard(updCreditCard)
		if err != nil {
			fmt.Println(err)
			goto updateCard
		}
	}
	return nil
}

func (svc *LocalService) UploadCreditCard(creditCard service.CreditCard) error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Please, enter card holder name")
	var err error
	creditCard.Holder, err = reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input: ", err)
	}
	if creditCard.Number == "" {
		fmt.Println("Please, enter card number")
		creditCard.Number, err = reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input: ", err)
		}
	}
	fmt.Println("Please, enter due date")
	creditCard.DueDate, err = reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input: ", err)
	}
	fmt.Println("Please, enter CVC/CVV code")
	creditCard.CVV, err = reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input: ", err)
	}
	fmt.Println("Please, enter description for the card")
	creditCard.Description, err = reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input: ", err)
	}

	creditCard.Holder, err = tools.EncryptString(creditCard.Holder, svc.key)
	if err != nil {
		return err
	}
	creditCard.DueDate, err = tools.EncryptString(creditCard.DueDate, svc.key)
	if err != nil {
		return err
	}
	creditCard.CVV, err = tools.EncryptString(creditCard.CVV, svc.key)
	if err != nil {
		return err
	}
	creditCard.UpdatedAt = time.Now()

	err = svc.storage.UpdateCreditCard(creditCard)
	if err != nil {
		return err
	}
	fmt.Println("Successfully saved to local storage")
	err = svc.api.PutCreditCard(creditCard)
	if err != nil {
		return err
	}
	fmt.Println("Successfully saved to remote")
	return nil
}

func (svc *LocalService) ShowBinaryList() error {
	binaryList, err := svc.storage.GetBinaryList()
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Description", "Last updated"})

	for _, binary := range binaryList {
		row := []string{strconv.FormatUint(uint64(binary.ID), 10), binary.Description, binary.UpdatedAt.Format(DateTimeLayout)}
		table.Append(row)
	}

	table.Render()

updateBinary:
	fmt.Print("If you want to update any credit card info enter it's ID\n" +
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
		var updBinary service.BinaryData
		updBinary.Overwrite = true
		id, err := strconv.ParseUint(choice, 10, 32)
		if err != nil {
			fmt.Println("I feel you bro, but i just dont get it, please try again")
			goto updateBinary
		}
		updBinary.ID = uint(id)
		updBinary.Description = ""
		for _, binary := range binaryList {
			if updBinary.ID == binary.ID {
				updBinary.Description = binary.Description
			}
		}
		if updBinary.Description == "" {
			fmt.Println("There is no such ID, try again")
			goto updateBinary
		}
		err = svc.UploadBinary(updBinary)
		if err != nil {
			fmt.Println(err)
			goto updateBinary
		}
	}
	return nil
}

func (svc *LocalService) UploadBinary(binary service.BinaryData) error {
	reader := bufio.NewReader(os.Stdin)
	var err error
	if binary.Description == "" {
		fmt.Println("Please, enter description for the binary")
		binary.Description, err = reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input: ", err)
		}
	}
	fmt.Println("Please, enter a path to upload your binary data")
	path, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input: ", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading file: ", err)
		return err
	}
	binary.Binary = string(content)

	binary.Binary, err = tools.EncryptString(binary.Binary, svc.key)
	if err != nil {
		return err
	}
	binary.UpdatedAt = time.Now()

	err = svc.api.PutBinary(binary)
	if err != nil {
		return err
	}
	fmt.Println("Successfully uploaded to remote")

	binary.Binary = ""
	err = svc.storage.UpdateBinaryList(binary)
	if err != nil {
		return err
	}
	fmt.Println("Successfully saved to local storage")
	return nil
}
