package task

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	authMW "github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/middleware/auth"
	teamAccessMW "github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/middleware/teamaccess"
	"github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/model/task"
	"github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/pkg/errorrenderer"
	"github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/pkg/log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

const defaultLimit uint64 = 20

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

	request.Title = strings.TrimSpace(request.Title)
	request.Description = strings.TrimSpace(request.Description)
	request.Status = task.Status(strings.ToLower(strings.TrimSpace(string(request.Status))))
	if count := utf8.RuneCountInString(request.Title); count < 1 || count > 255 {
		errorrenderer.Render(w, http.StatusBadRequest, errorrenderer.BadRequest, "incorrect task title length")
		return
	}
	if request.AssigneeID == uuid.Nil {
		errorrenderer.Render(w, http.StatusBadRequest, errorrenderer.BadRequest, "assignee_id is required")
		return
	}

	actor, ok := authMW.ActorFromContext(r.Context())
	if !ok {
		errorrenderer.Render(w, http.StatusUnauthorized, errorrenderer.Unauthorized, "unauthorized")
		return
	}
	access, ok := teamAccessMW.FromContext(r.Context())
	if !ok {
		errorrenderer.Render(w, http.StatusInternalServerError, errorrenderer.Internal, "internal server error")
		return
	}

	createdTask, err := h.service.CreateTask(
		r.Context(),
		access.TeamID,
		request.Title,
		request.Description,
		request.Status,
		request.AssigneeID,
		actor.UserID,
	)
	if err != nil {
		if status, code, message, ok := mapTaskError(err); ok {
			errorrenderer.Render(w, status, code, message)
		} else {
			h.logger.Error("failed to create task", log.Err(err))
			errorrenderer.Render(w, http.StatusInternalServerError, errorrenderer.Internal, "internal server error")
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, newTaskResponse(createdTask))
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	access, ok := teamAccessMW.FromContext(r.Context())
	if !ok {
		errorrenderer.Render(w, http.StatusInternalServerError, errorrenderer.Internal, "internal server error")
		return
	}

	filter := task.Filter{
		TeamID: access.TeamID,
		Limit:  defaultLimit,
		Offset: 0,
	}
	if rawStatus := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("status"))); rawStatus != "" {
		status := task.Status(rawStatus)
		filter.Status = &status
	}
	if rawAssigneeID := strings.TrimSpace(r.URL.Query().Get("assignee_id")); rawAssigneeID != "" {
		assigneeID, err := uuid.Parse(rawAssigneeID)
		if err != nil || assigneeID == uuid.Nil {
			errorrenderer.Render(w, http.StatusBadRequest, errorrenderer.BadRequest, "invalid assignee ID")
			return
		}
		filter.AssigneeID = &assigneeID
	}
	if rawLimit := strings.TrimSpace(r.URL.Query().Get("limit")); rawLimit != "" {
		limit, err := strconv.ParseUint(rawLimit, 10, 64)
		if err != nil || limit == 0 {
			errorrenderer.Render(w, http.StatusBadRequest, errorrenderer.BadRequest, "invalid limit")
			return
		}
		filter.Limit = limit
	}
	if rawOffset := strings.TrimSpace(r.URL.Query().Get("offset")); rawOffset != "" {
		offset, err := strconv.ParseUint(rawOffset, 10, 64)
		if err != nil {
			errorrenderer.Render(w, http.StatusBadRequest, errorrenderer.BadRequest, "invalid offset")
			return
		}
		filter.Offset = offset
	}

	tasks, err := h.service.ListTasks(r.Context(), filter)
	if err != nil {
		if status, code, message, ok := mapTaskError(err); ok {
			errorrenderer.Render(w, status, code, message)
		} else {
			h.logger.Error("failed to list tasks", log.Err(err))
			errorrenderer.Render(w, http.StatusInternalServerError, errorrenderer.Internal, "internal server error")
		}
		return
	}

	response := make([]taskResponse, 0, len(tasks))
	for _, value := range tasks {
		response = append(response, newTaskResponse(value))
	}
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, response)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	var request updateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		errorrenderer.Render(w, http.StatusBadRequest, errorrenderer.BadRequest, "invalid request body")
		return
	}

	h.update(w, r, request)
}

