.PHONY: test test-coverage test-verbose mock clean run-local docker-build docker-run

# デフォルトのターゲット
all: test

# テストを実行
test:
	go test ./... -v

# カバレッジレポート付きでテストを実行
test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

# 詳細なテスト結果を表示
test-verbose:
	go test ./... -v

# モックを生成
mock:
	mockery

# 生成されたファイルをクリーンアップ
clean:
	rm -f coverage.out coverage.html
	rm -rf mocks/ 

# ローカルでAPIサーバーを実行
run-local:
	cd admin_api_server && go run cmd/local/main.go

# Dockerイメージをビルド
docker-build:
	docker build -t cms-api-local .

# Dockerコンテナを実行
docker-run:
	docker run -p 8080:8080 cms-api-local 
