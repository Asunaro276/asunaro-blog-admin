package usecase

import (
	"cms_api/internal/domain/entity"
)

type getContents interface {
	GetArticles() ([]model.Article, error)
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
func (u *contentUsecase) GetArticles() ([]model.Article, error) {
	contents, err := u.contentRepository.GetArticles()
	if err != nil {
		return nil, err
	}

	result := contents
	return result, nil
}
