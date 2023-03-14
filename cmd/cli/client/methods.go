package client

import "fmt"

func Auth(login, password string) error {
	fmt.Print("Authorization successful, updating, please wait")
	err := UpdateAll()
	if err != nil {
		return err
	}
	return nil
}

func Register(login, password string) error {
	return nil
}

func UpdateAll() error {
	return nil
}

func ShowLogoPasses() error {
	return nil
}

func ShowTexts() error {
	return nil
}
func ShowCreditCards() error {
	return nil
}
func ShowBinaryList() error {
	return nil
}
func UploadLogoPass() error {
	return nil
}
func UploadText() error {
	return nil
}
func UploadCreditCard() error {
	return nil
}
func UploadBinary() error {
	return nil
}
