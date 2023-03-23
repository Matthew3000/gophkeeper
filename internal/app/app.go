package app

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"gophkeeper/internal/config"
	"gophkeeper/internal/storage"
	"gophkeeper/internal/tools"
	"log"
	"net/http"
)

type App struct {
	config        config.Config
	UserStorage   storage.UserStorage
	cookieStorage sessions.CookieStore
}

const (
	RegisterEndpoint       = "/api/user/register"
	LoginEndpoint          = "/api/user/login"
	PutLogoPassEndpoint    = "/api/user/upload/logopass"
	PutTextEndpoint        = "/api/user/upload/text"
	PutCreditCardEndpoint  = "/api/user/upload/credit-card"
	PutBinaryEndpoint      = "/api/user/upload/binary"
	GetLogoPassesEndpoint  = "/api/user/download/logopasses"
	GetTextsEndpoint       = "/api/user/download/texts"
	GetCreditCardsEndpoint = "/api/user/download/credit-cards"
	GetBinaryListEndpoint  = "/api/user/download/binary-list"
	GetBinaryEndpoint      = "/api/user/download/binary"
)

func NewApp(cfg config.Config, userStorage storage.UserStorage, cookieStorage sessions.CookieStore) *App {
	return &App{config: cfg, UserStorage: userStorage, cookieStorage: cookieStorage}
}

func (app *App) Run() {
	router := mux.NewRouter()
	router.Use(tools.GzipMiddleware, app.AddContext)
	router.HandleFunc(RegisterEndpoint, app.Register).Methods(http.MethodPost)
	router.HandleFunc(LoginEndpoint, app.Login).Methods(http.MethodPost)
	router.HandleFunc(PutLogoPassEndpoint, app.IsAuthorized(app.UploadLogoPass)).Methods(http.MethodPost)
	router.HandleFunc(PutTextEndpoint, app.IsAuthorized(app.UploadText)).Methods(http.MethodPost)
	router.HandleFunc(PutCreditCardEndpoint, app.IsAuthorized(app.UploadCreditCard)).Methods(http.MethodPost)
	router.HandleFunc(PutBinaryEndpoint, app.IsAuthorized(app.UploadBinary)).Methods(http.MethodPost)
	router.HandleFunc(GetLogoPassesEndpoint, app.IsAuthorized(app.BatchDownloadLogoPasses)).Methods(http.MethodGet)
	router.HandleFunc(GetTextsEndpoint, app.IsAuthorized(app.BatchDownloadTexts)).Methods(http.MethodGet)
	router.HandleFunc(GetCreditCardsEndpoint, app.IsAuthorized(app.BatchDownloadCreditCards)).Methods(http.MethodGet)
	router.HandleFunc(GetBinaryListEndpoint, app.IsAuthorized(app.DownloadBinaryList)).Methods(http.MethodGet)
	router.HandleFunc(GetBinaryEndpoint, app.IsAuthorized(app.DownloadBinary)).Methods(http.MethodPost)

	router.HandleFunc("/", app.handleDefault)

	log.Fatal(http.ListenAndServe(app.config.ServerAddress, router))
}
