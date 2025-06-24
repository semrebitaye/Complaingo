package models

type Role string

const (
	AdminRole Role = "admin"
	UerRole   Role = "user"
)

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Role      string `json:"role"`
	RoleID    int    `json:"role_id"`
}
