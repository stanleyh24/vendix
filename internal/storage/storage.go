package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"vendix/internal/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Client wraps the MinIO/S3 client
type Client struct {
	client *minio.Client
	bucket string
	cfg    *config.Config
}

// NewClient creates a new S3/MinIO storage client with a default bucket
func NewClient(cfg *config.Config) (*Client, error) {
	return NewClientWithBucket(cfg, cfg.S3Bucket)
}

// NewClientWithBucket creates a new S3/MinIO storage client with a specific bucket
func NewClientWithBucket(cfg *config.Config, bucketName string) (*Client, error) {
	// Parse endpoint
	endpoint := cfg.S3Endpoint
	useSSL := cfg.S3UseSSL

	// Initialize minio client
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.S3AccessKey, cfg.S3SecretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	client := &Client{
		client: minioClient,
		bucket: bucketName,
		cfg:    cfg,
	}

	// Ensure bucket exists
	ctx := context.Background()
	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if !exists {
		if err := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{
			Region: cfg.S3Region,
		}); err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	return client, nil
}

// GetBucketName returns the current bucket name
func (c *Client) GetBucketName() string {
	return c.bucket
}

// SetBucket changes the bucket for this client instance
func (c *Client) SetBucket(bucketName string) error {
	ctx := context.Background()
	exists, err := c.client.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if !exists {
		if err := c.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{
			Region: c.cfg.S3Region,
		}); err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	c.bucket = bucketName
	return nil
}

// UploadFile uploads a file to S3/MinIO and returns the object key
func (c *Client) UploadFile(ctx context.Context, reader io.Reader, objectKey string, contentType string, size int64) error {
	_, err := c.client.PutObject(ctx, c.bucket, objectKey, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	return nil
}

// UploadBytes uploads bytes as a file
func (c *Client) UploadBytes(ctx context.Context, data []byte, objectKey string, contentType string) error {
	reader := bytes.NewReader(data)
	return c.UploadFile(ctx, reader, objectKey, contentType, int64(len(data)))
}

// DownloadFile downloads a file from S3/MinIO
func (c *Client) DownloadFile(ctx context.Context, objectKey string) ([]byte, error) {
	obj, err := c.client.GetObject(ctx, c.bucket, objectKey, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}
	defer obj.Close()

	data, err := io.ReadAll(obj)
	if err != nil {
		return nil, fmt.Errorf("failed to read object: %w", err)
	}

	return data, nil
}

// GetFileStream returns a reader for the file (useful for streaming large files)
func (c *Client) GetFileStream(ctx context.Context, objectKey string) (io.Reader, error) {
	obj, err := c.client.GetObject(ctx, c.bucket, objectKey, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}

	return obj, nil
}

// DeleteFile deletes a file from S3/MinIO
func (c *Client) DeleteFile(ctx context.Context, objectKey string) error {
	err := c.client.RemoveObject(ctx, c.bucket, objectKey, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// PresignedURL generates a presigned URL for temporary access to a file
func (c *Client) PresignedURL(ctx context.Context, objectKey string, expiry time.Duration) (string, error) {
	url, err := c.client.PresignedGetObject(ctx, c.bucket, objectKey, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return url.String(), nil
}

// FileExists checks if a file exists in S3/MinIO
func (c *Client) FileExists(ctx context.Context, objectKey string) (bool, error) {
	_, err := c.client.StatObject(ctx, c.bucket, objectKey, minio.StatObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return false, nil
		}
		return false, fmt.Errorf("failed to stat object: %w", err)
	}

	return true, nil
}

// BuildObjectKey builds an object key for a file based on tenant and type
func BuildObjectKey(tenantSchema, documentType, documentID, filename string) string {
	// Format: tenants/{tenant}/documents/{type}/{id}/{filename}
	// Example: tenants/tenant_demo/invoices/abc-123/invoice.pdf
	return filepath.Join("tenants", tenantSchema, "documents", documentType, documentID, filename)
}

// BuildInvoicePDFKey builds the object key for an invoice PDF
func BuildInvoicePDFKey(tenantSchema, invoiceID, invoiceNumber string) string {
	filename := fmt.Sprintf("invoice-%s.pdf", invoiceNumber)
	return BuildObjectKey(tenantSchema, "invoices", invoiceID, filename)
}

// BuildInvoiceXMLKey builds the object key for an invoice XML (DGII)
func BuildInvoiceXMLKey(tenantSchema, invoiceID, invoiceNumber string) string {
	filename := fmt.Sprintf("invoice-%s.xml", invoiceNumber)
	return BuildObjectKey(tenantSchema, "invoices", invoiceID, filename)
}

// BuildReportKey builds the object key for a report
// reportType is used as part of the document type (e.g., "monthly", "annual")
func BuildReportKey(tenantSchema, reportType, reportID, filename string) string {
	// Format: tenants/{tenant}/documents/reports/{type}/{id}/{filename}
	documentType := filepath.Join("reports", reportType)
	return filepath.Join("tenants", tenantSchema, "documents", documentType, reportID, filename)
}

// BuildTenantBucketName generates a bucket name for a tenant
// Format: vendix-{tenant_schema} (e.g., vendix-tenant_demo)
// MinIO bucket names must be lowercase and can contain hyphens
func BuildTenantBucketName(tenantSchema string) string {
	// Convert to lowercase and replace underscores with hyphens
	bucketName := fmt.Sprintf("vendix-%s", strings.ToLower(strings.ReplaceAll(tenantSchema, "_", "-")))
	// Ensure bucket name is valid (3-63 chars, lowercase alphanumeric and hyphens)
	// Remove any invalid characters
	return bucketName
}

