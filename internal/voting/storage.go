package voting

import (
	"bytes"
	"context"
	"fmt"
	"os"

	storage_go "github.com/supabase-community/storage-go"
)

// signatureStorage handles Supabase Storage operations for digital signatures
type signatureStorage struct {
	client *storage_go.Client
	url    string
}

func newSignatureStorage() (*signatureStorage, error) {
	url := os.Getenv("SUPABASE_URL")
	key := os.Getenv("SUPABASE_SECRET_KEY")
	if url == "" || key == "" {
		return nil, fmt.Errorf("SUPABASE_URL and SUPABASE_SECRET_KEY required")
	}

	headers := map[string]string{
		"apikey": key,
	}
	client := storage_go.NewClient(url+"/storage/v1", key, headers)
	return &signatureStorage{client: client, url: url}, nil
}

func (s *signatureStorage) Upload(ctx context.Context, bucket, path string, data []byte, contentType string) (string, error) {
	reader := bytes.NewReader(data)
	_, err := s.client.UploadFile(bucket, path, reader)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/storage/v1/object/public/%s/%s", s.url, bucket, path), nil
}

func getSignatureBucket() string {
	bucket := os.Getenv("SUPABASE_MEDIA_BUCKET")
	if bucket == "" {
		return "pemira"
	}
	return bucket
}
