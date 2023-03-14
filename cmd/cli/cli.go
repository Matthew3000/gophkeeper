package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"gophkeeper/cmd/cli/client"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

var (
	serverURL  string
	store      bool
	retrieve   bool
	dataFile   string
	outputFile string
)

func main() {
	var cfg client.Config

	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}
	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "Server address")
	flag.StringVar(&cfg.OutputFolder, "o", cfg.OutputFolder, "Output folder for files")
	flag.Parse()

	fmt.Print("Welcome at Gophkeeper")
	printBuildData()

	err := client.Communicate()
	if err != nil {
		log.Fatal(err)
	}

	if store && retrieve {
		fmt.Fprintln(os.Stderr, "Error: cannot store and retrieve at the same time")
		os.Exit(1)
	}

	if store {
		err := storeData(dataFile, serverURL+"/store")
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error storing data:", err)
			os.Exit(1)
		}
		fmt.Println("Data stored successfully")
	}

	if retrieve {
		err := retrieveData(outputFile, serverURL+"/retrieve")
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error retrieving data:", err)
			os.Exit(1)
		}
		fmt.Println("Data retrieved successfully")
	}
}

func storeData(filename string, url string) error {
	// Read the data from the file
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	// Send the data in a POST request to the server
	resp, err := http.Post(url, "application/octet-stream", bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status code %d", resp.StatusCode)
	}

	return nil
}

func retrieveData(filename string, url string) error {
	// Send a GET request to the server to retrieve the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status code %d", resp.StatusCode)
	}

	// Read the binary data from the response body
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Write the binary data to the output file
	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
