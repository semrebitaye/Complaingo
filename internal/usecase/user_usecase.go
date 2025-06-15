package usecase

import (
	"context"
	"crud_api/internal/domain/models"
	"crud_api/internal/repository"
	"crud_api/internal/utility"

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
		return err
	}
	u.Password = string(hash)
	if err = uc.repo.CreateUser(ctx, u); err != nil {
		return err
	}
	return nil
}

func (uc *UserUsecase) GetAllUser(ctx context.Context) ([]*models.User, error) {
	users, err := uc.repo.GetAllUser(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (uc *UserUsecase) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	user, err := uc.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *UserUsecase) UpdateUser(ctx context.Context, u *models.User) error {
	if err := uc.repo.UpdateUser(ctx, u); err != nil {
		return err
	}
	return nil
}

func (uc *UserUsecase) DeleteUser(ctx context.Context, id int) error {
	if err := uc.repo.DeleteUser(ctx, id); err != nil {
		return err
	}
	return nil
}

func (uc *UserUsecase) Login(ctx context.Context, email string, password string) (string, error) {
	user, err := uc.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", err
	}

	token, err := utility.GenerateJWT(user.ID, user.Email, user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}
