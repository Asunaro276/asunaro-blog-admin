package main

import (
	"cms_api/internal/config"
	route "cms_api/internal/di"
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"
)

var echoLambda *echoadapter.EchoLambda

func init() {
	// Lambda用の初期化
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Lambda設定の読み込みに失敗しました: %v", err)
	}

	log.Printf("Lambda関数を初期化します: DB=%s:%d/%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)

	// Echo サーバーとルーティングの初期化
	e := route.RouteHandler(cfg)
	echoLambda = echoadapter.New(e)

	log.Printf("Lambda関数の初期化が完了しました")
}

func main() {
	lambda.Start(Handler)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Lambda Handler実行
	response, err := echoLambda.ProxyWithContext(ctx, req)
	if err != nil {
		log.Printf("Lambda Handler実行中にエラーが発生しました: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error": "Internal Server Error"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, err
	}
	return response, nil
}
