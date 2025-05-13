package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"goproject/internal/storage"
	"goproject/internal/storage/postres/models"
	"time"

	_ "github.com/lib/pq"
)

type PostgresStorage struct {
	db *sql.DB
}

// New создает новое подключение к PostgreSQL
func New(connString string) (*PostgresStorage, error) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	return &PostgresStorage{db: db}, nil
}

// Init создает таблицу если ее нет
func (s *PostgresStorage) Init(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS records (
			id SERIAL PRIMARY KEY,
			data TEXT NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)
	`

	_, err := s.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	return nil
}

// Save сохраняет данные и возвращает ID
func (s *PostgresStorage) Save(ctx context.Context, prj models.Project, task models.Task, rep models.Report, dev models.Developer) (int, error) {
	// Начинаем транзакцию
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Откатываем при ошибке

	// 1. Сохраняем разработчика
	devQuery := `
        INSERT INTO developers (firstname, last_name, created_at, modified_at)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `
	err = tx.QueryRowContext(ctx, devQuery,
		dev.Firstname,
		dev.LastName,
		time.Now(), // CreatedAt
		time.Now(), // ModifiedAt
	).Scan(&dev.ID)
	if err != nil {
		return 0, fmt.Errorf("failed to save developer: %w", err)
	}

	// 2. Сохраняем проект
	prjQuery := `
        INSERT INTO projects (name, description, created_at, modified_at)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `
	err = tx.QueryRowContext(ctx, prjQuery,
		prj.Name,
		prj.Description,
		time.Now(), // CreatedAt
		time.Now(), // ModifiedAt
	).Scan(&prj.ID)
	if err != nil {
		return 0, fmt.Errorf("failed to save project: %w", err)
	}

	// 3. Сохраняем отчет
	repQuery := `
        INSERT INTO reports (developer_id, created_at)
        VALUES ($1, $2)
        RETURNING id
    `
	err = tx.QueryRowContext(ctx, repQuery,
		dev.ID,     // Связь с разработчиком
		time.Now(), // CreatedAt
	).Scan(&rep.ID)
	if err != nil {
		return 0, fmt.Errorf("failed to save report: %w", err)
	}

	// 4. Сохраняем задачу
	taskQuery := `
        INSERT INTO tasks (
            report_id, project_id, name, developer_note,
            estimate_planed, estimate_progress,
            start_timestamp, end_timestamp, created_at
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING id
    `
	err = tx.QueryRowContext(ctx, taskQuery,
		rep.ID, // Связь с отчетом
		prj.ID, // Связь с проектом
		task.Name,
		task.DeveloperNote,
		task.EstimatePlaned,
		task.EstimateProgress,
		task.StartTimestamp,
		task.EndTimestamp,
		time.Now(), // CreatedAt
	).Scan(&task.ID)
	if err != nil {
		return 0, fmt.Errorf("failed to save task: %w", err)
	}

	// Фиксируем транзакцию
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Возвращаем ID задачи как результат
	return int(task.ID), nil
}

// GetByID возвращает запись по ID
func (s *PostgresStorage) GetByID(ctx context.Context, id int) (*storage.Record, error) {
	query := `
		SELECT id, data, created_at
		FROM records
		WHERE id = $1
	`

	var record storage.Record
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&record.ID,
		&record.Data,
		&record.Time,
	)

	if err == sql.ErrNoRows {
		return nil, storage.ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get record: %w", err)
	}

	return &record, nil
}

// Close закрывает соединение с БД
func (s *PostgresStorage) Close() error {
	return s.db.Close()
}
