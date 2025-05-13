package models

import "time"

// Developer - модель разработчика
type Developer struct {
	ID         uint
	Firstname  string
	LastName   string
	CreatedAt  time.Time
	ModifiedAt time.Time
	DeletedAt  *time.Time

	Reports []Report
}

// Report - отчёт разработчика
type Report struct {
	ID          uint
	DeveloperID uint
	CreatedAt   time.Time

	Developer *Developer
	Tasks     []Task
}

// Task - задача в отчёте
type Task struct {
	ID               uint
	ReportID         uint
	ProjectID        uint
	Name             string
	DeveloperNote    string
	EstimatePlaned   int
	EstimateProgress int
	StartTimestamp   time.Time
	EndTimestamp     time.Time
	CreatedAt        time.Time

	Report  *Report
	Project *Project
}

// Project - проект
type Project struct {
	ID          uint
	Name        string
	Description string
	CreatedAt   time.Time
	ModifiedAt  time.Time

	Tasks []Task
}
