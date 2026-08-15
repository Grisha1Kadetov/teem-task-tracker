package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/mail"
	"strings"
	"unicode/utf8"

	"github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/model/user"
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

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var request registerRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		errorrenderer.Render(w, http.StatusBadRequest, errorrenderer.BadRequest, "invalid request body")
		return
	}
	request.Email = strings.ToLower(strings.TrimSpace(request.Email))
	request.Name = strings.ToLower(strings.TrimSpace(request.Name))
	request.Password = strings.TrimSpace(request.Password)

	if !validEmail(request.Email) {
		errorrenderer.Render(w, http.StatusBadRequest, errorrenderer.BadRequest, "incorrect email")
		return
	}

	count := utf8.RuneCountInString(request.Name)
	if count < 1 || count > 30 {
		errorrenderer.Render(w, http.StatusBadRequest, errorrenderer.BadRequest, "incorrect name lenght")
		return
	}

	count = utf8.RuneCountInString(request.Password)
	if count < 8 || count > 30 {
		errorrenderer.Render(w, http.StatusBadRequest, errorrenderer.BadRequest, "incorrect password lenght")
		return
	}

	registeredUser, token, err := h.service.Register(
		r.Context(),
		request.Email,
		request.Password,
		request.Name,
	)

	if errors.Is(err, user.ErrEmailAlreadyExists) {
		errorrenderer.Render(w, http.StatusConflict, errorrenderer.Conflict, "email already exists")
		return
	}
	if err != nil {
		h.logger.Error("failed to register user", log.Err(err))
		errorrenderer.Render(w, http.StatusInternalServerError, errorrenderer.Internal, "internal server error")
		return
	}

	response := registerResponse{
		ID:          registeredUser.ID,
		Email:       registeredUser.Email,
		Name:        registeredUser.Name,
		AccessToken: token,
	}
	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, response)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var request loginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		errorrenderer.Render(w, http.StatusBadRequest, errorrenderer.BadRequest, "invalid request body")
		return
	}

	if !validEmail(request.Email) {
		errorrenderer.Render(w, http.StatusBadRequest, errorrenderer.BadRequest, "incorrect email")
		return
	}

	token, err := h.service.Login(r.Context(), request.Email, request.Password)
	if errors.Is(err, user.ErrUserNotFound) || errors.Is(err, user.ErrIncorrectPassword) {
		errorrenderer.Render(w, http.StatusUnauthorized, errorrenderer.Unauthorized, "invalid email or password")
		return
	}
	if err != nil {
		h.logger.Error("failed to log in user", log.Err(err))
		errorrenderer.Render(w, http.StatusInternalServerError, errorrenderer.Internal, "internal server error")
		return
	}

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, loginResponse{AccessToken: token})
}

func validEmail(value string) bool {
	value = strings.ToLower(value)
	address, err := mail.ParseAddress(value)
	return err == nil && address.Address == value
}
