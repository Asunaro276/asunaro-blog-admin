package database

import (
	"cms_api/internal/config"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// PostgresDB はPostgreSQLデータベース接続を管理する構造体です
type PostgresDB struct {
	db     *gorm.DB
	config *config.Config
}

// NewPostgresDB は新しいPostgreSQLデータベース接続を作成します
func NewPostgresDB(cfg *config.Config) (*PostgresDB, error) {
	dsn := cfg.GetDSN()

	// Gormのロガー設定
	gormLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	// Gormデータベース接続を開く
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("データベース接続の作成に失敗しました: %w", err)
	}

	// 基盤となるsql.DBを取得して接続プールを設定
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("sql.DBの取得に失敗しました: %w", err)
	}

	// 接続プールの設定
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	sqlDB.SetConnMaxIdleTime(1 * time.Minute)

	// 接続テスト
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("データベースへのpingが失敗しました: %w", err)
	}

	log.Printf("PostgreSQL接続が確立されました: %s:%d/%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)

	return &PostgresDB{
		db:     db,
		config: cfg,
	}, nil
}

// GetDB はgorm.DBインスタンスを返します
func (p *PostgresDB) GetDB() *gorm.DB {
	return p.db
}

// Close はデータベース接続を閉じます
func (p *PostgresDB) Close() error {
	if p.db != nil {
		log.Println("PostgreSQL接続を閉じます")
		sqlDB, err := p.db.DB()
		if err != nil {
			return fmt.Errorf("sql.DBの取得に失敗しました: %w", err)
		}
		return sqlDB.Close()
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

	// 基盤となるsql.DBを取得してPingを実行
	sqlDB, err := p.db.DB()
	if err != nil {
		return fmt.Errorf("sql.DBの取得に失敗しました: %w", err)
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("データベースヘルスチェックに失敗しました: %w", err)
	}

	// Gormを使用した基本的なクエリテスト
	var count int64
	if err := p.db.WithContext(ctx).Table("content_types").Where("is_active = ?", true).Count(&count).Error; err != nil {
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
	
	sqlDB, err := p.db.DB()
	if err != nil {
		return sql.DBStats{}
	}
	
	return sqlDB.Stats()
}
