package repository

import (
	model "cms_api/internal/domain/entity"
	"cms_api/internal/infrastructure"
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// contentRepository はContentRepositoryの実装です
type ContentRepository interface {
	GetArticles(ctx context.Context) ([]model.Article, error)
	CreateContent(content *model.Article) error
	UpdateContent(content *model.Article) error
	DeleteContent(id string) error
}

type contentRepository struct {
	dbClient  *infrastructure.DynamoDBClient
	tableName string
}

// ContentItem はDynamoDBに保存するコンテンツアイテムの構造体です
type ArticleItem struct {
	PK           string    `dynamodbav:"PK"`
	SK           string    `dynamodbav:"SK"`
	Type         string    `dynamodbav:"type"`
	Title        string    `dynamodbav:"title"`
	Description  string    `dynamodbav:"description"`
	Body         string    `dynamodbav:"body"`
	CoverImage   string    `dynamodbav:"coverImage"`
	PublishedAt  time.Time `dynamodbav:"publishedAt"`
	UpdatedAt    time.Time `dynamodbav:"updatedAt"`
	Status       string    `dynamodbav:"status"`
	CategoryID   string    `dynamodbav:"categoryID"`
	CategoryName string    `dynamodbav:"categoryName"`
	TagName      string    `dynamodbav:"tagName"`
	ArticleCount int       `dynamodbav:"articleCount"`
	GSI1PK       string    `dynamodbav:"GSI1PK"`
	GSI1SK       string    `dynamodbav:"GSI1SK"`
}

// NewContentRepository は新しいContentRepositoryインスタンスを作成します
func NewContentRepository(dbClient *infrastructure.DynamoDBClient) ContentRepository {
	return &contentRepository{
		dbClient:  dbClient,
		tableName: "Contents",
	}
}

// GetArticles はコンテンツのリストを取得します
func (cr *contentRepository) GetArticles(ctx context.Context) ([]model.Article, error) {
	input := &dynamodb.ScanInput{
		TableName:        aws.String(cr.tableName),
		FilterExpression: aws.String("PK = :articleType"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":articleType": &types.AttributeValueMemberS{Value: "ARTICLE"},
		},
	}

	result, err := cr.dbClient.Client.Scan(ctx, input)
	if err != nil {
		return nil, err
	}

	articles := []model.Article{}
	for _, item := range result.Items {
		article := model.Article{}
		err = attributevalue.UnmarshalMap(item, &article)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}

	return articles, nil
}

// CreateContent は新しいコンテンツを作成します
func (cr *contentRepository) CreateArticle(content *model.Article) error {
	ctx := context.Background()

	item := ArticleItem{
		PK:          content.ID,
		SK:          "ARTICLE#" + content.ID,
		Type:        "ARTICLE",
		Title:       content.Title,
		Description: content.Description,
		Body:        content.Body,
		CoverImage:  content.CoverImage,
		Status:      content.Status,
		CategoryID:  content.CategoryID,
		UpdatedAt:   time.Now(),
		GSI1PK:      "ARTICLE",
		GSI1SK:      time.Now().Format(time.RFC3339),
	}

	if content.PublishedAt != "" {
		publishedTime, err := time.Parse(time.RFC3339, content.PublishedAt)
		if err == nil {
			item.PublishedAt = publishedTime
		}
	}

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(cr.tableName),
		Item:      av,
	}

	_, err = cr.dbClient.Client.PutItem(ctx, input)
	return err
}

// UpdateContent はコンテンツを更新します
func (cr *contentRepository) UpdateContent(content *model.Article) error {
	panic("not implemented")
}

// DeleteContent はコンテンツを削除します
func (cr *contentRepository) DeleteContent(id string) error {
	panic("not implemented")
}
