package usecase

import (
	"context"
	"crud_api/internal/domain/models"
	"crud_api/internal/repository"
	"crud_api/internal/utility"
	"crud_api/internal/validation"

	"github.com/joomcode/errorx"

	appErrors "crud_api/internal/errors"
)

type UserUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(r repository.UserRepository) *UserUsecase {
	return &UserUsecase{repo: r}
}

type LoginResponse struct {
	Token string
	User  *models.User
}

func (uc *UserUsecase) RegisterUser(ctx context.Context, u *models.User) error {
	// validate user input
	if err := validation.ValidateUser(u); err != nil {
		return appErrors.ErrInvalidPayload.Wrap(err, "usecase: validation failed")
	}

	// get role_id from role name
	roleID, err := uc.repo.GetRoleByName(ctx, u.Role)
	if err != nil {
		return appErrors.ErrInvalidPayload.Wrap(err, "usecase: role not found")
	}
	u.RoleID = roleID

	// hash the password
	hashed, err := utility.HashPassword(u.Password)
	if err != nil {
		return appErrors.ErrInvalidPayload.Wrap(err, "password hashing failed")
	}
	u.Password = hashed

	// create user in db
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
	// check validation
	if err := validation.ValidateId(id); err != nil {
		return nil, appErrors.ErrInvalidPayload.New("usecase: Invalid user id")
	}

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
	if err := validation.ValidateUser(u); err != nil {
		return appErrors.ErrDbFailure.Wrap(err, "usecase: validation failed")
	}

	if err := uc.repo.UpdateUser(ctx, u); err != nil {
		return appErrors.ErrDbFailure.Wrap(err, "failed to update user in usecase")
	}
	return nil
}

func (uc *UserUsecase) DeleteUser(ctx context.Context, id int) error {
	if err := validation.ValidateId(id); err != nil {
		return appErrors.ErrInvalidPayload.New("validation of id failed")
	}

	if err := uc.repo.DeleteUser(ctx, id); err != nil {
		return appErrors.ErrDbFailure.Wrap(err, "usecase: failed to get user to be deleted")
	}
	return nil
}

func (uc *UserUsecase) Login(ctx context.Context, email string, password string) (*LoginResponse, error) {
	// validate email and password input
	input := validation.LoginInput{
		Email:    email,
		Password: password,
	}

	if err := input.ValidateLoginInput(); err != nil {
		return nil, appErrors.ErrInvalidPayload.Wrap(err, "usecase: Invalid login input")
	}

	// fetch user by email
	user, err := uc.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, appErrors.ErrUserNotFound.Wrap(err, "usecase: login failed, email not found")
	}

	// compare hashed password with the input password
	err = utility.ComparePassword(user.Password, password)
	if err != nil {
		return nil, appErrors.ErrUnauthorized.New("Invalid credential")
	}

	// generate jwt token
	token, err := utility.GenerateJWT(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, appErrors.ErrDbFailure.Wrap(err, "failed to generate jwt")
	}

	return &LoginResponse{
		Token: token,
		User:  user,
	}, nil
}
