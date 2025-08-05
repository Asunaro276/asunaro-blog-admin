# コードスタイルと規約

## 全般的なスタイル
- Go標準のコーディング規約に従う
- `gofmt`でフォーマット
- `golangci-lint`でリンティング（v2.2.1使用）

## 命名規約
- パッケージ名: 小文字、短く、わかりやすい名前
- 構造体: PascalCase（例：`Article`）
- フィールド: PascalCase（例：`CreatedAt`）
- 関数・メソッド: PascalCase（public）、camelCase（private）

## ディレクトリ構造
Clean Architectureパターン:
```
internal/
├── domain/entity/     # エンティティ
├── usecase/          # ユースケース
└── infrastructure/   # インフラ層
    ├── repository/   # データアクセス
    └── controller/   # HTTP処理
```

## JSON構造体タグ
- JSON出力用に`json`タグを必須で付与
- 例：`json:"created_at"`

## テスト規約
- テーブル駆動テストパターンを使用
- モックは`mockery`で自動生成
- テストファイル名：`*_test.go`
- Testcontainersを使用した統合テスト

## インポート順序
1. 標準ライブラリ
2. サードパーティライブラリ
3. 内部パッケージ