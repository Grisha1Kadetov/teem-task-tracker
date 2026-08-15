package errorrenderer

import (
	"encoding/json"
	"net/http"
)

type Code string

const (
	BadRequest   Code = "BAD_REQUEST"
	Conflict     Code = "CONFLICT"
	Unauthorized Code = "UNAUTHORIZED"
	Forbidden    Code = "FORBIDDEN"
	NotFound     Code = "NOT_FOUND"
	Internal     Code = "INTERNAL_ERROR"
)

type response struct {
	Error errorBody `json:"error"`
}

type errorBody struct {
	Code    Code   `json:"code"`
	Message string `json:"message"`
}

func Render(w http.ResponseWriter, status int, code Code, message string) {
	body, err := json.Marshal(response{
		Error: errorBody{
			Code:    code,
			Message: message,
		},
	})
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(append(body, '\n'))
}
