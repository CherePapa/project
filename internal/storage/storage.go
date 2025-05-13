package storage

import (
	"context"
	"time"
)

type Record struct {
	ID   int
	Data string
	Time time.Time
}

type Storage interface {
	Init(ctx context.Context) error
	Save(ctx context.Context, data string) (int, error)
	GetByID(ctx context.Context, id int) (*Record, error)
	Close() error
}
