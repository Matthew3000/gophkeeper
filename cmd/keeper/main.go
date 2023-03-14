package main

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/gorilla/sessions"
	"gophkeeper/internal/app"
	"gophkeeper/internal/config"
	"gophkeeper/internal/service"
	"gophkeeper/internal/storage"
	"log"
)

var (
	// To get desired build credentials use commands:
	// For win powershell:
	// go run -ldflags "-X main.buildVersion=v1.0.1 -X main.buildCommit=07fa3a5 -X 'main.buildDate=$(Get-Date -uformat %Y/%m/%d-%H:%M)'" main.go
	// For unix:
	// go run -ldflags "-X main.buildVersion=v1.0.1 -X main.buildCommit=07fa3a5 -X 'main.buildDate=$(date +'%Y/%m/%d-%H:%M')'" main.go
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func printBuildData() {
	fmt.Println("Build version:", buildVersion)
	fmt.Println("Build date:", buildDate)
	fmt.Println("Build commit:", buildCommit)
}

func main() {
	printBuildData()
	var cfg config.Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}
	flag.StringVar(&cfg.DatabaseDSN, "d", cfg.DatabaseDSN, "File Storage Path")
	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "Server address")
	flag.Parse()

	fmt.Println(cfg.DatabaseDSN)

	userStorage := storage.NewUserStorage(cfg.DatabaseDSN)
	cookieStorage := sessions.NewCookieStore([]byte(service.SecretKey))
	var application = app.NewApp(cfg, userStorage, *cookieStorage)
	application.Run()
}
