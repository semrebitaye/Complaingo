package repository

import (
	"context"
	"crud_api/internal/domain/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, u *models.User) error
	GetAllUser(ctx context.Context) ([]*models.User, error)
}
