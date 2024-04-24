package entity

import "github.com/google/uuid"

//nolint:musttag
type URL struct {
	ID       string
	UUID     uuid.UUID
	Short    string
	Original string
}
