package data

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// TransferProgress represents progress of an S3 transfer operation
type TransferProgress struct {
	BytesTransferred int64
	TotalBytes       int64
	Percentage       float64
	Speed            int64 // bytes per second
	ETA              time.Duration
	StartTime        time.Time
	LastUpdate       time.Time
}

// ProgressCallback is called during transfer operations to report progress
type ProgressCallback func(progress TransferProgress)

// S3Manager handles optimized S3 transfer operations
type S3Manager struct {
	client     *s3.Client
	uploader   *manager.Uploader
	downloader *manager.Downloader
	region     string
	
	// Transfer configuration
	partSize      int64
	concurrency   int
	maxRetries    int
	
	// Progress tracking
	progressMu    sync.RWMutex
	activeTransfers map[string]*TransferProgress
}

// S3ManagerConfig holds configuration for S3Manager
type S3ManagerConfig struct {
	PartSize      int64 // Default: 16MB
	Concurrency   int   // Default: 10
	MaxRetries    int   // Default: 3
}

// NewS3Manager creates a new S3 transfer manager with optimizations
func NewS3Manager(client *s3.Client, region string, config *S3ManagerConfig) *S3Manager {
	if config == nil {
		config = &S3ManagerConfig{
			PartSize:    16 * 1024 * 1024, // 16MB
			Concurrency: 10,
			MaxRetries:  3,
		}
	}

	uploader := manager.NewUploader(client, func(u *manager.Uploader) {
		u.PartSize = config.PartSize
		u.Concurrency = config.Concurrency
	})

	downloader := manager.NewDownloader(client, func(d *manager.Downloader) {
		d.PartSize = config.PartSize
		d.Concurrency = config.Concurrency
	})

	return &S3Manager{
		client:          client,
		uploader:        uploader,
		downloader:      downloader,
		region:          region,
		partSize:        config.PartSize,
		concurrency:     config.Concurrency,
		maxRetries:      config.MaxRetries,
		activeTransfers: make(map[string]*TransferProgress),
	}
}

// UploadFile uploads a file to S3 with progress tracking and optimization
func (sm *S3Manager) UploadFile(ctx context.Context, bucket, key, filePath string, callback ProgressCallback) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat file %s: %w", filePath, err)
	}

	transferID := fmt.Sprintf("upload:%s:%s", bucket, key)
	progress := &TransferProgress{
		TotalBytes: stat.Size(),
		StartTime:  time.Now(),
		LastUpdate: time.Now(),
	}

	sm.progressMu.Lock()
	sm.activeTransfers[transferID] = progress
	sm.progressMu.Unlock()

	defer func() {
		sm.progressMu.Lock()
		delete(sm.activeTransfers, transferID)
		sm.progressMu.Unlock()
	}()

	// Create progress reader wrapper
	progressReader := &progressReader{
		reader:     file,
		progress:   progress,
		callback:   callback,
		progressMu: &sm.progressMu,
	}

	_, err = sm.uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   progressReader,
	})

	if err != nil {
		return fmt.Errorf("failed to upload file to s3://%s/%s: %w", bucket, key, err)
	}

	return nil
}

// DownloadFile downloads a file from S3 with progress tracking and optimization
func (sm *S3Manager) DownloadFile(ctx context.Context, bucket, key, filePath string, callback ProgressCallback) error {
	// Get object info for progress tracking
	headResult, err := sm.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to get object info for s3://%s/%s: %w", bucket, key, err)
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filePath, err)
	}
	defer file.Close()

	transferID := fmt.Sprintf("download:%s:%s", bucket, key)
	progress := &TransferProgress{
		TotalBytes: *headResult.ContentLength,
		StartTime:  time.Now(),
		LastUpdate: time.Now(),
	}

	sm.progressMu.Lock()
	sm.activeTransfers[transferID] = progress
	sm.progressMu.Unlock()

	defer func() {
		sm.progressMu.Lock()
		delete(sm.activeTransfers, transferID)
		sm.progressMu.Unlock()
	}()

	// Create progress writer wrapper
	progressWriter := &progressWriterAt{
		writerAt:   file,
		progress:   progress,
		callback:   callback,
		progressMu: &sm.progressMu,
	}

	_, err = sm.downloader.Download(ctx, progressWriter, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return fmt.Errorf("failed to download file from s3://%s/%s: %w", bucket, key, err)
	}

	return nil
}

// ListObjects lists objects in an S3 bucket with pagination support
func (sm *S3Manager) ListObjects(ctx context.Context, bucket, prefix string, maxKeys int32) ([]types.Object, error) {
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	}
	
	if prefix != "" {
		input.Prefix = aws.String(prefix)
	}
	
	if maxKeys > 0 {
		input.MaxKeys = aws.Int32(maxKeys)
	}

	result, err := sm.client.ListObjectsV2(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to list objects in bucket %s: %w", bucket, err)
	}

	return result.Contents, nil
}

