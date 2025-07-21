package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"time"

	storage "github.com/supabase-community/storage-go"
)

type ImageUploader interface {
	UploadImage(ctx context.Context, fileHeader *multipart.FileHeader, path string) (string, error)
	UploadImages(ctx context.Context, files []*multipart.FileHeader, path string) ([]string, error)
}

type SupabaseUploader struct {
	client *storage.Client
	bucket string
}

func NewSupabaseUploader(url, key, bucket string) (*SupabaseUploader, error) {
	client := storage.NewClient(url, key, nil)
	return &SupabaseUploader{
		client: client,
		bucket: bucket,
	}, nil
}

func (s *SupabaseUploader) UploadImage(ctx context.Context, fileHeader *multipart.FileHeader, path string) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		return "", fmt.Errorf("failed to copy file to buffer: %w", err)
	}

	fileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), fileHeader.Filename)
	fullPath := fmt.Sprintf("%s/%s", path, fileName)
	reader := bytes.NewReader(buf.Bytes())

	if _, err = s.client.UploadFile(s.bucket, fullPath, reader); err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	publicURL := fmt.Sprintf("%s/storage/v1/object/public/%s/%s", s.client.GetPublicUrl(s.bucket, fullPath), s.bucket, fullPath)
	return publicURL, nil
}

func (s *SupabaseUploader) UploadImages(ctx context.Context, files []*multipart.FileHeader, path string) ([]string, error) {
	urls := make([]string, 0, len(files))
	for _, file := range files {
		url, err := s.UploadImage(ctx, file, path)
		if err != nil {
			return nil, fmt.Errorf("failed to upload image %s: %w", file.Filename, err)
		}
		urls = append(urls, url)
	}
	return urls, nil
}
