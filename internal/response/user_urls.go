package response

// UserURLResponse.
type UserURLResponse struct {
	ShortURL    string `json:"short_url"`    //nolint:tagliatelle
	OriginalURL string `json:"original_url"` //nolint:tagliatelle
}

// NewUserURLResponse Constructor for UserURLResponse.
func NewUserURLResponse(shortURL string, originalURL string) UserURLResponse {
	return UserURLResponse{
		ShortURL:    shortURL,
		OriginalURL: originalURL,
	}
}