// CreateBucket creates an S3 bucket with appropriate configuration
func (sm *S3Manager) CreateBucket(ctx context.Context, bucket string) error {
	input := &s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	}

	// Add location constraint for regions other than us-east-1
	if sm.region != "us-east-1" {
		input.CreateBucketConfiguration = &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(sm.region),
		}
	}

	_, err := sm.client.CreateBucket(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to create bucket %s: %w", bucket, err)
	}

	return nil
}

// BucketExists checks if a bucket exists and is accessible
func (sm *S3Manager) BucketExists(ctx context.Context, bucket string) (bool, error) {
	_, err := sm.client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})
	
	if err != nil {
		// Check if it's a "not found" error
		if strings.Contains(err.Error(), "NoSuchBucket") || strings.Contains(err.Error(), "NotFound") {
			return false, nil
		}
		return false, fmt.Errorf("failed to check bucket %s: %w", bucket, err)
	}
	
	return true, nil
}

// GetActiveTransfers returns current active transfer progress
func (sm *S3Manager) GetActiveTransfers() map[string]TransferProgress {
	sm.progressMu.RLock()
	defer sm.progressMu.RUnlock()

	transfers := make(map[string]TransferProgress)
	for id, progress := range sm.activeTransfers {
		transfers[id] = *progress
	}

	return transfers
}

// progressReader wraps an io.Reader to track upload progress
type progressReader struct {
	reader     io.Reader
	progress   *TransferProgress
	callback   ProgressCallback
	progressMu *sync.RWMutex
}

func (pr *progressReader) Read(p []byte) (int, error) {
	n, err := pr.reader.Read(p)
	
	if n > 0 {
		pr.progressMu.Lock()
		pr.progress.BytesTransferred += int64(n)
		pr.updateProgress()
		pr.progressMu.Unlock()
	}
	
	return n, err
}

func (pr *progressReader) updateProgress() {
	now := time.Now()
	elapsed := now.Sub(pr.progress.StartTime)
	
	if pr.progress.TotalBytes > 0 {
		pr.progress.Percentage = float64(pr.progress.BytesTransferred) / float64(pr.progress.TotalBytes) * 100
	}
	
	if elapsed.Seconds() > 0 {
		pr.progress.Speed = pr.progress.BytesTransferred / int64(elapsed.Seconds())
		
		if pr.progress.Speed > 0 && pr.progress.TotalBytes > pr.progress.BytesTransferred {
			remaining := pr.progress.TotalBytes - pr.progress.BytesTransferred
			pr.progress.ETA = time.Duration(remaining/pr.progress.Speed) * time.Second
		}
	}
	
	pr.progress.LastUpdate = now
	
	if pr.callback != nil {
		pr.callback(*pr.progress)
	}
}

// progressWriterAt wraps an io.WriterAt to track download progress
type progressWriterAt struct {
	writerAt   io.WriterAt
	progress   *TransferProgress
	callback   ProgressCallback
	progressMu *sync.RWMutex
}

func (pw *progressWriterAt) WriteAt(p []byte, off int64) (int, error) {
	n, err := pw.writerAt.WriteAt(p, off)
	
	if n > 0 {
		pw.progressMu.Lock()
		pw.progress.BytesTransferred += int64(n)
		pw.updateProgress()
		pw.progressMu.Unlock()
	}
	
	return n, err
}

func (pw *progressWriterAt) Write(p []byte) (int, error) {
	// This is a fallback for io.Writer interface compatibility
	return pw.WriteAt(p, 0)
}

func (pw *progressWriterAt) updateProgress() {
	now := time.Now()
	elapsed := now.Sub(pw.progress.StartTime)
	
	if pw.progress.TotalBytes > 0 {
		pw.progress.Percentage = float64(pw.progress.BytesTransferred) / float64(pw.progress.TotalBytes) * 100
	}
	
	if elapsed.Seconds() > 0 {
		pw.progress.Speed = pw.progress.BytesTransferred / int64(elapsed.Seconds())
		
		if pw.progress.Speed > 0 && pw.progress.TotalBytes > pw.progress.BytesTransferred {
			remaining := pw.progress.TotalBytes - pw.progress.BytesTransferred
			pw.progress.ETA = time.Duration(remaining/pw.progress.Speed) * time.Second
		}
	}
	
	pw.progress.LastUpdate = now
	
	if pw.callback != nil {
		pw.callback(*pw.progress)
	}
}