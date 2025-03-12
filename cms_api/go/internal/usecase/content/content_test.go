package content

import (
	"errors"
	"net/http"
	"testing"

	"cms/internal/testutil"
)

func TestContentUseCase_GetContent(t *testing.T) {
	// テストケースを定義
	tests := []struct {
		name           string
		id             string
		setupMock      func() *testutil.MockContentRepository
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "正常系: コンテンツを取得できる",
			id:   "1",
			setupMock: func() *testutil.MockContentRepository {
				mock := &testutil.MockContentRepository{}
				mock.GetContentFunc = func(id string) (string, error) {
					return "Content for ID: " + id, nil
				}
				return mock
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "Content for ID: 1",
		},
		{
			name: "異常系: リポジトリからのエラー",
			id:   "error",
			setupMock: func() *testutil.MockContentRepository {
				mock := &testutil.MockContentRepository{}
				mock.GetContentFunc = func(id string) (string, error) {
					return "", errors.New("repository error")
				}
				return mock
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"Failed to get content"}`,
		},
	}

	// 各テストケースを実行
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックリポジトリをセットアップ
			mockRepo := tt.setupMock()

			// テスト対象のユースケースを作成
			uc := NewContentUseCase(mockRepo)

			// テスト用のコンテキストを作成
			pathParams := map[string]string{"id": tt.id}
			c, rec := testutil.CreateEchoContext(http.MethodGet, "/content/"+tt.id, nil, pathParams, nil)

			// テスト対象の関数を呼び出し
			err := uc.GetContent(c)

			// エラーがないことを確認
			if err != nil {
				t.Errorf("GetContent returned an error: %v", err)
			}

			// ステータスコードが期待通りであることを確認
			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status code %d but got %d", tt.expectedStatus, rec.Code)
			}

			// レスポンスボディが期待通りであることを確認
			if rec.Body.String() != tt.expectedBody {
				t.Errorf("expected response body %q but got %q", tt.expectedBody, rec.Body.String())
			}
		})
	}
}
