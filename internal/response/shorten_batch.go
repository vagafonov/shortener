package response

// ShortenBatchResponse.
type ShortenBatchResponse struct {
	CorrelationID string `json:"correlation_id"` //nolint:tagliatelle
	ShortURL      string `json:"short_url"`      //nolint:tagliatelle
}
