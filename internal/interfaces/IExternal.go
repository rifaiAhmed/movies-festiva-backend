package interfaces

import (
	"context"
	"movie-festival/internal/models"
)

type IExternal interface {
	ValidateToken(ctx context.Context, token string) (models.TokenData, error)
}
