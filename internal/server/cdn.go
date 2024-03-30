package server

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"cloud.google.com/go/storage"
)

type CDN interface {
	// UploadFile uploads a file to the CDN
	UploadFolder(path string) error

	// RootURL returns the root URL of the CDN
	RootURL() string
}

type LocalCDN struct {
	rootURL string
}

func NewLocalCDN(rootURL string) *LocalCDN {
	return &LocalCDN{rootURL: rootURL}
}

func (l *LocalCDN) RootURL() string {
	return l.rootURL
}

func (l *LocalCDN) UploadFolder(path string) error {
	return nil
}

type GoogleCloudBucketCDN struct {
	Bucket string
}

func NewGoogleCloudBucketCDN(bucket string) *GoogleCloudBucketCDN {
	return &GoogleCloudBucketCDN{Bucket: bucket}
}

func (g *GoogleCloudBucketCDN) RootURL() string {
	return fmt.Sprintf("https://storage.googleapis.com/%s/", g.Bucket)
}

func (g *GoogleCloudBucketCDN) UploadFolder(path string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	return uploadFolder(ctx, client, g.Bucket, path)
}

// uploadFolder uploads all files within a specified folder to a GCS bucket.
func uploadFolder(ctx context.Context, client *storage.Client, bucket, folderPath string) error {
	// Iterate through all files in the folder.
	return filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			// Construct the object name and ensure it reflects the folder structure.
			objectName := filepath.Join(filepath.Base(folderPath), path[len(folderPath):])
			if err := uploadFile(ctx, client, bucket, objectName, path); err != nil {
				return err
			}
		}
		return nil
	})
}

// uploadFile uploads an individual file to a GCS bucket.
func uploadFile(ctx context.Context, client *storage.Client, bucket, objectName, filePath string) error {
	bkt := client.Bucket(bucket)

	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("os.Open: %w", err)
	}
	defer f.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	o := bkt.Object(objectName)

	// Optional: set a generation-match precondition to avoid potential race
	// conditions and data corruptions. The request to upload is aborted if the
	// object's generation number does not match your precondition.
	// For an object that does not yet exist, set the DoesNotExist precondition.
	o = o.If(storage.Conditions{DoesNotExist: true})
	// If the live object already exists in your bucket, set instead a
	// generation-match precondition using the live object's generation number.
	// attrs, err := o.Attrs(ctx)
	// if err != nil {
	//      return fmt.Errorf("object.Attrs: %w", err)
	// }
	// o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})

	// Upload an object with storage.Writer.
	wc := o.NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return fmt.Errorf("io.Copy: %w", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %w", err)
	}
	return nil
}
