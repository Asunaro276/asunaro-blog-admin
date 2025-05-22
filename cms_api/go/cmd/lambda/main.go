package main

import (
	"cms_api/internal/di"
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"
)

var echoLambda *echoadapter.EchoLambda

func main() {
	lambda.Start(Handler)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	e := route.RouteHandler()
	echoLambda = echoadapter.New(e)
	return echoLambda.ProxyWithContext(ctx, req)
}
