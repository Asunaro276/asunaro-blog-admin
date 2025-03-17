package repository

import (
	"admin/model"
	"time"
)

// contentRepository はContentRepositoryの実装です
type contentRepository struct {
	// ここにデータベース接続などの依存関係を追加することができます
	// db *sql.DB
}

// NewContentRepository は新しいContentRepositoryインスタンスを作成します
func NewContentRepository() *contentRepository {
	return &contentRepository{
		// ここでデータベース接続などの依存関係を注入することができます
		// db: db,
	}
}

// GetContent はコンテンツを取得します
func (cr *contentRepository) GetContent() (*model.Content, error) {
	// 実際の実装ではデータベースからデータを取得します
	// ここではサンプルデータを返します
	now := time.Now()
	content := &model.Content{
		ID:        "1",
		Title:     "サンプルコンテンツ",
		Body:      "これはサンプルコンテンツです。",
		Author:    "管理者",
		CreatedAt: now.Add(-24 * time.Hour),
		UpdatedAt: now,
	}

	return content, nil
}

// ListContents はコンテンツのリストを取得します
func (cr *contentRepository) ListContents() ([]*model.Content, error) {
	// 実際の実装ではデータベースからデータを取得します
	// ここではサンプルデータを返します
	now := time.Now()
	contents := []*model.Content{
		{
			ID:        "1",
			Title:     "サンプルコンテンツ1",
			Body:      "これはサンプルコンテンツ1です。",
			Author:    "管理者",
			CreatedAt: now.Add(-24 * time.Hour),
			UpdatedAt: now,
		},
		{
			ID:        "2",
			Title:     "サンプルコンテンツ2",
			Body:      "これはサンプルコンテンツ2です。",
			Author:    "管理者",
			CreatedAt: now.Add(-48 * time.Hour),
			UpdatedAt: now.Add(-24 * time.Hour),
		},
	}

	return contents, nil
}

// CreateContent は新しいコンテンツを作成します
func (cr *contentRepository) CreateContent(content *model.Content) error {
	// 実際の実装ではデータベースにデータを保存します
	return nil
}

// UpdateContent はコンテンツを更新します
func (cr *contentRepository) UpdateContent(content *model.Content) error {
	// 実際の実装ではデータベースのデータを更新します
	return nil
}

// DeleteContent はコンテンツを削除します
func (cr *contentRepository) DeleteContent(id string) error {
	// 実際の実装ではデータベースからデータを削除します
	return nil
}
