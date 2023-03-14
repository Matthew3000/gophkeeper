package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/render"
	"gophkeeper/internal/service"
	"gophkeeper/internal/storage"
	"log"
	"net/http"
	"time"
)

func (app *App) IsAuthorized(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := app.cookieStorage.Get(r, "session.id")
		authenticated := session.Values["authenticated"]
		if authenticated != nil && authenticated != false {
			handler.ServeHTTP(w, r)
			return
		}

		http.Redirect(w, r, "/login", http.StatusUnauthorized)
	}
}

func (app *App) AddContext(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		r = r.WithContext(ctx)
		handler.ServeHTTP(w, r)
	})
}

func (app *App) Register(w http.ResponseWriter, r *http.Request) {
	var user service.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("register err: json parse error: %s", err)
		http.Error(w, fmt.Sprintf("json parse error: %s", err), http.StatusBadRequest)
		return
	}

	err = app.userStorage.RegisterUser(user, r.Context())
	if err != nil {
		log.Printf("register err: %s for user: %s, password: %s", err, user.Login, user.Password)
		if errors.Is(err, storage.ErrUserExists) {
			http.Error(w, fmt.Sprintf("register error: %s", err), http.StatusConflict)
			return
		}
		http.Error(w, fmt.Sprintf("register error: %s", err), http.StatusInternalServerError)
		return
	}

	var authDetails service.Authentication
	authDetails.Login = user.Login
	authDetails.Password = user.Password
	err = app.userStorage.CheckUserAuth(authDetails, r.Context())
	if err != nil {
		log.Printf("register then auth err: %s for user: %s, password: %s", err, authDetails.Login, authDetails.Password)
		if errors.Is(err, storage.ErrInvalidCredentials) {
			http.Error(w, fmt.Sprintf("auth error: %s", err), http.StatusUnauthorized)
			return
		}
		http.Error(w, fmt.Sprintf("auth error: %s", err), http.StatusInternalServerError)
		return
	}

	session, _ := app.cookieStorage.Get(r, "session.id")
	session.Values["authenticated"] = true
	session.Values["login"] = user.Login
	session.Save(r, w)
	w.WriteHeader(http.StatusOK)
}

func (app *App) Login(w http.ResponseWriter, r *http.Request) {
	var authDetails service.Authentication
	err := json.NewDecoder(r.Body).Decode(&authDetails)
	if err != nil {
		log.Printf("auth err: json parse error: %s", err)
		http.Error(w, fmt.Sprintf("json parse error: %s", err), http.StatusBadRequest)
		return
	}

	err = app.userStorage.CheckUserAuth(authDetails, r.Context())
	if err != nil {
		log.Printf("auth err: %s for user: %s, password: %s", err, authDetails.Login, authDetails.Password)
		if errors.Is(err, storage.ErrInvalidCredentials) {
			http.Error(w, fmt.Sprintf("auth error: %s", err), http.StatusUnauthorized)
			return
		}
		http.Error(w, fmt.Sprintf("auth error: %s", err), http.StatusInternalServerError)
		return
	}

	session, _ := app.cookieStorage.Get(r, "session.id")
	session.Values["authenticated"] = true
	session.Values["login"] = authDetails.Login
	session.Save(r, w)
	w.WriteHeader(http.StatusOK)
}

func (app *App) UploadLogoPass(w http.ResponseWriter, r *http.Request) {
	var logoPass service.LogoPass

	session, _ := app.cookieStorage.Get(r, "session.id")
	logoPass.Login = session.Values["login"].(string)

	err := json.NewDecoder(r.Body).Decode(&logoPass)
	if err != nil {
		log.Printf("upload logopass pair: json parse error: %s", err)
		http.Error(w, "json parse error:", http.StatusBadRequest)
		return
	}

	err = app.userStorage.PutLogoPass(logoPass, r.Context())
	if err != nil {
		log.Printf("put logopass pair: save to db: %s for user: %s", err, logoPass.Login)
		return
	}
	http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
}

func (app *App) BatchDownloadLogoPasses(w http.ResponseWriter, r *http.Request) {
	var logoPass service.LogoPass

	session, _ := app.cookieStorage.Get(r, "session.id")
	logoPass.Login = session.Values["login"].(string)

	listLogoPasses, err := app.userStorage.BatchGetLogoPasses(logoPass.Login, r.Context())
	if err != nil {
		log.Printf("get logpass pairs: %s for user: %s", err, logoPass.Login)
		if errors.Is(err, storage.ErrEmpty) {
			w.WriteHeader(http.StatusNoContent)
		} else {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}
	}
	render.JSON(w, r, listLogoPasses)
}

