// Package main initializes the cli of Gophkeeper secrets manager service
package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"gophkeeper/cmd/cli/client"
	"log"
	"net"
	"time"
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

	fmt.Println("Welcome to Gophkeeper")
	printBuildData()

	var api = client.NewApi(cfg.ServerAddress)
	var storage, err = client.NewStorage(cfg.OutputFolder)
	if err != nil {
		log.Fatal(err)
	}
	var service = client.NewService(cfg, api, storage)

	tickerUpdate := time.NewTicker(client.UpdateDataTimer)
	go func() {
		for range tickerUpdate.C {
			log.Printf("Updating data from server")
			err = service.UpdateAll()
			if err != nil {
				if errors.Is(err, &net.OpError{}) {
					fmt.Println("Update failed due to poor internet connection, continuing offline")
				}
				log.Printf("update data from server: %s", err)
			}
		}
	}()

	err = service.StartCommunicate()
	if err != nil {
		log.Fatal(err)
	}

}
