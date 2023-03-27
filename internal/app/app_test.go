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

	app.UserStorage.DeleteAll()

	cookies := AuthTest(t, app)
	PutLogoPassTest(t, app, cookies)
	PutTextTest(t, app, cookies)
	PutCreditCardTest(t, app, cookies)
	PutBinaryTest(t, app, cookies)
	GetLogoPassesTest(t, app, cookies)
	GetTextsTest(t, app, cookies)
	GetCreditCardsTest(t, app, cookies)
	GetBinaryListTest(t, app, cookies)
	GetBinaryTest(t, app, cookies)

	//app.UserStorage.DeleteAll()
}

func AuthTest(t *testing.T, app *App) []http.Cookie {
	var cookies []http.Cookie
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
			name: "register ok. empty user",
			addr: RegisterEndpoint,
			user: service.User{
				Login:    "you know the rules",
				Password: "and so do I",
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
		{
			name: "login OK: empty user",
			addr: LoginEndpoint,
			user: service.User{
				Login:    "you know the rules",
				Password: "and so do I",
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
					cookies = append(cookies, *s[0])
				}
			}
		})
	}
	return cookies
}

func PutLogoPassTest(t *testing.T, app *App, cookies []http.Cookie) {

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
				SetBody(tt.data).SetCookie(&cookies[0])

			result, err := request.Post("http://" + app.config.ServerAddress + tt.addr)
			require.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, result.StatusCode())
			assert.Equal(t, tt.want.contentType, result.Header().Get("Content-Type"))
		})
	}
}

func PutTextTest(t *testing.T, app *App, cookies []http.Cookie) {

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
				SetBody(tt.data).SetCookie(&cookies[0])

			result, err := request.Post("http://" + app.config.ServerAddress + tt.addr)
			require.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, result.StatusCode())
			assert.Equal(t, tt.want.contentType, result.Header().Get("Content-Type"))
		})
	}
}

func PutCreditCardTest(t *testing.T, app *App, cookies []http.Cookie) {

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
				SetBody(tt.data).SetCookie(&cookies[0])

			result, err := request.Post("http://" + app.config.ServerAddress + tt.addr)
			require.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, result.StatusCode())
			assert.Equal(t, tt.want.contentType, result.Header().Get("Content-Type"))
		})
	}
}

func PutBinaryTest(t *testing.T, app *App, cookies []http.Cookie) {

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
				SetBody(tt.data).SetCookie(&cookies[0])

			result, err := request.Post("http://" + app.config.ServerAddress + tt.addr)
			require.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, result.StatusCode())
			assert.Equal(t, tt.want.contentType, result.Header().Get("Content-Type"))
		})
	}
}

func GetLogoPassesTest(t *testing.T, app *App, cookies []http.Cookie) {
	type want struct {
		statusCode  int
		contentType string
		bodyLen     int
	}
	tests := []struct {
		name      string
		addr      string
		cookieNum int
		data      []service.LogoPass
		want      want
	}{
		{
			name:      "logopass list download ok",
			addr:      GetLogoPassesEndpoint,
			cookieNum: 0,
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
		{
			name:      "logopass list download fail: no content",
			addr:      GetLogoPassesEndpoint,
			cookieNum: 1,
			want: want{
				statusCode: http.StatusNotFound,
				bodyLen:    0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := resty.New().R().SetHeader("Content-Type", "application/json").SetCookie(&cookies[tt.cookieNum])

			result, err := request.Get("http://" + app.config.ServerAddress + tt.addr)
			require.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, result.StatusCode())
			assert.Equal(t, tt.want.contentType, result.Header().Get("Content-Type"))

			if result.StatusCode() == http.StatusOK {
				var resp []service.LogoPass
				err = json.Unmarshal(result.Body(), &resp)
				require.NoError(t, err)
				tt.data[0].Model = resp[0].Model
				assert.Equal(t, tt.data, resp)
			} else {
				assert.Equal(t, tt.want.bodyLen, len(result.Body()))
			}
		})
	}
}

func GetTextsTest(t *testing.T, app *App, cookies []http.Cookie) {
	type want struct {
		statusCode  int
		contentType string
		bodyLen     int
	}
	tests := []struct {
		name      string
		addr      string
		cookieNum int
		data      []service.TextData
		want      want
	}{
		{
			name:      "text list download ok",
			addr:      GetTextsEndpoint,
			cookieNum: 0,
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
		{
			name:      "text list download fail: no content",
			addr:      GetTextsEndpoint,
			cookieNum: 1,
			want: want{
				statusCode: http.StatusNotFound,
				bodyLen:    0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := resty.New().R().SetHeader("Content-Type", "application/json").SetCookie(&cookies[tt.cookieNum])

			result, err := request.Get("http://" + app.config.ServerAddress + tt.addr)
			require.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, result.StatusCode())
			assert.Equal(t, tt.want.contentType, result.Header().Get("Content-Type"))

			if result.StatusCode() == http.StatusOK {
				var resp []service.TextData
				err = json.Unmarshal(result.Body(), &resp)
				require.NoError(t, err)
				tt.data[0].Model = resp[0].Model
				assert.Equal(t, tt.data, resp)
			} else {
				assert.Equal(t, tt.want.bodyLen, len(result.Body()))
			}
		})
	}
}

