package application

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/vagafonov/shortener/internal/middleware"
	"github.com/vagafonov/shortener/internal/validate"
	"github.com/vagafonov/shortener/pkg/entity"
)

type Application struct {
	cnt     *Container
	service Service
}

func NewApplication(cnt *Container) *Application {
	return &Application{
		cnt: cnt,
		service: NewService(
			cnt.logger,
			cnt.GetStorage(),
			cnt.GetBackupStorage(),
			cnt.GetHasher(),
		),
	}
}

func (a *Application) Serve() error {
	restored, err := a.service.RestoreURLs(a.cnt.cfg.FileStoragePath)
	if err != nil {
		return err
	}
	a.cnt.logger.Info().Msgf("restored urls %v", restored)

	a.cnt.logger.Info().Msgf("server started and listen %s", a.cnt.cfg.ServerURL)
	err = http.ListenAndServe(a.cnt.cfg.ServerURL, a.Routes()) //nolint:gosec
	if err != nil {
		return err
	}

	return nil
}

func (a *Application) Routes() *chi.Mux {
	r := chi.NewRouter()
	// Middleware для логирования запросов
	mw := middleware.NewMiddleware(a.cnt.logger)
	r.Use(mw.WithLogging)
	r.Use(mw.WithCompress)
	r.Get("/{short_url}", a.getShortURL)
	r.Post("/", a.createShortURL)
	r.Post("/api/", a.createShortURL)
	r.Get("/ping", a.ping)

	r.Route("/api", func(r chi.Router) {
		r.Post("/shorten", a.shorten)
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
	shortURL, err := NewService(
		a.cnt.logger,
		a.cnt.GetStorage(),
		a.cnt.GetBackupStorage(),
		a.cnt.GetHasher(),
	).MakeShortURL(string(body), a.cnt.cfg.ShortURLLength)
	// TODO check error type
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
	res.WriteHeader(http.StatusCreated)

	if _, err := fmt.Fprintf(res, "%s/%s", a.cnt.cfg.ResultURL, shortURL.Short); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func (a *Application) shorten(res http.ResponseWriter, req *http.Request) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		a.cnt.logger.Warn().Str("error", err.Error()).Msg("cannot read body")
		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	validatedRequest := validate.NewValidator(a.cnt.logger).ShortenRequest(buf)
	if validatedRequest == nil {
		res.WriteHeader(http.StatusBadRequest)

		return
	}

	svc := NewService(
		a.cnt.logger,
		a.cnt.GetStorage(),
		a.cnt.GetBackupStorage(),
		a.cnt.GetHasher(),
	)
	shortURL, err := svc.MakeShortURL(validatedRequest.URL, a.cnt.cfg.ShortURLLength)
	if err != nil {
		a.cnt.logger.Warn().Str("error", err.Error()).Msg("cannot make short")
		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	jsonRes, err := json.Marshal(entity.ShortenResponse{
		Result: fmt.Sprintf("%s/%s", a.cnt.cfg.ResultURL, shortURL.Short),
	})
	if err != nil {
		a.cnt.logger.Warn().Str("error", err.Error()).Msg("cannot encode response to JSON")
		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	res.Header().Set("content-type", "application/json")
	res.WriteHeader(http.StatusCreated)
	_, err = res.Write(jsonRes)
	if err != nil {
		a.cnt.logger.Warn().Str("error", err.Error()).Msg("cannot encode response to JSON")

		return
	}
}

func (a *Application) getShortURL(res http.ResponseWriter, req *http.Request) {
	shortURL, err := NewService(
		a.cnt.logger,
		a.cnt.GetStorage(),
		a.cnt.GetBackupStorage(),
		a.cnt.GetHasher(),
	).GetShortURL(chi.URLParam(req, "short_url"))
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	if shortURL == nil {
		res.WriteHeader(http.StatusNotFound)

		return
	}
	res.Header().Set("Location", shortURL.Full)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func (a *Application) ping(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := a.cnt.db.PingContext(ctx); err != nil { //nolint:contextcheck
		a.cnt.logger.Error().Err(err).Send()
		res.WriteHeader(http.StatusInternalServerError)
	}
	res.WriteHeader(http.StatusOK)
}
