package teamaccess

import (
	"errors"
	"net/http"

	"github.com/Grisha1Kadetov/TeemTaskTrackerService/internal/middleware/auth"
	"github.com/Grisha1Kadetov/TeemTaskTrackerService/internal/model/role"
	"github.com/Grisha1Kadetov/TeemTaskTrackerService/internal/pkg/errorrenderer"
)

type Checker struct {
	service service
	resolve Resolver
	allowed map[role.Role]struct{}
}

func New(service service, resolve Resolver, roles ...role.Role) *Checker {
	allowed := make(map[role.Role]struct{}, len(roles))
	for _, value := range roles {
		allowed[value] = struct{}{}
	}

	return &Checker{
		service: service,
		resolve: resolve,
		allowed: allowed,
	}
}

func (c *Checker) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actor, ok := auth.ActorFromContext(r.Context())
			if !ok {
				errorrenderer.Render(w, http.StatusUnauthorized, errorrenderer.Unauthorized, "authentication required")
				return
			}

			teamID, err := c.resolve(r)
			if err != nil {
				renderResolveError(w, err)
				return
			}

			teamRole, found, err := c.service.Role(r.Context(), actor.UserID, teamID)
			if err != nil {
				errorrenderer.Render(w, http.StatusInternalServerError, errorrenderer.Internal, "internal server error")
				return
			}
			if !found {
				errorrenderer.Render(w, http.StatusForbidden, errorrenderer.Forbidden, "team access denied")
				return
			}
			if !teamRole.IsValid() {
				errorrenderer.Render(w, http.StatusInternalServerError, errorrenderer.Internal, "internal server error")
				return
			}
			if _, ok := c.allowed[teamRole]; !ok {
				errorrenderer.Render(w, http.StatusForbidden, errorrenderer.Forbidden, "team access denied")
				return
			}

			ctx := withAccess(r.Context(), Access{TeamID: teamID, Role: teamRole})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func renderResolveError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrInvalidID):
		errorrenderer.Render(w, http.StatusBadRequest, errorrenderer.BadRequest, "invalid resource ID")
	case errors.Is(err, ErrNotFound):
		errorrenderer.Render(w, http.StatusNotFound, errorrenderer.NotFound, "resource not found")
	default:
		errorrenderer.Render(w, http.StatusInternalServerError, errorrenderer.Internal, "internal server error")
	}
}
