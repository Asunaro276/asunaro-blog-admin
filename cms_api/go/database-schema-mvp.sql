-- =============================================================================
-- CMS API データベーススキーマ設計 v2.0 - MVP版
-- Aurora PostgreSQL 15 対応
-- 基本的なコンテンツ管理機能のみを含む最小実行可能バージョン
-- =============================================================================

-- 拡張機能の有効化
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- =============================================================================
-- コンテンツタイプ管理テーブル
-- =============================================================================

/**
 * コンテンツタイプ定義テーブル
 * ユーザーが定義可能なコンテンツタイプ（ブログ記事、商品、イベントなど）
 */
CREATE TABLE content_types (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL UNIQUE,
    display_name VARCHAR(200) NOT NULL,
    description TEXT,
    icon VARCHAR(50),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(100) NOT NULL
);

-- =============================================================================
-- コンテンツマスターテーブル（MVP版）
-- =============================================================================

/**
 * コンテンツマスターテーブル
 * CMSの中核となるコンテンツ情報を格納
 */
CREATE TABLE contents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    content_type_id UUID NOT NULL REFERENCES content_types(id),
    title VARCHAR(500) NOT NULL,
    slug VARCHAR(200),
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    published_at TIMESTAMP WITH TIME ZONE,
    author_id VARCHAR(100) NOT NULL,
    version INTEGER NOT NULL DEFAULT 1,
    UNIQUE(content_type_id, slug),
    CONSTRAINT chk_contents_status 
        CHECK (status IN ('draft', 'published', 'archived', 'trash'))
);

-- =============================================================================
-- ブロックベースコンテンツ管理テーブル（MVP版）
-- =============================================================================

/**
 * コンテンツブロック定義テーブル
 * 各コンテンツが持つブロック（段落、画像、埋め込みなど）を定義
 */
CREATE TABLE content_blocks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    content_id UUID NOT NULL REFERENCES contents(id) ON DELETE CASCADE,
    block_type VARCHAR(50) NOT NULL,
    block_order INTEGER NOT NULL DEFAULT 0,
    is_visible BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_block_type 
        CHECK (block_type IN ('richtext', 'plaintext', 'number', 'image', 'video', 'audio', 'file', 'content_reference', 'embed', 'code', 'quote', 'divider'))
);

/**
 * ブロックデータテーブル
 * 各ブロックの実際のデータを格納
 */
CREATE TABLE content_block_data (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    block_id UUID NOT NULL REFERENCES content_blocks(id) ON DELETE CASCADE,
    data_type VARCHAR(50) NOT NULL,
    content_text TEXT,
    content_richtext JSONB,
    content_number DECIMAL,
    content_url VARCHAR(1000),
    content_json JSONB,
    referenced_content_id UUID REFERENCES contents(id),
    settings JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_data_type 
        CHECK (data_type IN ('text', 'richtext', 'number', 'url', 'json', 'reference')),
    CONSTRAINT chk_reference_consistency
        CHECK (
            (data_type = 'reference' AND referenced_content_id IS NOT NULL) OR
            (data_type != 'reference' AND referenced_content_id IS NULL)
        )
);

-- =============================================================================
-- インデックス設計（MVP版）
-- =============================================================================

-- コンテンツテーブルのインデックス
CREATE INDEX idx_contents_content_type_id ON contents(content_type_id);
CREATE INDEX idx_contents_status ON contents(status);
CREATE INDEX idx_contents_author_id ON contents(author_id);
CREATE INDEX idx_contents_created_at ON contents(created_at DESC);
CREATE INDEX idx_contents_updated_at ON contents(updated_at DESC);
CREATE INDEX idx_contents_published_at ON contents(published_at DESC) 
    WHERE published_at IS NOT NULL;
CREATE INDEX idx_contents_slug ON contents(content_type_id, slug);

CREATE INDEX idx_contents_status_published_at ON contents(status, published_at DESC) 
    WHERE status = 'published';
CREATE INDEX idx_contents_status_created_at ON contents(status, created_at DESC);
CREATE INDEX idx_contents_author_status ON contents(author_id, status);
CREATE INDEX idx_contents_type_status ON contents(content_type_id, status);

CREATE INDEX idx_contents_title_gin ON contents USING GIN (title gin_trgm_ops);

-- コンテンツタイプ関連のインデックス
CREATE INDEX idx_content_types_name ON content_types(name);
CREATE INDEX idx_content_types_is_active ON content_types(is_active);

-- ブロック関連のインデックス
CREATE INDEX idx_content_blocks_content_id ON content_blocks(content_id);
CREATE INDEX idx_content_blocks_content_order ON content_blocks(content_id, block_order);
CREATE INDEX idx_content_blocks_type ON content_blocks(block_type);
CREATE INDEX idx_content_blocks_visible ON content_blocks(content_id, is_visible, block_order);

CREATE INDEX idx_content_block_data_block_id ON content_block_data(block_id);
CREATE INDEX idx_content_block_data_data_type ON content_block_data(data_type);
CREATE INDEX idx_content_block_data_text_gin ON content_block_data USING GIN (content_text gin_trgm_ops);
CREATE INDEX idx_content_block_data_richtext_gin ON content_block_data USING GIN (content_richtext);
CREATE INDEX idx_content_block_data_referenced_content ON content_block_data(referenced_content_id) WHERE referenced_content_id IS NOT NULL;
CREATE INDEX idx_content_block_data_json_gin ON content_block_data USING GIN (content_json);

-- =============================================================================
-- ビュー定義（MVP版）
-- =============================================================================

/**
 * コンテンツ詳細ビュー
 * コンテンツとそのタイプ情報を結合
 */