func GetCreditCardsTest(t *testing.T, app *App, cookies []http.Cookie) {
	type want struct {
		statusCode  int
		contentType string
		bodyLen     int
	}
	tests := []struct {
		name      string
		addr      string
		cookieNum int
		data      []service.CreditCard
		want      want
	}{
		{
			name:      "credit cards download ok",
			addr:      GetCreditCardsEndpoint,
			cookieNum: 0,
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
		{
			name:      "credit cards download fail: no content",
			addr:      GetCreditCardsEndpoint,
			cookieNum: 1,
			want: want{
				statusCode: http.StatusNotFound,
				bodyLen:    0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := resty.New().R().SetHeader("Content-Type", "application/json").
				SetCookie(&cookies[tt.cookieNum])

			result, err := request.Get("http://" + app.config.ServerAddress + tt.addr)
			require.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, result.StatusCode())
			assert.Equal(t, tt.want.contentType, result.Header().Get("Content-Type"))

			if result.StatusCode() == http.StatusOK {
				var resp []service.CreditCard
				err = json.Unmarshal(result.Body(), &resp)
				require.NoError(t, err)
				tt.data[0].Model = resp[0].Model
				assert.Equal(t, tt.data, resp)
			} else {
				assert.Equal(t, tt.want.bodyLen, len(result.Body()))
			}
		})
	}
}

func GetBinaryListTest(t *testing.T, app *App, cookies []http.Cookie) {
	type want struct {
		statusCode  int
		contentType string
		bodyLen     int
	}
	tests := []struct {
		name      string
		addr      string
		cookieNum int
		data      []service.BinaryData
		want      want
	}{
		{
			name:      "binary list download ok",
			addr:      GetBinaryListEndpoint,
			cookieNum: 0,
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
		{
			name:      "binary list download fail: no content",
			addr:      GetBinaryListEndpoint,
			cookieNum: 1,
			want: want{
				statusCode: http.StatusNotFound,
				bodyLen:    0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := resty.New().R().SetHeader("Content-Type", "application/json").SetCookie(&cookies[tt.cookieNum])

			result, err := request.Get("http://" + app.config.ServerAddress + tt.addr)
			require.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, result.StatusCode())
			assert.Equal(t, tt.want.contentType, result.Header().Get("Content-Type"))

			if result.StatusCode() == http.StatusOK {
				var resp []service.BinaryData
				err = json.Unmarshal(result.Body(), &resp)
				require.NoError(t, err)
				tt.data[0].Model = resp[0].Model
				assert.Equal(t, tt.data, resp)
			} else {
				assert.Equal(t, tt.want.bodyLen, len(result.Body()))
			}
		})
	}
}

func GetBinaryTest(t *testing.T, app *App, cookies []http.Cookie) {
	type want struct {
		statusCode  int
		contentType string
		bodyLen     int
	}
	tests := []struct {
		name      string
		addr      string
		cookieNum int
		reqData   service.BinaryData
		data      service.BinaryData
		want      want
	}{
		{
			name:      "binary download ok",
			addr:      GetBinaryEndpoint,
			cookieNum: 0,
			reqData: service.BinaryData{
				Description: "aaa",
			},
			data: service.BinaryData{
				Binary:      "mЕYђ8+012gЎZQСBБ00516МЖЪQ”dм™±cgѕџфДлИдЏ'уЮ",
				Description: "aaa",
			},
			want: want{
				statusCode:  http.StatusOK,
				contentType: "application/json; charset=utf-8",
			},
		},
		{
			name:      "binary download fail: no content",
			addr:      GetBinaryEndpoint,
			cookieNum: 1,
			want: want{
				statusCode: http.StatusNotFound,
				bodyLen:    0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.reqData)
			request := resty.New().R().SetHeader("Content-Type", "application/json").
				SetBody(body).SetCookie(&cookies[tt.cookieNum])

			result, err := request.Post("http://" + app.config.ServerAddress + tt.addr)

			require.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, result.StatusCode())
			assert.Equal(t, tt.want.contentType, result.Header().Get("Content-Type"))

			if result.StatusCode() == http.StatusOK {
				var resp service.BinaryData
				err = json.Unmarshal(result.Body(), &resp)
				require.NoError(t, err)
				tt.data.Model = resp.Model
				assert.Equal(t, tt.data, resp)
			} else {
				assert.Equal(t, tt.want.bodyLen, len(result.Body()))
			}
		})
	}
}
