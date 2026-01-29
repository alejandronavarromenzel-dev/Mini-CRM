package models

import "time"

type Task struct {
	ID        int64
	ClientID  int64
	Title     string
	Status    string
	Priority  string
	Owner     string
	Progress  int
	DueDate   *time.Time
	CreatedAt time.Time
}
