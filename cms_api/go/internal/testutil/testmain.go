package testutil

import (
	"os"
	"testing"
)

// TestMain は、すべてのテストの前後に実行される関数です。
// テスト環境のセットアップとクリーンアップを行います。
func TestMain(m *testing.M) {
	// テスト前の準備
	setup()

	// テストを実行
	code := m.Run()

	// テスト後のクリーンアップ
	teardown()

	// テスト結果のコードで終了
	os.Exit(code)
}

// setup はテスト環境をセットアップします。
func setup() {
	// テスト用のデータベース接続などをセットアップ
}

// teardown はテスト環境をクリーンアップします。
func teardown() {
	// テスト用のデータベース接続などをクリーンアップ
}
