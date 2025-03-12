package content

import (
	"net/http"

	"cms/internal/repository"

	"github.com/labstack/echo/v4"
)

// ContentUseCase はコンテンツユースケースの構造体です。
type ContentUseCase struct {
	repo repository.ContentRepository
}

// NewContentUseCase は新しいContentUseCaseを作成します。
func NewContentUseCase(repo repository.ContentRepository) *ContentUseCase {
	return &ContentUseCase{
		repo: repo,
	}
}

// GetContent はコンテンツを取得するハンドラーです。
func (uc *ContentUseCase) GetContent(c echo.Context) error {
	// パスパラメータからIDを取得
	id := c.Param("id")
	if id == "" {
		id = "default"
	}

	// リポジトリからコンテンツを取得
	content, err := uc.repo.GetContent(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get content",
		})
	}

	return c.String(http.StatusOK, content)
}
