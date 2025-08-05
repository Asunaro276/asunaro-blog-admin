# コードベース構造

## ディレクトリ構造
```
cms_api/go/
├── cmd/                    # エントリーポイント
│   ├── main.go            # APIサーバーメイン
│   └── lambda/main.go     # Lambda用エントリーポイント
├── internal/               # 内部パッケージ
│   ├── domain/entity/     # ドメインエンティティ
│   │   └── content.go     # Article構造体
│   ├── usecase/           # ユースケース層
│   │   ├── content/       # コンテンツ関連ユースケース
│   │   └── healthcheck/   # ヘルスチェック
│   ├── infrastructure/    # インフラ層
│   │   ├── repository/    # データアクセス
│   │   └── controller/    # HTTP処理
│   └── di/               # 依存性注入
│       └── route.go      # ルーティング設定
├── volume/               # ローカル開発用ボリューム
├── output/               # ビルド出力
├── go.mod               # Go modules設定
├── go.sum               # 依存関係ロック
├── Makefile             # ビルド・テストコマンド
├── Dockerfile           # コンテナ設定
├── .mockery.yaml        # モック生成設定
└── README.md            # プロジェクト説明
```

## 主要ファイル
- `cmd/main.go`: Echo サーバーを:8080で起動
- `internal/domain/entity/content.go`: Article構造体定義
- `internal/di/route.go`: ルーティング設定とDI
- `.mockery.yaml`: モック生成設定（testify/mock使用）

## テスト構成
- `*_test.go`: 単体テスト
- `mocks/`: mockeryで生成されたモック
- testcontainers使用でDynamoDB統合テスト