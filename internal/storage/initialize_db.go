package storage

import (
	"gophkeeper/internal/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type DBStorage struct {
	db *gorm.DB
}

func NewUserStorage(DatabaseURL string) *DBStorage {
	connection, err := gorm.Open(postgres.Open(DatabaseURL), &gorm.Config{})
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

	InitializeTables(connection)

	return &DBStorage{
		db: connection,
	}
}

func InitializeTables(connection *gorm.DB) {
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
	err = connection.AutoMigrate(service.UserBinaryList{})
	if err != nil {
		log.Fatalf("database failed to create user table: %s", err)
	}
	err = connection.AutoMigrate(service.BinaryData{})
	if err != nil {
		log.Fatalf("database failed to create user table: %s", err)
	}
}
