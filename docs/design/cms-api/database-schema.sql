-- =============================================================================
-- CMS API データベーススキーマ設計 v2.0
-- Aurora PostgreSQL 15 対応
-- 柔軟なコンテンツタイプ & ブロックベースCMS
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

/**
 * コンテンツタイプフィールド定義テーブル
 * 各コンテンツタイプに属するフィールドの定義
 */
CREATE TABLE content_type_fields (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    content_type_id UUID NOT NULL REFERENCES content_types(id) ON DELETE CASCADE,
    field_name VARCHAR(100) NOT NULL,
    display_name VARCHAR(200) NOT NULL,
    field_type VARCHAR(50) NOT NULL,
    is_required BOOLEAN NOT NULL DEFAULT false,
    is_unique BOOLEAN NOT NULL DEFAULT false,
    default_value TEXT,
    validation_rules JSONB DEFAULT '{}',
    field_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(content_type_id, field_name),
    CONSTRAINT chk_field_type 
        CHECK (field_type IN ('text', 'richtext', 'number', 'date', 'datetime', 'boolean', 'select', 'multiselect', 'url', 'email', 'json'))
);

-- =============================================================================
-- コンテンツマスターテーブル（拡張版）
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

/**
 * コンテンツフィールド値テーブル
 * 各コンテンツのフィールド値を格納
 */
CREATE TABLE content_field_values (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    content_id UUID NOT NULL REFERENCES contents(id) ON DELETE CASCADE,
    field_id UUID NOT NULL REFERENCES content_type_fields(id) ON DELETE CASCADE,
    value_text TEXT,
    value_number DECIMAL,
    value_date DATE,
    value_datetime TIMESTAMP WITH TIME ZONE,
    value_boolean BOOLEAN,
    value_json JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(content_id, field_id)
);

-- =============================================================================
-- ブロックベースコンテンツ管理テーブル
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
-- コンテンツ関係管理テーブル
-- =============================================================================

/**
 * コンテンツリレーションテーブル
 * コンテンツ間の関係性を管理（関連記事、カテゴリ分類など）
 */
CREATE TABLE content_relations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    source_content_id UUID NOT NULL REFERENCES contents(id) ON DELETE CASCADE,
    target_content_id UUID NOT NULL REFERENCES contents(id) ON DELETE CASCADE,
    relation_type VARCHAR(50) NOT NULL,
    relation_metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(100) NOT NULL,
    UNIQUE(source_content_id, target_content_id, relation_type),
    CONSTRAINT chk_no_self_reference 
        CHECK (source_content_id != target_content_id),
    CONSTRAINT chk_relation_type
        CHECK (relation_type IN ('related', 'category', 'tag', 'parent', 'child', 'reference', 'similar'))
);

/**
 * コンテンツ参照階層テーブル
 * 循環参照の検出と防止のための階層管理
 */
CREATE TABLE content_reference_hierarchy (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ancestor_id UUID NOT NULL REFERENCES contents(id) ON DELETE CASCADE,
    descendant_id UUID NOT NULL REFERENCES contents(id) ON DELETE CASCADE,
    depth INTEGER NOT NULL DEFAULT 0,
    path_length INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(ancestor_id, descendant_id),
    CONSTRAINT chk_no_self_hierarchy 
        CHECK (ancestor_id != descendant_id OR depth = 0),
    CONSTRAINT chk_positive_depth 
        CHECK (depth >= 0)
);

/**
 * コンテンツタグテーブル
 * フラットなタグ管理（従来のカテゴリ・タグ機能の代替）
 */
CREATE TABLE content_tags (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL UNIQUE,
    display_name VARCHAR(200) NOT NULL,
    description TEXT,
    color VARCHAR(7),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(100) NOT NULL
);

/**
 * コンテンツ-タグ関連テーブル
 * コンテンツとタグの多対多関係
 */
CREATE TABLE content_tag_assignments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    content_id UUID NOT NULL REFERENCES contents(id) ON DELETE CASCADE,
    tag_id UUID NOT NULL REFERENCES content_tags(id) ON DELETE CASCADE,
    assigned_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    assigned_by VARCHAR(100) NOT NULL,
    UNIQUE(content_id, tag_id)
);

-- =============================================================================
-- インデックス設計
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

CREATE INDEX idx_content_type_fields_content_type_id ON content_type_fields(content_type_id);
CREATE INDEX idx_content_type_fields_field_order ON content_type_fields(content_type_id, field_order);
CREATE INDEX idx_content_type_fields_field_type ON content_type_fields(field_type);

