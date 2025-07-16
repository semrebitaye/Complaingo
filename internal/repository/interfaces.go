package repository

import (
	"Complaingo/internal/domain/models"
	"context"
)

type MessageSaver interface {
	SaveMessage(ctx context.Context, msg *models.MessageEntity) error
}
