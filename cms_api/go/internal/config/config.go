package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
)

// Config はアプリケーション設定を管理する構造体です
type Config struct {
	Server   ServerConfig   `koanf:"server"`
	Database DatabaseConfig `koanf:"database"`
	AWS      AWSConfig      `koanf:"aws"`
}

// ServerConfig はサーバー関連の設定を管理します
type ServerConfig struct {
	Host string `koanf:"host"`
	Port string `koanf:"port"`
}

// DatabaseConfig はデータベース関連の設定を管理します
type DatabaseConfig struct {
	Host     string `koanf:"host"`
	Port     int    `koanf:"port"`
	User     string `koanf:"user"`
	Password string `koanf:"password"`
	DBName   string `koanf:"dbname"`
	SSLMode  string `koanf:"sslmode"`
}

// AWSConfig はAWS関連の設定を管理します
type AWSConfig struct {
	Region string `koanf:"region"`
}

// DefaultConfig はデフォルト設定を返します
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Host: "0.0.0.0",
			Port: "8080",
		},
		Database: DatabaseConfig{
			Host:    "localhost",
			Port:    5432,
			User:    "postgres",
			DBName:  "cms_api",
			SSLMode: "disable",
		},
		AWS: AWSConfig{
			Region: "ap-northeast-1",
		},
	}
}

// LoadConfig は環境変数からアプリケーション設定を読み込みます
func LoadConfig() (*Config, error) {
	// koanfインスタンスを作成
	k := koanf.New(".")

	// デフォルト設定をロード
	cfg := DefaultConfig()

	// 環境変数プロバイダーでkoanfを設定
	// 環境変数プレフィックスは "CMS_API_" を使用
	if err := k.Load(env.Provider("CMS_API_", ".", func(s string) string {
		// CMS_API_SERVER_PORT -> server.port のように変換
		return strings.ToLower(strings.Replace(s, "_", ".", -1))
	}), nil); err != nil {
		log.Printf("環境変数の読み込みでエラーが発生しました: %v", err)
		return cfg, nil // エラーがあってもデフォルト設定で続行
	}

	// 設定構造体にアンマーシャル
	if err := k.Unmarshal("", cfg); err != nil {
		return nil, fmt.Errorf("設定の解析に失敗しました: %w", err)
	}

	// 必須設定の検証
	if err := validateConfig(cfg); err != nil {
		return nil, fmt.Errorf("設定の検証に失敗しました: %w", err)
	}

	log.Printf("設定の読み込みが完了しました: Server=%s:%s, DB=%s:%d/%s",
		cfg.Server.Host, cfg.Server.Port,
		cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)

	return cfg, nil
}

// validateConfig は設定の妥当性を検証します
func validateConfig(cfg *Config) error {
	if cfg.Server.Port == "" {
		return fmt.Errorf("サーバーポートが設定されていません")
	}

	if cfg.Database.Host == "" {
		return fmt.Errorf("データベースホストが設定されていません")
	}

	if cfg.Database.User == "" {
		return fmt.Errorf("データベースユーザーが設定されていません")
	}

	if cfg.Database.DBName == "" {
		return fmt.Errorf("データベース名が設定されていません")
	}

	return nil
}

// GetDSN はPostgreSQLのデータソース名を生成します
func (cfg *Config) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode)
}
