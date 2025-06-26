package repository

import (
	"context"
	"crud_api/internal/domain/models"
	appErrors "crud_api/internal/errors"

	"github.com/jackc/pgx/v5"
)

type PgxUserRepo struct {
	db *pgx.Conn
}

func NewPgxUserRepo(db *pgx.Conn) *PgxUserRepo {
	return &PgxUserRepo{db: db}
}

func (r *PgxUserRepo) CreateUser(ctx context.Context, u *models.User) error {
	var roleID int
	err := r.db.QueryRow(ctx, `SELECT id FROM roles WHERE name=$1`, u.Role).Scan(&roleID)
	if err != nil {
		return appErrors.ErrInvalidPayload.New("Invalid role")
	}

	query := `
	INSERT INTO users (first_name, last_name, email, password, role_id)
	VALUES ($1, $2, $3, $4, (SELECT id FROM roles WHERE name = $5))
	RETURNING id`
	err = r.db.QueryRow(ctx, query, u.FirstName, u.LastName, u.Email, u.Password, u.Role).Scan(&u.ID)

	if err != nil {
		return appErrors.ErrDbFailure.Wrap(err, "failed to query the user")
	}

	u.RoleID = roleID

	return nil
}

func (r *PgxUserRepo) GetAllUser(ctx context.Context) ([]*models.User, error) {
	query := `
	SELECT u.id, u.first_name, u.last_name, u.email, u.password, r.name as role, u.role_id
	FROM users u
	LEFT JOIN roles r ON u.role_id = r.id
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, appErrors.ErrDbFailure.Wrap(err, "query failed")
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var u models.User
		err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Password, &u.Role, &u.RoleID)
		if err != nil {
			return nil, appErrors.ErrUserNotFound.New("failed to scan user row")
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
		if err == pgx.ErrNoRows {
			return nil, appErrors.ErrUserNotFound.New("user not found the required id")
		}
		return nil, appErrors.ErrDbFailure.New("query failed")
	}
	return &u, nil
}

func (r *PgxUserRepo) GetRoleByName(ctx context.Context, roleName string) (int, error) {
	var id int

	query := `SELECT id FROM roles WHERE name=$1`
	err := r.db.QueryRow(ctx, query, roleName).Scan(&id)
	if err != nil {
		return 0, appErrors.ErrInvalidPayload.Wrap(err, "Invalid Role Name")
	}
	return id, nil
}

func (r *PgxUserRepo) UpdateUser(ctx context.Context, u *models.User) error {
	query := `UPDATE users SET first_name=$1, last_name=$2, email=$3, password=$4, role=$5 WHERE id=$6`
	_, err := r.db.Exec(ctx, query, u.FirstName, u.LastName, u.Email, u.Password, u.Role, u.ID)
	if err != nil {
		return appErrors.ErrDbFailure.New("unable to update user")
	}
	return nil
}

func (r *PgxUserRepo) DeleteUser(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id=$1`
	res, err := r.db.Exec(ctx, query, id)
	if err != nil {
		if res.RowsAffected() == 0 {
			return appErrors.ErrUserNotFound.New("no user found to delete")
		}
		return appErrors.ErrDbFailure.Wrap(err, "failed to excute delete query")
	}
	return nil
}

func (r *PgxUserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}

	query := `SELECT u.id, u.email, u.password, r.name as role
	FROM users u
	LEFT JOIN roles r on u.role_id = r.id
	WHERE u.email=$1`
	err := r.db.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.Password, &user.Role)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, appErrors.ErrUserNotFound.New("user not found with email")
		}
		return nil, appErrors.ErrDbFailure.New("query failed")
	}

	return user, nil
}
