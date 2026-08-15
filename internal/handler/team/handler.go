package team

import (
	"encoding/json"
	"net/http"
	"strings"
	"unicode/utf8"

	authMW "github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/middleware/auth"
	"github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/pkg/errorrenderer"
	"github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/pkg/log"
	"github.com/go-chi/render"
)

type Handler struct {
	service service
	logger  log.Logger
}

func New(service service, logger log.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var request createRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		errorrenderer.Render(w, http.StatusBadRequest, errorrenderer.BadRequest, "invalid request body")
		return
	}

	request.Name = strings.TrimSpace(request.Name)
	count := utf8.RuneCountInString(request.Name)
	if count < 1 || count > 255 {
		errorrenderer.Render(w, http.StatusBadRequest, errorrenderer.BadRequest, "incorrect team name length")
		return
	}

	actor, ok := authMW.ActorFromContext(r.Context())
	if !ok {
		errorrenderer.Render(w, http.StatusUnauthorized, errorrenderer.Unauthorized, "unauthorized")
		return
	}

	createdTeam, err := h.service.Create(r.Context(), request.Name, actor.UserID)
	if err != nil {
		h.logger.Error("failed to create team", log.Err(err))
		errorrenderer.Render(w, http.StatusInternalServerError, errorrenderer.Internal, "internal server error")
		return
	}

	response := createResponse{
		ID:        createdTeam.ID,
		Name:      createdTeam.Name,
		CreatedBy: createdTeam.CreatedBy,
		CreatedAt: createdTeam.CreatedAt,
	}
	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, response)
}
