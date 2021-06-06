package storage

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"strings"
	"time"
)

type Storage struct {
	bucketName        *string
	clientS3          *s3.Client
	clientS3Presigned *s3.PresignClient
	Logger            *zap.Logger
}

const (
	ErrUnknown = iota
	ErrObjectNotFound
	ErrObjectNotAccessible
)

type StorageError struct {
	Message    string
	Code       int
	BucketName string
	Key        string
	Err        error
}

func (s StorageError) Error() string {
	panic(s.Message)
}

const PresignUrlDuration = time.Duration(10) * time.Minute

func NewStorage(bucketName string, logger *zap.Logger) (*Storage, error) {
	awsCfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		logger.Error("error loading aws config", zap.Error(err))
		return nil, fmt.Errorf("error loading aws config")
	}

	clientS3 := s3.NewFromConfig(awsCfg)
	clientS3Presigned := s3.NewPresignClient(clientS3)

	return &Storage{
		clientS3:          s3.NewFromConfig(awsCfg),
		clientS3Presigned: clientS3Presigned,
		bucketName:        &bucketName,
		Logger:            logger,
	}, nil
}

func (svc *Storage) ListDirs(ctxId string, key string) (*[]string, *StorageError) {
	svc.Logger.Debug(fmt.Sprintf("%s storageSvc.ListDirs() called", ctxId),
		zap.String("key", key),
	)

	if !strings.HasSuffix(key, "/") {
		key = fmt.Sprintf("%s/", key)
	}

	delimiter := "/"
	params := &s3.ListObjectsV2Input{
		Bucket:    svc.bucketName,
		Prefix:    &key,
		Delimiter: &delimiter,
	}
	resp, err := svc.clientS3.ListObjectsV2(context.Background(), params)
	if err != nil {
		return nil, svc.handleError(ctxId, err, "storageSvc.ListDirs", key,
			zap.String("key", key),
			zap.Reflect("params", params),
		)
	}
	svc.Logger.Debug(fmt.Sprintf("%s resp", ctxId),
		zap.Reflect("resp", resp),
	)

	var result []string

	for _, path := range resp.CommonPrefixes {
		result = append(result, *path.Prefix)
	}

	svc.Logger.Debug(fmt.Sprintf("%s storageSvc.ListDirs() return", ctxId),
		zap.String("key", key),
		zap.Reflect("result", result),
	)

	return &result, nil
}

func (svc *Storage) GetDownloadUrl(ctxId string, key string) (*string, *StorageError) {
	svc.Logger.Debug(fmt.Sprintf("%s storageSvc.GetDownloadUrl() called", ctxId),
		zap.String("key", key),
	)

	paramsHead := &s3.HeadObjectInput{
		Bucket: svc.bucketName,
		Key:    &key,
	}

	_, err := svc.clientS3.HeadObject(context.TODO(), paramsHead)
	if err != nil {
		return nil, svc.handleError(ctxId, err, "storageSvc.GetDownloadUrl.HeadObject", key,
			zap.Reflect("bucket", svc.bucketName),
			zap.String("key", key),
			zap.Reflect("params", paramsHead),
		)
	}

	paramsSign := &s3.GetObjectInput{
		Bucket: svc.bucketName,
		Key:    &key,
	}

	resp, err := svc.clientS3Presigned.PresignGetObject(context.TODO(), paramsSign, s3.WithPresignExpires(PresignUrlDuration))
	if err != nil {
		return nil, svc.handleError(ctxId, err, "storageSvc.GetDownloadUrl.PresignObject", key,
			zap.Reflect("bucket", svc.bucketName),
			zap.String("key", key),
			zap.Reflect("params", paramsSign),
		)
	}

	svc.Logger.Debug(fmt.Sprintf("%s storageSvc.ListDirs() return", ctxId),
		zap.String("key", key),
		zap.Reflect("result", &resp.URL),
	)

	return &resp.URL, nil
}

func (svc *Storage) handleError(ctxId string, err error, method string, key string, fields ...zap.Field) *StorageError {
	var oe *smithy.OperationError
	var errApi *smithy.GenericAPIError
	if errors.As(err, &oe) && oe.Service() == "S3" {
		if errors.As(err, &errApi) {
			if errApi.Code == "NotFound" {
				fields = append(fields, zap.NamedError("errApi", errApi))
				svc.Logger.Warn(fmt.Sprintf("%s storageSvc.%s() S3.NotFound", ctxId, method),
					fields...,
				)
				return &StorageError{
					Message:    fmt.Sprintf("Error #%d Object Not Found while generating signed url for %s", ErrObjectNotFound, key),
					Code:       ErrObjectNotFound,
					BucketName: *svc.bucketName,
					Key:        key,
					Err:        err,
				}
			}
		}
	}

	fields = append(fields, zap.NamedError("err", err))
	svc.Logger.Warn(fmt.Sprintf("%s storageSvc.%s() Unknown", ctxId, method),
		fields...,
	)
	return &StorageError{
		Message:    fmt.Sprintf("Error #%d while generating signed url for %s", ErrUnknown, key),
		Code:       ErrUnknown,
		BucketName: *svc.bucketName,
		Key:        key,
		Err:        err,
	}
}
