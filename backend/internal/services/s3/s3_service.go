// Package s3 provides file storage services with support for both
// AWS S3 and local filesystem storage.
package s3

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/research-paper-analyzer/backend/internal/config"
)

// Service handles file storage operations.
// It supports both AWS S3 and local filesystem storage.
type Service struct {
	cfg      *config.Config
	s3Client *s3.Client
}

// NewService creates a new S3/storage service instance.
// If the storage provider is "s3", it initializes the AWS S3 client.
// If "local", it ensures the upload directory exists.
func NewService(cfg *config.Config) (*Service, error) {
	svc := &Service{
		cfg: cfg,
	}

	if cfg.StorageProvider == "s3" {
		// Initialize AWS S3 client
		var awsOpts []func(*awsconfig.LoadOptions) error

		// Set the AWS region
		awsOpts = append(awsOpts, awsconfig.WithRegion(cfg.AWSRegion))

		// Use explicit credentials if provided, otherwise fall back to IAM roles
		if cfg.AWSAccessKey != "" && cfg.AWSSecretKey != "" {
			awsOpts = append(awsOpts, awsconfig.WithCredentialsProvider(
				credentials.NewStaticCredentialsProvider(cfg.AWSAccessKey, cfg.AWSSecretKey, ""),
			))
		}

		awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(), awsOpts...)
		if err != nil {
			return nil, fmt.Errorf("failed to load AWS config: %w", err)
		}

		svc.s3Client = s3.NewFromConfig(awsCfg)
		log.Println("✅ S3 storage initialized")
	} else {
		// Ensure local upload directory exists
		if err := os.MkdirAll(cfg.UploadDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create upload directory: %w", err)
		}
		log.Printf("✅ Local storage initialized at %s\n", cfg.UploadDir)
	}

	return svc, nil
}

// UploadFile stores a file and returns the file URL/path and S3 key.
// For S3: uploads to the configured bucket and returns the S3 URL.
// For local: saves to the upload directory and returns the local path.
func (svc *Service) UploadFile(file io.Reader, originalFilename string) (fileURL string, s3Key string, err error) {
	// Generate a unique filename to prevent collisions
	ext := filepath.Ext(originalFilename)
	uniqueFilename := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().Unix(), ext)

	if svc.cfg.StorageProvider == "s3" {
		return svc.uploadToS3(file, uniqueFilename)
	}
	return svc.uploadToLocal(file, uniqueFilename)
}

// uploadToS3 uploads a file to the configured S3 bucket.
func (svc *Service) uploadToS3(file io.Reader, filename string) (string, string, error) {
	key := fmt.Sprintf("papers/%s", filename)

	_, err := svc.s3Client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:      aws.String(svc.cfg.S3BucketName),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String("application/pdf"),
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to upload to S3: %w", err)
	}

	// Construct the S3 URL
	fileURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s",
		svc.cfg.S3BucketName, svc.cfg.AWSRegion, key)

	return fileURL, key, nil
}

// uploadToLocal saves a file to the local upload directory.
func (svc *Service) uploadToLocal(file io.Reader, filename string) (string, string, error) {
	filePath := filepath.Join(svc.cfg.UploadDir, filename)

	// Create the destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return "", "", fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	// Copy file content
	if _, err := io.Copy(dst, file); err != nil {
		return "", "", fmt.Errorf("failed to write file: %w", err)
	}

	return filePath, "", nil
}

// GetPresignedURL generates a pre-signed URL for downloading a file from S3.
// For local storage, it returns the local file path.
func (svc *Service) GetPresignedURL(s3Key string) (string, error) {
	if svc.cfg.StorageProvider == "s3" && s3Key != "" {
		presignClient := s3.NewPresignClient(svc.s3Client)

		presignedURL, err := presignClient.PresignGetObject(context.Background(), &s3.GetObjectInput{
			Bucket: aws.String(svc.cfg.S3BucketName),
			Key:    aws.String(s3Key),
		}, s3.WithPresignExpires(1*time.Hour))
		if err != nil {
			return "", fmt.Errorf("failed to generate presigned URL: %w", err)
		}
		return presignedURL.URL, nil
	}

	// For local storage, return the path directly
	return s3Key, nil
}
