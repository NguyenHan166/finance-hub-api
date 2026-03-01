package utils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
)

type R2Storage struct {
	client        *s3.S3
	bucket        string
	publicBaseURL string
}

// NewR2Storage creates a new R2 storage client
func NewR2Storage(endpoint, accessKeyID, secretAccessKey, bucket, publicBaseURL string) (*R2Storage, error) {
	// Create AWS session for R2 (R2 is S3-compatible)
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("auto"), // R2 uses "auto" region
		Endpoint:    aws.String(endpoint),
		Credentials: credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
		S3ForcePathStyle: aws.Bool(true), // Required for R2
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create R2 session: %w", err)
	}

	return &R2Storage{
		client:        s3.New(sess),
		bucket:        bucket,
		publicBaseURL: publicBaseURL,
	}, nil
}

// UploadFile uploads a file to R2 and returns the public URL
func (r *R2Storage) UploadFile(ctx context.Context, file multipart.File, header *multipart.FileHeader, folder string) (string, error) {
	// Read file content
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%s-%s%s", uuid.New().String(), time.Now().Format("20060102150405"), ext)
	
	// Create full path
	key := filename
	if folder != "" {
		key = fmt.Sprintf("%s/%s", strings.Trim(folder, "/"), filename)
	}

	// Determine content type
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Upload to R2
	_, err = r.client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(r.bucket),
		Key:           aws.String(key),
		Body:          bytes.NewReader(fileBytes),
		ContentType:   aws.String(contentType),
		ContentLength: aws.Int64(header.Size),
		// Make the object publicly accessible
		ACL: aws.String("public-read"),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to R2: %w", err)
	}

	// Generate public URL
	publicURL := fmt.Sprintf("%s/%s", strings.TrimSuffix(r.publicBaseURL, "/"), key)
	
	return publicURL, nil
}

// DeleteFile deletes a file from R2
func (r *R2Storage) DeleteFile(ctx context.Context, fileURL string) error {
	// Extract key from URL
	key := strings.TrimPrefix(fileURL, r.publicBaseURL+"/")
	
	_, err := r.client.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete from R2: %w", err)
	}

	return nil
}

// ValidateFile validates file size and type
func ValidateFile(file multipart.File, header *multipart.FileHeader, maxSize int64, allowedTypes []string) error {
	// Check file size
	if header.Size > maxSize {
		return fmt.Errorf("file size exceeds limit of %d bytes", maxSize)
	}

	// Check file type
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		return fmt.Errorf("content type not specified")
	}

	// Validate against allowed types
	allowed := false
	for _, t := range allowedTypes {
		if strings.HasPrefix(contentType, t) || contentType == t {
			allowed = true
			break
		}
	}
	if !allowed {
		return fmt.Errorf("file type %s not allowed", contentType)
	}

	return nil
}
