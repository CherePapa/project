package main

import (
	"context"
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

	// Инициализация таблицы
	ctx := context.Background()
	if err := storage.Init(ctx); err != nil {
		log.Fatalf("failed to init storage: %v", err)
	}

	// Тут нужно save поменять, ибо мы в там его поменяли
	id, err := storage.Save(ctx, "test data")
	if err != nil {
		log.Printf("failed to save: %v", err)
	} else {
		log.Printf("saved with id: %d", id)
	}

	// Пример получения данных
	record, err := storage.GetByID(ctx, id)
	if err != nil {
		log.Printf("failed to get record: %v", err)
	} else {
		log.Printf("got record: %+v", record)
	}
}
