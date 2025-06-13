package repository

import (
	"context"
	"crud_api/internal/domain/models"

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
