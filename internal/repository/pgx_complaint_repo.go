package repository

import (
	"Complaingo/internal/domain/models"
	appErrors "Complaingo/internal/errors"
	"Complaingo/internal/middleware"
	"Complaingo/internal/utility"
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/joomcode/errorx"
)

type PgxComplaintRepo struct {
	db *pgx.Conn
}

type PgxComplaintMessageRepo struct {
	db *pgx.Conn
}

func NewPgxComplaintRepo(db *pgx.Conn) *PgxComplaintRepo {
	return &PgxComplaintRepo{db: db}
}

func NewPgxComplaintMessageRepo(db *pgx.Conn) *PgxComplaintMessageRepo {
	return &PgxComplaintMessageRepo{
		db: db,
	}
}

func (r *PgxComplaintRepo) CreateComplaint(ctx context.Context, c *models.Complaints) error {
	if middleware.IsAdmin(ctx) {
		return appErrors.ErrInvalidPayload.New("users only have permission to create complaint")
	}

	query := `INSERT INTO complaints (user_id, subject, message, status) VALUES($1, $2, $3, $4) RETURNING id`

	err := r.db.QueryRow(ctx, query, c.UserID, c.Subject, c.Message, c.Status).Scan(&c.ID)
	if err != nil {
		return appErrors.ErrDbFailure.Wrap(err, "failed to query the user")
	}

	return nil
}

func (r *PgxComplaintRepo) GetComplaintByRole(ctx context.Context, UserID int, param utility.FilterParam) ([]*models.Complaints, error) {
	query := `SELECT * FROM complaints WHERE 1=1`

	args := []interface{}{UserID}
	query += " AND user_id=$1"
	argIdx := 2

	// add filters
	for _, f := range param.Filters {
		query += fmt.Sprintf(" AND %s %s $%d", f.ColumnName, f.Operator, argIdx)
		args = append(args, f.Value)
		argIdx++
	}

	// add search
	if param.Search != "" {
		query += fmt.Sprintf(" AND (user_id ILIKE $%d OR subject ILIKE $%d)", argIdx, argIdx+1)
		searchVal := "%" + param.Search + "%"
		args = append(args, searchVal, searchVal)
		argIdx += 2
	}

	// add sort
	sortCol := param.Sort.ColumnName
	sortOrder := param.Sort.Value

	if sortCol == "" {
		sortCol = "id"
	}
	if sortOrder == "" {
		sortOrder = "asc"
	}
	query += fmt.Sprintf(" ORDER BY %s %s", sortCol, sortOrder)

	//add pagination
	offset := (param.Page - 1) * param.PerPage
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, param.PerPage, offset)

	rows, err := r.db.Query(ctx, query, args)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, appErrors.ErrUserNotFound.Wrap(err, "Complaint not found")
		}
		return nil, appErrors.ErrDbFailure.New("quey failed")
	}
	defer rows.Close()

	var complaint []*models.Complaints
	for rows.Next() {
		var c models.Complaints
		err := rows.Scan(&c.ID, &c.UserID, &c.Subject, &c.Message, &c.Status, &c.CreatedAt)
		if err != nil {
			return nil, appErrors.ErrUserNotFound.New("failed to scan complaint row")
		}
		complaint = append(complaint, &c)
	}

	return complaint, nil
}

func (r *PgxComplaintRepo) GetComplaintByID(ctx context.Context, complaintID int) (*models.Complaints, error) {
	var c models.Complaints
	query := `SELECT * FROM complaints WHERE id=$1`
	err := r.db.QueryRow(ctx, query, complaintID).Scan(&c.ID, &c.UserID, &c.Subject, &c.Message, &c.Status, &c.CreatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, appErrors.ErrUserNotFound.New("complaint not found")
		}
		return nil, appErrors.ErrDbFailure.Wrap(err, "query failed")
	}

	return &c, nil
}

func (r *PgxComplaintRepo) UpdateComplaints(ctx context.Context, ComplaintId int, status string) error {
	query := `UPDATE complaints SET status=$1 WHERE id=$2`

	_, err := r.db.Exec(ctx, query, status, ComplaintId)
	if err != nil {
		return appErrors.ErrDbFailure.Wrap(err, "failed to update complaint")
	}

	return nil
}

