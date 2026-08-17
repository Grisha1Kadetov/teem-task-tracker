package taskhistory

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/model/taskhistory"
	"github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/pkg/errorrenderer"
	"github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/pkg/log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
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
	taskID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil || taskID == uuid.Nil {
		errorrenderer.Render(w, http.StatusBadRequest, errorrenderer.BadRequest, "invalid task ID")
		return
	}

	history, err := h.service.ListTaskHistory(r.Context(), taskID)
	if err != nil {
		if errors.Is(err, taskhistory.ErrTaskNotFound) {
			errorrenderer.Render(w, http.StatusNotFound, errorrenderer.NotFound, "task not found")
			return
		}

		h.logger.Error("failed to list task history", log.Err(err))
		errorrenderer.Render(w, http.StatusInternalServerError, errorrenderer.Internal, "internal server error")
		return
	}

	response := make([]historyResponse, 0, len(history))
	for _, value := range history {
		response = append(response, historyResponse{
			ID:        value.ID,
			TaskID:    value.TaskID,
			ChangedBy: value.ChangedBy,
			Changes:   json.RawMessage(value.Changes),
			CreatedAt: value.CreatedAt,
		})
	}

	render.JSON(w, r, response)
}
