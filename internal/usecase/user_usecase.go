package usecase

import (
	"context"
	"crud_api/internal/domain/models"
	"crud_api/internal/repository"
	"crud_api/internal/utility"

	"github.com/joomcode/errorx"
	"golang.org/x/crypto/bcrypt"

	appErrors "crud_api/internal/errors"
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
		return appErrors.ErrDbFailure.New("usecase: Failed to generate password")
	}
	u.Password = string(hash)
	err = uc.repo.CreateUser(ctx, u)
	if err != nil {
		if errorx.IsOfType(err, appErrors.ErrUserDuplicate) {
			return err
		}
		return appErrors.ErrDbFailure.New("usecase: failed to register in usecase")
	}
	return nil
}

func (uc *UserUsecase) GetAllUser(ctx context.Context) ([]*models.User, error) {
	users, err := uc.repo.GetAllUser(ctx)
	if err != nil {
		return nil, appErrors.ErrDbFailure.Wrap(err, "usecase: failed to get all user on usecase")
	}
	return users, nil
}

func (uc *UserUsecase) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	user, err := uc.repo.GetUserByID(ctx, id)
	if err != nil {
		if errorx.IsOfType(err, appErrors.ErrUserNotFound) {
			return nil, appErrors.ErrUserNotFound.New("user not found in usecase")
		}
		return nil, appErrors.ErrDbFailure.Wrap(err, "usecase: unexpected db error")
	}

	return user, nil
}

func (uc *UserUsecase) UpdateUser(ctx context.Context, u *models.User) error {
	if err := uc.repo.UpdateUser(ctx, u); err != nil {
		return appErrors.ErrDbFailure.Wrap(err, "failed to update user in usecase")
	}
	return nil
}

func (uc *UserUsecase) DeleteUser(ctx context.Context, id int) error {
	if err := uc.repo.DeleteUser(ctx, id); err != nil {
		return appErrors.ErrDbFailure.Wrap(err, "usecase: failed to get user to be deleted")
	}
	return nil
}

func (uc *UserUsecase) Login(ctx context.Context, email string, password string) (string, error) {
	user, err := uc.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", appErrors.ErrUserNotFound.Wrap(err, "usecase: login failed, email not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", appErrors.ErrUnauthorized.New("Invalid credential")
	}

	token, err := utility.GenerateJWT(user.ID, user.Email, user.Role)
	if err != nil {
		return "", appErrors.ErrDbFailure.Wrap(err, "failed to generate jwt")
	}

	return token, nil
}
