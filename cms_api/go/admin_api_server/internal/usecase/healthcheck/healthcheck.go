package healthcheck

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Healthcheck はヘルスチェックエンドポイントのハンドラーです
func Healthcheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "ok",
	})
}
