package config

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"   envDefault:"localhost:8080"`
	DatabaseDSN   string `env:"DATABASE_URI"     envDefault:"postgres://matt:pvtjoker@localhost:5432/gophkeeper?sslmode=disable"`
}
