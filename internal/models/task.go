package models

import "time"

type Task struct {
	Id        int
	Title     string
	CreatedAt time.Time
	DeletedAt *time.Time
}
