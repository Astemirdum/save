package repository

import (
	"context"

	"github.com/Astemirdum/save/server/models"
)

type Repository struct {
	FileContentRepo
}

func NewRepository(filePath string) *Repository {
	return &Repository{
		FileContentRepo: NewFileRepo(filePath),
	}
}

// FileContentRepo is a local store for file content
type FileContentRepo interface {
	Create(_ context.Context, file *models.File) error
	Append(context.Context, *models.File, []byte) error
	Download(context.Context, *models.File) ([]byte, error)
}
