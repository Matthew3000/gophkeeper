package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gophkeeper/internal/app"
	"gophkeeper/internal/service"
	"net/http"
)

type Api interface {
	Register(user service.User) error
	Login(user service.User) error
	GetLogoPasses() ([]service.LogoPass, error)
	GetTexts() ([]service.TextData, error)
	GetCreditCards() ([]service.CreditCard, error)
	GetBinaryList() (service.UserBinaryList, error)
	GetBinary() (service.BinaryData, error)
	PutLogoPass(logoPass service.LogoPass) error
	PutText(text service.TextData) error
	PutCreditCard(card service.CreditCard) error
	PutBinary(binary service.BinaryData) error
}

type ServerApi struct {
	BaseURL string
}

func NewApi(url string) *ServerApi {
	return &ServerApi{BaseURL: url}
}

func (api ServerApi) Register(user service.User) error {
	req, err := json.Marshal(user)
	if err != nil {
		return err
	}

	resp, err := http.Post(api.BaseURL+app.RegisterEndpoint, "application/octet-stream", bytes.NewReader(req))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status code %d", resp.StatusCode)
	}
	return nil
}

func (api ServerApi) Login(user service.User) error {
	req, err := json.Marshal(user)
	if err != nil {
		return err
	}

	resp, err := http.Post(api.BaseURL+app.LoginEndpoint, "application/octet-stream", bytes.NewReader(req))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status code %d", resp.StatusCode)
	}
	return nil
}

func (api ServerApi) GetLogoPasses() ([]service.LogoPass, error) {
	var listLogoPasses []service.LogoPass
	resp, err := http.Get(api.BaseURL + app.GetLogoPassesEndpoint)
	if err != nil {
		return listLogoPasses, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&listLogoPasses)
	if err != nil {
		return listLogoPasses, err
	}
	return listLogoPasses, nil
}

func (api ServerApi) GetTexts() ([]service.TextData, error) {
	var listTexts []service.TextData
	resp, err := http.Get(api.BaseURL + app.GetTextsEndpoint)
	if err != nil {
		return listTexts, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&listTexts)
	if err != nil {
		return listTexts, err
	}
	return listTexts, nil
}

func (api ServerApi) GetCreditCards() ([]service.CreditCard, error) {
	var listCreditCards []service.CreditCard
	resp, err := http.Get(api.BaseURL + app.GetCreditCardsEndpoint)
	if err != nil {
		return listCreditCards, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&listCreditCards)
	if err != nil {
		return listCreditCards, err
	}
	return listCreditCards, nil
}

func (api ServerApi) GetBinaryList() (service.UserBinaryList, error) {
	var binaryList service.UserBinaryList
	resp, err := http.Get(api.BaseURL + app.GetBinaryListEndpoint)
	if err != nil {
		return binaryList, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&binaryList)
	if err != nil {
		return binaryList, err
	}
	return binaryList, nil
}

func (api ServerApi) GetBinary() (service.BinaryData, error) {
	var binary service.BinaryData
	resp, err := http.Get(api.BaseURL + app.GetBinaryEndpoint)
	if err != nil {
		return binary, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&binary)
	if err != nil {
		return binary, err
	}
	return binary, nil
}

func (api ServerApi) PutLogoPass(logoPass service.LogoPass) error {
	jsonBody, err := json.Marshal(logoPass)
	if err != nil {
		return err
	}

	resp, err := http.Post(api.BaseURL+app.PutLogoPassEndpoint, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status code %d", resp.StatusCode)
	}

	fmt.Printf("login password pair has been succesfully updated")
	return nil
}

func (api ServerApi) PutText(text service.TextData) error {
	jsonBody, err := json.Marshal(text)
	if err != nil {
		return err
	}

	resp, err := http.Post(api.BaseURL+app.PutTextEndpoint, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status code %d", resp.StatusCode)
	}

	fmt.Printf("login password pair has been succesfully updated")
	return nil
}

func (api ServerApi) PutCreditCard(card service.CreditCard) error {
	jsonBody, err := json.Marshal(card)
	if err != nil {
		return err
	}

	resp, err := http.Post(api.BaseURL+app.PutCreditCardEndpoint, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status code %d", resp.StatusCode)
	}

	fmt.Printf("login password pair has been succesfully updated")
	return nil
}

func (api ServerApi) PutBinary(binary service.BinaryData) error {
	jsonBody, err := json.Marshal(binary)
	if err != nil {
		return err
	}

	resp, err := http.Post(api.BaseURL+app.PutBinaryEndpoint, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status code %d", resp.StatusCode)
	}

	fmt.Printf("login password pair has been succesfully updated")
	return nil
}
