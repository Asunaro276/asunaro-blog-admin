package usecase

import (
	"admin/model"
)

type getContents interface {
	GetContents() ([]model.Content, error)
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
func (u *contentUsecase) GetContents() ([]model.Content, error) {
	contents, err := u.contentRepository.GetContents()
	if err != nil {
		return nil, err
	}

	result := contents
	return result, nil
}
