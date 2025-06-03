package server

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/olegtemek/mini/internal/repository"
)

type App struct {
	server *http.Server
	repo   *repository.Repository
}

func New(repository *repository.Repository) *App {
	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	app := &App{
		server: server,
		repo:   repository,
	}

	mux.HandleFunc("/products", app.GetAll)
	mux.HandleFunc("/product", app.Create)

	return app
}

func (a *App) Start() (err error) {
	err = a.server.ListenAndServe()
	return
}

func (a *App) Stop(ctx context.Context) (err error) {
	err = a.server.Shutdown(ctx)
	return
}

func (a *App) GetAll(w http.ResponseWriter, r *http.Request) {
	products, err := a.repo.GetAll(r.Context())
	if err != nil {
		slog.Error("cannot get all products", "error", err)
		return
	}

	res, err := json.Marshal(products)
	if err != nil {
		slog.Error("cannot marshal products")
		return
	}

	w.Write(res)
}

func (a *App) Create(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	title := params.Get("title")
	if title == "" {
		slog.Error("cannot get title, cause title is null")
		return
	}

	ok, err := a.repo.Create(r.Context(), title)
	if err != nil {
		slog.Error("cannot create product", "error", err)
		return
	}

	res, err := json.Marshal(ok)
	if err != nil {
		slog.Error("cannot marshal status")
		return
	}

	w.Write(res)
}
