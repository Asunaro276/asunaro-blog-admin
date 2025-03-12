package testutil

import (
	"net/http/httptest"

	"github.com/labstack/echo/v4"
)

// CreateEchoContext は、テスト用のecho.Contextを作成するヘルパー関数です。
// method: HTTPメソッド
// path: リクエストパス
// body: リクエストボディ（省略可能）
// pathParams: パスパラメータのマップ（省略可能）
// queryParams: クエリパラメータのマップ（省略可能）
func CreateEchoContext(method, path string, body []byte, pathParams map[string]string, queryParams map[string]string) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(method, path, nil)
	if body != nil {
		req = httptest.NewRequest(method, path, nil)
	}

	// クエリパラメータを設定
	if queryParams != nil {
		q := req.URL.Query()
		for key, value := range queryParams {
			q.Add(key, value)
		}
		req.URL.RawQuery = q.Encode()
	}

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// パスパラメータを設定
	if pathParams != nil {
		for key, value := range pathParams {
			c.SetParamNames(key)
			c.SetParamValues(value)
		}
	}

	return c, rec
}

// AssertStatusCode は、レスポンスのステータスコードを検証するヘルパー関数です。
func AssertStatusCode(t interface {
	Errorf(format string, args ...interface{})
}, rec *httptest.ResponseRecorder, expected int) {
	if rec.Code != expected {
		t.Errorf("expected status code %d but got %d", expected, rec.Code)
	}
}

// AssertResponseBody は、レスポンスボディを検証するヘルパー関数です。
func AssertResponseBody(t interface {
	Errorf(format string, args ...interface{})
}, rec *httptest.ResponseRecorder, expected string) {
	if rec.Body.String() != expected {
		t.Errorf("expected response body %q but got %q", expected, rec.Body.String())
	}
}
