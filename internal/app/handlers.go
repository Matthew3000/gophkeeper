package app

// Here are all the handler functions of the App

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/render"
	"gophkeeper/internal/config"
	"gophkeeper/internal/service"
	"gophkeeper/internal/storage"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"time"
)

// isAuthorized is a middleware used to find out if the user authorized.
//
// Returns:
//   - `401` upon unsuccessful cookie check.
func (app *App) isAuthorized(handler http.HandlerFunc) http.HandlerFunc {
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

// addContext is a middleware that adds context.Context to all the incoming requests.
func (app *App) addContext(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		r = r.WithContext(ctx)
		handler.ServeHTTP(w, r)
	})
}

// register handles registration.
//
// Accepts json.Marshalled service.User struct with 'login' and 'password' fields obligatory.
// Also authorizes the user right away if the registration is successful.
//
// Returns:
//   - `400` if json is corrupted
//   - `409` if 'login' already exists in storage
//   - `500` if storage methods fail to comprehend the request
//   - `200` and the access cookie in the header - if everything works out
func (app *App) register(w http.ResponseWriter, r *http.Request) {
	var user service.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("register err: json parse error: %s", err)
		http.Error(w, fmt.Sprintf("json parse error: %s", err), http.StatusBadRequest)
		return
	}

	err = app.UserStorage.RegisterUser(user, r.Context())
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
	err = app.UserStorage.CheckUserAuth(authDetails, r.Context())
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

// login handles signing in.
//
// Accepts json.Marshalled service.Authentication struct with 'login' and 'password' fields obligatory.
//
// Returns:
//   - `400` if json is corrupted
//   - `401` if user does not exist
//   - `500` if storage methods fail to comprehend the request
//   - `200` and the access cookie in the header - if everything works out
func (app *App) login(w http.ResponseWriter, r *http.Request) {
	var authDetails service.Authentication
	err := json.NewDecoder(r.Body).Decode(&authDetails)
	if err != nil {
		log.Printf("auth err: json parse error: %s", err)
		http.Error(w, fmt.Sprintf("json parse error: %s", err), http.StatusBadRequest)
		return
	}

	err = app.UserStorage.CheckUserAuth(authDetails, r.Context())
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

// uploadLogoPass handles uploading logo-pass pairs via http.Post request.
//
// Accepts json.Marshalled service.LogoPass struct with 'description' field obligatory.
// 'secret_login', 'secret' are optional but are recommended in the sake of common sense.
// 'overwrite' field is obligatory if the request is meant to update existing information.
//
// Returns:
//   - `400` if json is corrupted
//   - `409` if this data already exist and 'overwrite' flag is false or omitted
//   - `500` if storage methods fail to comprehend the request
//   - `201` if everything is OK
func (app *App) uploadLogoPass(w http.ResponseWriter, r *http.Request) {
	var logoPass service.LogoPass

	session, _ := app.cookieStorage.Get(r, "session.id")
	logoPass.Login = session.Values["login"].(string)

	err := json.NewDecoder(r.Body).Decode(&logoPass)
	if err != nil {
		log.Printf("upload logopass pair: json parse error: %s", err)
		http.Error(w, "json parse error:", http.StatusBadRequest)
		return
	}

	err = app.UserStorage.PutLogoPass(logoPass, r.Context())
	if err != nil {
		if errors.Is(err, storage.ErrOldData) || errors.Is(err, storage.ErrAlreadyExists) {
			http.Error(w, fmt.Sprint(err), http.StatusConflict)
			return
		}
		log.Printf("put logopass pair: save to db: %s for user: %s", err, logoPass.Login)
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return

	}
	w.WriteHeader(http.StatusCreated)
}

// batchDownloadLogoPasses handles sending the list of all user's logo-pass pairs via http.Get request.
//
// All it needs to run is a user login from http.Cookie.
//
// Returns:
//   - `404` if user has no such data in storage
//   - `500` if storage methods fail to comprehend the request
//   - json.Marshalled struct of []service.LogoPass type that contains all fields got from storage
//     except 'overwrite' field by virtue of its needlessness
func (app *App) batchDownloadLogoPasses(w http.ResponseWriter, r *http.Request) {
	var logoPass service.LogoPass

	session, _ := app.cookieStorage.Get(r, "session.id")
	logoPass.Login = session.Values["login"].(string)

	listLogoPasses, err := app.UserStorage.BatchGetLogoPasses(logoPass.Login, r.Context())
	if err != nil {
		log.Printf("get logpass pairs: %s for user: %s", err, logoPass.Login)
		if errors.Is(err, storage.ErrEmpty) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}
	}
	render.JSON(w, r, listLogoPasses)
}

