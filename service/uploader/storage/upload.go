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
		return "", err
	}
	defer file.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		return "", err
	}

	fileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), fileHeader.Filename)
	fullPath := fmt.Sprintf("%s/%s", path, fileName)
	reader := bytes.NewReader(buf.Bytes())
	_, err = s.client.UploadFile(s.bucket, fullPath, reader)
	if err != nil {
		return "", err
	}

	// Return public URL
	publicURL := fmt.Sprintf("%s/storage/v1/object/public/%s/%s", s.client.GetPublicUrl(s.bucket, fullPath), s.bucket, fullPath)
	return publicURL, nil
}
