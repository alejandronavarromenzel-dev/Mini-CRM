package models

type Task struct {
	ID         int64
	ClientID   int64
	ClientName string

	Title       string
	Status      string
	Owner       string
	Tag         string
	Archived    bool
}
