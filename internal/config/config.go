// Package config holds the configuration credentials needed for the App
package config

// Config struct holds server address of where the App is working and a db access url
type Config struct {
	ServerAddress  string `env:"SERVER_ADDRESS"  envDefault:"localhost:8080"`
	DatabaseDSN    string `env:"DATABASE_URI"    envDefault:"postgres://matt:pvtjoker@localhost:5432/gophkeeper?sslmode=disable"`
	DownloadFolder string `env:"DOWNLOAD_FOLDER" envDefault:"C:/repo/GO/gophkeeper/internal/bin/"`
}

const (
	WinFileName = "gophkeeper_windows.zip"
	MacFileName = "gophkeeper_mac.zip"
	LinFileName = "gophkeeper_linux.zip"
)