func (h *Handler) Replace(w http.ResponseWriter, r *http.Request) {
	var request replaceRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		errorrenderer.Render(w, http.StatusBadRequest, errorrenderer.BadRequest, "invalid request body")
		return
	}

	h.update(w, r, updateRequest{
		Title:       &request.Title,
		Description: &request.Description,
		Status:      &request.Status,
		AssigneeID:  &request.AssigneeID,
	})
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request, request updateRequest) {
	if request.Title != nil {
		title := strings.TrimSpace(*request.Title)
		if count := utf8.RuneCountInString(title); count < 1 || count > 255 {
			errorrenderer.Render(w, http.StatusBadRequest, errorrenderer.BadRequest, "incorrect task title length")
			return
		}
		request.Title = &title
	}
	if request.Description != nil {
		description := strings.TrimSpace(*request.Description)
		request.Description = &description
	}
	if request.Status != nil {
		status := task.Status(strings.ToLower(strings.TrimSpace(string(*request.Status))))
		request.Status = &status
	}
	if request.AssigneeID != nil && *request.AssigneeID == uuid.Nil {
		errorrenderer.Render(w, http.StatusBadRequest, errorrenderer.BadRequest, "invalid assignee ID")
		return
	}

	taskID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil || taskID == uuid.Nil {
		errorrenderer.Render(w, http.StatusBadRequest, errorrenderer.BadRequest, "invalid task ID")
		return
	}
	actor, ok := authMW.ActorFromContext(r.Context())
	if !ok {
		errorrenderer.Render(w, http.StatusUnauthorized, errorrenderer.Unauthorized, "unauthorized")
		return
	}

	err = h.service.UpdateTask(
		r.Context(),
		request.Title,
		request.Description,
		request.Status,
		request.AssigneeID,
		taskID,
		actor.UserID,
	)
	if err != nil {
		if status, code, message, ok := mapTaskError(err); ok {
			errorrenderer.Render(w, status, code, message)
		} else {
			h.logger.Error("failed to update task", log.Err(err))
			errorrenderer.Render(w, http.StatusInternalServerError, errorrenderer.Internal, "internal server error")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func newTaskResponse(value task.Task) taskResponse {
	return taskResponse{
		ID:          value.ID,
		TeamID:      value.TeamID,
		Title:       value.Title,
		Description: value.Description,
		Status:      value.Status,
		CreatedBy:   value.CreatedBy,
		AssigneeID:  value.AssigneeID,
		CreatedAt:   value.CreatedAt,
		UpdatedAt:   value.UpdatedAt,
		ClosedAt:    value.ClosedAt,
		Version:     value.Version,
	}
}

func mapTaskError(err error) (status int, code errorrenderer.Code, message string, ok bool) {
	if errors.Is(err, task.ErrInvalidStatus) {
		return http.StatusBadRequest, errorrenderer.BadRequest, "invalid task status", true
	}
	if errors.Is(err, task.ErrInvalidAssigneeID) {
		return http.StatusBadRequest, errorrenderer.BadRequest, "bad assignee_id", true
	}
	if errors.Is(err, task.ErrAssigneeNotTeamMember) {
		return http.StatusBadRequest, errorrenderer.BadRequest, "assignee is not a team member", true
	}
	if errors.Is(err, task.ErrEditorNotTeamMember) {
		return http.StatusForbidden, errorrenderer.Forbidden, "editor is not a team member", true
	}
	if errors.Is(err, task.ErrTaskNotFound) {
		return http.StatusNotFound, errorrenderer.NotFound, "task not found", true
	}
	if errors.Is(err, task.ErrNoChanges) {
		return http.StatusBadRequest, errorrenderer.BadRequest, "task has no changes", true
	}
	if errors.Is(err, task.ErrNoPermissionToUpdateTask) {
		return http.StatusForbidden, errorrenderer.Forbidden, "no permission to update task", true
	}
	if errors.Is(err, task.ErrInsufficientPermissions) {
		return http.StatusForbidden, errorrenderer.Forbidden, "insufficient permissions to update task", true
	}

	return 0, "", "", false
}
