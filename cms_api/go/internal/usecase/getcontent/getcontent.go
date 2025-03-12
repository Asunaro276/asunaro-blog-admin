package getcontent

import (
	"net/http"

	"github.com/labstack/echo/v4"
)


func GetContent(c echo.Context) error {
	return c.String(http.StatusOK, string("Content"))
}
