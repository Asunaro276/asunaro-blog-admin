package getcontent

import (
	"net/http"
	"testing"

	"cms/internal/testutil"
)

func TestGetContent(t *testing.T) {
	// テストケースを定義
	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "正常系: コンテンツを取得できる",
			method:         http.MethodGet,
			path:           "/content",
			expectedStatus: http.StatusOK,
			expectedBody:   "Content",
		},
	}

	// 各テストケースを実行
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// テスト用のコンテキストを作成
			c, rec := testutil.CreateEchoContext(tt.method, tt.path, nil, nil, nil)

			// テスト対象の関数を呼び出し
			err := GetContent(c)

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
