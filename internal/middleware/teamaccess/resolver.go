package teamaccess

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidID   = errors.New("invalid resource ID")
	ErrInvalidBody = errors.New("invalid request body")
	ErrNotFound    = errors.New("resource not found")
)

type Resolver func(r *http.Request) (uuid.UUID, error)

func Path(param string) Resolver {
	return func(r *http.Request) (uuid.UUID, error) {
		return parseID(chi.URLParam(r, param))
	}
}

func Query(param string) Resolver {
	return func(r *http.Request) (uuid.UUID, error) {
		return parseID(r.URL.Query().Get(param))
	}
}

func Body(param string) Resolver {
	return func(r *http.Request) (uuid.UUID, error) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return uuid.Nil, ErrInvalidBody
		}
		r.Body = io.NopCloser(bytes.NewReader(body))

		var values map[string]json.RawMessage
		if err := json.Unmarshal(body, &values); err != nil {
			return uuid.Nil, ErrInvalidBody
		}

		var value string
		if err := json.Unmarshal(values[param], &value); err != nil {
			return uuid.Nil, ErrInvalidID
		}

		return parseID(value)
	}
}

func PathWithService(param string, findTeam FindTeamFunc) Resolver {
	return func(r *http.Request) (uuid.UUID, error) {
		resourceID, err := parseID(chi.URLParam(r, param))
		if err != nil {
			return uuid.Nil, err
		}

		teamID, found, err := findTeam(r.Context(), resourceID)
		if err != nil {
			return uuid.Nil, fmt.Errorf("find resource team: %w", err)
		}
		if !found {
			return uuid.Nil, ErrNotFound
		}
		if teamID == uuid.Nil {
			return uuid.Nil, errors.New("resource has empty team ID")
		}

		return teamID, nil
	}
}

func parseID(value string) (uuid.UUID, error) {
	id, err := uuid.Parse(value)
	if err != nil || id == uuid.Nil {
		return uuid.Nil, ErrInvalidID
	}
	return id, nil
}
