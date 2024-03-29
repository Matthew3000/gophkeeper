package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gophkeeper/internal/app"
	"gophkeeper/internal/service"
	"io"
	"log"
	"net/http"
)

// Api is an interface of all api interactions needed for Gophkeeper
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

// ServerApi holds the url of remote and user cookie for requests
type ServerApi struct {
	BaseURL string
	cookie  *http.Cookie
}

// NewApi creates a new instance of ServerApi according to the settings
func NewApi(url string) *ServerApi {
	return &ServerApi{BaseURL: url, cookie: &http.Cookie{Name: "session.id", Value: ""}}
}

// downloadData sends a get request and returns []bytes of response
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

// uploadData sends a post request marshalling input data and checks for proper response
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

// Register sends post request with service.User with 'login' and 'password' fields
// returns any errors occurred in the process
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

// Login sends post request with service.User with 'login' and 'password' fields
// returns any errors occurred in the process
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

// GetLogoPasses sends a http.Get request and returns the list of service.LogoPass acquired from the remote
func (api *ServerApi) GetLogoPasses() ([]service.LogoPass, error) {
	var listLogoPasses []service.LogoPass

	resp, err := api.downloadData(app.GetLogoPassesEndpoint)
	if err != nil {
		return nil, err
	}

	if len(resp) != 0 {
		err = json.Unmarshal(resp, &listLogoPasses)
		if err != nil {
			return nil, fmt.Errorf("json unmarshall: %w", err)
		}
	}
	return listLogoPasses, nil
}

// GetTexts sends a http.Get request and returns the list of service.TextData acquired from the remote
func (api *ServerApi) GetTexts() ([]service.TextData, error) {
	var listTexts []service.TextData

	resp, err := api.downloadData(app.GetTextsEndpoint)
	if err != nil {
		return nil, err
	}

	if len(resp) != 0 {
		err = json.Unmarshal(resp, &listTexts)
		if err != nil {
			return nil, fmt.Errorf("json unmarshall: %w", err)
		}
	}
	return listTexts, nil
}

// GetCreditCards sends a http.Get request and returns the list of service.CreditCard acquired from the remote
func (api *ServerApi) GetCreditCards() ([]service.CreditCard, error) {
	var listCreditCards []service.CreditCard

	resp, err := api.downloadData(app.GetCreditCardsEndpoint)
	if err != nil {
		return nil, err
	}

	if len(resp) != 0 {
		err = json.Unmarshal(resp, &listCreditCards)
		if err != nil {
			return nil, fmt.Errorf("json unmarshall: %w", err)
		}
	}
	return listCreditCards, nil
}

// GetBinaryList sends a http.Get request and returns the list of service.BinaryData acquired from the remote
func (api *ServerApi) GetBinaryList() ([]service.BinaryData, error) {
	var binaryList []service.BinaryData

	resp, err := api.downloadData(app.GetBinaryListEndpoint)
	if err != nil {
		return nil, err
	}

	if len(resp) != 0 {
		err = json.Unmarshal(resp, &binaryList)
		if err != nil {
			return nil, fmt.Errorf("json unmarshall: %w", err)
		}
	}
	return binaryList, nil
}

// GetBinary sends a http.Get request and returns service.BinaryData acquired from the remote
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

// UploadLogoPass sends post request that contains service.LogoPass
func (api *ServerApi) UploadLogoPass(logoPass service.LogoPass) error {
	err := api.uploadData(logoPass, app.PutLogoPassEndpoint)
	if err != nil {
		return err
	}
	log.Println("login password pair has been successfully updated")
	return nil
}

// UploadText sends post request that contains service.TextData
func (api *ServerApi) UploadText(text service.TextData) error {
	err := api.uploadData(text, app.PutTextEndpoint)
	if err != nil {
		return err
	}
	log.Println("The secret text has been successfully updated. What's it about, I wonder")
	return nil
}

// UploadCreditCard sends post request that contains service.CreditCard
func (api *ServerApi) UploadCreditCard(card service.CreditCard) error {
	err := api.uploadData(card, app.PutCreditCardEndpoint)
	if err != nil {
		return err
	}
	log.Println("The credit card info has been successfully updated")
	return nil
}

// UploadBinary sends post request that contains service.BinaryData
func (api *ServerApi) UploadBinary(binary service.BinaryData) error {
	err := api.uploadData(binary, app.PutBinaryEndpoint)
	if err != nil {
		return err
	}
	log.Println("Your binary data has been successfully updated")
	return nil
}
