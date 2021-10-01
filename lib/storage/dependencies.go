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

func (svc *Storage) ListDependenciesByParent(ctxId string, parent *string) (*[]DependencyDto, *StorageErrorRest) {
	if parent == nil {
		parent = ptr.String(RootParent)
	}

	svc.Logger.Debug(fmt.Sprintf("%s ListDependenciesByParent() called", ctxId),
		zap.String("parent", *parent),
	)

	var consistentRead = false
	params := &dynamodb.QueryInput{
		TableName:              svc.Config.DependenciesTableName,
		ConsistentRead:         &consistentRead,
		KeyConditionExpression: aws.String("Parent = :parent"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":parent": &types.AttributeValueMemberS{Value: *parent},
		},
	}
	paginator := dynamodb.NewQueryPaginator(svc.DynamoDb, params)

	var result []DependencyDto
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			return nil, svc.handleError(ctxId, err, "ListDependenciesByParent",
				map[string]string{
					"parent": *parent,
				},
				zap.String("parent", *parent),
			)
		}

		var depsResp []DependencyDto
		err = attributevalue.UnmarshalListOfMaps(page.Items, &depsResp)
		if err != nil {
			return nil, svc.handleError(ctxId, err, "ListDependenciesByParent",
				map[string]string{
					"parent": *parent,
				},
				zap.String("parent", *parent),
			)
		}
		result = append(result, depsResp...)
	}

	svc.Logger.Debug(fmt.Sprintf("%s ListDependenciesByParent() result", ctxId),
		zap.String("parent", *parent),
		zap.Reflect("result", &result),
	)

	return &result, nil
}

func (svc *Storage) ListDependenciesByRepo(ctxId string, repo string, ref string) (*[]StorageDto, *StorageErrorRest) {
	svc.Logger.Debug(fmt.Sprintf("%s ListDependenciesByRepo() called", ctxId),
		zap.String("repo", repo),
		zap.String("ref", ref),
	)

	var consistentRead = false
	params := &dynamodb.QueryInput{
		TableName:              svc.Config.StorageTableName,
		IndexName:				ptr.String("Repository"),
		ConsistentRead:         &consistentRead,
		KeyConditionExpression: aws.String("#repo = :repo and #ref = :ref"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":repo": &types.AttributeValueMemberS{Value: repo},
			":ref":  &types.AttributeValueMemberS{Value: ref},
		},
		ExpressionAttributeNames: map[string]string{
			"#repo": "Repo",
			"#ref": "Ref",
		},
		Select: types.SelectAllAttributes,
	}
	paginator := dynamodb.NewQueryPaginator(svc.DynamoDb, params)

	var result []StorageDto
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			return nil, svc.handleError(ctxId, err, "ListDependenciesByRepo",
				map[string]string{
					"repo": repo,
					"ref":  ref,
				},
				zap.String("repo", repo),
				zap.String("ref", ref),
			)
		}

		var depsResp []StorageDto
		err = attributevalue.UnmarshalListOfMaps(page.Items, &depsResp)
		if err != nil {
			return nil, svc.handleError(ctxId, err, "ListDependenciesByRepo",
				map[string]string{
					"repo": repo,
					"ref":  ref,
				},
				zap.String("repo", repo),
				zap.String("ref", ref),
			)
		}
		result = append(result, depsResp...)
	}

	svc.Logger.Debug(fmt.Sprintf("%s ListDependenciesByRepo() result", ctxId),
		zap.String("repo", repo),
		zap.String("ref", ref),
		zap.Reflect("result", &result),
	)

	return &result, nil
}
