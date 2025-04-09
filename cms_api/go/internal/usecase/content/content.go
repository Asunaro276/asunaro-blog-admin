package usecase

import (
	model "cms_api/internal/domain/entity"
	"context"
)

type getContents interface {
	GetArticles(ctx context.Context) ([]model.Article, error)
}

type contentUsecase struct {
	contentRepository getContents
}

// NewContentUsecase は新しいContentUsecaseインスタンスを作成します
func NewContentUsecase(getContents getContents) *contentUsecase {
	return &contentUsecase{
		contentRepository: getContents,
	}
}

// GetContent はコンテンツを取得します
func (u *contentUsecase) GetArticles(ctx context.Context) ([]model.Article, error) {
	contents, err := u.contentRepository.GetArticles(ctx)
	if err != nil {
		return nil, err
	}

	result := contents
	return result, nil
}
