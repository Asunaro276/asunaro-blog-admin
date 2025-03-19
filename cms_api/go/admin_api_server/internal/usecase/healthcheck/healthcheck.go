package healthcheck

import (
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
