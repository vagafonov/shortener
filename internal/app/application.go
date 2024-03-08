package app

import (
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

const host = "http://localhost:8080/"

type Application struct {
	container *Container
}

func NewApplication(cnt *Container) *Application {
	return &Application{container: cnt}
}

func (a *Application) Serve() {
	err := http.ListenAndServe(`:8080`, a.Routes())
	if err != nil {
		panic(err)
	}
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

	if string(body) == "" {
		res.WriteHeader(http.StatusBadRequest)
	} else {
		shortURL, err := NewService(a.container.GetStorage()).MakeShortURL(string(body))
		// TODO check error type
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
		res.WriteHeader(http.StatusCreated)

		if _, err := res.Write([]byte(host + shortURL.Short)); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (a *Application) getShortURL(res http.ResponseWriter, req *http.Request) {
	shortURL := NewService(a.container.GetStorage()).GetShortURL(chi.URLParam(req, "short_url"))

	if shortURL == nil {
		res.WriteHeader(http.StatusBadRequest)
	} else {
		res.Header().Set("Location", shortURL.Full)
		res.WriteHeader(http.StatusTemporaryRedirect)
	}
}
