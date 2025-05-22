package route

import (
	"cms_api/internal/infrastructure/controller"
	"cms_api/internal/infrastructure/repository"
	usecase "cms_api/internal/usecase/content"
	"cms_api/internal/usecase/healthcheck"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RouteHandler() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	dbClient, _ := repository.NewDynamoDBClient()
	contentRepository := repository.NewContentRepository(dbClient)

	contentUsecase := usecase.NewContentUsecase(contentRepository)

	contentController := controller.NewContentController(contentUsecase)

	e.GET("/", contentController.GetContent)
	e.GET("/healthcheck", healthcheck.Healthcheck)

	return e
}
