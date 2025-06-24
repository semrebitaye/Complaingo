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

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (l LoginInput) ValidateLoginInput() error {
	return validation.ValidateStruct(&l,
		validation.Field(&l.Email, validation.Required, is.Email),
		validation.Field(&l.Password, validation.Required, validation.Length(6, 100)),
	)
}
