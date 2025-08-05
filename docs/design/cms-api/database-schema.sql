-- =============================================================================
-- CMS API データベーススキーマ設計
-- Aurora PostgreSQL 15 対応
-- =============================================================================

-- 拡張機能の有効化
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- =============================================================================
-- コンテンツテーブル
-- =============================================================================

/**
 * コンテンツマスターテーブル
 * CMSの中核となるコンテンツ情報を格納
 */
CREATE TABLE contents (
    -- 基本情報
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(500) NOT NULL,
    body TEXT NOT NULL,
    
    -- ステータス管理
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    
    -- 日時管理
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    published_at TIMESTAMP WITH TIME ZONE,
    
    -- 作成者情報
    author_id VARCHAR(100) NOT NULL,
    
    -- メタデータ（JSONB形式）
    metadata JSONB DEFAULT '{}',
    
    -- 制約
    CONSTRAINT chk_contents_status 
        CHECK (status IN ('draft', 'published', 'archived')),
    CONSTRAINT chk_contents_title_length 
        CHECK (char_length(title) >= 1 AND char_length(title) <= 500),
    CONSTRAINT chk_contents_body_length 
        CHECK (char_length(body) >= 1),
    CONSTRAINT chk_contents_published_at_when_published
        CHECK (
            (status = 'published' AND published_at IS NOT NULL) OR 
            (status != 'published')
        )
);

-- =============================================================================
-- インデックス設計
-- =============================================================================

-- 基本検索用インデックス
CREATE INDEX idx_contents_status ON contents(status);
CREATE INDEX idx_contents_author_id ON contents(author_id);
CREATE INDEX idx_contents_created_at ON contents(created_at DESC);
CREATE INDEX idx_contents_updated_at ON contents(updated_at DESC);
CREATE INDEX idx_contents_published_at ON contents(published_at DESC) 
    WHERE published_at IS NOT NULL;

-- 複合インデックス（よく使われる組み合わせ）
CREATE INDEX idx_contents_status_published_at ON contents(status, published_at DESC) 
    WHERE status = 'published';
CREATE INDEX idx_contents_status_created_at ON contents(status, created_at DESC);
CREATE INDEX idx_contents_author_status ON contents(author_id, status);

-- 全文検索用インデックス
CREATE INDEX idx_contents_title_gin ON contents USING GIN (title gin_trgm_ops);
CREATE INDEX idx_contents_body_gin ON contents USING GIN (body gin_trgm_ops);

-- JSONBメタデータ用インデックス
CREATE INDEX idx_contents_metadata_gin ON contents USING GIN (metadata);
CREATE INDEX idx_contents_metadata_category ON contents USING GIN ((metadata->>'category'));
CREATE INDEX idx_contents_metadata_tags ON contents USING GIN ((metadata->'tags'));

-- =============================================================================
-- トリガー関数
-- =============================================================================

/**
 * updated_at自動更新トリガー関数
 */
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

/**
 * published_at自動設定トリガー関数
 */
CREATE OR REPLACE FUNCTION set_published_at_on_publish()
RETURNS TRIGGER AS $$
BEGIN
    -- ステータスがdraftからpublishedに変更された場合
    IF OLD.status != 'published' AND NEW.status = 'published' AND NEW.published_at IS NULL THEN
        NEW.published_at = CURRENT_TIMESTAMP;
    END IF;
    
    -- ステータスがpublishedから他に変更された場合はpublished_atをクリア
    IF OLD.status = 'published' AND NEW.status != 'published' THEN
        NEW.published_at = NULL;
    END IF;
    
    RETURN NEW;
END;
$$ language 'plpgsql';

-- =============================================================================
-- トリガー設定
-- =============================================================================

-- updated_at自動更新トリガー
CREATE TRIGGER trigger_contents_updated_at
    BEFORE UPDATE ON contents
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- published_at自動設定トリガー
CREATE TRIGGER trigger_contents_published_at
    BEFORE INSERT OR UPDATE ON contents
    FOR EACH ROW
    EXECUTE FUNCTION set_published_at_on_publish();

-- =============================================================================
-- ビュー定義
-- =============================================================================

/**
 * 公開コンテンツビュー
 * 公開状態のコンテンツのみを抽出
 */
CREATE VIEW published_contents AS
SELECT 
    id,
    title,
    body,
    status,
    created_at,
    updated_at,
    published_at,
    author_id,
    metadata,
    -- 計算フィールド
    metadata->>'seoTitle' as seo_title,
    metadata->>'seoDescription' as seo_description,
    metadata->>'excerpt' as excerpt,
    metadata->>'category' as category,
    metadata->'tags' as tags,
    metadata->>'featuredImageUrl' as featured_image_url
FROM contents
WHERE status = 'published' 
    AND published_at IS NOT NULL 
    AND published_at <= CURRENT_TIMESTAMP;

/**
 * コンテンツサマリービュー
 * 一覧表示用の軽量ビュー
 */
CREATE VIEW content_summaries AS
SELECT 
    id,
    title,
    status,
    created_at,
    updated_at,
    published_at,
    author_id,
    -- メタデータから抽出
    COALESCE(metadata->>'excerpt', LEFT(body, 200) || '...') as excerpt,
    COALESCE(metadata->>'category', 'uncategorized') as category,
    COALESCE(metadata->'tags', '[]'::jsonb) as tags,
    metadata->>'featuredImageUrl' as featured_image_url
FROM contents;

-- =============================================================================
-- パフォーマンス最適化
-- =============================================================================

-- 統計情報更新の自動化設定
ALTER TABLE contents SET (autovacuum_analyze_scale_factor = 0.05);
ALTER TABLE contents SET (autovacuum_vacuum_scale_factor = 0.1);

