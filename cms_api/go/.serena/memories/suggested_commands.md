# 推奨コマンド

## 開発コマンド
- `make run` または `go run cmd/main.go`: APIサーバーを起動（:8080）
- `make test` または `go test ./... -v`: 全テストを実行
- `make mock` または `mockery`: モックファイルを生成

## コード品質チェック
- `golangci-lint run`: リンターを実行（go.modにツールとして定義済み）
- `go fmt ./...`: コードフォーマット
- `go mod tidy`: 依存関係の整理

## テスト関連
- `go test ./... -v`: 詳細なテスト実行
- `go test ./... -cover`: カバレッジ付きテスト実行
- `go test -run TestName`: 特定のテストを実行

## ビルド・デプロイ
- `go build cmd/main.go`: バイナリビルド
- `go build cmd/lambda/main.go`: Lambda用バイナリビルド

## 依存関係管理
- `go mod download`: 依存関係のダウンロード
- `go mod verify`: 依存関係の検証