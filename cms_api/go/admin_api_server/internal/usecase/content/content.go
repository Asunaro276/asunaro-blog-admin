package usecase

import (
	"admin/model"
)

// ContentUsecase はコンテンツ取得のユースケースインターフェースです
type getContent interface {
	GetContent() (*model.Content, error)
}

// contentUsecase はContentUsecaseの実装です
type contentUsecase struct {
	contentRepository getContent
}

// NewContentUsecase は新しいContentUsecaseインスタンスを作成します
func NewContentUsecase(getContent getContent) *contentUsecase {
	return &contentUsecase{
		contentRepository: getContent,
	}
}

// GetContent はコンテンツを取得します
func (u *contentUsecase) GetContent() (map[string]interface{}, error) {
	content, err := u.contentRepository.GetContent()
	if err != nil {
		return nil, err
	}

	// ドメインモデルをレスポンス用のマップに変換
	result := map[string]interface{}{
		"id":         content.ID,
		"title":      content.Title,
		"body":       content.Body,
		"author":     content.Author,
		"created_at": content.CreatedAt,
		"updated_at": content.UpdatedAt,
	}

	return result, nil
}
