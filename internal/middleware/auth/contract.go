package auth

import (
	"context"

	"github.com/google/uuid"
)

type service interface {
	ParseToken(ctx context.Context, token string) (uuid.UUID, error)
}
