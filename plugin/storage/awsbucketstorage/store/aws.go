package store

import (
	"bytes"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// AWSStorage represents a Google Cloud Storage object.
type AWSStorage struct {
	bucket string
	object string
}

// NewAWSStorage returns a Google Cloud storage item given a bucket and an object
// path.
func NewAWSStorage(bucket string, object string) *AWSStorage {
	return &AWSStorage{
		bucket: bucket,
		object: object,
	}
}

// Load downloads an object from a bucket and returns an error if it cannot
// be read.
func (s *AWSStorage) Load() ([]byte, error) {
	// Load config.
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	// Create the buffer to save the file.
	var buf []byte
	f := manager.NewWriteAtBuffer(buf)

	// Create an uploader with the session and default options.
	client := s3.NewFromConfig(cfg)
	downloader := manager.NewDownloader(client)
	_, err = downloader.Download(context.TODO(), f, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(s.object),
	})
	if err != nil {
		return nil, err
	}

	return f.Bytes(), nil
}

// Save uploads an object to a bucket and returns an error if it cannot be
// written.
func (s *AWSStorage) Save(b []byte) error {
	// Load config.
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}

	// Create an uploader with the session and default options.
	client := s3.NewFromConfig(cfg)
	uploader := manager.NewUploader(client)
	_, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(s.object),
		Body:   bytes.NewBuffer(b),
	})
	if err != nil {
		return err
	}

	return nil
}
