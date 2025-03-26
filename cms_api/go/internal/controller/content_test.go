package controller

import (
	"cms_api/internal/domain/entity"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"cms_api/internal/usecase/content/mocks"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

type contentsControllerTestSuite struct {
	suite.Suite
	echo        *echo.Echo
	controller  *ContentController
	mockUsecase *mocks.GetContents
}

// TestContentsControllerを実行（テストメインエントリーポイント）
func TestContentsController(t *testing.T) {
	suite.Run(t, new(contentsControllerTestSuite))
}

// スイート全体のセットアップ
func (s *contentsControllerTestSuite) SetupSuite() {
	s.echo = echo.New()
}

// 各サブテスト実行前のセットアップ
func (s *contentsControllerTestSuite) SetupSubTest() {
	s.mockUsecase = new(mocks.GetContents)
	s.controller = NewContentController(s.mockUsecase)
}

// セットアップ関数の型
type setupFunc func(s *contentsControllerTestSuite)

// 複数のセットアップ関数を実行
func (s *contentsControllerTestSuite) setup(fs ...setupFunc) {
	for _, f := range fs {
		f(s)
	}
}

// GetContentsのテスト
func (s *contentsControllerTestSuite) TestGetContents() {
	testCases := []struct {
		name           string
		setup          setupFunc
		expectedStatus int
		expectedBody   []model.Article
		expectError    bool
	}{
		{
			name: "正常系：コンテンツが正常に取得できる場合",
			setup: func(s *contentsControllerTestSuite) {
				contents := []model.Article{
					{
						ID:        "1",
						Title:     "テストタイトル1",
						Body:      "テスト本文1",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					{
						ID:        "2",
						Title:     "テストタイトル2",
						Body:      "テスト本文2",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
				}
				s.mockUsecase.EXPECT().GetArticles().Return(contents, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: []model.Article{
				{
					ID:     "1",
					Title:  "テストタイトル1",
					Body:   "テスト本文1",
				},
				{
					ID:     "2",
					Title:  "テストタイトル2",
					Body:   "テスト本文2",
				},
			},
			expectError: false,
		},
		{
			name: "異常系：コンテンツ取得でエラーが発生する場合",
			setup: func(s *contentsControllerTestSuite) {
				s.mockUsecase.EXPECT().GetArticles().Return(nil, errors.New("取得エラー"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   nil,
			expectError:    true,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// サブテストのセットアップ
			s.SetupSubTest()

			// テストケース固有のセットアップを実行
			s.setup(tc.setup)

			// リクエストとレコーダーを作成
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := s.echo.NewContext(req, rec)

			// テスト対象のメソッドを実行
			err := s.controller.GetContent(c)

			// アサーション
			if tc.expectError {
				// エラーケースのレスポンスを検証
				s.Equal(tc.expectedStatus, rec.Code)

				// JSONレスポンスがエラーを含むか検証
				var responseBody map[string]interface{}
				err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
				s.NoError(err)
				s.Contains(responseBody, "error")
			} else {
				// 正常系のレスポンスを検証
				s.NoError(err)
				s.Equal(tc.expectedStatus, rec.Code)

				// レスポンスボディを検証
				if tc.expectedBody != nil {
					var responseContents []model.Article
					err := json.Unmarshal(rec.Body.Bytes(), &responseContents)
					s.NoError(err)

					// 時間フィールドは動的に生成されるため、比較から除外する
					for i := range responseContents {
						responseContents[i].CreatedAt = time.Time{}
						responseContents[i].UpdatedAt = time.Time{}
					}
					for i := range tc.expectedBody {
						tc.expectedBody[i].CreatedAt = time.Time{}
						tc.expectedBody[i].UpdatedAt = time.Time{}
					}

					s.Equal(tc.expectedBody, responseContents)
				}
			}

			// モックの検証
			s.mockUsecase.AssertExpectations(s.T())
		})
	}
}
