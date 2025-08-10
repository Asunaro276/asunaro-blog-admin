.PHONY: test test-coverage test-verbose mock clean run-local docker-build docker-run build-lambda build-standalone

# デフォルトのターゲット
all: test

# APIサーバーを実行
run:
	go run cmd/main.go

# Lambda Handlerを実行（テスト用）
run-lambda:
	go run cmd/lambda/main.go

# スタンドアロン版をビルド
build-standalone:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o bin/cms-api-standalone cmd/main.go

# Lambda用バイナリをビルド
build-lambda:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o bin/cms-api-lambda cmd/lambda/main.go

# すべてのバイナリをビルド
build: build-standalone build-lambda

# テストを実行
test:
	go test ./... -v

# モックを生成
mock:
	mockery
