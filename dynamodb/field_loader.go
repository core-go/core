package dynamodb

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

type FieldLoader struct {
	database  *dynamodb.DynamoDB
	tableName string
	name      string
}

func NewFieldLoader(db *dynamodb.DynamoDB, tableName string, name string) *FieldLoader {
	return &FieldLoader{
		database:  db,
		tableName: tableName,
		name:      name,
	}
}

func (l *FieldLoader) Values(ctx context.Context, ids []string) ([]string, error) {
	var array []string
	projection := expression.NamesList(expression.Name(l.name))
	var filterConditions *expression.ConditionBuilder
	for _, id := range ids {
		c := expression.Equal(expression.Name(l.name), expression.Value(id))
		if filterConditions == nil {
			filterConditions = &c
		} else {
			and := filterConditions.Or(c)
			filterConditions = &and
		}
	}
	expr, err := expression.NewBuilder().
		WithFilter(*filterConditions).
		WithProjection(projection).
		Build()
	if err != nil {
		return nil, err
	}

	query := &dynamodb.ScanInput{
		TableName:                 aws.String(l.tableName),
		ProjectionExpression:      expr.Projection(),
		FilterExpression:          expr.Filter(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	}

	if err != nil {
		return nil, err
	}
	output, err := l.database.ScanWithContext(ctx, query)
	if err != nil {
		return array, err
	}
	if len(output.Items) == 0 {
		return array, nil
	}
	var result []map[string]interface{}
	err = dynamodbattribute.UnmarshalListOfMaps(output.Items, &result)
	if err != nil {
		return array, err
	}
	for _, model := range result {
		array = append(array, model[l.name].(string))
	}
	return array, err
}
