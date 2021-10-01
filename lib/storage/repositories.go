package storage

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/smithy-go/ptr"
	"go.uber.org/zap"
)

func (svc *Storage) ListRepositoriesByParent(ctxId string, parent *string) (*[]RepositoryDto, *StorageErrorRest) {
	if parent == nil {
		parent = ptr.String(RootParent)
	}

	svc.Logger.Debug(fmt.Sprintf("%s ListRepositoriesByParent() called", ctxId),
		zap.String("parent", *parent),
	)

	var consistentRead = false
	params := &dynamodb.QueryInput{
		TableName:              svc.Config.RepositoriesTableName,
		ConsistentRead:         &consistentRead,
		KeyConditionExpression: aws.String("Parent = :parent"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":parent": &types.AttributeValueMemberS{Value: *parent},
		},
	}
	paginator := dynamodb.NewQueryPaginator(svc.DynamoDb, params)

	var result []RepositoryDto
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			return nil, svc.handleError(ctxId, err, "ListRepositoriesByParent",
				map[string]string{
					"parent": *parent,
				},
				zap.String("parent", *parent),
			)
		}

		var repsResp []RepositoryDto
		err = attributevalue.UnmarshalListOfMaps(page.Items, &repsResp)
		if err != nil {
			return nil, svc.handleError(ctxId, err, "ListRepositoriesByParent",
				map[string]string{
					"parent": *parent,
				},
				zap.String("parent", *parent),
			)
		}
		result = append(result, repsResp...)
	}

	svc.Logger.Debug(fmt.Sprintf("%s ListRepositoriesByParent() result", ctxId),
		zap.String("parent", *parent),
		zap.Reflect("result", &result),
	)

	return &result, nil
}
