package storage

import (
	"context"
	"goproject/internal/storage/postres/models"
)

type ReportRepository interface {
	SaveReport(ctx context.Context, rep *models.Report) error
}

type TaskRepository interface {
	SaveTask(ctx context.Context, task *models.Task) error
}

type ProjectRepository interface {
	SaveProject(ctx context.Context, prj *models.Project) error
}

type DeveloperStorage interface {
	SaveDeveloper(ctx context.Context, dev *models.Developer) error
	GetDeveloper(ctx context.Context, id uint) (*models.Developer, error)
	PachDeveloper(ctx context.Context, id uint) error
	Close() error
}
