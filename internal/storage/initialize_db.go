package storage

import (
	"gophkeeper/internal/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

// DBStorage hold pointer to gorm type DB
type DBStorage struct {
	db *gorm.DB
}

// NewUserStorage is used to open a connection to a postgres db and migrate all the tables needed
func NewUserStorage(databaseURL string) *DBStorage {
	connection, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("database failed to open: %s", err)
	}
	sqlDB, err := connection.DB()
	if err != nil {
		log.Fatalf("database failed to connect: %s", err)
	}

	err = sqlDB.Ping()
	if err != nil {
		log.Fatalf("database failed to ping: %s", err)
	}
	log.Printf("Database connection successful")

	initializeTables(connection)

	return &DBStorage{
		db: connection,
	}
}

func initializeTables(connection *gorm.DB) {
	err := connection.AutoMigrate(service.User{})
	if err != nil {
		log.Fatalf("database failed to create user table: %s", err)
	}
	err = connection.AutoMigrate(service.LogoPass{})
	if err != nil {
		log.Fatalf("database failed to create user table: %s", err)
	}
	err = connection.AutoMigrate(service.TextData{})
	if err != nil {
		log.Fatalf("database failed to create user table: %s", err)
	}
	err = connection.AutoMigrate(service.CreditCard{})
	if err != nil {
		log.Fatalf("database failed to create user table: %s", err)
	}
	err = connection.AutoMigrate(service.BinaryData{})
	if err != nil {
		log.Fatalf("database failed to create user table: %s", err)
	}
}
