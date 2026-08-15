package member

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	authMW "github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/middleware/auth"
	teamAccessMW "github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/middleware/teamaccess"
	"github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/model/member"
	"github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/model/role"
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

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	actor, ok := authMW.ActorFromContext(r.Context())
	if !ok {
		errorrenderer.Render(w, http.StatusUnauthorized, errorrenderer.Unauthorized, "unauthorized")
		return
	}

	members, err := h.service.List(r.Context(), actor.UserID)
	if err != nil {
		h.logger.Error("failed to list user teams", log.Err(err))
		errorrenderer.Render(w, http.StatusInternalServerError, errorrenderer.Internal, "internal server error")
		return
	}

	response := make([]listResponse, 0, len(members))
	for _, teamMember := range members {
		response = append(response, listResponse{
			TeamID: teamMember.TeamID,
			Role:   teamMember.Role,
		})
	}

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, response)
}

func (h *Handler) Invite(w http.ResponseWriter, r *http.Request) {
	var request inviteRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		errorrenderer.Render(w, http.StatusBadRequest, errorrenderer.BadRequest, "invalid request body")
		return
	}

	request.Role = role.Role(strings.ToLower(strings.TrimSpace(string(request.Role))))

	access, ok := teamAccessMW.FromContext(r.Context())
	if !ok {
		errorrenderer.Render(w, http.StatusInternalServerError, errorrenderer.Internal, "internal server error")
		return
	}

	teamMember, err := h.service.Invite(r.Context(), access.TeamID, request.UserID, request.Role)
	if errors.Is(err, member.ErrUserNotFound) {
		errorrenderer.Render(w, http.StatusNotFound, errorrenderer.NotFound, "user not found")
		return
	}
	if errors.Is(err, member.ErrRoleMismatch) {
		errorrenderer.Render(w, http.StatusBadRequest, errorrenderer.BadRequest, "invalid role")
		return
	}
	if errors.Is(err, member.ErrAlreadyExists) {
		errorrenderer.Render(w, http.StatusConflict, errorrenderer.Conflict, "user is already a team member")
		return
	}
	if err != nil {
		h.logger.Error("failed to invite team member", log.Err(err))
		errorrenderer.Render(w, http.StatusInternalServerError, errorrenderer.Internal, "internal server error")
		return
	}

	response := inviteResponse{
		TeamID: teamMember.TeamID,
		UserID: teamMember.UserID,
		Role:   teamMember.Role,
	}
	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, response)
}
