package repository

import (
	model "cms_api/internal/domain/entity"
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type contentRepository struct {
	dbClient  *dynamoDBClient
	tableName string
}

// ContentItem はDynamoDBに保存するコンテンツアイテムの構造体です
type ArticleItem struct {
	PK           string    `dynamodbav:"PK"`
	SK           string    `dynamodbav:"SK"`
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
func NewContentRepository(dbClient *dynamoDBClient) *contentRepository {
	return &contentRepository{
		dbClient:  dbClient,
		tableName: "Contents",
	}
}

// GetArticles はコンテンツのリストを取得します
func (cr *contentRepository) GetArticles(ctx context.Context) ([]model.Article, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(cr.tableName),
		IndexName:              aws.String("GSI1"),
		KeyConditionExpression: aws.String("GSI1PK = :articleType"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":articleType": &types.AttributeValueMemberS{Value: "ARTICLE"},
		},
	}

	result, err := cr.dbClient.client.Query(ctx, input)
	if err != nil {
		return nil, err
	}

	articles := []model.Article{}
	for _, item := range result.Items {
		var articleItem ArticleItem
		err = attributevalue.UnmarshalMap(item, &articleItem)
		if err != nil {
			return nil, err
		}

		article := model.Article{
			ID:          articleItem.PK,
			Title:       articleItem.Title,
			Description: articleItem.Description,
			Body:        articleItem.Body,
			CoverImage:  articleItem.CoverImage,
			PublishedAt: articleItem.PublishedAt,
			UpdatedAt:   articleItem.UpdatedAt,
			Status:      articleItem.Status,
			CategoryID:  articleItem.CategoryID,
			// Tagsフィールドは現在のArticleItemには含まれていないため設定しない
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
		SK:          "a#" + content.ID,
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

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(cr.tableName),
		Item:      av,
	}

	_, err = cr.dbClient.client.PutItem(ctx, input)
	return err
}

// UpdateContent はコンテンツを更新します
func (cr *contentRepository) UpdateArticle(content *model.Article) error {
	panic("not implemented")
}

// DeleteArticle はコンテンツを削除します
func (cr *contentRepository) DeleteArticle(id string) error {
	panic("not implemented")
}
