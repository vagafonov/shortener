package entity

import "github.com/google/uuid"

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	Result string `json:"result"`
}

//nolint:musttag
type URL struct {
	UUID  uuid.UUID
	Short string
	Full  string
}
