package main

import (
	postgres "goproject/internal/storage/postres"
	"log"
)

func main() {
	// Подключение к PostgreSQL
	storage, err := postgres.New("postgres://user:password@localhost/dbname?sslmode=disable")
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer storage.Close()

}
