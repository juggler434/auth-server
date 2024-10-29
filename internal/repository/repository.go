package repository

import (
	"context"

	"github.com/juggler434/auth-server/pkg/models"
)

type Repository interface {
	CreateUser(ctx context.Context, user *models.User) error
	Close(ctx context.Context) error
}
