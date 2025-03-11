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

func Healthcheck(c echo.Context) error {
	msg := &HealthcheckMessage{
		Status:  http.StatusOK,
		Message: "Success to connect echo",
	}
	res, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
	}
	return c.String(http.StatusOK, string(res))
}
