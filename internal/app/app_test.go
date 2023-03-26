package app

import (
	"encoding/json"
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/go-resty/resty/v2"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gophkeeper/internal/config"
	"gophkeeper/internal/service"
	"gophkeeper/internal/storage"
	"gorm.io/gorm"
	"log"
	"net/http"
	"testing"
	"time"
)

// Чтобы запустить тест, необходимо поднять Postgres базу и написать url в cfg ниже
// для запуска обоих тестов сразу использовать go test -p 1 ./.../

func TestApp(t *testing.T) {
	var cfg config.Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}
	flag.StringVar(&cfg.DatabaseDSN, "d", cfg.DatabaseDSN, "File Storage Path")
	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "Server address")
	flag.Parse()

	userStorage := storage.NewUserStorage(cfg.DatabaseDSN)
	cookieStorage := sessions.NewCookieStore([]byte(service.SecretKey))
	var app = NewApp(cfg, userStorage, *cookieStorage)
	go app.Run()

	cookie := AuthTest(t, app)
	PutLogoPassTest(t, app, cookie)
	PutTextTest(t, app, cookie)
	PutCreditCardTest(t, app, cookie)
	PutBinaryTest(t, app, cookie)
	GetLogoPassesTest(t, app, cookie)
	GetTextsTest(t, app, cookie)
	GetCreditCardsTest(t, app, cookie)
	GetBinaryListTest(t, app, cookie)
	GetBinaryTest(t, app, cookie)

	app.UserStorage.DeleteAll()
}

