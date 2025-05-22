.PHONY: test test-coverage test-verbose mock clean run-local docker-build docker-run

# デフォルトのターゲット
all: test

# APIサーバーを実行
run:
	go run cmd/main.go

# テストを実行
test:
	go test ./... -v

# モックを生成
mock:
	mockery
