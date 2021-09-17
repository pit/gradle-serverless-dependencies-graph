package storage

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/smithy-go"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"time"
)

type Storage struct {
	tableName      *string
	clientDynamoDb *dynamodb.Client
	Logger         *zap.Logger
}

const (
	ErrUnknown = iota
	ErrObjectNotFound
)

type Dependencies struct {
	Dependencies []Dependency `json:"dependencies"`
}

type Dependency struct {
	Id      string `json:"id"`
	Version string `json:"version"`
}

type UpsertResult struct {
	UsedCapacity float64
}

type StorageError struct {
	Message string
	Code    int
	Repo    string
	Ref     string
	Id      string
	Version string
	Err     error
}

func (s StorageError) Error() string {
	panic(s.Message)
}

func NewStorage(tableName string, logger *zap.Logger) (*Storage, error) {
	awsCfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		logger.Error("error loading aws config", zap.Error(err))
		return nil, fmt.Errorf("error loading aws config")
	}

	clientDynamoDb := dynamodb.NewFromConfig(awsCfg)

	return &Storage{
		clientDynamoDb: clientDynamoDb,
		tableName:      &tableName,
		Logger:         logger,
	}, nil
}

func (svc *Storage) UpsertRepoInfo(ctxId string, repo string, ref string, deps Dependencies) (*UpsertResult, *StorageError) {
	svc.Logger.Debug(fmt.Sprintf("%s UpsertRepoInfo() called", ctxId),
		zap.String("repo", repo),
		zap.String("ref", ref),
		zap.Reflect("deps", deps),
	)

	updated := time.Now().Format(time.RFC3339)

	var itemsToInsert []types.WriteRequest
	for _, dep := range deps.Dependencies {
		Id := fmt.Sprintf("%s:%s:%s", repo, ref, dep.Id)
		itemsToInsert = append(itemsToInsert, types.WriteRequest{
			PutRequest: &types.PutRequest{
				Item: map[string]types.AttributeValue{
					"Id":         &types.AttributeValueMemberS{Value: Id},
					"Dependency": &types.AttributeValueMemberS{Value: dep.Id},
					"Version":    &types.AttributeValueMemberS{Value: dep.Version},
					"Repo":       &types.AttributeValueMemberS{Value: repo},
					"Ref":        &types.AttributeValueMemberS{Value: ref},
					"Updated":    &types.AttributeValueMemberS{Value: updated},
				},
			},
		})
	}

	retry := 5
	result := UpsertResult{
		0,
	}
	for len(itemsToInsert) > 0 && retry > 0 {
		params := &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
				*svc.tableName: itemsToInsert[:min(25, len(itemsToInsert))],
			},
			ReturnConsumedCapacity:      types.ReturnConsumedCapacityIndexes,
			ReturnItemCollectionMetrics: types.ReturnItemCollectionMetricsSize,
		}

		svc.Logger.Debug(fmt.Sprintf("%s UpsertRepoInfo()", ctxId),
			zap.String("repo", repo),
			zap.String("ref", ref),
			zap.Reflect("params", params),
		)

		resp, err := svc.clientDynamoDb.BatchWriteItem(context.Background(), params)
		if err != nil {
			return nil, svc.handleError(ctxId, err, "UpsertRepoInfo",
				map[string]string{
					"repo": repo,
					"ref":  ref,
				},
				zap.String("repo", repo),
				zap.String("ref", ref),
				zap.Reflect("params", params),
			)
		}
		svc.Logger.Debug(fmt.Sprintf("%s UpsertRepoInfo()", ctxId),
			zap.String("repo", repo),
			zap.String("ref", ref),
			zap.Reflect("resp", resp),
		)
		result.UsedCapacity = result.UsedCapacity + *resp.ConsumedCapacity[0].CapacityUnits

		if len(itemsToInsert) > 25 {
			itemsToInsert = itemsToInsert[25:]
		} else {
			itemsToInsert = make([]types.WriteRequest, 0)
		}

		if len(resp.UnprocessedItems) > 0 {
			retry--
			itemsToInsert = append(itemsToInsert, resp.UnprocessedItems[*svc.tableName]...)
		}
	}

	svc.Logger.Debug(fmt.Sprintf("%s UpsertRepoInfo()", ctxId),
		zap.String("repo", repo),
		zap.String("ref", ref),
	)

	return &result, nil
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func (svc *Storage) handleError(ctxId string, err error, method string, keys map[string]string, fields ...zap.Field) *StorageError {
	var oe *smithy.OperationError
	var errApi *smithy.GenericAPIError
	if errors.As(err, &oe) && oe.Service() == "DynamoDB" {
		if errors.As(err, &errApi) {
			if errApi.Code == "NotFound" {
				fields = append(fields, zap.NamedError("errApi", errApi))
				svc.Logger.Warn(fmt.Sprintf("%s storageSvc.%s() DynamoDB.NotFound", ctxId, method),
					fields...,
				)
				return &StorageError{
					Message: fmt.Sprintf("Error #%d Data Not Found", ErrObjectNotFound),
					Code:    ErrObjectNotFound,
					Repo:    keys["repo"],
					Ref:     keys["ref"],
					Id:      keys["id"],
					Version: keys["version"],
					Err:     err,
				}
			}
		}
	}

	fields = append(fields, zap.NamedError("err", err))
	svc.Logger.Error(fmt.Sprintf("%s storageSvc.%s() Unknown", ctxId, method),
		fields...,
	)
	return &StorageError{
		Message: fmt.Sprintf("Error #%d while quering data", ErrUnknown),
		Code:    ErrUnknown,
		Repo:    keys["repo"],
		Ref:     keys["ref"],
		Id:      keys["id"],
		Version: keys["version"],
		Err:     err,
	}
}
