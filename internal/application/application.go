package application

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	_ "net/http/pprof" //nolint:gosec
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/vagafonov/shortener/internal/config"
	"github.com/vagafonov/shortener/internal/container"
	"github.com/vagafonov/shortener/internal/cookie"
	"github.com/vagafonov/shortener/internal/customerror"
	"github.com/vagafonov/shortener/internal/middleware"
	"github.com/vagafonov/shortener/internal/response"
	"github.com/vagafonov/shortener/internal/validate"
	"github.com/vagafonov/shortener/pkg/encrypting"
	"github.com/vagafonov/shortener/pkg/entity"
)

// Application Contains routes and starts the server.
type Application struct {
	cnt *container.Container
}

// Constructor for application.
func NewApplication(cnt *container.Container) *Application {
	return &Application{
		cnt: cnt,
	}
}

// Serve run server.
func (a *Application) Serve(ctx context.Context) error {
	a.restoreURLs(ctx)
	a.cnt.GetLogger().Info().Msgf("server started and listen %s", a.cnt.GetConfig().ServerURL)
	err := http.ListenAndServe(a.cnt.GetConfig().ServerURL, a.Routes()) //nolint:gosec
	if err != nil {
		return err
	}

	return nil
}

// ServeHTTPS Serve run HTTPS server.
func (a *Application) ServeHTTPS(ctx context.Context) error {
	a.restoreURLs(ctx)
	_, err := encrypting.GenerateCertificate("localhost") // TODO
	if err != nil {
		a.cnt.GetLogger().Err(err).Msg("failed to generate certificate")
	}

	a.cnt.GetLogger().Info().Msgf("HTPS server started and listen URL: %s", a.cnt.GetConfig().ServerURL)
	//nolint:gosec
	err = http.ListenAndServeTLS(a.cnt.GetConfig().ServerURL, "certs/server.crt", "certs/server.key", a.Routes())
	if err != nil {
		a.cnt.GetLogger().Err(err).Msg("fail to start TLS server")

		return err
	}

	return nil
}

func (a *Application) restoreURLs(ctx context.Context) {
	if a.cnt.GetConfig().FileStoragePath != "" {
		restored, err := a.cnt.GetServiceURL().RestoreURLs(ctx, a.cnt.GetConfig().FileStoragePath)
		if err != nil {
			a.cnt.GetLogger().Info().Msgf("cannot restore URLs: %s", err.Error())
		}
		a.cnt.GetLogger().Info().Msgf("restored urls %v", restored)
	}
}

// Routes register routes and middlewares.
func (a *Application) Routes() *chi.Mux {
	r := chi.NewRouter()
	// Middleware для логирования запросов
	mw := middleware.NewMiddleware(a.cnt.GetLogger())
	r.Use(mw.WithLogging)
	r.Use(mw.WithCompress)
	r.Use(func(handler http.Handler) http.Handler {
		return mw.WithUserIDCookie(handler, a.cnt.GetConfig().CryptoKey)
	})
	if a.cnt.GetConfig().Mode == config.ModeDev {
		r.Mount("/debug", chimiddleware.Profiler())
	}
	r.Get("/{short_url}", a.getShortURL)
	r.Post("/", a.createShortURL)
	r.Post("/api/", a.createShortURL)
	r.Get("/ping", a.ping)

	r.Route("/api", func(r chi.Router) {
		r.Post("/shorten", a.shorten)
		r.Post("/shorten/batch", a.shortenBatch)
		r.Get("/user/urls", a.userUrls)
		r.Delete("/user/urls", a.deleteUserURLs)
	})

	return r
}

