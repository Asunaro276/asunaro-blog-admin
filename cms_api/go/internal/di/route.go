package route

import (
	"cms_api/internal/config"
	"cms_api/internal/infrastructure/controller"
	"cms_api/internal/infrastructure/database"
	"cms_api/internal/infrastructure/repository"
	usecase "cms_api/internal/usecase/content"
	"cms_api/internal/usecase/healthcheck"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// RouteHandler は設定を受け取ってEchoインスタンスを構築します
func RouteHandler(cfg *config.Config) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	// PostgreSQLデータベース接続の初期化
	postgresDB, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("PostgreSQL接続の初期化に失敗しました: %v", err)
	}

	// DynamoDB接続も保持（段階的移行のため）
	dynamoClient, err := repository.NewDynamoDBClient()
	if err != nil {
		log.Printf("DynamoDB接続の初期化に失敗しました（PostgreSQLで続行）: %v", err)
	}

	// リポジトリの初期化（現在はDynamoDBを使用、将来的にPostgreSQLに移行予定）
	contentRepository := repository.NewContentRepository(dynamoClient)

	// ユースケースの初期化
	contentUsecase := usecase.NewContentUsecase(contentRepository)

	// コントローラーの初期化
	contentController := controller.NewContentController(contentUsecase)

	// ルーティング設定
	e.GET("/", contentController.GetContent)
	e.GET("/healthcheck", func(c echo.Context) error {
		return healthcheck.HealthcheckWithDB(c, postgresDB)
	})

	return e
}
