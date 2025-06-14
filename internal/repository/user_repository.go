package repository

import (
	"context"
	"crud_api/internal/domain/models"
	"log"

	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	db *pgx.Conn
}

func NewUserRepository(db *pgx.Conn) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, u *models.User) error {
	query := `INSERT INTO users (first_name, last_name, email, password, role) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return r.db.QueryRow(ctx, query, u.FirstName, u.LastName, u.Email, u.Password, u.Role).Scan(&u.ID)
}

func (r *UserRepository) GetAllUser(ctx context.Context) ([]*models.User, error) {
	query := `SELECT * FROM users`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var u models.User
		err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Password, &u.Role)
		if err != nil {
			return nil, err
		}
		users = append(users, &u)
	}

	return users, nil

}

func (r *UserRepository) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	var u models.User
	query := `SELECT * FROM users WHERE id = $1`
	err := r.db.QueryRow(ctx, query, id).Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Password, &u.Role)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, u *models.User) error {
	query := `UPDATE users SET first_name=$1, last_name=$2, email=$3, password=$4, role=$5 WHERE id=$6`
	_, err := r.db.Exec(ctx, query, u.FirstName, u.LastName, u.Email, u.Password, u.Role, u.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id=$1`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}

	query := `SELECT email, password FROM users where email=$1`
	err := r.db.QueryRow(ctx, query, email).Scan(&user.Email, &user.Password)

	if err != nil {
		log.Println("the requested email does not exist")
		return nil, err
	}

	return user, nil
}