func (app *App) UploadSecret(w http.ResponseWriter, r *http.Request) {
	var text service.TextData

	session, _ := app.cookieStorage.Get(r, "session.id")
	text.Login = session.Values["login"].(string)

	err := json.NewDecoder(r.Body).Decode(&text)
	if err != nil {
		log.Printf("upload secret text: json parse error: %s", err)
		http.Error(w, "json parse error:", http.StatusBadRequest)
		return
	}

	err = app.userStorage.PutText(text, r.Context())
	if err != nil {
		log.Printf("put secret text: save to db: %s for user: %s", err, text.Login)
		return
	}
	http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
}

func (app *App) BatchDownloadSecrets(w http.ResponseWriter, r *http.Request) {
	var text service.TextData

	session, _ := app.cookieStorage.Get(r, "session.id")
	text.Login = session.Values["login"].(string)

	listTexts, err := app.userStorage.BatchGetTexts(text.Login, r.Context())
	if err != nil {
		log.Printf("get secrets: %s for user: %s", err, text.Login)
		if errors.Is(err, storage.ErrEmpty) {
			w.WriteHeader(http.StatusNoContent)
		} else {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}
	}
	render.JSON(w, r, listTexts)
}

func (app *App) UploadCreditCard(w http.ResponseWriter, r *http.Request) {
	var card service.CreditCard

	session, _ := app.cookieStorage.Get(r, "session.id")
	card.Login = session.Values["login"].(string)

	err := json.NewDecoder(r.Body).Decode(&card)
	if err != nil {
		log.Printf("upload credit card: json parse error: %s", err)
		http.Error(w, "json parse error:", http.StatusBadRequest)
		return
	}

	err = app.userStorage.PutCreditCard(card, r.Context())
	if err != nil {
		log.Printf("put credit card: save to db: %s for user: %s", err, card.Login)
		return
	}
	http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
}

func (app *App) BatchDownloadCreditCards(w http.ResponseWriter, r *http.Request) {
	var card service.CreditCard

	session, _ := app.cookieStorage.Get(r, "session.id")
	card.Login = session.Values["login"].(string)

	listCards, err := app.userStorage.BatchGetCreditCards(card.Login, r.Context())
	if err != nil {
		log.Printf("get logpass pairs: %s for user: %s", err, card.Login)
		if errors.Is(err, storage.ErrEmpty) {
			w.WriteHeader(http.StatusNoContent)
		} else {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}
	}
	render.JSON(w, r, listCards)
}

func (app *App) UploadBinary(w http.ResponseWriter, r *http.Request) {
	var binary service.BinaryData

	session, _ := app.cookieStorage.Get(r, "session.id")
	binary.Login = session.Values["login"].(string)

	err := json.NewDecoder(r.Body).Decode(&binary)
	if err != nil {
		log.Printf("upload credit card: json parse error: %s", err)
		http.Error(w, "json parse error:", http.StatusBadRequest)
		return
	}

	err = app.userStorage.PutBinary(binary, r.Context())
	if err != nil {
		log.Printf("put credit card: save to db: %s for user: %s", err, binary.Login)
		return
	}
	http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
}

func (app *App) DownloadBinaryList(w http.ResponseWriter, r *http.Request) {
	var binaryList service.UserBinaryList

	session, _ := app.cookieStorage.Get(r, "session.id")
	binaryList.Login = session.Values["login"].(string)

	var err error
	binaryList, err = app.userStorage.GetBinaryList(binaryList.Login, r.Context())
	if err != nil {
		log.Printf("get binary list: %s for user: %s", err, binaryList.Login)
		if errors.Is(err, storage.ErrEmpty) {
			w.WriteHeader(http.StatusNoContent)
		} else {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}
	}
	render.JSON(w, r, binaryList)
}

func (app *App) DownloadBinary(w http.ResponseWriter, r *http.Request) {
	var binary service.CreditCard

	session, _ := app.cookieStorage.Get(r, "session.id")
	binary.Login = session.Values["login"].(string)

	err := json.NewDecoder(r.Body).Decode(&binary)
	if err != nil {
		log.Printf("get binary: json parse error: %s", err)
		http.Error(w, "json parse error:", http.StatusBadRequest)
		return
	}

	binary, err = app.userStorage.GetBinary(binary.Login, r.Context())
	if err != nil {
		log.Printf("get binary: %s for user: %s", err, binary.Login)
		if errors.Is(err, storage.ErrEmpty) {
			w.WriteHeader(http.StatusNoContent)
		} else {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}
	}
	render.JSON(w, r, binary)
}

func (app *App) handleDefault(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://www.youtube.com/watch?v=dQw4w9WgXcQ", http.StatusTemporaryRedirect)
}
