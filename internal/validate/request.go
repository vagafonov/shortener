package validate

import (
	"bytes"
	"encoding/json"

	"github.com/rs/zerolog"
	"github.com/vagafonov/shortener/pkg/entity"
)

type validator struct {
	logger zerolog.Logger
}

func NewValidator(l zerolog.Logger) *validator {
	return &validator{logger: l}
}

func (v *validator) ShortenRequest(buf bytes.Buffer) *entity.ShortenRequest {
	var shortenReq entity.ShortenRequest
	if err := json.Unmarshal(buf.Bytes(), &shortenReq); err != nil {
		v.logger.Warn().Str("error", err.Error()).Str("request", buf.String()).Msg("cannot unmarshal request")

		return nil
	}

	return &shortenReq
}
