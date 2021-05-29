package lib

import(
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	log "github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

type Storage struct {
	bucketName *string
	clientS3   *s3.Client
	logger     *zap.Logger
}

func NewStorage(bucketName *string, loggerObj *zap.Logger) (*Storage, error) {
	awsCfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.WithField("err", err).Error("error loading aws config")
		return nil,fmt.Errorf("error loading aws config")
	}

	return &Storage{
		clientS3: s3.NewFromConfig(awsCfg),
		bucketName: bucketName,
		logger: loggerObj,
	}, nil
}

func (svc *Storage) ListDirs(key string) ([]string, error) {
	svc.logger.Debug("storage.ListDirs() called",
		zap.String("key", key),
	)

	delimiter := "/"
	params := &s3.ListObjectsV2Input{
		Bucket: svc.bucketName,
		Prefix: &key,
		Delimiter: &delimiter,
	}
	resp,err := svc.clientS3.ListObjectsV2(context.Background(), params)
	if err != nil {
		svc.logger.Error("Error calling s3.ListObjectsV2",
			zap.Reflect("key", params))
		return nil,fmt.Errorf("error calling s3.ListObjectsV2 for %s", key)
	}

	var result []string
	for _,path := range resp.CommonPrefixes{
		result = append(result, *path.Prefix)
	}

	svc.logger.Debug("storage.ListDirs() return",
		zap.String("key", key),
		zap.Reflect("result", result),
	)

	return result,nil
}