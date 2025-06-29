package repository

import (
	"context"
	"crud_api/internal/domain/models"

	appErrors "crud_api/internal/errors"

	"github.com/jackc/pgx/v5"
)

type MessageRepository struct {
	db *pgx.Conn
}

func NewMessageRepository(db *pgx.Conn) *MessageRepository {
	return &MessageRepository{
		db: db,
	}
}

func (r *MessageRepository) SaveMessage(ctx context.Context, msg *models.MessageEntity) error {
	query := `INSERT INTO messages(from_user_id, to_user_id, to_role, channel, message) VALUES($1, $2, $3, $4, $5)`

	_, err := r.db.Exec(ctx, query, msg.FromUserID, msg.ToUserID, msg.ToRole, msg.Channel, msg.Message)
	if err != nil {
		return appErrors.ErrDbFailure.Wrap(err, "Query failed")
	}

	return nil
}
