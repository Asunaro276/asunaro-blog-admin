.PHONY: test test-coverage test-verbose mock clean

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
	go install github.com/vektra/mockery/v2@latest
	mockery

# 生成されたファイルをクリーンアップ
clean:
	rm -f coverage.out coverage.html
	rm -rf mocks/ 
