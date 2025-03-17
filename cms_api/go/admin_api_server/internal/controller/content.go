package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ContentController struct {
	contentUsecase interface {
		GetContent() (map[string]interface{}, error)
	}
}

func NewContentController(cu interface {
	GetContent() (map[string]interface{}, error)
}) *ContentController {
	return &ContentController{
		contentUsecase: cu,
	}
}

// GetContent godoc
// @Summary Get content
// @Description Get content
// @Tags content
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router / [get]
func (cc *ContentController) GetContent(c echo.Context) error {
	content, err := cc.contentUsecase.GetContent()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, content)
}
