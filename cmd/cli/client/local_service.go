package client

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"gophkeeper/internal/service"
	"gophkeeper/internal/tools"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type LocalService struct {
	config  Config
	Api     Api
	storage Storage
	key     string
}

func NewService(config Config, api Api, storage Storage) *LocalService {
	return &LocalService{config: config, Api: api, storage: storage}
}

func (svc *LocalService) getAnswer(ask string) string {
	reader := bufio.NewReader(os.Stdin)
ask:
	fmt.Println(ask)
	answer, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input: ", err)
		goto ask
	}
	answer = strings.TrimRight(answer, "\r\n")
	return answer
}

func (svc *LocalService) StartCommunicate() error {
auth:
	choice := svc.getAnswer("login: type 1\nRegister: type 2")
	login := svc.getAnswer("Enter your login")
	password := svc.getAnswer("Enter your password")

	var err error
	switch choice {
	case "1":
		err = svc.Auth(login, password)
	case "2":
		err = svc.Register(login, password)
	default:
		fmt.Println("Houston we got a problem")
		goto auth
	}
	if err != nil {
		fmt.Println(err)
		goto auth
	}
	svc.getActionFromUser()
	return nil
}

func (svc *LocalService) getActionFromUser() {
initialActionChoice:

	choice := svc.getAnswer("What are we going to do today?\n" +
		"View all stored login password pairs: type 1\n" +
		"Create a new login password pair:     type 2\n" +
		"View all texts:                       type 3\n" +
		"Create a new text secret:             type 4\n" +
		"View all credit cards' info:          type 5\n" +
		"Create a new credit card entry:       type 6\n" +
		"Review all binary data:               type 7\n" +
		"Upload a new binary:                  type 8")

	var err error
	switch choice {
	case "1":
		err = svc.showLogoPasses()
	case "2":
		var newLogoPass service.LogoPass
		newLogoPass.Overwrite = false
		err = svc.PutLogoPass(newLogoPass)
	case "3":
		err = svc.showTexts()
	case "4":
		var newText service.TextData
		newText.Overwrite = false
		err = svc.PutText(newText)
	case "5":
		err = svc.showCreditCards()
	case "6":
		var newCreditCard service.CreditCard
		newCreditCard.Overwrite = false
		err = svc.PutCreditCard(newCreditCard)
	case "7":
		err = svc.showBinaryList()
	case "8":
		var newBinary service.BinaryData
		newBinary.Overwrite = false
		err = svc.putBinary(newBinary)
	default:
		fmt.Println("Houston we got a problem!")
		goto initialActionChoice
	}
	if err != nil {
		fmt.Println(err)
	}
	goto initialActionChoice
}

func (svc *LocalService) Auth(login, password string) error {
	var user service.User
	user.Login = login
	user.Password = password
	err := svc.Api.Login(user)
	if err != nil {
		return err
	}
	svc.key = password

	err = svc.storage.UpdatePath(fmt.Sprintf("%s/%s/", svc.config.OutputFolder, login))
	if err != nil {
		return err
	}
	fmt.Println("Authorization successful, updating, please wait")
	fmt.Printf("Output folder is %s/%s/\n", svc.config.OutputFolder, login)

	err = svc.UpdateAll()
	if err != nil {
		if errors.Is(err, &net.OpError{}) {
			fmt.Println("Update failed due to poor internet connection, continuing offline")
			return nil
		}
		return err
	}
	fmt.Println("Update successful")
	return nil
}

func (svc *LocalService) Register(login, password string) error {
	var user service.User
	user.Login = login
	user.Password = password
	err := svc.Api.Register(user)
	if err != nil {
		return err
	}
	svc.key = password
	err = svc.storage.UpdatePath(fmt.Sprintf("%s/%s/", svc.config.OutputFolder, login))
	if err != nil {
		return err
	}
	fmt.Println("Registration successful, please proceed")
	fmt.Printf("Output folder is %s/%s/\n", svc.config.OutputFolder, login)
	return nil
}

func (svc *LocalService) UpdateAll() error {
	// means no auth yet, so no update required
	if svc.key == "" {
		return nil
	}

	listLogoPasses, err := svc.Api.GetLogoPasses()
	if err != nil {
		return fmt.Errorf("get logopass: %w", err)
	}
	updLogoPasses, err := svc.storage.StoreLogoPasses(listLogoPasses)
	if err != nil {
		return fmt.Errorf("store logopass: %w", err)
	}
	for _, logoPass := range updLogoPasses {
		err = svc.Api.UploadLogoPass(logoPass)
		if err != nil {
			return err
		}
	}

	listTexts, err := svc.Api.GetTexts()
	if err != nil {
		return fmt.Errorf("get texts: %w", err)
	}
	updTexts, err := svc.storage.StoreTexts(listTexts)
	if err != nil {
		return fmt.Errorf("store cards: %w", err)
	}
	for _, text := range updTexts {
		err = svc.Api.UploadText(text)
		if err != nil {
			return err
		}
	}

	listCreditCards, err := svc.Api.GetCreditCards()
	if err != nil {
		return fmt.Errorf("get cards: %w", err)
	}
	updCreditCards, err := svc.storage.StoreCreditCards(listCreditCards)
	if err != nil {
		return fmt.Errorf("store texts: %w", err)
	}
	for _, card := range updCreditCards {
		err = svc.Api.UploadCreditCard(card)
		if err != nil {
			return err
		}
	}

	binaryList, err := svc.Api.GetBinaryList()
	if err != nil {
		return fmt.Errorf("get binarylist: %w", err)
	}
	err = svc.storage.StoreBinaries(binaryList)
	if err != nil {
		return fmt.Errorf("store binarylist: %w", err)
	}

	return nil
}

