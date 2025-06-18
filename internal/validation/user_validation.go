package validation

import (
	"crud_api/internal/domain/models"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func ValidateUser(u *models.User) error {
	return validation.ValidateStruct(u,
		validation.Field(&u.FirstName, validation.Required),
		validation.Field(&u.LastName, validation.Required),
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password, validation.Required, validation.Length(6, 100)),
		validation.Field(&u.Role, validation.Required, validation.In("admin", "user")),
	)
}

func ValidateId(id int) error {
	return validation.Validate(&id,
		validation.Required,
		validation.Min(1),
	)
}

func ValidateLoginInput(email, password string) error {
	return validation.ValidateStruct(&struct {
		Email    string
		Password string
	}{
		Email:    email,
		Password: password,
	},
		validation.Field(&email, validation.Required, is.Email),
		validation.Field(&password, validation.Required, validation.Length(6, 100)),
	)
}
