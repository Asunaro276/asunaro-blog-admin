package infrastructure

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/testcontainers/testcontainers-go"
	dynamodbLocal "github.com/testcontainers/testcontainers-go/modules/dynamodb"
)

// DynamoDBContainer はDynamoDB Localコンテナを管理する構造体
type DynamoDBContainer struct {
	Container testcontainers.Container
	Client    *dynamodb.Client
	Endpoint  string
}

var dynamoDBContainer *DynamoDBContainer

// GetTestDynamoDBClient はテスト用のDynamoDBクライアントを返します
func GetTestDynamoDBClient() *DynamoDBClient {
	return NewDynamoDBClientFromContainer(dynamoDBContainer)
}

// SetupDynamoDBContainer はDynamoDB Localコンテナをセットアップします
func SetupDynamoDBContainer(ctx context.Context) (*DynamoDBContainer, error) {
	// DynamoDB Localコンテナの設定
	container, err := dynamodbLocal.Run(context.Background(), "amazon/dynamodb-local:latest")
	if err != nil {
		return nil, fmt.Errorf("DynamoDB Localコンテナの起動に失敗しました: %w", err)
	}

	// コンテナのエンドポイントを取得
	mappedPort, err := container.MappedPort(ctx, "8000")
	if err != nil {
		return nil, fmt.Errorf("ポートマッピングの取得に失敗しました: %w", err)
	}

	hostIP, err := container.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("ホストIPの取得に失敗しました: %w", err)
	}

	endpoint := fmt.Sprintf("http://%s:%s", hostIP, mappedPort.Port())

	// DynamoDBクライアントを設定
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("ap-northeast-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("dummy", "dummy", "dummy")),
	)
	if err != nil {
		return nil, fmt.Errorf("AWS設定の読み込みに失敗しました: %w", err)
	}

	client := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	})

	// コンテナとクライアントを保持
	dynamodbContainer := &DynamoDBContainer{
		Container: container,
		Client:    client,
		Endpoint:  endpoint,
	}

	return dynamodbContainer, nil
}

// CreateTable はDynamoDBテーブルを作成します
func (d *DynamoDBContainer) CreateTable(ctx context.Context, tableName string) error {
	// テーブル作成
	input := &dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("PK"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("SK"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("GSI1PK"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("GSI1SK"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("PK"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("SK"),
				KeyType:       types.KeyTypeRange,
			},
		},
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
			{
				IndexName: aws.String("GSI1"),
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("GSI1PK"),
						KeyType:       types.KeyTypeHash,
					},
					{
						AttributeName: aws.String("GSI1SK"),
						KeyType:       types.KeyTypeRange,
					},
				},
				Projection: &types.Projection{
					ProjectionType: types.ProjectionTypeAll,
				},
				ProvisionedThroughput: &types.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(5),
					WriteCapacityUnits: aws.Int64(5),
				},
			},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
	}

	_, err := d.Client.CreateTable(ctx, input)
	if err != nil {
		return fmt.Errorf("テーブル %s の作成に失敗しました: %w", tableName, err)
	}

	// テーブルが作成されるまで待機
	waiter := dynamodb.NewTableExistsWaiter(d.Client)
	err = waiter.Wait(ctx, &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	}, 5*time.Second)

	if err != nil {
		return fmt.Errorf("テーブル %s の作成完了を待機中にエラーが発生しました: %w", tableName, err)
	}

	log.Printf("テーブル %s が正常に作成されました", tableName)
	return nil
}

// Teardown はコンテナを停止して削除します
func (d *DynamoDBContainer) Teardown(ctx context.Context) error {
	return d.Container.Terminate(ctx)
}

// NewDynamoDBClientFromContainer はDynamoDBコンテナからDynamoDBClientを作成します
func NewDynamoDBClientFromContainer(container *DynamoDBContainer) *DynamoDBClient {
	return &DynamoDBClient{
		Client: container.Client,
	}
}
