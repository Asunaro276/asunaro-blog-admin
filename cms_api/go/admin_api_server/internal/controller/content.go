package controller

import (
	"admin/model"
	"net/http"

	"github.com/labstack/echo/v4"
)

type getContents interface {
	GetContents() ([]model.Content, error)
}

type ContentController struct {
	contentUsecase getContents
}

func NewContentController(cu getContents) *ContentController {
	return &ContentController{
		contentUsecase: cu,
	}
}

// GetContent godoc
// @Summary コンテンツ一覧の取得
// @Description コンテンツ一覧を取得します
// @Tags content
// @Accept json
// @Produce json
// @Success 200 {array} model.Content
// @Failure 500 {object} map[string]string
// @Router / [get]
func (cc *ContentController) GetContent(c echo.Context) error {
	contents, err := cc.contentUsecase.GetContents()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, contents)
}