-- パーティション設定（将来的な大容量対応）
-- 作成日時ベースでの月次パーティション（コメントアウト - 必要に応じて有効化）
/*
-- マスターテーブルをパーティションテーブルに変更
CREATE TABLE contents_partitioned (
    LIKE contents INCLUDING ALL
) PARTITION BY RANGE (created_at);

-- 月次パーティション作成例
CREATE TABLE contents_202401 PARTITION OF contents_partitioned
    FOR VALUES FROM ('2024-01-01') TO ('2024-02-01');
CREATE TABLE contents_202402 PARTITION OF contents_partitioned
    FOR VALUES FROM ('2024-02-01') TO ('2024-03-01');
*/

-- =============================================================================
-- セキュリティ設定
-- =============================================================================

-- RLS (Row Level Security) の有効化準備
-- ALTER TABLE contents ENABLE ROW LEVEL SECURITY;

-- サンプルRLSポリシー（コメントアウト - 認証実装時に使用）
/*
-- 作成者のみが自分のコンテンツを編集可能
CREATE POLICY contents_author_policy ON contents
    FOR ALL
    TO cms_api_role
    USING (author_id = current_setting('app.current_user_id'));

-- 公開コンテンツは全員が閲覧可能
CREATE POLICY contents_public_read_policy ON contents
    FOR SELECT
    TO public
    USING (status = 'published' AND published_at <= CURRENT_TIMESTAMP);
*/

-- =============================================================================
-- インデックス使用統計監視用ビュー
-- =============================================================================

/**
 * インデックス使用統計ビュー
 * パフォーマンス監視用
 */
CREATE VIEW index_usage_stats AS
SELECT 
    schemaname,
    tablename,
    indexname,
    idx_tup_read,
    idx_tup_fetch,
    idx_scan,
    CASE 
        WHEN idx_scan = 0 THEN 'UNUSED'
        WHEN idx_scan < 10 THEN 'LOW_USAGE'
        WHEN idx_scan < 100 THEN 'MODERATE_USAGE'
        ELSE 'HIGH_USAGE'
    END as usage_level
FROM pg_stat_user_indexes
WHERE schemaname = 'public' AND tablename = 'contents'
ORDER BY idx_scan DESC;

-- =============================================================================
-- 初期データ投入（テスト用）
-- =============================================================================

-- テスト用データ（開発環境でのみ実行）
-- INSERT文は実際の運用では削除してください

INSERT INTO contents (
    title, 
    body, 
    status, 
    author_id, 
    metadata,
    published_at
) VALUES 
-- 公開済みコンテンツ
(
    'CMS API システムの概要',
    'このCMS APIシステムは、AWS Lambda と Aurora Serverless v2を使用したサーバーレスアーキテクチャを採用しています。高い可用性と拡張性を実現し、コンテンツ管理を効率的に行うことができます。',
    'published',
    'admin',
    '{
        "seoTitle": "CMS API システム - 高性能なコンテンツ管理",
        "seoDescription": "AWS Lambda と Aurora Serverlessを使用した高性能CMS API",
        "category": "technology",
        "tags": ["CMS", "API", "AWS", "Lambda", "Aurora"],
        "excerpt": "AWS Lambda と Aurora Serverless v2を使用したサーバーレスCMSシステム"
    }',
    CURRENT_TIMESTAMP - INTERVAL '1 day'
),
-- ドラフトコンテンツ
(
    'パフォーマンス最適化ガイド',
    'このガイドでは、CMS APIのパフォーマンスを最適化するための具体的な手法について説明します。データベースインデックスの活用、クエリの最適化、キャッシュ戦略などを含みます。',
    'draft',
    'admin',
    '{
        "category": "guide",
        "tags": ["performance", "optimization", "database"],
        "excerpt": "CMS APIのパフォーマンス最適化手法"
    }',
    NULL
),
-- アーカイブコンテンツ
(
    '旧バージョンの仕様書',
    '以前のバージョンで使用されていた仕様書です。現在は使用されていませんが、参考資料として保管されています。',
    'archived',
    'admin',
    '{
        "category": "archive",
        "tags": ["legacy", "specification"],
        "excerpt": "旧バージョンの仕様書（アーカイブ）"
    }',
    NULL
);

-- =============================================================================
-- 運用・監視用クエリ
-- =============================================================================

-- よく使用される運用クエリをコメントで記載

/*
-- コンテンツ統計取得
SELECT 
    status,
    COUNT(*) as count,
    AVG(char_length(body)) as avg_body_length
FROM contents 
GROUP BY status;

-- 月別コンテンツ作成数
SELECT 
    DATE_TRUNC('month', created_at) as month,
    COUNT(*) as content_count
FROM contents
WHERE created_at >= CURRENT_DATE - INTERVAL '12 months'
GROUP BY DATE_TRUNC('month', created_at)
ORDER BY month DESC;

-- 人気カテゴリTop10
SELECT 
    metadata->>'category' as category,
    COUNT(*) as content_count
FROM contents
WHERE status = 'published'
    AND metadata->>'category' IS NOT NULL
GROUP BY metadata->>'category'
ORDER BY content_count DESC
LIMIT 10;

-- インデックス使用状況確認
SELECT * FROM index_usage_stats;

-- テーブルサイズ確認
SELECT 
    pg_size_pretty(pg_total_relation_size('contents')) as total_size,
    pg_size_pretty(pg_relation_size('contents')) as table_size,
    pg_size_pretty(pg_total_relation_size('contents') - pg_relation_size('contents')) as index_size;
*/