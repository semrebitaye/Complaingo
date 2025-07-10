package repository

import (
	"Complaingo/internal/domain/models"
	"Complaingo/internal/utility"
	"context"
)

type UserRepository interface {
	CreateUser(ctx context.Context, u *models.User) error
	GetAllUser(ctx context.Context, param utility.FilterParam) ([]*models.User, error)
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	GetRoleByName(ctx context.Context, roleName string) (int, error)
	UpdateUser(ctx context.Context, u *models.User) error
	DeleteUser(ctx context.Context, id int) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}
