package validate

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"

	"github.com/rs/zerolog"
	"github.com/vagafonov/shortener/internal/request"
)

// ErrValidateEmpty error for empty input.
var ErrValidateEmpty = errors.New("empty")

type validator struct {
	logger *zerolog.Logger
}

// NewValidator Constructor for Validator.
func NewValidator(l *zerolog.Logger) *validator {
	return &validator{logger: l}
}

// ShortenRequest create ShortenRequest from input.
func (v *validator) ShortenRequest(buf bytes.Buffer) *request.ShortenRequest {
	var shortenReq request.ShortenRequest
	if err := json.Unmarshal(buf.Bytes(), &shortenReq); err != nil {
		v.logger.Warn().Str("error", err.Error()).Str("request", buf.String()).Msg("cannot unmarshal request")

		return nil
	}

	return &shortenReq
}

// ShortenBatchRequest create ShortenBatchRequest from input.
func (v *validator) ShortenBatchRequest(ctx context.Context, buf bytes.Buffer) ([]request.ShortenBatchRequest, error) {
	var req []request.ShortenBatchRequest
	if err := json.Unmarshal(buf.Bytes(), &req); err != nil {
		v.logger.Warn().Str("error", err.Error()).Str("request", buf.String()).Msg("cannot unmarshal shorten batch request")

		return nil, err
	}
	if len(req) == 0 {
		return nil, ErrValidateEmpty
	}
	for _, v := range req {
		if v.CorrelationID == "" {
			return nil, ErrValidateEmpty
		}

		if v.OriginalURL == "" {
			return nil, ErrValidateEmpty
		}
	}

	return req, nil
}

// DeleteUserURLsRequest create slice of string for DeleteUserURLsRequest.
func (v *validator) DeleteUserURLsRequest(ctx context.Context, buf bytes.Buffer) ([]string, error) {
	var req []string
	if err := json.Unmarshal(buf.Bytes(), &req); err != nil {
		v.logger.
			Warn().
			Str("error", err.Error()).
			Str("request", buf.String()).
			Msg("cannot unmarshal delete user URLs request")

		return nil, err
	}
	if len(req) == 0 {
		return nil, ErrValidateEmpty
	}

	return req, nil
}