func (svc *LocalService) showLogoPasses() error {
updateLogoPass:
	listLogoPasses, err := svc.storage.GetLogoPasses()
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "login", "Password", "Description", "Last updated"})

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
			logoPass.Description, logoPass.UpdatedAt.Format(dateTimeLayout)}
		table.Append(row)
	}
	table.Render()

	choice := svc.getAnswer("If you want to update any pair enter it's ID\notherwise type exit")
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
		err = svc.PutLogoPass(updLogoPass)
		if err != nil {
			fmt.Println(err)
			goto updateLogoPass
		}
	}
	return nil
}

func (svc *LocalService) PutLogoPass(logoPass service.LogoPass) error {
	logoPass.SecretLogin = svc.getAnswer("Please, enter login")
	logoPass.SecretPass = svc.getAnswer("Please, enter password")

	if logoPass.Description == "" {
		logoPass.Description = svc.getAnswer("Please, enter description for the pair")
	}

	var err error
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
	err = svc.Api.UploadLogoPass(logoPass)
	if err != nil {
		return err
	}
	fmt.Println("Successfully saved to remote")
	return nil
}

func (svc *LocalService) showTexts() error {
updateText:
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

		row := []string{strconv.FormatUint(uint64(text.ID), 10), text.Description, text.Text, text.UpdatedAt.Format(dateTimeLayout)}
		table.Append(row)
	}
	table.Render()

	choice := svc.getAnswer("If you want to update any text enter it's ID\notherwise type exit")
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
		err = svc.PutText(updText)
		if err != nil {
			fmt.Println(err)
			goto updateText
		}
	}

	return nil
}

func (svc *LocalService) PutText(text service.TextData) error {
	text.Text = svc.getAnswer("Please, enter text")
	if text.Description == "" {
		text.Description = svc.getAnswer("Please, enter description for the text")
	}

	var err error
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
	err = svc.Api.UploadText(text)
	if err != nil {
		return err
	}
	fmt.Println("Successfully saved to remote")
	return nil
}

func (svc *LocalService) showCreditCards() error {
updateCard:
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
			card.Description, card.UpdatedAt.Format(dateTimeLayout)}
		table.Append(row)
	}
	table.Render()

	choice := svc.getAnswer("If you want to update any credit card info enter it's ID\notherwise type exit")
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
		err = svc.PutCreditCard(updCreditCard)
		if err != nil {
			fmt.Println(err)
			goto updateCard
		}
	}
	return nil
}

func (svc *LocalService) PutCreditCard(creditCard service.CreditCard) error {
	creditCard.Holder = svc.getAnswer("Please, enter card holder name")
	if creditCard.Number == "" {
		creditCard.Number = svc.getAnswer("Please, enter card number")
	}
	creditCard.DueDate = svc.getAnswer("Please, enter due date")
	creditCard.CVV = svc.getAnswer("Please, enter CVC/CVV code")
	creditCard.Description = svc.getAnswer("Please, enter description for the card")

	var err error
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
	err = svc.Api.UploadCreditCard(creditCard)
	if err != nil {
		return err
	}
	fmt.Println("Successfully saved to remote")
	return nil
}

func (svc *LocalService) showBinaryList() error {
updateBinary:
	binaryList, err := svc.storage.GetBinaryList()
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Description", "Last updated"})

	for _, binary := range binaryList {
		row := []string{strconv.FormatUint(uint64(binary.ID), 10), binary.Description, binary.UpdatedAt.Format(dateTimeLayout)}
		table.Append(row)
	}
	table.Render()

	choice := svc.getAnswer("If you want to update or download any binary enter it's ID\notherwise type exit")
	switch choice {
	case "exit":
		return nil
	default:
		var updBinary service.BinaryData
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

		choice = svc.getAnswer("If you want to update: type 1\n" +
			"if you want to download: type 2\n" +
			"otherwise: type exit")
		switch choice {
		case "exit":
			return nil
		case "1":
			updBinary.Overwrite = true
			err = svc.putBinary(updBinary)
			if err != nil {
				fmt.Println(err)
				goto updateBinary
			}
		case "2":
			err = svc.downloadBinary(updBinary)
			if err != nil {
				fmt.Println(err)
				goto updateBinary
			}
		default:
		}
	}
	return nil
}

func (svc *LocalService) putBinary(binary service.BinaryData) error {
	if binary.Description == "" {
		binary.Description = svc.getAnswer("Please, enter description for the binary")
	}
	path := svc.getAnswer("Please, enter a path to upload your binary data")

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

	err = svc.Api.UploadBinary(binary)
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

func (svc *LocalService) downloadBinary(binary service.BinaryData) error {
	path := svc.getAnswer("Please enter a path to folder where you want to save a binary")

	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}

	name := svc.getAnswer("Please enter a name for a file")

	binary, err = svc.Api.GetBinary(binary)
	if err != nil {
		return err
	}

	data, err := tools.DecryptString(binary.Binary, svc.key)
	if err != nil {
		return err
	}

	err = os.WriteFile(path+"/"+name, []byte(data), 0644)
	if err != nil {
		return err
	}
	return nil
}
