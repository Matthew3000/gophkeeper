package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gophkeeper/internal/app"
	"gophkeeper/internal/service"
	"io"
	"net/http"
)

type Api interface {
	Register(user service.User) error
	Login(user service.User) error
	GetLogoPasses() ([]service.LogoPass, error)
	GetTexts() ([]service.TextData, error)
	GetCreditCards() ([]service.CreditCard, error)
	GetBinaryList() ([]service.BinaryData, error)
	GetBinary(binary service.BinaryData) (service.BinaryData, error)
	UploadLogoPass(logoPass service.LogoPass) error
	UploadText(text service.TextData) error
	UploadCreditCard(card service.CreditCard) error
	UploadBinary(binary service.BinaryData) error
}

type ServerApi struct {
	BaseURL string
	cookie  *http.Cookie
}

func NewApi(url string) *ServerApi {
	return &ServerApi{BaseURL: url, cookie: &http.Cookie{Name: "session.id", Value: ""}}
}

func (api *ServerApi) Register(user service.User) error {
	req, err := json.Marshal(user)
	if err != nil {
		return err
	}

	resp, err := http.Post(api.BaseURL+app.RegisterEndpoint, "application/octet-stream", bytes.NewReader(req))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	cookies := resp.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "session.id" {
			api.cookie = cookie
			break
		}
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status code %d", resp.StatusCode)
	}
	return nil
}

func (api *ServerApi) Login(user service.User) error {
	req, err := json.Marshal(user)
	if err != nil {
		return err
	}

	resp, err := http.Post(api.BaseURL+app.LoginEndpoint, "application/octet-stream", bytes.NewReader(req))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	cookies := resp.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "session.id" {
			api.cookie = cookie
			break
		}
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status code %d", resp.StatusCode)
	}
	return nil
}

func (api *ServerApi) GetLogoPasses() ([]service.LogoPass, error) {
	var listLogoPasses []service.LogoPass

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, api.BaseURL+app.GetLogoPassesEndpoint, nil)
	req.AddCookie(api.cookie)
	resp, err := client.Do(req)
	if err != nil {
		return listLogoPasses, fmt.Errorf("request: %w", err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&listLogoPasses)
	if err != nil {
		if !errors.Is(err, io.EOF) {
			return listLogoPasses, fmt.Errorf("json decoder: %w", err)
		}
	}
	return listLogoPasses, nil
}

func (api *ServerApi) GetTexts() ([]service.TextData, error) {
	var listTexts []service.TextData

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, api.BaseURL+app.GetTextsEndpoint, nil)
	req.AddCookie(api.cookie)
	resp, err := client.Do(req)
	if err != nil {
		return listTexts, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&listTexts)
	if !errors.Is(err, io.EOF) {
		if err != nil {
			return listTexts, err
		}
	}
	return listTexts, nil
}

func (api *ServerApi) GetCreditCards() ([]service.CreditCard, error) {
	var listCreditCards []service.CreditCard

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, api.BaseURL+app.GetCreditCardsEndpoint, nil)
	req.AddCookie(api.cookie)
	resp, err := client.Do(req)
	if err != nil {
		return listCreditCards, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&listCreditCards)
	if err != nil {
		if !errors.Is(err, io.EOF) {
			return listCreditCards, err
		}
	}
	return listCreditCards, nil
}

func (api *ServerApi) GetBinaryList() ([]service.BinaryData, error) {
	var binaryList []service.BinaryData

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, api.BaseURL+app.GetBinaryListEndpoint, nil)
	req.AddCookie(api.cookie)
	resp, err := client.Do(req)
	if err != nil {
		return binaryList, fmt.Errorf("request: %w", err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&binaryList)
	if err != nil {
		if !errors.Is(err, io.EOF) {
			return binaryList, fmt.Errorf("json decoder: %w", err)
		}
	}
	return binaryList, nil
}

func (api *ServerApi) GetBinary(binary service.BinaryData) (service.BinaryData, error) {
	jsonBody, err := json.Marshal(binary)
	if err != nil {
		return binary, err
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, api.BaseURL+app.GetBinaryEndpoint, bytes.NewBuffer(jsonBody))
	req.AddCookie(api.cookie)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
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

func (api *ServerApi) UploadLogoPass(logoPass service.LogoPass) error {
	jsonBody, err := json.Marshal(logoPass)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, api.BaseURL+app.PutLogoPassEndpoint, bytes.NewBuffer(jsonBody))
	req.AddCookie(api.cookie)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("server returned status code %d", resp.StatusCode)
	}

	fmt.Println("login password pair has been succesfully updated")
	return nil
}

func (api *ServerApi) UploadText(text service.TextData) error {
	jsonBody, err := json.Marshal(text)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, api.BaseURL+app.PutTextEndpoint, bytes.NewBuffer(jsonBody))
	req.AddCookie(api.cookie)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("server returned status code %d", resp.StatusCode)
	}

	fmt.Println("login password pair has been succesfully updated")
	return nil
}

func (api *ServerApi) UploadCreditCard(card service.CreditCard) error {
	jsonBody, err := json.Marshal(card)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, api.BaseURL+app.PutCreditCardEndpoint, bytes.NewBuffer(jsonBody))
	req.AddCookie(api.cookie)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("server returned status code %d", resp.StatusCode)
	}

	fmt.Println("login password pair has been succesfully updated")
	return nil
}

func (api *ServerApi) UploadBinary(binary service.BinaryData) error {
	jsonBody, err := json.Marshal(binary)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, api.BaseURL+app.PutBinaryEndpoint, bytes.NewBuffer(jsonBody))
	req.AddCookie(api.cookie)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("server returned status code %d", resp.StatusCode)
	}

	fmt.Println("login password pair has been succesfully updated")
	return nil
}
