package repository

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// DynamoDBClient はDynamoDBとの接続を管理する構造体です
type dynamoDBClient struct {
	client *dynamodb.Client
}

// NewDynamoDBClient は新しいDynamoDBクライアントを作成します
func NewDynamoDBClient() (*dynamoDBClient, error) {
	// AWS SDK設定を読み込み
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("AWS設定の読み込みに失敗しました: %v", err)
		return nil, err
	}

	// DynamoDBクライアントを作成
	client := dynamodb.NewFromConfig(cfg)

	return &dynamoDBClient{
		client: client,
	}, nil
}

// NewDynamoDBClientWithEndpoint はカスタムエンドポイントを使用した新しいDynamoDBクライアントを作成します
// ローカル開発やテスト環境で使用します
func NewDynamoDBClientWithEndpoint(region, endpoint string) (*dynamoDBClient, error) {
	// AWS SDK設定を読み込み
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
	)
	if err != nil {
		log.Printf("AWS設定の読み込みに失敗しました: %v", err)
		return nil, err
	}

	// カスタムエンドポイントでDynamoDBクライアントを作成
	client := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	})

	return &dynamoDBClient{
		client: client,
	}, nil
}

// GetClient はDynamoDBクライアントを返します
func (d *dynamoDBClient) GetClient() *dynamodb.Client {
	return d.client
}
