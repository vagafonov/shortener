package application

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/vagafonov/shortener/internal/application/services"
	"github.com/vagafonov/shortener/pkg/entity"
)

type Application struct {
	cnt    *Container
	logger zerolog.Logger
}

func NewApplication(cnt *Container) *Application {
	// Инициализация логгера zerolog
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	// human-friendly и цветной output
	logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr}) //nolint:exhaustruct
	// Уровень логирования
	zerolog.SetGlobalLevel(cnt.cfg.LogLevel)

	return &Application{cnt: cnt, logger: logger}
}

func (a *Application) Serve() error {
	a.logger.Info().Msgf("server started and listen %s", a.cnt.cfg.ServerURL)
	err := http.ListenAndServe(a.cnt.cfg.ServerURL, a.Routes()) //nolint:gosec
	if err != nil {
		return err
	}

	return nil
}

func (a *Application) Routes() *chi.Mux {
	r := chi.NewRouter()
	// Middleware для логирования запросов
	r.Use(a.withLogging)
	r.Get("/{short_url}", a.getShortURL)
	r.Post("/", a.createShortURL)
	r.Post("/api/", a.createShortURL)

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
	shortURL, err := services.NewService(a.cnt.GetStorage()).MakeShortURL(string(body), a.cnt.cfg.ShortURLLength)
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
		a.logger.Warn().Str("error", err.Error()).Msg("cannot read body")
		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	var shortenReq entity.ShortenRequest
	if err = json.Unmarshal(buf.Bytes(), &shortenReq); err != nil {
		a.logger.Warn().Str("error", err.Error()).Str("request", buf.String()).Msg("cannot unmarshal request")
		res.WriteHeader(http.StatusBadRequest)

		return
	}

	shortURL, err := services.NewService(a.cnt.GetStorage()).MakeShortURL(shortenReq.URL, a.cnt.cfg.ShortURLLength)
	if err != nil {
		a.logger.Warn().Str("error", err.Error()).Msg("cannot make short")
		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	jsonRes, err := json.Marshal(entity.ShortenResponse{
		Result: fmt.Sprintf("%s/%s", a.cnt.cfg.ResultURL, shortURL.Short),
	})
	if err != nil {
		a.logger.Warn().Str("error", err.Error()).Msg("cannot encode response to JSON")
		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	res.Header().Set("content-type", "application/json")
	res.WriteHeader(http.StatusCreated)
	_, err = res.Write(jsonRes)
	if err != nil {
		a.logger.Warn().Str("error", err.Error()).Msg("cannot encode response to JSON")

		return
	}
}

func (a *Application) getShortURL(res http.ResponseWriter, req *http.Request) {
	shortURL := services.NewService(a.cnt.GetStorage()).GetShortURL(chi.URLParam(req, "short_url"))
	if shortURL == nil {
		res.WriteHeader(http.StatusNotFound)

		return
	}
	res.Header().Set("Location", shortURL.Full)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

// WithLogging middleware для логирования.
func (a *Application) withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lw := loggingResponseWriter{
			ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
			responseData: &responseData{
				status: 0,
				size:   0,
			},
		}
		l := a.logger.Info().Str("URI", r.RequestURI)
		next.ServeHTTP(&lw, r)
		l.Dur("duration", time.Since(start))
		l.Int("status", lw.responseData.status)
		l.Int("size", lw.responseData.size).Send()
	})
}
