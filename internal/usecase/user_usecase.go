package usecase

import (
	"context"
	"crud_api/internal/domain/models"
	"crud_api/internal/repository"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(r repository.UserRepository) *UserUsecase {
	return &UserUsecase{repo: r}
}

func (uc *UserUsecase) RegisterUser(ctx context.Context, u *models.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	u.Password = string(hash)
	return uc.repo.CreateUser(ctx, u)
}
