package entity

import (
	"time"

	"github.com/google/uuid"
)

// URL entity.
type URL struct {
	ID        string     `json:"id"`
	UUID      uuid.UUID  `json:"uuid"`
	Short     string     `json:"short"`
	Original  string     `json:"original"`
	UserID    uuid.UUID  `json:"userId"`
	DeletedAt *time.Time `json:"deletedAt"`
}
