package request

// ShortenBatchRequest.
type ShortenBatchRequest struct {
	CorrelationID string `json:"correlation_id"` //nolint:tagliatelle
	OriginalURL   string `json:"original_url"`   //nolint:tagliatelle
}
