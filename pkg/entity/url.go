package entity

import (
	"time"

	"github.com/google/uuid"
)

//nolint:musttag
type URL struct {
	ID        string
	UUID      uuid.UUID
	Short     string
	Original  string
	UserID    uuid.UUID
	DeletedAt *time.Time
}
