package application

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/vagafonov/shortener/internal/container"
	"github.com/vagafonov/shortener/internal/contract"
	"github.com/vagafonov/shortener/internal/middleware"
	"github.com/vagafonov/shortener/internal/response"
	"github.com/vagafonov/shortener/internal/validate"
)

type Application struct {
	cnt *container.Container
}

func NewApplication(cnt *container.Container) *Application {
	return &Application{
		cnt: cnt,
	}
}

func (a *Application) Serve() error {
	if a.cnt.GetConfig().FileStoragePath != "" {
		restored, err := a.cnt.GetServiceURL().RestoreURLs(a.cnt.GetConfig().FileStoragePath)
		if err != nil {
			return fmt.Errorf("cannot restore URLs: %w", err)
		}
		a.cnt.GetLogger().Info().Msgf("restored urls %v", restored)
	}

	a.cnt.GetLogger().Info().Msgf("server started and listen %s", a.cnt.GetConfig().ServerURL)
	err := http.ListenAndServe(a.cnt.GetConfig().ServerURL, a.Routes()) //nolint:gosec
	if err != nil {
		return err
	}

	return nil
}

func (a *Application) Routes() *chi.Mux {
	r := chi.NewRouter()
	// Middleware для логирования запросов
	mw := middleware.NewMiddleware(a.cnt.GetLogger())
	r.Use(mw.WithLogging)
	r.Use(mw.WithCompress)
	r.Get("/{short_url}", a.getShortURL)
	r.Post("/", a.createShortURL)
	r.Post("/api/", a.createShortURL)
	r.Get("/ping", a.ping)

	r.Route("/api", func(r chi.Router) {
		r.Post("/shorten", a.shorten)
		r.Post("/shorten/batch", a.shortenBatch)
	})

	return r
}

func (a *Application) createShortURL(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)

		return
	}

	if len(body) == 0 {
		res.WriteHeader(http.StatusBadRequest)

		return
	}
	shortURL, err := a.cnt.GetServiceURL().MakeShortURL(string(body), a.cnt.GetConfig().ShortURLLength)
	statusCode := http.StatusCreated
	if err != nil {
		if errors.Is(err, contract.ErrURLAlreadyExists) {
			statusCode = http.StatusConflict
		} else {
			a.cnt.GetLogger().Err(err).Msg("cannot read body")
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
	}

	res.WriteHeader(statusCode)
	if _, err := fmt.Fprintf(res, "%s/%s", a.cnt.GetConfig().ResultURL, shortURL.Short); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func (a *Application) shorten(res http.ResponseWriter, req *http.Request) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		a.cnt.GetLogger().Warn().Str("error", err.Error()).Msg("cannot read body")
		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	validatedRequest := validate.NewValidator(a.cnt.GetLogger()).ShortenRequest(buf)
	if validatedRequest == nil {
		res.WriteHeader(http.StatusBadRequest)

		return
	}

	shortURL, err := a.cnt.GetServiceURL().MakeShortURL(validatedRequest.URL, a.cnt.GetConfig().ShortURLLength)
	statusCode := http.StatusCreated
	if err != nil {
		if errors.Is(err, contract.ErrURLAlreadyExists) {
			statusCode = http.StatusConflict
		} else {
			a.cnt.GetLogger().Err(err).Msg("cannot read body")
			http.Error(res, err.Error(), http.StatusInternalServerError)

			return
		}
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(statusCode)

	jsonRes, err := json.Marshal(response.ShortenResponse{
		Result: fmt.Sprintf("%s/%s", a.cnt.GetConfig().ResultURL, shortURL.Short),
	})
	if err != nil {
		a.cnt.GetLogger().Warn().Str("error", err.Error()).Msg("cannot encode response to JSON")
		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	_, err = res.Write(jsonRes)
	if err != nil {
		a.cnt.GetLogger().Warn().Str("error", err.Error()).Msg("cannot encode response to JSON")

		return
	}
}

func (a *Application) getShortURL(res http.ResponseWriter, req *http.Request) {
	shortURL, err := a.cnt.GetServiceURL().GetShortURL(chi.URLParam(req, "short_url"))
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	if shortURL == nil {
		res.WriteHeader(http.StatusNotFound)

		return
	}
	res.Header().Set("Location", shortURL.Original)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func (a *Application) ping(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := a.cnt.GetDB().PingContext(ctx); err != nil { //nolint:contextcheck
		a.cnt.GetLogger().Error().Err(err).Send()
		res.WriteHeader(http.StatusInternalServerError)
	}
	res.WriteHeader(http.StatusOK)
}

func (a *Application) shortenBatch(res http.ResponseWriter, req *http.Request) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		a.cnt.GetLogger().Warn().Str("error", err.Error()).Msg("cannot read body")
		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	validatedRequest, err := validate.NewValidator(a.cnt.GetLogger()).ShortenBatchRequest(buf)
	if err != nil {
		if errors.Is(err, validate.ErrValidateEmpty) {
			res.WriteHeader(http.StatusBadRequest)

			return
		}

		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	shortenBatchResponse, err := a.cnt.GetServiceURL().MakeShortURLBatch(
		validatedRequest,
		a.cnt.GetConfig().ShortURLLength,
		a.cnt.GetConfig().ResultURL,
	)
	if err != nil {
		a.cnt.GetLogger().Warn().Str("error", err.Error()).Msg("cannot make shorten batch")
		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	jsonRes, err := json.Marshal(shortenBatchResponse)
	if err != nil {
		a.cnt.GetLogger().Warn().Str("error", err.Error()).Msg("cannot encode response to JSON")
		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	res.Header().Set("content-type", "application/json")
	res.WriteHeader(http.StatusCreated)
	_, err = res.Write(jsonRes)
	if err != nil {
		a.cnt.GetLogger().Warn().Str("error", err.Error()).Msg("cannot encode response to JSON")

		return
	}
}
