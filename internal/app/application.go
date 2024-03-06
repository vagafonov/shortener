package app

import (
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
	a.Routes()
	err := http.ListenAndServe(`:8080`, a.mux)
	if err != nil {
		panic(err)
	}
}

func (a *Application) Routes() {
	a.mux.HandleFunc(`/`, a.Route)
}

func (a *Application) Route(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		a.createShortURL(res, req)
	case http.MethodGet:
		a.getShortURL(res, req)
	default:
		http.Error(res, "Only POST or GET requests are allowed!", http.StatusMethodNotAllowed)
	}
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
