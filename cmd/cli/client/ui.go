package client

import (
	"bufio"
	"fmt"
	"os"
)

func (svc *LocalService) Communicate() error {
	reader := bufio.NewReader(os.Stdin)

auth:
	fmt.Print("Login: type 1\nRegister: type 2")
	choice, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input: ", err)
	}

	fmt.Print("Enter your login")
	login, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input: ", err)
	}

	fmt.Print("Enter your password")
	password, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input: ", err)
	}

	switch choice {
	case "1":
		err = svc.Auth(login, password)
	case "2":
		err = svc.Register(login, password)
	}
	if err != nil {
		fmt.Print(err)
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
		err = svc.UploadLogoPass()
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
	}
	if err != nil {
		fmt.Print(err)
	}
	goto initialActionChoice
}
