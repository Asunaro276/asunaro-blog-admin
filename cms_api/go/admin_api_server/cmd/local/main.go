package main

import (
	"admin/infrastructure/repository"
	"admin_api_server/internal/controller"
	usecase "admin_api_server/internal/usecase/content"
	"admin_api_server/internal/usecase/healthcheck"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// ローカル開発用のサーバーを起動する
func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// リポジトリの初期化
	contentRepo := repository.NewContentRepository()

	// ユースケースの初期化
	contentUsecase := usecase.NewContentUsecase(contentRepo)

	// コントローラーの初期化
	contentController := controller.NewContentController(contentUsecase)

	// ルーティング
	e.GET("/", contentController.GetContent)
	e.GET("/healthcheck", healthcheck.Healthcheck)

	// サーバー起動
	log.Println("ローカルサーバーを起動しています... http://localhost:8080")
	e.Logger.Fatal(e.Start(":8080"))
}
