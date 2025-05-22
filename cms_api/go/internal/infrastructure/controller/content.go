package controller

import (
	model "cms_api/internal/domain/entity"
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

type getContents interface {
	GetArticles(ctx context.Context) ([]model.Article, error)
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
// @Success 200 {array} model.Article
// @Failure 500 {object} map[string]string
// @Router / [get]
func (cc *ContentController) GetContent(c echo.Context) error {
	contents, err := cc.contentUsecase.GetArticles(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, contents)
}
