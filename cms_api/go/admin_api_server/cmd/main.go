package main

import (
	"admin/infrastructure/repository"
	"admin_api_server/internal/controller"
	usecase "admin_api_server/internal/usecase/content"
	"admin_api_server/internal/usecase/healthcheck"
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var echoLambda *echoadapter.EchoLambda

func init() {
	log.Printf("echo cold start")

	e := echo.New()
	e.Use(middleware.Recover())

	// リポジトリの初期化
	contentRepository := repository.NewContentRepository()

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
