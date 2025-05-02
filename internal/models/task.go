package models

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Status string

const (
	NEW         Status = "new"
	IN_PROGRESS Status = "in_progress"
	DONE        Status = "done"
)

type Task struct {
	ID          int64       `json:"id"`
	Title       string      `json:"title"`
	Description pgtype.Text `json:"description"`
	Status      Status      `json:"status"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}
