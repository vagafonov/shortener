package application

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

type Application struct {
	cnt *Container
}

func NewApplication(cnt *Container) *Application {
	return &Application{cnt: cnt}
}

func (a *Application) Serve() {
	fmt.Println(a.cnt.config.ServerURL)
	err := http.ListenAndServe(a.cnt.config.ServerURL, a.Routes())
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
		shortURL, err := NewService(a.cnt.GetStorage()).MakeShortURL(string(body), a.cnt.config.ShortURLLength)
		// TODO check error type
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
		res.WriteHeader(http.StatusCreated)

		if _, err := res.Write([]byte(a.cnt.config.ResultURL + "/" + shortURL.Short)); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (a *Application) getShortURL(res http.ResponseWriter, req *http.Request) {
	shortURL := NewService(a.cnt.GetStorage()).GetShortURL(chi.URLParam(req, "short_url"))

	if shortURL == nil {
		res.WriteHeader(http.StatusBadRequest)
	} else {
		res.Header().Set("Location", shortURL.Full)
		res.WriteHeader(http.StatusTemporaryRedirect)
	}
}
