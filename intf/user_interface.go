package intf

import (
	"context"
	"crud_api/internal/domain/models"
)

type UserInterface interface {
	CreateUser(ctx context.Context, u *models.User) error
	GetAllUser(ctx context.Context) ([]*models.User, error)
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	UpdateUser(ctx context.Context, u *models.User) error
	DeleteUser(ctx context.Context, id int) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}
