package main

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"gophkeeper/cmd/cli/client"
	"log"
)

var (
	// To get desired build credentials use commands:
	// For win powershell:
	// go run -ldflags "-X main.buildVersion=v1.0.1 -X main.buildCommit=07fa3a5 -X 'main.buildDate=$(Get-Date -uformat %Y/%m/%d-%H:%M)'" cli.go
	// For unix:
	// go run -ldflags "-X main.buildVersion=v1.0.1 -X main.buildCommit=07fa3a5 -X 'main.buildDate=$(date +'%Y/%m/%d-%H:%M')'" cli.go
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
	var cfg client.Config

	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}
	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "Server address")
	flag.StringVar(&cfg.OutputFolder, "o", cfg.OutputFolder, "Output folder for files")
	flag.Parse()

	fmt.Print("Welcome to Gophkeeper")
	printBuildData()

	var api = client.NewApi(cfg.ServerAddress)
	var storage = client.NewStorage(cfg.OutputFolder)
	var service = client.NewService(cfg, api, storage)

	err := service.Communicate()
	if err != nil {
		log.Fatal(err)
	}

}
