package repository

import (
	"context"
	"Complaingo/internal/domain/models"
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

func (r *DocumentRepository) GetDocumentByUser(ctx context.Context, user_id int) ([]*models.Document, error) {
	query := `SELECT id, user_id, file_name, file_path, uploaded_at FROM documents WHERE user_id=$1`

	rows, err := r.db.Query(ctx, query, user_id)
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
