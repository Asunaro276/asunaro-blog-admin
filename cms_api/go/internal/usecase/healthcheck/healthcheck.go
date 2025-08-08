package healthcheck

import (
	"cms_api/internal/infrastructure/database"
	"encoding/json"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type HealthcheckMessage struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// Healthcheck はヘルスチェックエンドポイントのハンドラーです
func Healthcheck(c echo.Context) error {
	msg := &HealthcheckMessage{
		Status:  http.StatusOK,
		Message: "Success to connect echo server",
	}
	res, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Failed to marshal healthcheck message: %v", err)
	}
	return c.String(http.StatusOK, string(res))
}

// HealthcheckWithDB はデータベース接続を含むヘルスチェックを実行します
func HealthcheckWithDB(c echo.Context, db *database.PostgresDB) error {
	response := map[string]interface{}{
		"status":  "OK",
		"message": "CMS API is running",
	}

	// データベースヘルスチェック
	if err := db.HealthCheck(); err != nil {
		response["status"] = "ERROR"
		response["database"] = "DISCONNECTED"
		response["error"] = err.Error()
		return c.JSON(http.StatusServiceUnavailable, response)
	}

	// データベース統計情報を追加
	stats := db.GetStats()
	response["database"] = map[string]interface{}{
		"status":         "CONNECTED",
		"open_conns":     stats.OpenConnections,
		"in_use_conns":   stats.InUse,
		"idle_conns":     stats.Idle,
		"max_open_conns": stats.MaxOpenConnections,
	}

	return c.JSON(http.StatusOK, response)
}
