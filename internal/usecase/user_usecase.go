package usecase

import (
	"context"
	"crud_api/internal/domain/models"
	"crud_api/intf"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	repo intf.UserInterface
}

func NewUserUsecase(r intf.UserInterface) *UserUsecase {
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

func (uc *UserUsecase) GetAllUser(ctx context.Context) ([]*models.User, error) {
	return uc.repo.GetAllUser(ctx)
}

func (uc *UserUsecase) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	return uc.repo.GetUserByID(ctx, id)
}

func (uc *UserUsecase) UpdateUser(ctx context.Context, u *models.User) error {
	return uc.repo.UpdateUser(ctx, u)
}

func (uc *UserUsecase) DeleteUser(ctx context.Context, id int) error {
	return uc.repo.DeleteUser(ctx, id)
}