func AuthTest(t *testing.T, app *App) http.Cookie {
	var cookie http.Cookie
	type want struct {
		statusCode  int
		contentType string
	}
	tests := []struct {
		name string
		addr string
		user service.User
		want want
	}{
		{
			name: "404",
			addr: "/iamlost",
			user: service.User{},
			want: want{
				statusCode:  http.StatusNotFound,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "register ok",
			addr: RegisterEndpoint,
			user: service.User{
				Login:    "nevergonna",
				Password: "giveyouup",
			},
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			name: "register conflict",
			addr: RegisterEndpoint,
			user: service.User{
				Login:    "nevergonna",
				Password: "giveyouup",
			},
			want: want{
				statusCode:  http.StatusConflict,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "login fail: no such user",
			addr: LoginEndpoint,
			user: service.User{
				Login:    "imgona",
				Password: "giveyouup",
			},
			want: want{
				statusCode:  http.StatusUnauthorized,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "login fail: wrong pass",
			addr: LoginEndpoint,
			user: service.User{
				Login:    "nevergonna",
				Password: "letyoudown",
			},
			want: want{
				statusCode:  http.StatusUnauthorized,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "login OK",
			addr: LoginEndpoint,
			user: service.User{
				Login:    "nevergonna",
				Password: "giveyouup",
			},
			want: want{
				statusCode: http.StatusOK,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.user)
			request := resty.New().R().SetHeader("Content-Type", "application/json").SetBody(body)

			result, err := request.Post("http://" + app.config.ServerAddress + tt.addr)
			require.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, result.StatusCode())
			assert.Equal(t, tt.want.contentType, result.Header().Get("Content-Type"))
			if tt.addr == LoginEndpoint {
				if result.StatusCode() == http.StatusOK {
					s := result.Cookies()
					cookie = *s[0]
				}
			}
		})
	}
	return cookie
}

func PutLogoPassTest(t *testing.T, app *App, cookie http.Cookie) {

	type want struct {
		statusCode  int
		contentType string
	}
	tests := []struct {
		name string
		addr string
		data service.LogoPass
		want want
	}{
		{
			name: "logoPass upload ok",
			addr: PutLogoPassEndpoint,
			data: service.LogoPass{
				SecretLogin: "aaa",
				SecretPass:  "aaa",
				Description: "aaa",
			},
			want: want{
				statusCode: http.StatusCreated,
			},
		},
		{
			name: "logoPass upload conflict: already exists",
			addr: PutLogoPassEndpoint,
			data: service.LogoPass{
				SecretLogin: "aaa",
				SecretPass:  "aaa",
				Description: "aaa",
			},
			want: want{
				statusCode:  http.StatusConflict,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "logoPass update ok",
			addr: PutLogoPassEndpoint,
			data: service.LogoPass{
				SecretLogin: "aaa",
				SecretPass:  "aaa",
				Description: "aaa",
				Overwrite:   true,
				Model: gorm.Model{
					UpdatedAt: time.Now().Add(time.Minute * 5),
				},
			},
			want: want{
				statusCode: http.StatusCreated,
			},
		},
		{
			name: "logoPass update conflict: old data",
			addr: PutLogoPassEndpoint,
			data: service.LogoPass{
				SecretLogin: "aaa",
				SecretPass:  "aaa",
				Description: "aaa",
				Overwrite:   true,
				Model: gorm.Model{
					UpdatedAt: time.Now().Add(time.Minute * -10),
				},
			},
			want: want{
				statusCode:  http.StatusConflict,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := resty.New().R().SetHeader("Content-Type", "application/json").
				SetBody(tt.data).SetCookie(&cookie)

			result, err := request.Post("http://" + app.config.ServerAddress + tt.addr)
			require.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, result.StatusCode())
			assert.Equal(t, tt.want.contentType, result.Header().Get("Content-Type"))
		})
	}
}

func PutTextTest(t *testing.T, app *App, cookie http.Cookie) {

	type want struct {
		statusCode  int
		contentType string
	}
	tests := []struct {
		name string
		addr string
		data service.TextData
		want want
	}{
		{
			name: "text upload ok",
			addr: PutTextEndpoint,
			data: service.TextData{
				Text:        "never gonna run around, desert you",
				Description: "aaa",
			},
			want: want{
				statusCode: http.StatusCreated,
			},
		},
		{
			name: "text upload conflict: already exists",
			addr: PutTextEndpoint,
			data: service.TextData{
				Text:        "never gonna run around, desert you",
				Description: "aaa",
			},
			want: want{
				statusCode:  http.StatusConflict,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "text update ok",
			addr: PutTextEndpoint,
			data: service.TextData{
				Text:        "never gonna run around, desert you",
				Description: "aaa",
				Overwrite:   true,
				Model: gorm.Model{
					UpdatedAt: time.Now().Add(time.Minute * 5),
				},
			},
			want: want{
				statusCode: http.StatusCreated,
			},
		},
		{
			name: "text update conflict: old data",
			addr: PutTextEndpoint,
			data: service.TextData{
				Text:        "never gonna run around, desert you",
				Description: "aaa",
				Overwrite:   true,
				Model: gorm.Model{
					UpdatedAt: time.Now().Add(time.Minute * -10),
				},
			},
			want: want{
				statusCode:  http.StatusConflict,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := resty.New().R().SetHeader("Content-Type", "application/json").
				SetBody(tt.data).SetCookie(&cookie)

			result, err := request.Post("http://" + app.config.ServerAddress + tt.addr)
			require.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, result.StatusCode())
			assert.Equal(t, tt.want.contentType, result.Header().Get("Content-Type"))
		})
	}
}

func PutCreditCardTest(t *testing.T, app *App, cookie http.Cookie) {

	type want struct {
		statusCode  int
		contentType string
	}
	tests := []struct {
		name string
		addr string
		data service.CreditCard
		want want
	}{
		{
			name: "Credit card upload ok",
			addr: PutCreditCardEndpoint,
			data: service.CreditCard{
				Number:      "1111",
				Holder:      "Rick Astley",
				DueDate:     "12/2025",
				CVV:         "123",
				Description: "aaa",
			},
			want: want{
				statusCode: http.StatusCreated,
			},
		},
		{
			name: "Credit card upload conflict: already exists",
			addr: PutCreditCardEndpoint,
			data: service.CreditCard{
				Number:      "1111",
				Holder:      "Rick Astley",
				DueDate:     "12/2025",
				CVV:         "123",
				Description: "aaa",
			},
			want: want{
				statusCode:  http.StatusConflict,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Credit card update ok",
			addr: PutCreditCardEndpoint,
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
			want: want{
				statusCode: http.StatusCreated,
			},
		},
		{
			name: "Credit card update conflict: old data",
			addr: PutCreditCardEndpoint,
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
			want: want{
				statusCode:  http.StatusConflict,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := resty.New().R().SetHeader("Content-Type", "application/json").
				SetBody(tt.data).SetCookie(&cookie)

			result, err := request.Post("http://" + app.config.ServerAddress + tt.addr)
			require.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, result.StatusCode())
			assert.Equal(t, tt.want.contentType, result.Header().Get("Content-Type"))
		})
	}
}

func PutBinaryTest(t *testing.T, app *App, cookie http.Cookie) {

	type want struct {
		statusCode  int
		contentType string
	}
	tests := []struct {
		name string
		addr string
		data service.BinaryData
		want want
	}{
		{
			name: "binary upload ok",
			addr: PutBinaryEndpoint,
			data: service.BinaryData{
				Binary:      "mЕYђ8+012gЎZQСBБ00516МЖЪQ”dм™±cgѕџфДлИдЏ'уЮ",
				Description: "aaa",
			},
			want: want{
				statusCode: http.StatusCreated,
			},
		},
		{
			name: "binary upload conflict: already exists",
			addr: PutBinaryEndpoint,
			data: service.BinaryData{
				Login:       "nevergonna",
				Binary:      "mЕYђ8+012gЎZQСBБ00516МЖЪQ”dм™±cgѕџфДлИдЏ'уЮ",
				Description: "aaa",
			},
			want: want{
				statusCode:  http.StatusConflict,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "binary update ok",
			addr: PutBinaryEndpoint,
			data: service.BinaryData{
				Login:       "nevergonna",
				Binary:      "mЕYђ8+012gЎZQСBБ00516МЖЪQ”dм™±cgѕџфДлИдЏ'уЮ",
				Description: "aaa",
				Overwrite:   true,
				Model: gorm.Model{
					UpdatedAt: time.Now().Add(time.Minute * 5),
				},
			},
			want: want{
				statusCode: http.StatusCreated,
			},
		},
		{
			name: "binary update conflict: old data",
			addr: PutBinaryEndpoint,
			data: service.BinaryData{
				Login:       "nevergonna",
				Binary:      "mЕYђ8+012gЎZQСBБ00516МЖЪQ”dм™±cgѕџфДлИдЏ'уЮ",
				Description: "aaa",
				Overwrite:   true,
				Model: gorm.Model{
					UpdatedAt: time.Now().Add(time.Minute * -10),
				},
			},
			want: want{
				statusCode:  http.StatusConflict,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := resty.New().R().SetHeader("Content-Type", "application/json").
				SetBody(tt.data).SetCookie(&cookie)

			result, err := request.Post("http://" + app.config.ServerAddress + tt.addr)
			require.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, result.StatusCode())
			assert.Equal(t, tt.want.contentType, result.Header().Get("Content-Type"))
		})
	}
}

func GetLogoPassesTest(t *testing.T, app *App, cookie http.Cookie) {
	type want struct {
		statusCode  int
		contentType string
	}
	tests := []struct {
		name string
		addr string
		data []service.LogoPass
		want want
	}{
		{
			name: "logopass list download ok",
			addr: GetLogoPassesEndpoint,
			data: []service.LogoPass{
				{
					SecretLogin: "aaa",
					SecretPass:  "aaa",
					Description: "aaa",
				},
			},
			want: want{
				statusCode:  http.StatusOK,
				contentType: "application/json; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := resty.New().R().SetHeader("Content-Type", "application/json").SetCookie(&cookie)

			result, err := request.Get("http://" + app.config.ServerAddress + tt.addr)
			require.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, result.StatusCode())
			assert.Equal(t, tt.want.contentType, result.Header().Get("Content-Type"))
		})
	}
}

func GetTextsTest(t *testing.T, app *App, cookie http.Cookie) {
	type want struct {
		statusCode  int
		contentType string
	}
	tests := []struct {
		name string
		addr string
		data []service.TextData
		want want
	}{
		{
			name: "text list download ok",
			addr: GetTextsEndpoint,
			data: []service.TextData{
				{
					Text:        "never gonna run around, desert you",
					Description: "aaa",
				},
			},
			want: want{
				statusCode:  http.StatusOK,
				contentType: "application/json; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := resty.New().R().SetHeader("Content-Type", "application/json").SetCookie(&cookie)

			result, err := request.Get("http://" + app.config.ServerAddress + tt.addr)
			require.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, result.StatusCode())
			assert.Equal(t, tt.want.contentType, result.Header().Get("Content-Type"))
		})
	}
}