CREATE VIEW content_details AS
SELECT 
    c.id,
    c.content_type_id,
    ct.name as content_type_name,
    ct.display_name as content_type_display_name,
    c.title,
    c.slug,
    c.status,
    c.created_at,
    c.updated_at,
    c.published_at,
    c.author_id,
    c.version
FROM contents c
JOIN content_types ct ON c.content_type_id = ct.id
WHERE ct.is_active = true;

/**
 * 公開コンテンツビュー（MVP版）
 * 公開状態のコンテンツのみを表示（タグ機能は除外）
 */
CREATE VIEW published_contents AS
SELECT 
    cd.*
FROM content_details cd
WHERE cd.status = 'published' 
    AND cd.published_at IS NOT NULL 
    AND cd.published_at <= CURRENT_TIMESTAMP;

/**
 * コンテンツブロック詳細ビュー
 * コンテンツのブロック構造を表示
 */
CREATE VIEW content_block_details AS
SELECT 
    cb.id as block_id,
    cb.content_id,
    cb.block_type,
    cb.block_order,
    cb.is_visible,
    cbd.data_type,
    cbd.content_text,
    cbd.content_richtext,
    cbd.content_number,
    cbd.content_url,
    cbd.content_json,
    cbd.referenced_content_id,
    rc.title as referenced_content_title,
    cbd.settings,
    cb.created_at as block_created_at,
    cbd.created_at as data_created_at
FROM content_blocks cb
LEFT JOIN content_block_data cbd ON cb.id = cbd.block_id
LEFT JOIN contents rc ON cbd.referenced_content_id = rc.id
WHERE cb.is_visible = true
ORDER BY cb.content_id, cb.block_order;

-- パフォーマンス最適化設定
ALTER TABLE contents SET (autovacuum_analyze_scale_factor = 0.05);
ALTER TABLE contents SET (autovacuum_vacuum_scale_factor = 0.1);

-- =============================================================================
-- インデックス使用統計監視用ビュー（MVP版）
-- =============================================================================

/**
 * インデックス使用統計ビュー
 * パフォーマンス監視用
 */
CREATE VIEW index_usage_stats AS
SELECT 
    schemaname,
    relname as tablename,
    indexrelname as indexname,
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
WHERE schemaname = 'public' AND relname IN ('contents', 'content_types', 'content_blocks', 'content_block_data')
ORDER BY idx_scan DESC;

-- =============================================================================
-- 初期データ投入（MVP版テスト用）
-- =============================================================================

-- テスト用データ（開発環境でのみ実行）

-- コンテンツタイプの作成
INSERT INTO content_types (id, name, display_name, description, icon, created_by) VALUES 
('550e8400-e29b-41d4-a716-446655440001', 'blog_post', 'ブログ記事', 'ブログ記事用のコンテンツタイプ', 'article', 'admin'),
('550e8400-e29b-41d4-a716-446655440002', 'page', '固定ページ', '固定ページ用のコンテンツタイプ', 'page', 'admin'),
('550e8400-e29b-41d4-a716-446655440003', 'product', '商品', '商品情報用のコンテンツタイプ', 'shopping-cart', 'admin');

-- コンテンツの作成
INSERT INTO contents (id, content_type_id, title, slug, status, author_id, published_at) VALUES 
('550e8400-e29b-41d4-a716-446655440201', '550e8400-e29b-41d4-a716-446655440001', 'CMS API システムの概要', 'cms-api-overview', 'published', 'admin', CURRENT_TIMESTAMP - INTERVAL '1 day'),
('550e8400-e29b-41d4-a716-446655440202', '550e8400-e29b-41d4-a716-446655440001', 'パフォーマンス最適化ガイド', 'performance-optimization-guide', 'draft', 'admin', NULL),
('550e8400-e29b-41d4-a716-446655440203', '550e8400-e29b-41d4-a716-446655440002', 'プライバシーポリシー', 'privacy-policy', 'published', 'admin', CURRENT_TIMESTAMP - INTERVAL '7 days');

-- コンテンツブロックの作成
INSERT INTO content_blocks (id, content_id, block_type, block_order) VALUES 
('550e8400-e29b-41d4-a716-446655440301', '550e8400-e29b-41d4-a716-446655440201', 'richtext', 1),
('550e8400-e29b-41d4-a716-446655440302', '550e8400-e29b-41d4-a716-446655440201', 'plaintext', 2),
('550e8400-e29b-41d4-a716-446655440303', '550e8400-e29b-41d4-a716-446655440202', 'richtext', 1);

-- ブロックデータの作成
INSERT INTO content_block_data (block_id, data_type, content_richtext, content_text) VALUES 
('550e8400-e29b-41d4-a716-446655440301', 'richtext', 
'{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"このCMS APIシステムは、AWS Lambda と Aurora Serverless v2を使用したサーバーレスアーキテクチャを採用しています。"}]},{"type":"paragraph","content":[{"type":"text","marks":[{"type":"bold"}],"text":"主な特徴:"},{"type":"hard_break"},{"type":"text","text":"• 高い可用性と拡張性"},{"type":"hard_break"},{"type":"text","text":"• 効率的なコンテンツ管理"},{"type":"hard_break"},{"type":"text","text":"• 柔軟なブロックベースエディタ"}]}]}'::jsonb, 
NULL),
('550e8400-e29b-41d4-a716-446655440302', 'text', NULL, 
'技術スタック: AWS Lambda, Aurora Serverless v2, PostgreSQL 15'),
('550e8400-e29b-41d4-a716-446655440303', 'richtext',
'{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"このガイドでは、CMS APIのパフォーマンスを最適化するための具体的な手法について説明します。"}]}]}'::jsonb,
NULL);