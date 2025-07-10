package repository

import (
	"Complaingo/internal/domain/models"
	"Complaingo/internal/utility"
	"context"
	"fmt"
	"time"

	appErrors "Complaingo/internal/errors"

	"github.com/jackc/pgx/v5"
)

type DocumentRepository struct {
	db *pgx.Conn
}

func NewDocumentRepository(db *pgx.Conn) *DocumentRepository {
	return &DocumentRepository{
		db: db,
	}
}

func (r *DocumentRepository) SaveDocument(ctx context.Context, doc *models.Document) error {
	query := `INSERT INTO documents(user_id, file_name, file_path, uploaded_at) VALUES($1, $2, $3, $4) RETURNING id, uploaded_at`

	err := r.db.QueryRow(ctx, query, doc.UserID, doc.FileName, doc.FilePath, time.Now()).Scan(&doc.ID, &doc.UploadedAt)
	if err != nil {
		return appErrors.ErrDbFailure.Wrap(err, "qury failed")
	}

	return nil
}

func (r *DocumentRepository) GetDocumentByID(ctx context.Context, id int) (*models.Document, error) {
	query := `SELECT user_id, file_name, file_path, uploaded_at FROM documents WHERE id=$1`

	doc := &models.Document{}
	err := r.db.QueryRow(ctx, query, id).Scan(&doc.UserID, &doc.FileName, &doc.FilePath, &doc.UploadedAt)
	if err != nil {
		return nil, appErrors.ErrDbFailure.Wrap(err, "Query failed")
	}

	return doc, nil
}

func (r *DocumentRepository) GetDocumentByUser(ctx context.Context, user_id int, param utility.FilterParam) ([]*models.Document, error) {
	query := `SELECT id, user_id, file_name, file_path, uploaded_at FROM documents WHERE 1=1`

	args := []interface{}{user_id}
	query += " AND user_id=$1"
	argIdx := 2

	// add filters
	for _, f := range param.Filters {
		query += fmt.Sprintf(" AND %s %s $%d", f.ColumnName, f.Operator, argIdx)
		args = append(args, f.Value)
		argIdx++
	}

	// add search across user_id and file_name
	if param.Search != "" {
		query += fmt.Sprintf(" AND (user_id ILIKE $%d OR file_name ILIKE $%d)", argIdx, argIdx+1)
		searchVal := "%" + param.Search + "%"
		args = append(args, searchVal, searchVal)
	}

	// add sort
	sortCol := param.Sort.ColumnName
	sortOrder := param.Sort.Value

	if sortCol == "" {
		sortCol = "user_id"
	}
	if sortOrder == "" {
		sortOrder = "asc"
	}
	query += fmt.Sprintf(" ORDER BY %s %s", sortCol, sortOrder)

	// add pagination
	offset := (param.Page - 1) * param.PerPage
	query += fmt.Sprintf(" AND LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, param.PerPage, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, appErrors.ErrDbFailure.New("query failed")
	}
	defer rows.Close()

	var docs []*models.Document

	for rows.Next() {
		d := &models.Document{}
		err := rows.Scan(&d.ID, &d.UserID, &d.FileName, &d.FilePath, &d.UploadedAt)
		if err != nil {
			return nil, appErrors.ErrInvalidPayload.Wrap(err, "Invalid payload")
		}

		docs = append(docs, d)
	}
	return docs, nil
}

func (r *DocumentRepository) DeleteDocument(ctx context.Context, id int) error {
	query := `DELETE FROM documents WHERE id=$1`

	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return appErrors.ErrDbFailure.Wrap(err, "Filed to delete")
	}

	return nil
}
