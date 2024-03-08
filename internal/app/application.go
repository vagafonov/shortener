package app

import (
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"strings"
)

const host = "http://localhost:8080/"

type Application struct {
	container *Container
	mux       *http.ServeMux
}

func NewApplication(cnt *Container) *Application {
	return &Application{container: cnt}
}

func (a *Application) Serve() {
	a.mux = http.NewServeMux()
	routes := a.Routes()
	err := http.ListenAndServe(`:8080`, routes)
	if err != nil {
		panic(err)
	}
}

func (a *Application) Routes() *chi.Mux {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Get("/", a.getShortURL)
		r.Post("/", a.createShortURL)
	})
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
		shortURL := NewService(a.container.GetStorage()).MakeShortURL(string(body))
		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(host + shortURL))
	}
}

func (a *Application) getShortURL(res http.ResponseWriter, req *http.Request) {
	shortURL := NewService(a.container.GetStorage()).GetShortURL(strings.Trim(req.URL.String(), "/"))
	if shortURL == "" {
		res.WriteHeader(http.StatusBadRequest)
	} else {
		res.Header().Set("Location", shortURL)
		res.WriteHeader(http.StatusTemporaryRedirect)
	}
}
