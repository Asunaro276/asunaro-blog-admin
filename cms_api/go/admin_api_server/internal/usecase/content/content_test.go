package usecase

import (
	"admin/model"
	"errors"
	"fmt"
	"math/rand/v2"
	"testing"
	"time"

	"admin_api_server/internal/controller/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type contentsUsecaseTestSuite struct {
	suite.Suite
	usecase        *contentUsecase
	mockRepository *mocks.GetContents
}

func randomContent(timeInt int64) *model.Content {
	return &model.Content{
		ID:        uuid.New().String(),
		Title:     fmt.Sprintf("test title %d", timeInt),
		Body:      fmt.Sprintf("test body %d", timeInt),
		Author:    fmt.Sprintf("test author %d", timeInt),
		CreatedAt: time.Unix(timeInt, 0),
		UpdatedAt: time.Unix(timeInt, 0),
	}
}

// TestContentsControllerを実行（テストメインエントリーポイント）
func TestContentsUsecase(t *testing.T) {
	suite.Run(t, new(contentsUsecaseTestSuite))
}

// 各テスト実行前のセットアップ
func (s *contentsUsecaseTestSuite) SetupSubTest() {
	s.mockRepository = mocks.NewGetContents(s.T())
	s.usecase = NewContentUsecase(s.mockRepository)
}

// GetContentsのテスト
func (s *contentsUsecaseTestSuite) TestGetContents() {
	randomInt := rand.Int64()
	content1 := randomContent(randomInt)
	content2 := randomContent(randomInt + 1)
	testCases := []struct {
		name          string
		setup         func()
		expectedData  []model.Content
		expectedError error
	}{
		{
			name: "正常系：コンテンツが正常に取得できる場合",
			setup: func() {
				contents := []model.Content{*content1, *content2}
				s.mockRepository.EXPECT().GetContents().Return(contents, nil)
			},
			expectedData: []model.Content{
				*content1,
				*content2,
			},
			expectedError: nil,
		},
		{
			name: "異常系：コンテンツ取得でエラーが発生する場合",
			setup: func() {
				s.mockRepository.EXPECT().GetContents().Return(nil, errors.New("取得エラー"))
			},
			expectedData:  nil,
			expectedError: errors.New("取得エラー"),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// テストケースのセットアップを実行
			tc.setup()

			// テスト対象のメソッドを実行
			contents, err := s.usecase.GetContents()

			// エラーのアサーション
			assert.Equal(s.T(), tc.expectedError, err)

			// データのアサーション
			assert.Equal(s.T(), len(tc.expectedData), len(contents))
			for i, content := range contents {
				assert.Equal(s.T(), tc.expectedData[i].ID, content.ID)
				assert.Equal(s.T(), tc.expectedData[i].Title, content.Title)
				assert.Equal(s.T(), tc.expectedData[i].Body, content.Body)
				assert.Equal(s.T(), tc.expectedData[i].Author, content.Author)
				assert.True(s.T(), content.CreatedAt.Equal(tc.expectedData[i].CreatedAt))
				assert.True(s.T(), content.UpdatedAt.Equal(tc.expectedData[i].UpdatedAt))
			}
		})
	}
}
