package main

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestHandler_BasicRequest(t *testing.T) {
	// テスト用のリクエストを作成
	req := events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		Path:       "/healthcheck",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	// Lambda Handlerを実行
	response, err := Handler(context.Background(), req)

	// レスポンスを検証
	assert.NoError(t, err)
	assert.NotEqual(t, 0, response.StatusCode)
	assert.Contains(t, response.Headers, "Content-Type")
}

func TestHandler_InvalidRequest(t *testing.T) {
	// 無効なリクエストを作成
	req := events.APIGatewayProxyRequest{
		HTTPMethod: "INVALID",
		Path:       "/nonexistent",
	}

	// Lambda Handlerを実行
	response, err := Handler(context.Background(), req)

	// エラーハンドリングを検証
	// レスポンスが返却されることを確認（エラー処理が正常動作）
	assert.NoError(t, err) // Handler自体はエラーを返さない（内部処理でハンドリング）
	assert.NotEqual(t, 0, response.StatusCode)
}

func TestHandler_EmptyRequest(t *testing.T) {
	// 空のリクエストを作成
	req := events.APIGatewayProxyRequest{}

	// Lambda Handlerを実行
	response, err := Handler(context.Background(), req)

	// 基本的なレスポンスが返却されることを確認
	assert.NoError(t, err)
	assert.NotEqual(t, 0, response.StatusCode)
}
