package service

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/AndrewMislyuk/go-shop-backend/internal/transport/rest/domain"
	"github.com/AndrewMislyuk/go-shop-backend/internal/transport/rest/repository"
	"github.com/AndrewMislyuk/go-shop-backend/pkg/storage"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var folders = map[domain.FileType]string{
	domain.Image: "images",
}

type FileService struct {
	repo    repository.Files
	storage storage.Provider
}

func NewFileService(repo repository.Files, storage storage.Provider) *FileService {
	return &FileService{
		repo:    repo,
		storage: storage,
	}
}

func (f *FileService) Save(file domain.File) error {
	return f.repo.Create(file)
}

func (f *FileService) GetProductImage(productId string) (string, error) {
	return f.repo.GetProductImage(productId)
}

func (f *FileService) Upload(file domain.File) (string, error) {
	defer removeFile(file.Name)

	file.UploadStartedAt = time.Now()
	file.ID = uuid.New().String()

	c := make(chan string)
	go f.waitForUpload(file, c)
	url := <-c

	file.URL = url

	productImage, err := f.GetProductImage(file.ProductId)
	if err != nil {
		return "", err
	}

	if productImage != "" {
		filename := f.parseImageURL(productImage)

		err := f.storage.Delete(context.Background(), filename)
		if err != nil {
			return "", err
		}
	}

	if err := f.Save(file); err != nil {
		return "", err
	}

	return url, nil
}

func (f *FileService) waitForUpload(file domain.File, c chan string) {
	url, err := f.upload(file)
	if err != nil {
		logrus.Fatal(err)
	}

	c <- url
}

func (f *FileService) upload(file domain.File) (string, error) {
	fileData, err := os.Open(file.Name)
	if err != nil {
		return "", err
	}

	info, _ := fileData.Stat()
	logrus.Infof("file info: %+v", info)

	defer fileData.Close()

	return f.storage.Upload(context.Background(), storage.UploadInput{
		File:        fileData,
		Size:        file.Size,
		ContentType: file.ContentType,
		Name:        f.generateFilename(file),
	})
}

func (s *FileService) generateFilename(file domain.File) string {
	filename := fmt.Sprintf("%s.%s", uuid.New().String(), file.Name)
	folder := folders[file.Type]

	return fmt.Sprintf("%s/%s", folder, filename)
}

func (s *FileService) parseImageURL(url string) string {
	str := strings.Split(url, "/")

	return fmt.Sprintf("images/%s", str[len(str)-1])
}

func removeFile(filename string) {
	if err := os.Remove(filename); err != nil {
		logrus.Error("removeFile(): ", err)
	}
}
