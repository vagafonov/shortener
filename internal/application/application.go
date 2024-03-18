package application

import (
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vagafonov/shortener/internal/application/services"
)

type Application struct {
	cnt *Container
}

func NewApplication(cnt *Container) *Application {
	return &Application{cnt: cnt}
}

func (a *Application) Serve() error {
	err := http.ListenAndServe(a.cnt.cfg.ServerURL, a.Routes()) //nolint:gosec
	if err != nil {
		return err
	}

	return nil
}

func (a *Application) Routes() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/{short_url}", a.getShortURL)
	r.Post("/", a.createShortURL)

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

func (a *Application) getShortURL(res http.ResponseWriter, req *http.Request) {
	shortURL := services.NewService(a.cnt.GetStorage()).GetShortURL(chi.URLParam(req, "short_url"))
	if shortURL == nil {
		res.WriteHeader(http.StatusNotFound)

		return
	}
	res.Header().Set("Location", shortURL.Full)
	res.WriteHeader(http.StatusTemporaryRedirect)
}
