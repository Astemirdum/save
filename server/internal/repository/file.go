package repository

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/Astemirdum/save/server/models"
)

// FileRepo ...
type FileRepo struct {
	filePath string
}

func NewFileRepo(filePath string) *FileRepo {
	return &FileRepo{filePath: filePath}
}

// Create file
func (repo *FileRepo) Create(_ context.Context, file *models.File) error {
	if err := os.MkdirAll(repo.filePath, os.ModePerm); err != nil && !os.IsExist(err) {
		return fmt.Errorf("dir failed %w", err)
	}
	f, err := os.Create(repo.filePath + "/" + file.Filename + ".txt")
	if err != nil {
		log.Fatal(err)
	}
	return f.Close()
}

// Append file
func (repo *FileRepo) Append(_ context.Context, file *models.File, fileBody []byte) error {
	if len(fileBody) == 0 {
		return errors.New("no file body provided to upload")
	}
	f, err := os.OpenFile(repo.filePath+"/"+file.Filename+".txt", os.O_APPEND|os.O_WRONLY, 0o666)
	if err != nil {
		return fmt.Errorf("openFile %v", err)
	}
	defer f.Close()
	if _, err = fmt.Fprintf(f, "%s", string(fileBody)); err != nil {
		return fmt.Errorf("fprintf %v", err)
	}
	return nil
}

// Download file
func (repo *FileRepo) Download(_ context.Context, file *models.File) ([]byte, error) {
	return ioutil.ReadFile(repo.filePath + "/" + file.Filename + ".txt")
}