-- コンテンツフィールド値のインデックス
CREATE INDEX idx_content_field_values_content_id ON content_field_values(content_id);
CREATE INDEX idx_content_field_values_field_id ON content_field_values(field_id);
CREATE INDEX idx_content_field_values_text_gin ON content_field_values USING GIN (value_text gin_trgm_ops);
CREATE INDEX idx_content_field_values_number ON content_field_values(value_number) WHERE value_number IS NOT NULL;
CREATE INDEX idx_content_field_values_date ON content_field_values(value_date) WHERE value_date IS NOT NULL;
CREATE INDEX idx_content_field_values_datetime ON content_field_values(value_datetime) WHERE value_datetime IS NOT NULL;
CREATE INDEX idx_content_field_values_json_gin ON content_field_values USING GIN (value_json);

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

-- 関係性・参照関連のインデックス
CREATE INDEX idx_content_relations_source ON content_relations(source_content_id);
CREATE INDEX idx_content_relations_target ON content_relations(target_content_id);
CREATE INDEX idx_content_relations_type ON content_relations(relation_type);
CREATE INDEX idx_content_relations_source_type ON content_relations(source_content_id, relation_type);

CREATE INDEX idx_content_reference_hierarchy_ancestor ON content_reference_hierarchy(ancestor_id);
CREATE INDEX idx_content_reference_hierarchy_descendant ON content_reference_hierarchy(descendant_id);
CREATE INDEX idx_content_reference_hierarchy_depth ON content_reference_hierarchy(depth);

-- タグ関連のインデックス
CREATE INDEX idx_content_tags_name ON content_tags(name);
CREATE INDEX idx_content_tags_is_active ON content_tags(is_active);

CREATE INDEX idx_content_tag_assignments_content_id ON content_tag_assignments(content_id);
CREATE INDEX idx_content_tag_assignments_tag_id ON content_tag_assignments(tag_id);

-- =============================================================================
-- ビュー定義
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
 * 公開コンテンツビュー
 * 公開状態のコンテンツのみを表示
 */
CREATE VIEW published_contents AS
SELECT 
    cd.*,
    -- タグ情報を集約
    COALESCE(
        array_agg(
            json_build_object(
                'id', ct.id,
                'name', ct.name,
                'display_name', ct.display_name,
                'color', ct.color
            ) ORDER BY ct.display_name
        ) FILTER (WHERE ct.id IS NOT NULL),
        ARRAY[]::json[]
    ) as tags
FROM content_details cd
LEFT JOIN content_tag_assignments cta ON cd.id = cta.content_id
LEFT JOIN content_tags ct ON cta.tag_id = ct.id AND ct.is_active = true
WHERE cd.status = 'published' 
    AND cd.published_at IS NOT NULL 
    AND cd.published_at <= CURRENT_TIMESTAMP
GROUP BY cd.id, cd.content_type_id, cd.content_type_name, cd.content_type_display_name,
         cd.title, cd.slug, cd.status, cd.created_at, cd.updated_at, cd.published_at,
         cd.author_id, cd.version;

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

/**
 * コンテンツフィールド詳細ビュー
 * コンテンツのフィールド値とその定義を表示
 */
CREATE VIEW content_field_details AS
SELECT 
    cfv.content_id,
    ctf.field_name,
    ctf.display_name,
    ctf.field_type,
    ctf.is_required,
    cfv.value_text,
    cfv.value_number,
    cfv.value_date,
    cfv.value_datetime,
    cfv.value_boolean,
    cfv.value_json,
    ctf.field_order
FROM content_field_values cfv
JOIN content_type_fields ctf ON cfv.field_id = ctf.id
ORDER BY cfv.content_id, ctf.field_order;

/**
 * コンテンツ関係詳細ビュー
 * コンテンツ間の関係性を表示
 */
CREATE VIEW content_relation_details AS
SELECT 
    cr.id as relation_id,
    cr.source_content_id,
    sc.title as source_content_title,
    cr.target_content_id,
    tc.title as target_content_title,
    cr.relation_type,
    cr.relation_metadata,
    cr.created_at,
    cr.created_by
FROM content_relations cr
JOIN contents sc ON cr.source_content_id = sc.id
JOIN contents tc ON cr.target_content_id = tc.id;

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

-- コンテンツタイプの作成
INSERT INTO content_types (id, name, display_name, description, icon, created_by) VALUES 
('550e8400-e29b-41d4-a716-446655440001', 'blog_post', 'ブログ記事', 'ブログ記事用のコンテンツタイプ', 'article', 'admin'),
('550e8400-e29b-41d4-a716-446655440002', 'page', '固定ページ', '固定ページ用のコンテンツタイプ', 'page', 'admin'),
('550e8400-e29b-41d4-a716-446655440003', 'product', '商品', '商品情報用のコンテンツタイプ', 'shopping-cart', 'admin');