func (r *PgxComplaintRepo) GetAllComplaintByRole(ctx context.Context, param utility.FilterParam) ([]*models.Complaints, error) {
	query := `SELECT id, user_id, subject, message, status, created_at FROM complaints WHERE 1=1` //to add AND conditions later.
	var args []interface{}
	argIdx := 1

	// Filters
	for _, f := range param.Filters {
		query += fmt.Sprintf(" AND %s %s $%d", f.ColumnName, f.Operator, argIdx)
		args = append(args, f.Value)
		argIdx++
	}

	// search in subject or message
	if param.Search != "" {
		query += fmt.Sprintf(" AND (subject ILIKE $%d OR message ILIKE $%d)", argIdx, argIdx+1)
		args = append(args, "%"+param.Search+"%", "%"+param.Search+"%")
		argIdx += 2
	}

	// sorting
	query += fmt.Sprintf(" ORDER BY %s %s", param.Sort.ColumnName, param.Sort.Value)

	// pagination
	offset := (param.Page - 1) * param.PerPage
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, param.PerPage, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, appErrors.ErrDbFailure.Wrap(err, "query failed")
	}
	defer rows.Close()

	var complaints []*models.Complaints
	for rows.Next() {
		var c models.Complaints
		err := rows.Scan(&c.ID, &c.UserID, &c.Subject, &c.Message, &c.Status, &c.CreatedAt)
		if err != nil {
			return nil, appErrors.ErrDbFailure.New("Failed to scan row")
		}
		complaints = append(complaints, &c)
	}

	return complaints, nil
}

// using complaintMessage table
func (r *PgxComplaintMessageRepo) InsertCoplaintMessage(ctx context.Context, cm *models.ComplaintMessages) error {
	query := `INSERT INTO complaint_messages (complaint_id, sender_id, parent_id, message, file_url) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := r.db.QueryRow(ctx, query, cm.ComplaintID, cm.SenderID, cm.ParentID, cm.Message, cm.FileUrl).Scan(&cm.ID)
	if err != nil {
		return appErrors.ErrDbFailure.Wrap(err, "query failed")
	}

	return nil
}

func (r *PgxComplaintMessageRepo) AddMessage(ctx context.Context, cm *models.ComplaintMessages) error {
	query := `INSERT INTO complaint_messages (complaint_id, sender_id, parent_id, message, file_url) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := r.db.QueryRow(ctx, query, cm.ComplaintID, cm.SenderID, cm.ParentID, cm.Message, cm.FileUrl).Scan(&cm.ID)
	if err != nil {
		log.Printf("DB Insert failed: %v", err)
		return appErrors.ErrDbFailure.Wrap(err, "query failed")
	}

	return nil
}

func (r *PgxComplaintMessageRepo) GetMessagesByComplaint(ctx context.Context, complaintID int) ([]*models.ComplaintMessages, error) {
	query := `SELECT * FROM complaint_messages WHERE complaint_id=$1`
	rows, err := r.db.Query(ctx, query, complaintID)
	if err != nil {
		return nil, appErrors.ErrDbFailure.Wrap(err, "query failed")
	}
	defer rows.Close()

	var complaint_messages []*models.ComplaintMessages
	for rows.Next() {
		var cm models.ComplaintMessages
		err := rows.Scan(&cm.ID, &cm.ComplaintID, &cm.SenderID, &cm.ParentID, &cm.Message, &cm.FileUrl, &cm.CreatedAt)
		if err != nil {
			return nil, appErrors.ErrDbFailure.Wrap(err, "Failed to scan row of messages")
		}
		complaint_messages = append(complaint_messages, &cm)
	}

	return complaint_messages, nil
}

func (r *PgxComplaintMessageRepo) GetMessageByID(ctx context.Context, messageID int) (*models.ComplaintMessages, error) {
	var cm models.ComplaintMessages

	query := `SELECT * FROM complaint_messages WHERE id=$1`
	err := r.db.QueryRow(ctx, query, messageID).Scan(&cm.ID, &cm.ComplaintID, &cm.SenderID, &cm.ParentID, &cm.Message, &cm.FileUrl, &cm.CreatedAt)
	if err != nil {
		if errorx.IsOfType(err, appErrors.ErrUserNotFound) {
			return nil, appErrors.ErrUnauthorized.Wrap(err, "Message not found by the given id")
		}
		return nil, appErrors.ErrDbFailure.Wrap(err, "query failed")
	}

	return &cm, nil
}

// notifier interface
