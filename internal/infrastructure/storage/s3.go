package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	iface "todo-service/internal/domain/interface"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
)

type s3Repository struct {
	client *s3.Client
	bucket string
}

func NewS3Repository(ctx context.Context, bucket, endpoint string) (iface.S3Repository, error) {
	log.Info().Str("bucket", bucket).Str("endpoint", endpoint).Msg("Initializing S3 repository")
	if endpoint == "" {
		log.Error().Msg("S3 endpoint cannot be empty")
		return nil, errors.New("S3 endpoint cannot be empty")
	}
	if bucket == "" {
		log.Error().Msg("S3 bucket cannot be empty")
		return nil, errors.New("S3 bucket cannot be empty")
	}
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("test", "test", "")),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:               endpoint,
				SigningRegion:     "us-east-1",
				HostnameImmutable: true,
			}, nil
		})),
		config.WithRegion("us-east-1"),
		config.WithRetryMaxAttempts(5),
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to load AWS config")
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}
	client := s3.NewFromConfig(cfg)
	log.Info().Msg("S3 client initialized successfully")
	return &s3Repository{client: client, bucket: bucket}, nil
}

func (r *s3Repository) UploadFile(ctx context.Context, bucket, key string, reader io.Reader, size int64) (string, error) {
	log.Info().Str("bucket", bucket).Str("key", key).Int64("size", size).Msg("Uploading file to S3")
	u := manager.NewUploader(r.client)
	_, err := u.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   reader,
	})
	if err != nil {
		log.Error().Err(err).Str("bucket", bucket).Str("key", key).Msg("Failed to upload file")
		return "", err
	}
	log.Info().Str("bucket", bucket).Str("key", key).Msg("File uploaded successfully")
	return key, nil
}
