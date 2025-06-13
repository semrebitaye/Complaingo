package repository

import (
	"context"
	"crud_api/internal/domain/models"
	"log"

	"github.com/jackc/pgx/v5"
)

type PgxUserRepo struct {
	db *pgx.Conn
}

func NewPgxUserRepo(db *pgx.Conn) *PgxUserRepo {
	return &PgxUserRepo{db: db}
}

func (r *PgxUserRepo) CreateUser(ctx context.Context, u *models.User) error {
	query := `INSERT INTO users (first_name, last_name, email, password, role) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return r.db.QueryRow(ctx, query, u.FirstName, u.LastName, u.Email, u.Password, u.Role).Scan(&u.ID)
}

func (r *PgxUserRepo) GetAllUser(ctx context.Context) ([]*models.User, error) {
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

func (r *PgxUserRepo) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	var u models.User
	query := `SELECT * FROM users WHERE id = $1`
	err := r.db.QueryRow(ctx, query, id).Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Password, &u.Role)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *PgxUserRepo) UpdateUser(ctx context.Context, u *models.User) error {
	query := `UPDATE users SET first_name=$1, last_name=$2, email=$3, password=$4, role=$5 WHERE id=$6`
	_, err := r.db.Exec(ctx, query, u.FirstName, u.LastName, u.Email, u.Password, u.Role, u.ID)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
