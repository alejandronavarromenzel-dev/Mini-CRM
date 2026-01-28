package models

import "time"

type Client struct {
	ID        int64
	Name      string
	Status    string
	Owner     string
	Tags      string
	Notes     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
