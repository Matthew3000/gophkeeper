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
	userStorage   storage.UserStorage
	cookieStorage sessions.CookieStore
}

func NewApp(cfg config.Config, userStorage storage.UserStorage, cookieStorage sessions.CookieStore) *App {
	return &App{config: cfg, userStorage: userStorage, cookieStorage: cookieStorage}
}

func (app *App) Run() {
	router := mux.NewRouter()
	router.Use(tools.GzipMiddleware, app.AddContext)
	router.HandleFunc("/api/user/register", app.Register).Methods(http.MethodPost)
	router.HandleFunc("/api/user/login", app.Login).Methods(http.MethodPost)
	router.HandleFunc("/api/user/upload/logopass", app.IsAuthorized(app.UploadLogoPass)).Methods(http.MethodPost)
	router.HandleFunc("/api/user/upload/text", app.IsAuthorized(app.UploadText)).Methods(http.MethodPost)
	router.HandleFunc("/api/user/upload/creditcard", app.IsAuthorized(app.UploadCreditCard)).Methods(http.MethodPost)
	router.HandleFunc("/api/user/upload/binary", app.IsAuthorized(app.UploadBinary)).Methods(http.MethodPost)
	router.HandleFunc("/api/user/download/logopasses", app.IsAuthorized(app.BatchDownloadLogoPasses)).Methods(http.MethodGet)
	router.HandleFunc("/api/user/download/logopasses", app.IsAuthorized(app.BatchDownloadTexts)).Methods(http.MethodGet)
	router.HandleFunc("/api/user/download/logopasses", app.IsAuthorized(app.BatchDownloadCreditCards)).Methods(http.MethodGet)
	router.HandleFunc("/api/user/download/logopasses", app.IsAuthorized(app.DownloadBinaryList)).Methods(http.MethodGet)
	router.HandleFunc("/api/user/download/logopasses", app.IsAuthorized(app.DownloadBinary)).Methods(http.MethodPost)

	router.HandleFunc("/", app.handleDefault)

	log.Fatal(http.ListenAndServe(app.config.ServerAddress, router))
}
