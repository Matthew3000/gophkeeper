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
	router.HandleFunc("/api/user/register", app.handleRegister).Methods(http.MethodPost)
	router.HandleFunc("/api/user/login", app.handleLogin).Methods(http.MethodPost)
	router.HandleFunc("/api/user/orders", app.IsAuthorized(app.handleUploadLogoPass)).Methods(http.MethodPost)

	router.HandleFunc("/", app.handleDefault)

	log.Fatal(http.ListenAndServe(app.config.ServerAddress, router))
}