// uploadText handles uploading text data via http.Post request.
//
// Accepts json.Marshalled service.TextData struct with 'description' field obligatory.
// 'data' is optional but is recommended in the sake of common sense.
// 'overwrite' field is obligatory if the request is meant to update existing information.
//
// Returns:
//   - `400` if json is corrupted
//   - `409` if this data already exist and 'overwrite' flag is false or omitted
//   - `500` if storage methods fail to comprehend the request
//   - `201` if everything is OK
func (app *App) uploadText(w http.ResponseWriter, r *http.Request) {
	var text service.TextData

	session, _ := app.cookieStorage.Get(r, "session.id")
	text.Login = session.Values["login"].(string)

	err := json.NewDecoder(r.Body).Decode(&text)
	if err != nil {
		log.Printf("upload secret text: json parse error: %s", err)
		http.Error(w, "json parse error:", http.StatusBadRequest)
		return
	}

	err = app.UserStorage.PutText(text, r.Context())
	if err != nil {
		if errors.Is(err, storage.ErrOldData) || errors.Is(err, storage.ErrAlreadyExists) {
			http.Error(w, fmt.Sprint(err), http.StatusConflict)
			return
		}
		log.Printf("put secret text: save to db: %s for user: %s", err, text.Login)
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// batchDownloadTexts handles sending the list of all user's secret strings via http.Get request.
//
// All it needs to run is a user login from http.Cookie.
//
// Returns:
//   - `404` if user has no such data in storage
//   - `500` if storage methods fail to comprehend the request
//   - json.Marshalled struct of []service.TextData type that contains all fields got from storage
//     except 'overwrite' field by virtue of its needlessness
func (app *App) batchDownloadTexts(w http.ResponseWriter, r *http.Request) {
	var text service.TextData

	session, _ := app.cookieStorage.Get(r, "session.id")
	text.Login = session.Values["login"].(string)

	listTexts, err := app.UserStorage.BatchGetTexts(text.Login, r.Context())
	if err != nil {
		log.Printf("get secrets: %s for user: %s", err, text.Login)
		if errors.Is(err, storage.ErrEmpty) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}
	}
	render.JSON(w, r, listTexts)
}

// uploadCreditCard handles uploading credit card data via http.Post request.
//
// Accepts json.Marshalled service.CreditCard struct with 'number' field obligatory.
// 'holder', 'due_date', 'cvv', 'description' are optional but are recommended in the sake of common sense.
// 'overwrite' field is obligatory if the request is meant to update existing information.
//
// Returns:
//   - `400` if json is corrupted
//   - `409` if this data already exist and 'overwrite' flag is false or omitted
//   - `500` if storage methods fail to comprehend the request
//   - `201` if everything is OK
func (app *App) uploadCreditCard(w http.ResponseWriter, r *http.Request) {
	var card service.CreditCard

	session, _ := app.cookieStorage.Get(r, "session.id")
	card.Login = session.Values["login"].(string)

	err := json.NewDecoder(r.Body).Decode(&card)
	if err != nil {
		log.Printf("upload credit card: json parse error: %s", err)
		http.Error(w, "json parse error:", http.StatusBadRequest)
		return
	}

	err = app.UserStorage.PutCreditCard(card, r.Context())
	if err != nil {
		if errors.Is(err, storage.ErrOldData) || errors.Is(err, storage.ErrAlreadyExists) {
			http.Error(w, fmt.Sprint(err), http.StatusConflict)
			return
		}
		log.Printf("put credit card: save to db: %s for user: %s", err, card.Login)
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// batchDownloadCreditCards handles sending the list of all user's credit cards via http.Get request.
//
// All it needs to run is a user login from http.Cookie.
//
// Returns:
//   - `404` if user has no such data in storage
//   - `500` if storage methods fail to comprehend the request
//   - json.Marshalled struct of []service.CreditCard type that contains all fields got from storage
//     except 'overwrite' field by virtue of its needlessness
func (app *App) batchDownloadCreditCards(w http.ResponseWriter, r *http.Request) {
	var card service.CreditCard

	session, _ := app.cookieStorage.Get(r, "session.id")
	card.Login = session.Values["login"].(string)

	listCards, err := app.UserStorage.BatchGetCreditCards(card.Login, r.Context())
	if err != nil {
		log.Printf("get logpass pairs: %s for user: %s", err, card.Login)
		if errors.Is(err, storage.ErrEmpty) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}
	}
	render.JSON(w, r, listCards)
}

// uploadBinary handles uploading text data via http.Post request.
//
// Accepts json.Marshalled service.BinaryData struct with 'description' field obligatory.
// 'binary' is optional but is recommended in the sake of common sense.
// 'overwrite' field is obligatory if the request is meant to update existing information.
//
// Returns:
//   - `400` if json is corrupted
//   - `409` if this data already exist and 'overwrite' flag is false or omitted
//   - `500` if storage methods fail to comprehend the request
//   - `201` if everything is OK
func (app *App) uploadBinary(w http.ResponseWriter, r *http.Request) {
	var binary service.BinaryData

	session, _ := app.cookieStorage.Get(r, "session.id")
	binary.Login = session.Values["login"].(string)

	err := json.NewDecoder(r.Body).Decode(&binary)
	if err != nil {
		log.Printf("upload credit card: json parse error: %s", err)
		http.Error(w, "json parse error:", http.StatusBadRequest)
		return
	}

	err = app.UserStorage.PutBinary(binary, r.Context())
	if err != nil {
		log.Printf("put credit card: save to db: %s for user: %s", err, binary.Login)
		if errors.Is(err, storage.ErrOldData) || errors.Is(err, storage.ErrAlreadyExists) {
			http.Error(w, fmt.Sprint(err), http.StatusConflict)
			return
		}
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// downloadBinaryList handles sending the list of all user's binaries via http.Get request.
//
// All it needs to run is a user login from http.Cookie.
//
// Returns:
//   - `404` if user has no such data in storage
//   - `500` if storage methods fail to comprehend the request
//   - json.Marshalled struct of []service.BinaryData type that contains only the 'description' field
//     while all stored binary data itself can be too large to batch download all of them
func (app *App) downloadBinaryList(w http.ResponseWriter, r *http.Request) {
	var binary service.BinaryData

	session, _ := app.cookieStorage.Get(r, "session.id")
	binary.Login = session.Values["login"].(string)

	binaryList, err := app.UserStorage.GetBinaryList(binary.Login, r.Context())
	if err != nil {
		log.Printf("get binary list: %s for user: %s", err, binary.Login)
		if errors.Is(err, storage.ErrEmpty) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}
	}
	render.JSON(w, r, binaryList)
}

// downloadBinary handles sending a particular user binary data via http.Post request.
// Accepts json.Marshalled service.BinaryData struct with 'description' field obligatory.
// All other fields are optional but are recommended to be omitted in the sake of common sense.
//
// Returns:
//   - `400` if json is corrupted
//   - `404` if user has no such data in storage
//   - `500` if storage methods fail to comprehend the request
//   - json.Marshalled struct of service.BinaryData type that contains all fields got from storage
//     except 'overwrite' field by virtue of its needlessness
func (app *App) downloadBinary(w http.ResponseWriter, r *http.Request) {
	var binary service.BinaryData

	session, _ := app.cookieStorage.Get(r, "session.id")
	binary.Login = session.Values["login"].(string)

	err := json.NewDecoder(r.Body).Decode(&binary)
	if err != nil {
		log.Printf("get binary: json parse error: %s", err)
		http.Error(w, "json parse error:", http.StatusBadRequest)
		return
	}

	binary, err = app.UserStorage.GetBinary(binary, r.Context())
	if err != nil {
		log.Printf("get binary: %s for user: %s", err, binary.Login)
		if errors.Is(err, storage.ErrEmpty) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}
	}
	render.JSON(w, r, binary)
}

// handleDownloadExe lets user download an .zip with executable client file for specified platform
func (app *App) handleDownload(w http.ResponseWriter, r *http.Request) {
	platform := path.Base(r.URL.Path)
	var name string
	switch platform {
	case "linux":
		name = config.LinFileName
	case "mac":
		name = config.MacFileName
	default:
		name = config.WinFileName
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=gophkeeper.zip")

	zipFile, err := os.Open(app.config.DownloadFolder + name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer zipFile.Close()

	buffer := make([]byte, 1024)
	for {
		n, err := zipFile.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = w.Write(buffer[:n])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// handleDefault is meant to remind anyone that tries some kind of funny stuff with this server
// that they are not the only funny ones in the game
func (app *App) handleDefault(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://www.youtube.com/watch?v=dQw4w9WgXcQ", http.StatusTemporaryRedirect)
}