func (a *Application) createShortURL(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		a.cnt.GetLogger().Err(err).Msg("cannot get read body")
		http.Error(res, err.Error(), http.StatusInternalServerError)

		return
	}

	if len(body) == 0 {
		res.WriteHeader(http.StatusBadRequest)

		return
	}

	userID, err := a.getUserIDFromCookie(req)
	if err != nil {
		a.cnt.GetLogger().Err(err).Msg("cannot get cookie with userID")
		http.Error(res, err.Error(), http.StatusInternalServerError)

		return
	}

	shortURL, err := a.cnt.GetServiceURL().MakeShortURL(
		req.Context(),
		string(body),
		a.cnt.GetConfig().ShortURLLength,
		userID,
	)
	statusCode := http.StatusCreated
	if err != nil {
		if errors.Is(err, customerror.ErrURLAlreadyExists) {
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

//nolint:funlen
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

	userID, err := a.getUserIDFromCookie(req)
	if err != nil {
		a.cnt.GetLogger().Err(err).Msg("cannot get cookie with userID")
		http.Error(res, err.Error(), http.StatusInternalServerError)

		return
	}

	shortURL, err := a.cnt.GetServiceURL().MakeShortURL(
		req.Context(),
		validatedRequest.URL,
		a.cnt.GetConfig().ShortURLLength,
		userID,
	)
	statusCode := http.StatusCreated
	if err != nil {
		if errors.Is(err, customerror.ErrURLAlreadyExists) {
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
	shortURL, err := a.cnt.GetServiceURL().GetShortURL(req.Context(), chi.URLParam(req, "short_url"))
	if err != nil {
		if errors.Is(err, customerror.ErrURLDeleted) {
			a.cnt.GetLogger().Info().Msg("trying to get deleted address")
			res.WriteHeader(http.StatusGone)
		}

		a.cnt.GetLogger().Error().Err(err).Msg("cannot get short url")
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
	if err := a.cnt.GetServiceHealthCheck().Ping(ctx); err != nil { //nolint:contextcheck
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

	validatedRequest, err := validate.NewValidator(a.cnt.GetLogger()).ShortenBatchRequest(req.Context(), buf)
	if err != nil {
		if errors.Is(err, validate.ErrValidateEmpty) {
			res.WriteHeader(http.StatusBadRequest)

			return
		}

		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	URLs := make([]*entity.URL, len(validatedRequest))
	for k, v := range validatedRequest {
		URLs[k] = &entity.URL{
			ID:       v.CorrelationID,
			Short:    a.cnt.GetHasher().Hash(a.cnt.GetConfig().ShortURLLength),
			Original: v.OriginalURL,
		}
	}
	shortenBatchResponse, err := a.cnt.GetServiceURL().MakeShortURLBatch(
		req.Context(),
		URLs,
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

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	_, err = res.Write(jsonRes)
	if err != nil {
		a.cnt.GetLogger().Warn().Str("error", err.Error()).Msg("cannot encode response to JSON")

		return
	}
}

func (a *Application) userUrls(res http.ResponseWriter, req *http.Request) {
	userID, err := a.getUserIDFromCookie(req)
	if err != nil {
		a.cnt.GetLogger().Err(err).Msg("cannot get cookie with userID")
		http.Error(res, err.Error(), http.StatusInternalServerError)

		return
	}

	if userID == uuid.Nil {
		a.cnt.GetLogger().Err(err).Msg("cookie with userID is empty")
		res.WriteHeader(http.StatusUnauthorized)

		return
	}

	userURLs, err := a.cnt.GetServiceURL().GetUserURLs(req.Context(), userID, a.cnt.GetConfig().ResultURL)
	if err != nil {
		a.cnt.GetLogger().Warn().Str("error", err.Error()).Msg("cannot get user URLs")
		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	if len(userURLs) == 0 {
		res.WriteHeader(http.StatusUnauthorized)

		return
	}

	userURLsResp := make([]response.UserURLResponse, len(userURLs))
	for k, v := range userURLs {
		userURLsResp[k] = response.NewUserURLResponse(v.Short, v.Original)
	}

	jsonRes, err := json.Marshal(userURLsResp)
	if err != nil {
		a.cnt.GetLogger().Warn().Str("error", err.Error()).Msg("cannot encode get user URLs response to JSON")
		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	_, err = res.Write(jsonRes)
	if err != nil {
		a.cnt.GetLogger().Warn().Str("error", err.Error()).Msg("cannot write result to response")

		return
	}
}

func (a *Application) deleteUserURLs(res http.ResponseWriter, req *http.Request) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		a.cnt.GetLogger().Warn().Str("error", err.Error()).Msg("cannot read body")
		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	validatedRequest, err := validate.NewValidator(a.cnt.GetLogger()).DeleteUserURLsRequest(req.Context(), buf)
	if err != nil {
		if errors.Is(err, validate.ErrValidateEmpty) {
			res.WriteHeader(http.StatusBadRequest)

			return
		}

		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	userID, err := a.getUserIDFromCookie(req)
	if err != nil {
		a.cnt.GetLogger().Err(err).Msg("cannot get cookie with userID")
		http.Error(res, err.Error(), http.StatusInternalServerError)

		return
	}

	if userID == uuid.Nil {
		a.cnt.GetLogger().Err(err).Msg("cookie with userID is empty")
		res.WriteHeader(http.StatusUnauthorized)

		return
	}

	err = a.cnt.GetServiceURL().DeleteUserURLs(
		req.Context(),
		userID,
		validatedRequest,
		a.cnt.GetConfig().DeleteURLsBatchSize,
		a.cnt.GetConfig().DeleteURLsJobsCount,
	)
	if err != nil {
		a.cnt.GetLogger().Warn().Str("error", err.Error()).Msg("cannot delete user URLs")
		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	res.WriteHeader(http.StatusAccepted)
}

func (a *Application) getUserIDFromCookie(req *http.Request) (uuid.UUID, error) {
	userIDCoockie, err := req.Cookie("userID")
	if err != nil {
		if !errors.Is(err, http.ErrNoCookie) {
			a.cnt.GetLogger().Err(err).Msg("cannot get cookie with userID")

			return uuid.Nil, err
		}
	}

	if userIDCoockie == nil {
		return uuid.Nil, nil
	}

	decr, err := cookie.Decrypt(userIDCoockie.Value, a.cnt.GetConfig().CryptoKey)
	if err != nil {
		return uuid.Nil, err
	}

	decryptedUUID, err := uuid.Parse(*decr)
	if err != nil {
		return uuid.Nil, err
	}

	return decryptedUUID, err
}
