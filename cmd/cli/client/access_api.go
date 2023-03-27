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

func (api *ServerApi) downloadData(url string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, api.BaseURL+url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.AddCookie(api.cookie)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}
	defer resp.Body.Close()

	jsonBytes, err := io.ReadAll(resp.Body)
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("read esponse body: %w", err)
	}
	return jsonBytes, nil
}

func (api *ServerApi) uploadData(data interface{}, url string) error {
	jsonBody, err := json.Marshal(data)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, api.BaseURL+url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.AddCookie(api.cookie)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		if resp.StatusCode == http.StatusConflict {
			return ErrAlreadyExists
		}
		return fmt.Errorf("server returned status code %d", resp.StatusCode)
	}
	return nil
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
		if resp.StatusCode == http.StatusConflict {
			return ErrUserExists
		}
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
		if resp.StatusCode == http.StatusUnauthorized {
			return ErrInvalidCredentials
		}
		return fmt.Errorf("server returned status code %d", resp.StatusCode)
	}
	return nil
}

func (api *ServerApi) GetLogoPasses() ([]service.LogoPass, error) {
	var listLogoPasses []service.LogoPass

	resp, err := api.downloadData(app.GetLogoPassesEndpoint)
	if err != nil {
		return listLogoPasses, err
	}

	if len(resp) != 0 {
		err = json.Unmarshal(resp, &listLogoPasses)
		if err != nil {
			return nil, fmt.Errorf("json unmarshall: %w", err)
		}
	}
	return listLogoPasses, nil
}

func (api *ServerApi) GetTexts() ([]service.TextData, error) {
	var listTexts []service.TextData

	resp, err := api.downloadData(app.GetTextsEndpoint)
	if err != nil {
		return listTexts, err
	}

	if len(resp) != 0 {
		err = json.Unmarshal(resp, &listTexts)
		if err != nil {
			return nil, fmt.Errorf("json unmarshall: %w", err)
		}
	}
	return listTexts, nil
}

func (api *ServerApi) GetCreditCards() ([]service.CreditCard, error) {
	var listCreditCards []service.CreditCard

	resp, err := api.downloadData(app.GetCreditCardsEndpoint)
	if err != nil {
		return listCreditCards, err
	}

	if len(resp) != 0 {
		err = json.Unmarshal(resp, &listCreditCards)
		if err != nil {
			return nil, fmt.Errorf("json unmarshall: %w", err)
		}
	}
	return listCreditCards, nil
}

func (api *ServerApi) GetBinaryList() ([]service.BinaryData, error) {
	var binaryList []service.BinaryData

	resp, err := api.downloadData(app.GetBinaryListEndpoint)
	if err != nil {
		return binaryList, err
	}

	if len(resp) != 0 {
		err = json.Unmarshal(resp, &binaryList)
		if err != nil {
			return nil, fmt.Errorf("json unmarshall: %w", err)
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
	if err != nil {
		return binary, err
	}

	req.AddCookie(api.cookie)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return binary, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return binary, ErrEmpty
	}

	err = json.NewDecoder(resp.Body).Decode(&binary)
	if err != nil {
		return binary, err
	}
	return binary, nil
}

func (api *ServerApi) UploadLogoPass(logoPass service.LogoPass) error {
	err := api.uploadData(logoPass, app.PutLogoPassEndpoint)
	if err != nil {
		return err
	}

	fmt.Println("Login password pair has been successfully updated")
	return nil
}

func (api *ServerApi) UploadText(text service.TextData) error {
	err := api.uploadData(text, app.PutTextEndpoint)
	if err != nil {
		return err
	}

	fmt.Println("The secret text has been successfully updated. What's it about, I wonder")
	return nil
}

func (api *ServerApi) UploadCreditCard(card service.CreditCard) error {
	err := api.uploadData(card, app.PutCreditCardEndpoint)
	if err != nil {
		return err
	}

	fmt.Println("The credit card info has been successfully updated")
	return nil
}

func (api *ServerApi) UploadBinary(binary service.BinaryData) error {
	err := api.uploadData(binary, app.PutBinaryEndpoint)
	if err != nil {
		return err
	}

	fmt.Println("Your binary data has been successfully updated")
	return nil
}