func GetCreditCardsTest(t *testing.T, app *App, cookie http.Cookie) {
	type want struct {
		statusCode  int
		contentType string
	}
	tests := []struct {
		name string
		addr string
		data []service.CreditCard
		want want
	}{
		{
			name: "credit cards download ok",
			addr: GetCreditCardsEndpoint,
			data: []service.CreditCard{
				{
					Number:      "1111",
					Holder:      "Rick Astley",
					DueDate:     "12/2025",
					CVV:         "123",
					Description: "aaa",
				},
			},
			want: want{
				statusCode:  http.StatusOK,
				contentType: "application/json; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := resty.New().R().SetHeader("Content-Type", "application/json").SetCookie(&cookie)

			result, err := request.Get("http://" + app.config.ServerAddress + tt.addr)
			require.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, result.StatusCode())
			assert.Equal(t, tt.want.contentType, result.Header().Get("Content-Type"))
		})
	}
}

func GetBinaryListTest(t *testing.T, app *App, cookie http.Cookie) {
	type want struct {
		statusCode  int
		contentType string
	}
	tests := []struct {
		name string
		addr string
		data []service.BinaryData
		want want
	}{
		{
			name: "binary list download ok",
			addr: GetBinaryListEndpoint,
			data: []service.BinaryData{
				{
					Description: "aaa",
				},
			},
			want: want{
				statusCode:  http.StatusOK,
				contentType: "application/json; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := resty.New().R().SetHeader("Content-Type", "application/json").SetCookie(&cookie)

			result, err := request.Get("http://" + app.config.ServerAddress + tt.addr)
			require.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, result.StatusCode())
			assert.Equal(t, tt.want.contentType, result.Header().Get("Content-Type"))
		})
	}
}

func GetBinaryTest(t *testing.T, app *App, cookie http.Cookie) {
	type want struct {
		statusCode  int
		contentType string
	}
	tests := []struct {
		name    string
		addr    string
		reqData service.BinaryData
		data    []service.BinaryData
		want    want
	}{
		{
			name: "binary download ok",
			addr: GetBinaryEndpoint,
			reqData: service.BinaryData{
				Description: "aaa",
			},
			data: []service.BinaryData{
				{
					Binary:      "mЕYђ8+012gЎZQСBБ00516МЖЪQ”dм™±cgѕџфДлИдЏ'уЮ",
					Description: "aaa",
				},
			},
			want: want{
				statusCode:  http.StatusOK,
				contentType: "application/json; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.reqData)
			request := resty.New().R().SetHeader("Content-Type", "application/json").
				SetBody(body).SetCookie(&cookie)

			result, err := request.Post("http://" + app.config.ServerAddress + tt.addr)

			require.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, result.StatusCode())
			assert.Equal(t, tt.want.contentType, result.Header().Get("Content-Type"))
		})
	}
}
