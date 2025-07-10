package usecase

import (
	"Complaingo/internal/domain/models"
	"Complaingo/internal/repository"
	"Complaingo/internal/utility"
	"context"
)

type DocumentUsecase struct {
	repo *repository.DocumentRepository
}

func NewDocumentUsecase(repo *repository.DocumentRepository) *DocumentUsecase {
	return &DocumentUsecase{
		repo: repo,
	}
}

func (du *DocumentUsecase) Uplod(ctx context.Context, doc *models.Document) error {
	return du.repo.SaveDocument(ctx, doc)
}

func (du *DocumentUsecase) GetDocumentByID(ctx context.Context, id int) (*models.Document, error) {
	return du.repo.GetDocumentByID(ctx, id)
}

func (du *DocumentUsecase) GetDocumentByUser(ctx context.Context, user_id int, param utility.FilterParam) ([]*models.Document, error) {
	return du.repo.GetDocumentByUser(ctx, user_id, param)
}

func (du *DocumentUsecase) DeleteDocument(ctx context.Context, id int) error {
	return du.repo.DeleteDocument(ctx, id)
}
