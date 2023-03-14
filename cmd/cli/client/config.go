package client

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"   envDefault:"localhost:8080"`
	OutputFolder  string `env:"OUTPUT_FOLDER"     envDefault:"/AppData"`
}
