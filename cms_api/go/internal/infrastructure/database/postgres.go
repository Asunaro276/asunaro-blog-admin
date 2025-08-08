package database

import (
	"cms_api/internal/config"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

// PostgresDB はPostgreSQLデータベース接続を管理する構造体です
type PostgresDB struct {
	db     *sql.DB
	config *config.Config
}

// NewPostgresDB は新しいPostgreSQLデータベース接続を作成します
func NewPostgresDB(cfg *config.Config) (*PostgresDB, error) {
	dsn := cfg.GetDSN()

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("データベース接続の作成に失敗しました: %w", err)
	}

	// 接続プールの設定
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)

	// 接続テスト
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("データベースへのpingが失敗しました: %w", err)
	}

	log.Printf("PostgreSQL接続が確立されました: %s:%d/%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)

	return &PostgresDB{
		db:     db,
		config: cfg,
	}, nil
}

// GetDB はsql.DBインスタンスを返します
func (p *PostgresDB) GetDB() *sql.DB {
	return p.db
}

// Close はデータベース接続を閉じます
func (p *PostgresDB) Close() error {
	if p.db != nil {
		log.Println("PostgreSQL接続を閉じます")
		return p.db.Close()
	}
	return nil
}

// HealthCheck はデータベースの接続状態をチェックします
func (p *PostgresDB) HealthCheck() error {
	if p.db == nil {
		return fmt.Errorf("データベース接続が初期化されていません")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := p.db.PingContext(ctx); err != nil {
		return fmt.Errorf("データベースヘルスチェックに失敗しました: %w", err)
	}

	// 基本的なクエリテスト
	var count int
	query := "SELECT COUNT(*) FROM content_types WHERE is_active = true"
	if err := p.db.QueryRowContext(ctx, query).Scan(&count); err != nil {
		return fmt.Errorf("データベーステストクエリに失敗しました: %w", err)
	}

	log.Printf("データベースヘルスチェック成功: アクティブなコンテンツタイプ数 = %d", count)
	return nil
}

// GetStats はデータベース接続統計を返します
func (p *PostgresDB) GetStats() sql.DBStats {
	if p.db == nil {
		return sql.DBStats{}
	}
	return p.db.Stats()
}
