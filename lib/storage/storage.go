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
	"gradle-serverless-dependencies-graph/lib/helpers"
	"time"
)

type Storage struct {
	Config   *StorageConfig
	DynamoDb *dynamodb.Client
	Logger   *zap.Logger
}

type StorageConfig struct {
	DependenciesTableName *string
	RepositoriesTableName *string
	StorageTableName      *string
}

type InsertItem struct {
	Table string
	Item  types.WriteRequest
}

const (
	ErrUnknown = iota
	ErrObjectNotFound
)

const RootParent = "-"

func (s StorageErrorRest) Error() string {
	panic(s.Message)
}

func NewStorage(cfg StorageConfig, logger *zap.Logger) (*Storage, error) {
	awsCfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		logger.Error("error loading aws Config", zap.Error(err))
		return nil, fmt.Errorf("error loading aws Config")
	}
	clientDynamoDb := dynamodb.NewFromConfig(awsCfg)

	return &Storage{
		Config: &StorageConfig{
			DependenciesTableName: cfg.DependenciesTableName,
			RepositoriesTableName: cfg.RepositoriesTableName,
			StorageTableName:      cfg.StorageTableName,
		},
		DynamoDb: clientDynamoDb,
		Logger:   logger,
	}, nil
}

func (svc *Storage) UpsertRepositoryInfo(ctxId string, repo string, ref string, deps DependenciesRest) (*UpsertResultRest, *StorageErrorRest) {
	svc.Logger.Debug(fmt.Sprintf("%s UpsertRepositoryInfo() called", ctxId),
		zap.String("repo", repo),
		zap.String("ref", ref),
		zap.Reflect("deps", deps),
	)

	updated := time.Now().Format(time.RFC3339)

	var insertBatch []InsertItem

	// Add data to storage. Details with dependencies and versions per repo/ref
	var groupsToInsert map[string]bool
	for _, dep := range deps.Dependencies {
		Id := fmt.Sprintf("%s:%s:%s:%s", repo, ref, dep.Group, dep.Name)
		Dep := fmt.Sprintf("%s:%s", dep.Group, dep.Name)

		groupsToInsert[dep.Group] = true
		insertBatch = append(insertBatch,
			// Add info to dependencies table: groupId -> name
			InsertItem{
				Table: *svc.Config.DependenciesTableName,
				Item: types.WriteRequest{
					PutRequest: &types.PutRequest{
						Item: map[string]types.AttributeValue{
							"Parent":  &types.AttributeValueMemberS{Value: dep.Group},
							"Child":   &types.AttributeValueMemberS{Value: dep.Name},
							"Updated": &types.AttributeValueMemberS{Value: updated},
						},
					},
				},
			},
			InsertItem{
				Table: *svc.Config.StorageTableName,
				Item: types.WriteRequest{
					PutRequest: &types.PutRequest{
						Item: map[string]types.AttributeValue{
							"Id":         &types.AttributeValueMemberS{Value: Id},
							"Dependency": &types.AttributeValueMemberS{Value: Dep},
							"Version":    &types.AttributeValueMemberS{Value: dep.Version},
							"Parent":     &types.AttributeValueMemberS{Value: repo},
							"Child":      &types.AttributeValueMemberS{Value: ref},
							"Updated":    &types.AttributeValueMemberS{Value: updated},
						},
					},
				},
			},
		)
	}
	// Add info to dependencies table: root -> groupId
	for group, _ := range groupsToInsert {
		insertBatch = append(insertBatch,
			InsertItem{
				Table: *svc.Config.DependenciesTableName,
				Item: types.WriteRequest{
					PutRequest: &types.PutRequest{
						Item: map[string]types.AttributeValue{
							"Parent":  &types.AttributeValueMemberS{Value: RootParent},
							"Child":   &types.AttributeValueMemberS{Value: group},
							"Updated": &types.AttributeValueMemberS{Value: updated},
						},
					},
				},
			},
		)
	}

	// Add info to repositories table: root -> repo -> ref
	insertBatch = append(insertBatch,
		InsertItem{
			Table: *svc.Config.RepositoriesTableName,
			Item: types.WriteRequest{
				PutRequest: &types.PutRequest{
					Item: map[string]types.AttributeValue{
						"Parent":  &types.AttributeValueMemberS{Value: RootParent},
						"Child":   &types.AttributeValueMemberS{Value: repo},
						"Updated": &types.AttributeValueMemberS{Value: updated},
					},
				},
			},
		},
		InsertItem{
			Table: *svc.Config.RepositoriesTableName,
			Item: types.WriteRequest{
				PutRequest: &types.PutRequest{
					Item: map[string]types.AttributeValue{
						"Parent":  &types.AttributeValueMemberS{Value: repo},
						"Child":   &types.AttributeValueMemberS{Value: ref},
						"Updated": &types.AttributeValueMemberS{Value: updated},
					},
				},
			},
		},
	)

	retry := 5
	result := UpsertResultRest{
		0,
	}
	for len(insertBatch) > 0 && retry > 0 {
		params := &dynamodb.BatchWriteItemInput{
			RequestItems:                make(map[string][]types.WriteRequest),
			ReturnConsumedCapacity:      types.ReturnConsumedCapacityIndexes,
			ReturnItemCollectionMetrics: types.ReturnItemCollectionMetricsSize,
		}

		for _, item := range insertBatch[:helpers.Min(25, len(insertBatch))] {
			if _, ok := params.RequestItems[item.Table]; ok {
				params.RequestItems[item.Table] = append(params.RequestItems[item.Table], item.Item)
			} else {
				params.RequestItems[item.Table] = []types.WriteRequest{
					item.Item,
				}
			}
		}

		svc.Logger.Debug(fmt.Sprintf("%s UpsertRepositoryInfo()", ctxId),
			zap.String("repo", repo),
			zap.String("ref", ref),
			zap.Reflect("params", params),
		)

		resp, err := svc.DynamoDb.BatchWriteItem(context.Background(), params)
		if err != nil {
			return nil, svc.handleError(ctxId, err, "UpsertRepositoryInfo",
				map[string]string{
					"repo": repo,
					"ref":  ref,
				},
				zap.String("repo", repo),
				zap.String("ref", ref),
				zap.Reflect("params", params),
			)
		}
		svc.Logger.Debug(fmt.Sprintf("%s UpsertRepositoryInfo()", ctxId),
			zap.String("repo", repo),
			zap.String("ref", ref),
			zap.Reflect("resp", resp),
		)
		result.UsedCapacity = result.UsedCapacity + *resp.ConsumedCapacity[0].CapacityUnits

		if len(insertBatch) > 25 {
			insertBatch = insertBatch[25:]
		} else {
			insertBatch = make([]InsertItem, 0)
		}

		if len(resp.UnprocessedItems) > 0 {
			retry--
			for table, items := range resp.UnprocessedItems {
				for _, item := range items {
					insertBatch = append(insertBatch, InsertItem{
						Table: table,
						Item:  item,
					})
				}
			}
		}
	}

	svc.Logger.Debug(fmt.Sprintf("%s UpsertRepositoryInfo()", ctxId),
		zap.String("repo", repo),
		zap.String("ref", ref),
	)

	return &result, nil
}

func (svc *Storage) handleError(ctxId string, err error, method string, keys map[string]string, fields ...zap.Field) *StorageErrorRest {
	var oe *smithy.OperationError
	var errApi *smithy.GenericAPIError
	svc.Logger.Debug("handleError() err",
		zap.Error(err),
	)

	if errors.As(err, &oe) && oe.Service() == "DynamoDB" {
		svc.Logger.Debug("handleError() oe",
			zap.Reflect("oe",oe),
		)

		if errors.As(err, &errApi) {
			svc.Logger.Debug("handleError() errApi",
				zap.Reflect("errApi", errApi),
			)

			if errApi.Code == "NotFound" {
				fields = append(fields, zap.NamedError("errApi", errApi))
				svc.Logger.Warn(fmt.Sprintf("%s storageSvc.%s() DynamoDB.NotFound", ctxId, method),
					fields...,
				)
				return &StorageErrorRest{
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
	return &StorageErrorRest{
		Message: fmt.Sprintf("Error #%d while quering data", ErrUnknown),
		Code:    ErrUnknown,
		Repo:    keys["repo"],
		Ref:     keys["ref"],
		Id:      keys["id"],
		Version: keys["version"],
		Err:     err,
	}
}
