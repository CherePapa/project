package main

import (
	"context"
	postgres "goproject/internal/storage/postres"
	"goproject/internal/storage/postres/models"
	"log"
)

func main() {
	// Подключение к PostgreSQL
	storage, err := postgres.New("postgres://user:password@localhost/dbname?sslmode=disable")
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer storage.Close()

	// Инициализация таблицы
	ctx := context.Background()
	if err := storage.Init(ctx); err != nil {
		log.Fatalf("failed to init storage: %v", err)
	}
ы
	devService := .NewDeveloperService(storage)

	// Создание разработчика
	dev := models.Developer{
		Firstname: "John",
		LastName:  "Doe",
	}

	if err := devService.CreateDeveloper(context.Background(), &dev); err != nil {
		log.Fatal("Failed to create developer:", err)
	}

	log.Println("Developer created successfully!")
}
