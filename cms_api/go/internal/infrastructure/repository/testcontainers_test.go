package repository

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
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	dynamodbLocal "github.com/testcontainers/testcontainers-go/modules/dynamodb"
)

type dynamoDBContainer struct {
	container testcontainers.Container
	client    *dynamodb.Client
	endpoint  string
}


type dynamodbTestcontainersTestSuite struct {
	suite.Suite
	dynamoContainer *dynamoDBContainer
	ctx             context.Context
	contentRepository *contentRepository
}

// SetupDynamoDBContainer はDynamoDB Localコンテナをセットアップします
func setupDynamoDBContainer(ctx context.Context) (*dynamoDBContainer, error) {
	// DynamoDB Localコンテナの設定
	container, err := dynamodbLocal.Run(ctx, "amazon/dynamodb-local:latest")
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
	dynamodbContainer := &dynamoDBContainer{
		container: container,
		client:    client,
		endpoint:  endpoint,
	}

	return dynamodbContainer, nil
}

// CreateTable はDynamoDBテーブルを作成します
func (d *dynamoDBContainer) createTable(ctx context.Context, tableName string) error {
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

	_, err := d.client.CreateTable(ctx, input)
	if err != nil {
		return fmt.Errorf("テーブル %s の作成に失敗しました: %w", tableName, err)
	}

	// テーブルが作成されるまで待機
	waiter := dynamodb.NewTableExistsWaiter(d.client)
	err = waiter.Wait(ctx, &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	}, 5*time.Second)

	if err != nil {
		return fmt.Errorf("テーブル %s の作成完了を待機中にエラーが発生しました: %w", tableName, err)
	}

	log.Printf("テーブル %s が正常に作成されました", tableName)
	return nil
}

func (s *dynamodbTestcontainersTestSuite) SetupSuite() {
	s.ctx = context.Background()
	container, err := setupDynamoDBContainer(s.ctx)
	if err != nil {
		s.T().Fatalf("DynamoDB Localコンテナのセットアップに失敗しました: %v", err)
	}
	s.dynamoContainer = container
	s.contentRepository = NewContentRepository(&dynamoDBClient{
		client: container.client,
	})
	err = container.createTable(s.ctx, "Contents")
	if err != nil {
		s.T().Fatalf("テーブルの作成に失敗しました: %v", err)
	}
}

func (d *dynamodbTestcontainersTestSuite) TearDownSuite() {
	if d.dynamoContainer != nil {
		err := d.dynamoContainer.container.Terminate(d.ctx)
		if err != nil {
			d.T().Fatalf("DynamoDB Localコンテナの停止に失敗しました: %v", err)
		}
	}
}