-- ブログ記事タイプのフィールド定義
INSERT INTO content_type_fields (content_type_id, field_name, display_name, field_type, is_required, field_order) VALUES 
('550e8400-e29b-41d4-a716-446655440001', 'excerpt', '抜粋', 'text', false, 1),
('550e8400-e29b-41d4-a716-446655440001', 'featured_image', 'アイキャッチ画像', 'url', false, 2),
('550e8400-e29b-41d4-a716-446655440001', 'seo_title', 'SEOタイトル', 'text', false, 3),
('550e8400-e29b-41d4-a716-446655440001', 'seo_description', 'SEO説明文', 'text', false, 4);

-- 商品タイプのフィールド定義
INSERT INTO content_type_fields (content_type_id, field_name, display_name, field_type, is_required, field_order) VALUES 
('550e8400-e29b-41d4-a716-446655440003', 'price', '価格', 'number', true, 1),
('550e8400-e29b-41d4-a716-446655440003', 'stock_quantity', '在庫数', 'number', false, 2),
('550e8400-e29b-41d4-a716-446655440003', 'product_image', '商品画像', 'url', false, 3),
('550e8400-e29b-41d4-a716-446655440003', 'release_date', '発売日', 'date', false, 4);

-- タグの作成
INSERT INTO content_tags (id, name, display_name, description, color, created_by) VALUES 
('550e8400-e29b-41d4-a716-446655440101', 'technology', 'テクノロジー', '技術関連のコンテンツ', '#007bff', 'admin'),
('550e8400-e29b-41d4-a716-446655440102', 'cms', 'CMS', 'CMSに関する内容', '#28a745', 'admin'),
('550e8400-e29b-41d4-a716-446655440103', 'api', 'API', 'API関連の内容', '#dc3545', 'admin'),
('550e8400-e29b-41d4-a716-446655440104', 'performance', 'パフォーマンス', 'パフォーマンス最適化', '#ffc107', 'admin');

-- コンテンツの作成
INSERT INTO contents (id, content_type_id, title, slug, status, author_id, published_at) VALUES 
('550e8400-e29b-41d4-a716-446655440201', '550e8400-e29b-41d4-a716-446655440001', 'CMS API システムの概要', 'cms-api-overview', 'published', 'admin', CURRENT_TIMESTAMP - INTERVAL '1 day'),
('550e8400-e29b-41d4-a716-446655440202', '550e8400-e29b-41d4-a716-446655440001', 'パフォーマンス最適化ガイド', 'performance-optimization-guide', 'draft', 'admin', NULL),
('550e8400-e29b-41d4-a716-446655440203', '550e8400-e29b-41d4-a716-446655440002', 'プライバシーポリシー', 'privacy-policy', 'published', 'admin', CURRENT_TIMESTAMP - INTERVAL '7 days');

-- コンテンツのフィールド値
INSERT INTO content_field_values (content_id, field_id, value_text) VALUES 
('550e8400-e29b-41d4-a716-446655440201', (SELECT id FROM content_type_fields WHERE content_type_id = '550e8400-e29b-41d4-a716-446655440001' AND field_name = 'excerpt'), 'AWS Lambda と Aurora Serverless v2を使用したサーバーレスCMSシステムの概要'),
('550e8400-e29b-41d4-a716-446655440201', (SELECT id FROM content_type_fields WHERE content_type_id = '550e8400-e29b-41d4-a716-446655440001' AND field_name = 'seo_title'), 'CMS API システム - 高性能なコンテンツ管理'),
('550e8400-e29b-41d4-a716-446655440201', (SELECT id FROM content_type_fields WHERE content_type_id = '550e8400-e29b-41d4-a716-446655440001' AND field_name = 'seo_description'), 'AWS Lambda と Aurora Serverlessを使用した高性能CMS API');

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

-- タグ割り当て
INSERT INTO content_tag_assignments (content_id, tag_id, assigned_by) VALUES 
('550e8400-e29b-41d4-a716-446655440201', '550e8400-e29b-41d4-a716-446655440101', 'admin'),
('550e8400-e29b-41d4-a716-446655440201', '550e8400-e29b-41d4-a716-446655440102', 'admin'),
('550e8400-e29b-41d4-a716-446655440201', '550e8400-e29b-41d4-a716-446655440103', 'admin'),
('550e8400-e29b-41d4-a716-446655440202', '550e8400-e29b-41d4-a716-446655440104', 'admin');
