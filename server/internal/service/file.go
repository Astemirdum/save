package service

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"go.uber.org/atomic"

	"github.com/Astemirdum/save/enc"
	"github.com/Astemirdum/save/server/internal/repository"
	"github.com/Astemirdum/save/server/models"
)

// FileSvc ...
type FileSvc struct {
	repo        *repository.Repository
	fileCounter *atomic.Uint32

	fileName string
	mx       sync.Mutex

	rwmx sync.RWMutex
}

func NewFileSvc(repo *repository.Repository) *FileSvc {
	return &FileSvc{
		repo:        repo,
		fileCounter: atomic.NewUint32(0),
	}
}

// Create file
func (svc *FileSvc) Create(ctx context.Context) error {
	fileName := randStringRunes(rand.Intn(10) + 10)
	err := svc.repo.Create(ctx, &models.File{Filename: fileName})
	if err != nil {
		return fmt.Errorf("svc Create %w", err)
	}
	svc.setFilename(fileName)
	svc.fileCounter.Inc()
	return nil
}

func (svc *FileSvc) setFilename(fileName string) {
	svc.mx.Lock()
	svc.fileName = fileName
	svc.mx.Unlock()
}

// Append file content
func (svc *FileSvc) Append(ctx context.Context, req *models.WriteRequest) error {
	text, err := enc.Decrypt(req.Raw, req.Key)
	if err != nil {
		return fmt.Errorf("svc DecryptMessage %w", err)
	}
	row := fmt.Sprintf("%d-%s: %s", req.TimeStamp, req.ClientName, text)
	svc.rwmx.Lock()
	defer svc.rwmx.Unlock()
	if err = svc.repo.Append(ctx, &models.File{Filename: svc.fileName}, []byte(row)); err != nil {
		return fmt.Errorf("svc Append %w", err)
	}
	return nil
}

// Download file content
func (svc *FileSvc) Download(ctx context.Context) ([]byte, error) {
	svc.rwmx.RLock()
	defer svc.rwmx.RUnlock()
	fileContent, err := svc.repo.Download(ctx, &models.File{Filename: svc.fileName})
	if err != nil {
		return nil, fmt.Errorf("svc Download %w", err)
	}
	return fileContent, nil
}

// GetFileCount get file count
func (svc *FileSvc) GetFileCount() uint32 {
	return svc.fileCounter.Load()
}

// GetServerTime get server time
func (svc *FileSvc) GetServerTime() time.Time {
	return time.Now().UTC()
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
