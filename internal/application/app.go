package application

import (
	"context"
	"net/http"

	"github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/config"
	authHandler "github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/handler/auth"
	teamHandler "github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/handler/team"
	authMW "github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/middleware/auth"
	"github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/pkg/db"
	"github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/pkg/log"
	teamDB "github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/repository/team"
	userDB "github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/repository/user"
	authService "github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/service/auth"
	teamService "github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/service/team"
	"github.com/go-chi/chi/v5"
)

type App struct {
	ctx context.Context
	db  db.Database
	cnf config.Config
	l   log.Logger
}

func NewApplication(ctx context.Context, db db.Database, cnf config.Config, l log.Logger) *App {
	app := &App{
		ctx: ctx,
		db:  db,
		cnf: cnf,
		l:   l,
	}
	return app
}

func (a *App) NewRouter() http.Handler {
	r := chi.NewRouter()
	userRepo := userDB.New(a.db)
	authService := authService.New([]byte(a.cnf.JWTSecret), userRepo)
	authHandler := authHandler.New(authService, a.l)
	authRequired := authMW.New(authService)

	teamRepo := teamDB.New(a.db)
	teamService := teamService.New(teamRepo)
	teamHandler := teamHandler.New(teamService, a.l)

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)

		r.Group(func(r chi.Router) {
			r.Use(authRequired.Middleware())
			r.Post("/teams", teamHandler.Create)
		})
	})

	return r
}
