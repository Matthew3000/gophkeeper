package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/require"
	"gophkeeper/cmd/cli/client"
	"gophkeeper/internal/app"
	"gophkeeper/internal/config"
	"gophkeeper/internal/service"
	"gophkeeper/internal/storage"
	"gorm.io/gorm"
	"log"
	"os"
	"testing"
	"time"
)

// Чтобы запустить тест, необходимо поднять Postgres базу и задать настройки через флаги
// для запуска обоих тестов сразу использовать go test -p 1 ./.../

func TestClient(t *testing.T) {
	var serverCfg config.Config

	if err := env.Parse(&serverCfg); err != nil {
		log.Fatal(err)
	}
	flag.StringVar(&serverCfg.DatabaseDSN, "d", serverCfg.DatabaseDSN, "File Storage Path")
	flag.StringVar(&serverCfg.ServerAddress, "a", serverCfg.ServerAddress, "Server address")
	log.Println(serverCfg.DatabaseDSN)

	var secretKey = "watch?v=Qw4w9WgXcQ"
	userStorage := storage.NewUserStorage(serverCfg.DatabaseDSN)
	cookieStorage := sessions.NewCookieStore([]byte(secretKey))
	var application = app.NewApp(serverCfg, userStorage, *cookieStorage)
	go application.Run()

	var cfg client.Config

	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}
	flag.StringVar(&cfg.OutputFolder, "o", cfg.OutputFolder, "Output folder for files")
	cfg.ServerAddress = "http://" + serverCfg.ServerAddress
	flag.Parse()

	var api = client.NewApi(cfg.ServerAddress)
	var clientStorage, err = client.NewStorage(cfg.OutputFolder)
	if err != nil {
		log.Fatal(err)
	}
	var clientService = client.NewService(cfg, api, clientStorage)

	application.UserStorage.DeleteAll()
	err = clientStorage.ClearAll()
	if err != nil {
		log.Fatal(err)
	}

	// Suppress stdin and std out
	r, w, _ := os.Pipe()
	oldStdin := os.Stdin
	oldStdout := os.Stdout
	os.Stdin = r
	os.Stdout = w

	defer func() {
		os.Stdin = oldStdin
		os.Stdout = oldStdout
	}()

	RegisterTest(t, clientService)
	LoginTest(t, clientService)
	PutLogoPassTest(t, clientService)
	PutTextTest(t, clientService)
	PutCreditCardTest(t, clientService)
	PutBinaryTest(t, clientService)
	GetBinary(t, clientService)
	GetLogoPassesTest(t, clientService)
	GetTextsTest(t, clientService)
	GetCreditCardsTest(t, clientService)
	GetBinaryListTest(t, clientService)

	application.UserStorage.DeleteAll()
	err = clientStorage.ClearAll()
	if err != nil {
		log.Fatal(err)
	}
}

func RegisterTest(t *testing.T, svc *client.LocalService) {
	tests := []struct {
		name string
		data service.User
		want error
	}{
		{
			name: "register ok",
			data: service.User{
				Login:    "Major Tom",
				Password: "Ground control",
			},
		},
		{
			name: "register fail",
			data: service.User{
				Login:    "Major Tom",
				Password: "Ground control",
			},
			want: client.ErrUserExists,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.Register(tt.data.Login, tt.data.Password)
			if errors.Is(err, tt.want) {
				err = nil
			}
			require.NoError(t, err)
		})
	}
}

func LoginTest(t *testing.T, svc *client.LocalService) {
	tests := []struct {
		name string
		data service.User
		want error
	}{
		{
			name: "login ok",
			data: service.User{
				Login:    "Major Tom",
				Password: "Ground control",
			},
		},
		{
			name: "login fail",
			data: service.User{
				Login:    "You've really",
				Password: "made the grade",
			},
			want: client.ErrInvalidCredentials,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.Auth(tt.data.Login, tt.data.Password)
			if errors.Is(err, tt.want) {
				err = nil
			}
			require.NoError(t, err)
		})
	}
}

