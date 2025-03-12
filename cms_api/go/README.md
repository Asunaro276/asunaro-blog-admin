# CMS API

## テスト基盤

このプロジェクトでは、以下のテスト基盤を提供しています。

### テストの実行方法

```bash
# すべてのテストを実行
make test

# カバレッジレポート付きでテストを実行
make test-coverage

# 詳細なテスト結果を表示
make test-verbose
```

### モックの生成

```bash
# モックを生成
make mock
```

### テスト基盤の構成

- `internal/testutil`: テスト用のユーティリティ関数を提供
  - `testutil.go`: テスト用のヘルパー関数
  - `testmain.go`: テスト環境のセットアップとクリーンアップ
  - `mock_repository.go`: リポジトリのモック実装

### テーブル駆動テスト

テストは、テーブル駆動テストのパターンを使用して実装されています。これにより、複数のテストケースを簡潔に記述できます。

```go
tests := []struct {
    name           string
    input          string
    expectedOutput string
}{
    {
        name:           "ケース1",
        input:          "入力1",
        expectedOutput: "期待される出力1",
    },
    {
        name:           "ケース2",
        input:          "入力2",
        expectedOutput: "期待される出力2",
    },
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        // テストコード
    })
}
```

### モックの使用方法

リポジトリなどの外部依存をモック化するには、`testutil`パッケージのモック実装を使用します。

```go
// モックリポジトリをセットアップ
mockRepo := &testutil.MockContentRepository{}
mockRepo.GetContentFunc = func(id string) (string, error) {
    return "モックされたコンテンツ", nil
}

// テスト対象のユースケースを作成
uc := NewContentUseCase(mockRepo)

// テスト実行
// ...
``` 
