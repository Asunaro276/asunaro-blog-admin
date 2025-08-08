package main

import (
	"cms_api/internal/config"
	"cms_api/internal/di"
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
		log.Fatalf("設定の読み込みに失敗しました: %v", err)
	}

	log.Printf("Lambda関数を初期化します: DB=%s:%d/%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)

	e := route.RouteHandler(cfg)
	echoLambda = echoadapter.New(e)
}

func main() {
	lambda.Start(Handler)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return echoLambda.ProxyWithContext(ctx, req)
}
