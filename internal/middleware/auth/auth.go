package auth

import (
	"net/http"
	"strings"

	"github.com/Grisha1Kadetov/TeemTaskTrackerService/internal/pkg/errorrenderer"
	"github.com/google/uuid"
)

type Checker struct {
	service service
}

func New(service service) *Checker {
	return &Checker{service: service}
}

func (c *Checker) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, ok := bearerToken(r.Header.Get("Authorization"))
			if !ok {
				errorrenderer.Render(w, http.StatusUnauthorized, errorrenderer.Unauthorized, "invalid authorization header")
				return
			}

			userID, err := c.service.ParseToken(r.Context(), token)
			if err != nil || userID == uuid.Nil {
				errorrenderer.Render(w, http.StatusUnauthorized, errorrenderer.Unauthorized, "invalid authorization token")
				return
			}

			ctx := WithActor(r.Context(), Actor{UserID: userID})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func bearerToken(header string) (string, bool) {
	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return "", false
	}

	token := strings.TrimSpace(strings.TrimPrefix(header, prefix))
	return token, token != ""
}
