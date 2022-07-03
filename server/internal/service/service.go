package service

import (
	"context"
	"time"

	"github.com/Astemirdum/save/server/internal/repository"
	"github.com/Astemirdum/save/server/models"
)

type Service struct {
	FileContentService
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		FileContentService: NewFileSvc(repo),
	}
}

//go:generate go run github.com/golang/mock/mockgen -source=service.go -destination=mocks/mock.go
// FileContentService ...
type FileContentService interface {
	Create(ctx context.Context) error
	Append(ctx context.Context, req *models.WriteRequest) error
	Download(ctx context.Context) ([]byte, error)

	GetFileCount() uint32
	GetServerTime() time.Time
}
