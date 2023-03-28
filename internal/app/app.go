// Package app holds handlers and routing for Gophkeeper secrets manager
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

// App is a struct holding structures crucial for the working of the service
type App struct {
	config        config.Config
	UserStorage   storage.UserStorage
	cookieStorage sessions.CookieStore
}

// This holds all the routes available in App
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
	GetWindows             = "/download/windows"
	GetMac                 = "/download/mac"
	GetLinux               = "/download/linux"
)

// NewApp constructor for app
func NewApp(cfg config.Config, userStorage storage.UserStorage, cookieStorage sessions.CookieStore) *App {
	return &App{config: cfg, UserStorage: userStorage, cookieStorage: cookieStorage}
}

// Run creates routing and holds all the handlers.
func (app *App) Run() {
	router := mux.NewRouter()
	router.Use(tools.GzipMiddleware, app.addContext)
	router.HandleFunc(GetWindows, app.handleDownload).Methods(http.MethodGet)
	router.HandleFunc(GetLinux, app.handleDownload).Methods(http.MethodGet)
	router.HandleFunc(GetMac, app.handleDownload).Methods(http.MethodGet)
	router.HandleFunc(RegisterEndpoint, app.register).Methods(http.MethodPost)
	router.HandleFunc(LoginEndpoint, app.login).Methods(http.MethodPost)
	router.HandleFunc(PutLogoPassEndpoint, app.isAuthorized(app.uploadLogoPass)).Methods(http.MethodPost)
	router.HandleFunc(PutTextEndpoint, app.isAuthorized(app.uploadText)).Methods(http.MethodPost)
	router.HandleFunc(PutCreditCardEndpoint, app.isAuthorized(app.uploadCreditCard)).Methods(http.MethodPost)
	router.HandleFunc(PutBinaryEndpoint, app.isAuthorized(app.uploadBinary)).Methods(http.MethodPost)
	router.HandleFunc(GetLogoPassesEndpoint, app.isAuthorized(app.batchDownloadLogoPasses)).Methods(http.MethodGet)
	router.HandleFunc(GetTextsEndpoint, app.isAuthorized(app.batchDownloadTexts)).Methods(http.MethodGet)
	router.HandleFunc(GetCreditCardsEndpoint, app.isAuthorized(app.batchDownloadCreditCards)).Methods(http.MethodGet)
	router.HandleFunc(GetBinaryListEndpoint, app.isAuthorized(app.downloadBinaryList)).Methods(http.MethodGet)
	router.HandleFunc(GetBinaryEndpoint, app.isAuthorized(app.downloadBinary)).Methods(http.MethodPost)
	router.NotFoundHandler = http.HandlerFunc(app.handleDefault)

	log.Fatal(http.ListenAndServe(app.config.ServerAddress, router))
}
