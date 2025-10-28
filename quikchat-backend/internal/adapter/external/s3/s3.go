package s3

import (
	"context"
	"time"
)

// S3Client defines the interface for interacting with an S3-compatible service.
type S3Client interface {
	GeneratePresignedUploadURL(ctx context.Context, bucket, key string, expiry time.Duration) (string, error)
	GeneratePresignedDownloadURL(ctx context.Context, bucket, key string, expiry time.Duration) (string, error)
}

// MockS3Client provides a placeholder implementation for development.
type MockS3Client struct {
	// In a real scenario, this would hold the S3 client from the AWS SDK.
}

// NewMockS3Client creates a new mock S3 client.
func NewMockS3Client() S3Client {
	return &MockS3Client{}
}

// GeneratePresignedUploadURL returns a fake URL for uploading.
func (c *MockS3Client) GeneratePresignedUploadURL(ctx context.Context, bucket, key string, expiry time.Duration) (string, error) {
	// This is a placeholder. In a real implementation, you would generate a real presigned URL.
	return "https://s3.example.com/" + bucket + "/" + key + "?upload-signature=mock", nil
}

// GeneratePresignedDownloadURL returns a fake URL for downloading.
func (c *MockS3Client) GeneratePresignedDownloadURL(ctx context.Context, bucket, key string, expiry time.Duration) (string, error) {
	// This is a placeholder.
	return "https://s3.example.com/" + bucket + "/" + key + "?download-signature=mock", nil
}

