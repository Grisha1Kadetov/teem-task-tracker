package application

import (
	"context"
	"net/http"

	taskAdapter "github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/application/adapter/task"
	"github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/config"
	authHandler "github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/handler/auth"
	memberHandler "github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/handler/member"
	taskHandler "github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/handler/task"
	taskHistoryHandler "github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/handler/taskhistory"
	teamHandler "github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/handler/team"
	authMW "github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/middleware/auth"
	teamAccessMW "github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/middleware/teamaccess"
	"github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/model/role"
	"github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/pkg/db"
	"github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/pkg/log"
	memberDB "github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/repository/member"
	taskDB "github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/repository/task"
	taskHistoryDB "github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/repository/taskhistory"
	teamDB "github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/repository/team"
	userDB "github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/repository/user"
	authService "github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/service/auth"
	memberService "github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/service/member"
	taskService "github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/service/task"
	taskHistoryService "github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/service/taskhistory"
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

	memberRepo := memberDB.New(a.db)
	memberService := memberService.New(memberRepo)
	memberHandler := memberHandler.New(memberService, a.l)
	inviteAccess := teamAccessMW.New(
		memberService,
		teamAccessMW.Path("id"),
		role.Owner,
		role.Admin,
	)

	taskRepo := taskAdapter.New(taskDB.New(a.db))
	taskHistoryRepo := taskHistoryDB.New(a.db)
	taskHistoryService := taskHistoryService.New(taskHistoryRepo, taskRepo)
	taskHistoryHandler := taskHistoryHandler.New(taskHistoryService, a.l)
	taskService := taskService.New(taskRepo, memberService, taskHistoryService, a.db)
	taskHandler := taskHandler.New(taskAdapter.NewService(taskService), a.l)
	taskCreateAccess := teamAccessMW.New(
		memberService,
		teamAccessMW.Body("team_id"),
		role.Owner,
		role.Admin,
		role.Member,
	)
	taskListAccess := teamAccessMW.New(
		memberService,
		teamAccessMW.Query("team_id"),
		role.Owner,
		role.Admin,
		role.Member,
	)
	taskAccess := teamAccessMW.New(
		memberService,
		teamAccessMW.PathWithService("id", taskService.FindTaskTeam),
		role.Owner,
		role.Admin,
		role.Member,
	)

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)

		r.Group(func(r chi.Router) {
			r.Use(authRequired.Middleware())
			r.Post("/teams", teamHandler.Create)
			r.Get("/teams", memberHandler.List)
			r.With(inviteAccess.Middleware()).Post("/teams/{id}/invite", memberHandler.Invite)
			r.With(taskCreateAccess.Middleware()).Post("/tasks", taskHandler.Create)
			r.With(taskListAccess.Middleware()).Get("/tasks", taskHandler.List)
			r.With(taskAccess.Middleware()).Patch("/tasks/{id}", taskHandler.Update)
			r.With(taskAccess.Middleware()).Put("/tasks/{id}", taskHandler.Replace)
			r.With(taskAccess.Middleware()).Get("/tasks/{id}/history", taskHistoryHandler.List)
		})
	})

	return r
}
