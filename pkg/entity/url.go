package entity

import "github.com/google/uuid"

//nolint:musttag
type URL struct {
	UUID     uuid.UUID
	Short    string
	Original string
}
