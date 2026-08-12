package application

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type App struct {
	ctx context.Context
}

func NewApplication(ctx context.Context) *App {
	app := &App{}
	app.ctx = ctx

	return app
}

func (a *App) NewRouter() http.Handler {
	r := chi.NewRouter()

	return r
}
