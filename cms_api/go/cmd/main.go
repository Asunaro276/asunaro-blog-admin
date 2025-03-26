package main

import (
	"cms_api/internal/infrastructure"
	"cms_api/internal/infrastructure/repository"
	"cms_api/internal/controller"
	usecase "cms_api/internal/usecase/content"
	"cms_api/internal/usecase/healthcheck"
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var echoLambda *echoadapter.EchoLambda

func init() {
	e := echo.New()
	e.Use(middleware.Recover())

	dbClient, _ := infrastructure.NewDynamoDBClient()
	// リポジトリの初期化
	contentRepository := repository.NewContentRepository(dbClient)

	// ユースケースの初期化
	contentUsecase := usecase.NewContentUsecase(contentRepository)

	// コントローラーの初期化
	contentController := controller.NewContentController(contentUsecase)

	// ルーティング
	e.GET("/", contentController.GetContent)
	e.GET("/healthcheck", healthcheck.Healthcheck)

	echoLambda = echoadapter.New(e)
}

func main() {
	lambda.Start(Handler)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return echoLambda.ProxyWithContext(ctx, req)
}