func PutLogoPassTest(t *testing.T, svc *client.LocalService) {
	tests := []struct {
		name string
		data service.LogoPass
		want error
	}{
		{
			name: "put logoPass ok",
			data: service.LogoPass{
				SecretLogin: "Colonel",
				SecretPass:  "No one writes to",
				Description: "And nobody waits for him",
			},
		},
		{
			name: "put logoPass fail: conflict",
			data: service.LogoPass{
				SecretLogin: "Colonel",
				SecretPass:  "No one writes to",
				Description: "And nobody waits for him",
				Overwrite:   false,
			},
			want: client.ErrAlreadyExists,
		},
		{
			name: "update logoPass fail: old data",
			data: service.LogoPass{
				SecretLogin: "Colonel",
				SecretPass:  "No one writes to",
				Description: "And nobody waits for him",
				Overwrite:   true,
				Model: gorm.Model{
					UpdatedAt: time.Now().Add(time.Minute * -20),
				},
			},
			want: client.ErrAlreadyExists,
		},
		{
			name: "update logoPass ok",
			data: service.LogoPass{
				SecretLogin: "Colonel",
				SecretPass:  "No one writes to",
				Description: "And nobody waits for him",
				Overwrite:   true,
				Model: gorm.Model{
					UpdatedAt: time.Now().Add(time.Minute * 5),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.PutLogoPass(tt.data)
			if errors.Is(err, tt.want) {
				err = nil
			}
			if tt.name == "update logoPass ok" {
				err = svc.Api.UploadLogoPass(tt.data)
			}
			require.NoError(t, err)
		})
	}
}

func PutTextTest(t *testing.T, svc *client.LocalService) {
	tests := []struct {
		name string
		data service.TextData
		want error
	}{
		{
			name: "put text ok",
			data: service.TextData{
				Text:        "sublieutenant",
				Description: "young fellow",
			},
		},
		{
			name: "put text fail: conflict",
			data: service.TextData{
				Text:        "sublieutenant",
				Description: "young fellow",
				Overwrite:   false,
			},
			want: client.ErrAlreadyExists,
		},
		{
			name: "update text fail: old data",
			data: service.TextData{
				Text:        "sublieutenant",
				Description: "young fellow",
				Overwrite:   true,
				Model: gorm.Model{
					UpdatedAt: time.Now().Add(time.Minute * -10),
				},
			},
			want: client.ErrAlreadyExists,
		},
		{
			name: "update text ok",
			data: service.TextData{
				Text:        "sublieutenant",
				Description: "young fellow",
				Overwrite:   true,
				Model: gorm.Model{
					UpdatedAt: time.Now().Add(time.Minute * 5),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.PutText(tt.data)
			if errors.Is(err, tt.want) {
				err = nil
			}
			if tt.name == "update text ok" {
				err = svc.Api.UploadText(tt.data)
			}
			require.NoError(t, err)
		})
	}
}

func PutCreditCardTest(t *testing.T, svc *client.LocalService) {
	tests := []struct {
		name string
		data service.CreditCard
		want error
	}{
		{
			name: "put CreditCard ok",
			data: service.CreditCard{
				Number:      "1111",
				Holder:      "Rick Astley",
				DueDate:     "12/2025",
				CVV:         "123",
				Description: "aaa",
			},
		},
		{
			name: "put CreditCard fail: conflict",
			data: service.CreditCard{
				Number:      "1111",
				Holder:      "Rick Astley",
				DueDate:     "12/2025",
				CVV:         "123",
				Description: "aaa",
				Overwrite:   false,
			},
			want: client.ErrAlreadyExists,
		},
		{
			name: "update CreditCard fail: old data",
			data: service.CreditCard{
				Number:      "1111",
				Holder:      "Rick Astley",
				DueDate:     "12/2025",
				CVV:         "123",
				Description: "aaa",
				Overwrite:   true,
				Model: gorm.Model{
					UpdatedAt: time.Now().Add(time.Minute * -10),
				},
			},
			want: client.ErrAlreadyExists,
		},
		{
			name: "update CreditCard ok",
			data: service.CreditCard{
				Number:      "1111",
				Holder:      "Rick Astley",
				DueDate:     "12/2025",
				CVV:         "123",
				Description: "aaa",
				Overwrite:   true,
				Model: gorm.Model{
					UpdatedAt: time.Now().Add(time.Minute * 5),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.PutCreditCard(tt.data)
			if errors.Is(err, tt.want) {
				err = nil
			}
			if tt.name == "update CreditCard ok" {
				err = svc.Api.UploadCreditCard(tt.data)
			}
			require.NoError(t, err)
		})
	}
}

func PutBinaryTest(t *testing.T, svc *client.LocalService) {
	tests := []struct {
		name string
		data service.BinaryData
		want error
	}{
		{
			name: "put binary ok",
			data: service.BinaryData{
				Binary:      "mЕYђ8+012gЎZQСBБ00516МЖЪQ”dм™±cgѕџфДлИдЏ'уЮ",
				Description: "aaa",
			},
		},
		{
			name: "put binary fail: conflict",
			data: service.BinaryData{
				Binary:      "mЕYђ8+012gЎZQСBБ00516МЖЪQ”dм™±cgѕџфДлИдЏ'уЮ",
				Description: "aaa",
				Overwrite:   false,
			},
			want: client.ErrAlreadyExists,
		},
		{
			name: "update binary fail: old data",
			data: service.BinaryData{
				Binary:      "mЕYђ8+012gЎZQСBБ00516МЖЪQ”dм™±cgѕџфДлИдЏ'уЮ",
				Description: "aaa",
				Overwrite:   true,
				Model: gorm.Model{
					UpdatedAt: time.Now().Add(time.Minute * -10),
				},
			},
			want: client.ErrAlreadyExists,
		},
		{
			name: "update binary ok",
			data: service.BinaryData{
				Binary:      "mЕYђ8+012gЎZQСBБ00516МЖЪQ”dм™±cgѕџфДлИдЏ'уЮ",
				Description: "aaa",
				Overwrite:   true,
				Model: gorm.Model{
					UpdatedAt: time.Now().Add(time.Minute * 5),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.Api.UploadBinary(tt.data)
			if errors.Is(err, tt.want) {
				err = nil
			}
			if tt.name == "update binary ok" {
				err = svc.Api.UploadBinary(tt.data)
			}
			require.NoError(t, err)
		})
	}
}

func GetBinary(t *testing.T, svc *client.LocalService) {
	tests := []struct {
		name     string
		data     service.BinaryData
		want     error
		respData service.BinaryData
	}{
		{
			name: "get binary ok",
			data: service.BinaryData{
				Description: "aaa",
			},
			respData: service.BinaryData{
				Binary:      "mЕYђ8+012gЎZQСBБ00516МЖЪQ”dм™±cgѕџфДлИдЏ'уЮ",
				Description: "aaa",
			},
		},
		{
			name: "get binary fail: not found",
			data: service.BinaryData{
				Description: "bbb",
			},
			want: client.ErrEmpty,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			binary, err := svc.Api.GetBinary(tt.data)
			if errors.Is(err, tt.want) {
				err = nil
			}
			if binary.Binary != tt.respData.Binary {
				err = fmt.Errorf("wrong data")
			}
			require.NoError(t, err)
		})
	}
}

func GetLogoPassesTest(t *testing.T, svc *client.LocalService) {
	tests := []struct {
		name string
		data []service.LogoPass
		want error
	}{
		{
			name: "get logoPasses ok",
			data: []service.LogoPass{
				{
					SecretLogin: "Colonel",
					SecretPass:  "No one writes to",
					Description: "And nobody waits for him",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logoPasses, err := svc.Api.GetLogoPasses()
			if errors.Is(err, tt.want) {
				err = nil
			}

			if logoPasses[0].SecretLogin != tt.data[0].SecretLogin ||
				logoPasses[0].SecretPass != tt.data[0].SecretPass {
				err = fmt.Errorf("wrong data")
			}
			require.NoError(t, err)
		})
	}
}

func GetTextsTest(t *testing.T, svc *client.LocalService) {
	tests := []struct {
		name string
		data []service.TextData
		want error
	}{
		{
			name: "get texts ok",
			data: []service.TextData{
				{
					Text:        "sublieutenant",
					Description: "young fellow",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			texts, err := svc.Api.GetTexts()
			if errors.Is(err, tt.want) {
				err = nil
			}

			if texts[0].Text != tt.data[0].Text {
				err = fmt.Errorf("wrong data")
			}
			require.NoError(t, err)
		})
	}
}

func GetCreditCardsTest(t *testing.T, svc *client.LocalService) {
	tests := []struct {
		name string
		data []service.CreditCard
		want error
	}{
		{
			name: "get credit cards ok",
			data: []service.CreditCard{
				{
					Number:      "1111",
					Holder:      "Rick Astley",
					DueDate:     "12/2025",
					CVV:         "123",
					Description: "aaa",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			creditCards, err := svc.Api.GetCreditCards()
			if errors.Is(err, tt.want) {
				err = nil
			}

			if creditCards[0].Holder != tt.data[0].Holder ||
				creditCards[0].DueDate != tt.data[0].DueDate ||
				creditCards[0].CVV != tt.data[0].CVV {
				err = fmt.Errorf("wrong data")
			}
			require.NoError(t, err)
		})
	}
}

func GetBinaryListTest(t *testing.T, svc *client.LocalService) {
	tests := []struct {
		name string
		data []service.BinaryData
		want error
	}{
		{
			name: "get binary list ok",
			data: []service.BinaryData{
				{
					Description: "aaa",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			texts, err := svc.Api.GetBinaryList()
			if errors.Is(err, tt.want) {
				err = nil
			}
			if texts[0].Description != tt.data[0].Description {
				err = fmt.Errorf("wrong data")
			}
			require.NoError(t, err)
		})
	}
}
