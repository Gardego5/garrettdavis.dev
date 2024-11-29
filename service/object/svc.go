package object

import (
	"context"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Service struct {
	presignClient *s3.PresignClient
	logger        *slog.Logger
	bucket        string
}

func New(
	ctx context.Context,
	bucket string,
	logger *slog.Logger,
) (*Service, error) {
	sdkConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		logger.Error("failed to load aws config", "error", err)
		return nil, err
	}

	// Create S3 service client
	svc := s3.NewFromConfig(sdkConfig, func(o *s3.Options) {
		o.BaseEndpoint = aws.String("https://fly.storage.tigris.dev")
		o.Region = "auto"
	})

	// Presigning a request
	ps := s3.NewPresignClient(svc)

	return &Service{presignClient: ps, logger: logger, bucket: bucket}, nil
}

// GetObject makes a presigned request that can be used to get an object from a bucket.
func (svc *Service) GetObject(
	ctx context.Context, key string, expireSecs int64,
) (*v4.PresignedHTTPRequest, error) {
	request, err := svc.presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(svc.bucket),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(expireSecs * int64(time.Second))
	})
	if err != nil {
		svc.logger.Error("failed to create presigned GET",
			"bucket", svc.bucket, "key", key, "error", err)
	}
	return request, err
}

// PutObject makes a presigned request that can be used to put an object in a bucket.
func (svc *Service) PutObject(
	ctx context.Context, key string, expireSecs int64,
) (*v4.PresignedHTTPRequest, error) {
	request, err := svc.presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: &svc.bucket, Key: &key,
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(expireSecs * int64(time.Second))
	})
	if err != nil {
		svc.logger.Error("failed to create presigned PUT",
			"bucket", svc.bucket, "key", key, "error", err)
	}
	return request, err
}

// DeleteObject makes a presigned request that can be used to delete an object from a bucket.
func (svc *Service) DeleteObject(
	ctx context.Context, key string, expireSecs int64,
) (*v4.PresignedHTTPRequest, error) {
	request, err := svc.presignClient.PresignDeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(svc.bucket),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(expireSecs * int64(time.Second))
	})
	if err != nil {
		svc.logger.Error("failed to create presigned DELETE",
			"bucket", svc.bucket, "key", key, "error", err)
	}
	return request, err
}
