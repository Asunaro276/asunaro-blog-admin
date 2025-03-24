package repository

import (
	"admin/infrastructure"
	"admin/model"
	"time"
)

// contentRepository はContentRepositoryの実装です
type contentRepository struct {
	dbClient  *infrastructure.DynamoDBClient
	tableName string
}

// ContentItem はDynamoDBに保存するコンテンツアイテムの構造体です
type ContentItem struct {
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
func NewContentRepository(dbClient *infrastructure.DynamoDBClient) *contentRepository {
	return &contentRepository{
		dbClient:  dbClient,
		tableName: "Contents",
	}
}

// GetArticles はコンテンツのリストを取得します
func (cr *contentRepository) GetArticles() ([]model.Article, error) {
	panic("not implemented")
}

// CreateContent は新しいコンテンツを作成します
func (cr *contentRepository) CreateContent(content *model.Article) error {
	panic("not implemented")
}

// UpdateContent はコンテンツを更新します
func (cr *contentRepository) UpdateContent(content *model.Article) error {
	panic("not implemented")
}

// DeleteContent はコンテンツを削除します
func (cr *contentRepository) DeleteContent(id string) error {
	panic("not implemented")
}
