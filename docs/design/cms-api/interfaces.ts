// CMS API TypeScript インターフェース定義

// =============================================================================
// エンティティ定義
// =============================================================================

/**
 * コンテンツエンティティ
 */
export interface Content {
  /** コンテンツID (UUID) */
  id: string;
  /** タイトル */
  title: string;
  /** コンテンツ本文 */
  body: string;
  /** 公開状態 */
  status: ContentStatus;
  /** 作成日時 */
  createdAt: Date;
  /** 更新日時 */
  updatedAt: Date;
  /** 公開日時 */
  publishedAt: Date | null;
  /** 作成者ID */
  authorId: string;
  /** メタデータ */
  metadata: ContentMetadata;
}

/**
 * コンテンツステータス
 */
export type ContentStatus = 'draft' | 'published' | 'archived';

/**
 * コンテンツメタデータ
 */
export interface ContentMetadata {
  /** SEO タイトル */
  seoTitle?: string;
  /** SEO 説明文 */
  seoDescription?: string;
  /** タグ */
  tags: string[];
  /** カテゴリ */
  category?: string;
  /** 概要/抜粋 */
  excerpt?: string;
  /** アイキャッチ画像URL */
  featuredImageUrl?: string;
}

// =============================================================================
// APIリクエスト/レスポンス型定義
// =============================================================================

/**
 * 共通APIレスポンス型
 */
export interface ApiResponse<T = unknown> {
  /** 成功フラグ */
  success: boolean;
  /** レスポンスデータ */
  data?: T;
  /** エラー情報 */
  error?: ApiError;
  /** メタ情報 (ページネーション等) */
  meta?: ResponseMeta;
}

/**
 * APIエラー情報
 */
export interface ApiError {
  /** エラーコード */
  code: string;
  /** エラーメッセージ */
  message: string;
  /** 詳細エラー情報 */
  details?: Record<string, unknown>;
  /** タイムスタンプ */
  timestamp: string;
  /** リクエストID (トレーシング用) */
  requestId: string;
}

/**
 * レスポンスメタ情報
 */
export interface ResponseMeta {
  /** リクエストID */
  requestId: string;
  /** タイムスタンプ */
  timestamp: string;
  /** 処理時間 (ミリ秒) */
  processingTimeMs: number;
}

// =============================================================================
// コンテンツ詳細取得API
// =============================================================================

/**
 * コンテンツ詳細取得リクエストパラメータ
 */
export interface GetContentDetailParams {
  /** コンテンツID */
  id: string;
}

/**
 * コンテンツ詳細取得レスポンス
 */
export type GetContentDetailResponse = ApiResponse<Content>;

// =============================================================================
// コンテンツ一覧取得API
// =============================================================================

/**
 * コンテンツ一覧取得クエリパラメータ
 */
export interface GetContentListParams {
  /** 取得件数 (デフォルト: 20, 最大: 100) */
  limit?: number;
  /** オフセット (デフォルト: 0) */
  offset?: number;
  /** ステータスフィルタ */
  status?: ContentStatus;
  /** カテゴリフィルタ */
  category?: string;
  /** タグフィルタ (カンマ区切り) */
  tags?: string;
  /** 検索キーワード */
  search?: string;
  /** ソート順 */
  sort?: ContentSortOrder;
  /** ソート方向 */
  order?: SortDirection;
}

/**
 * コンテンツソート順
 */
export type ContentSortOrder = 
  | 'createdAt' 
  | 'updatedAt' 
  | 'publishedAt' 
  | 'title';

/**
 * ソート方向
 */
export type SortDirection = 'asc' | 'desc';

/**
 * コンテンツ一覧取得レスポンスデータ
 */
export interface ContentListData {
  /** コンテンツ一覧 */
  contents: Content[];
  /** ページネーション情報 */
  pagination: PaginationInfo;
}

/**
 * ページネーション情報
 */
export interface PaginationInfo {
  /** 現在のページ番号 (1から開始) */
  currentPage: number;
  /** 1ページあたりの件数 */
  perPage: number;
  /** 総件数 */
  totalCount: number;
  /** 総ページ数 */
  totalPages: number;
  /** 前のページが存在するか */
  hasPrev: boolean;
  /** 次のページが存在するか */
  hasNext: number;
  /** 前のページ番号 */
  prevPage: number | null;
  /** 次のページ番号 */
  nextPage: number | null;
}

/**
 * コンテンツ一覧取得レスポンス
 */
export type GetContentListResponse = ApiResponse<ContentListData>;

// =============================================================================
// エラーレスポンス型定義
// =============================================================================

/**
 * 400 Bad Request エラー
 */
export interface BadRequestError extends ApiError {
  code: 'INVALID_PARAMETER' | 'MISSING_PARAMETER' | 'INVALID_FORMAT';
}

/**
 * 404 Not Found エラー
 */
export interface NotFoundError extends ApiError {
  code: 'CONTENT_NOT_FOUND' | 'RESOURCE_NOT_FOUND';
}

/**
 * 500 Internal Server Error エラー
 */
export interface InternalServerError extends ApiError {
  code: 'INTERNAL_ERROR' | 'DATABASE_ERROR' | 'TIMEOUT_ERROR';
}

// =============================================================================
// Lambda関数内部型定義
// =============================================================================

/**
 * Lambda関数のイベント型（API Gateway）
 */
export interface ApiGatewayEvent {
  httpMethod: string;
  path: string;
  pathParameters: Record<string, string> | null;
  queryStringParameters: Record<string, string> | null;
  headers: Record<string, string>;
  body: string | null;
  requestContext: {
    requestId: string;
    httpMethod: string;
    path: string;
    stage: string;
    identity: {
      sourceIp: string;
      userAgent: string;
    };
  };
}

/**
 * Lambda関数のレスポンス型
 */
export interface ApiGatewayResponse {
  statusCode: number;
  headers: Record<string, string>;
  body: string;
}

/**
 * データベース接続設定
 */
export interface DatabaseConfig {
  /** クラスターARN */
  clusterArn: string;
  /** Secrets Manager ARN */
  secretArn: string;
  /** データベース名 */
  database: string;
  /** 接続タイムアウト (秒) */
  timeoutSeconds: number;
}

/**
 * ログレベル
 */
export type LogLevel = 'DEBUG' | 'INFO' | 'WARN' | 'ERROR';

/**
 * 構造化ログ
 */
export interface StructuredLog {
  /** ログレベル */
  level: LogLevel;
  /** メッセージ */
  message: string;
  /** タイムスタンプ */
  timestamp: string;
  /** リクエストID */
  requestId: string;
  /** 関数名 */
  functionName: string;
  /** 処理時間 (ミリ秒) */
  duration?: number;
  /** エラー情報 */
  error?: {
    name: string;
    message: string;
    stack?: string;
  };
  /** 追加のメタデータ */
  metadata?: Record<string, unknown>;
}

// =============================================================================
// バリデーション型定義
// =============================================================================

/**
 * バリデーションエラー詳細
 */
export interface ValidationError {
  /** フィールド名 */
  field: string;
  /** エラーメッセージ */
  message: string;
  /** 受信した値 */
  receivedValue: unknown;
  /** 期待される値の説明 */
  expectedFormat: string;
}

/**
 * バリデーション結果
 */
export interface ValidationResult {
  /** バリデーション成功フラグ */
  isValid: boolean;
  /** エラー一覧 */
  errors: ValidationError[];
}

// =============================================================================
// 設定・定数型定義
// =============================================================================

/**
 * API設定
 */
export interface ApiConfig {
  /** デフォルトページサイズ */
  defaultPageSize: number;
  /** 最大ページサイズ */
  maxPageSize: number;
  /** APIバージョン */
  version: string;
  /** タイムゾーン */
  timezone: string;
}

/**
 * HTTPステータスコード
 */
export const HTTP_STATUS = {
  OK: 200,
  BAD_REQUEST: 400,
  NOT_FOUND: 404,
  INTERNAL_SERVER_ERROR: 500,
  SERVICE_UNAVAILABLE: 503,
  GATEWAY_TIMEOUT: 504,
} as const;

/**
 * エラーコード定義
 */
export const ERROR_CODES = {
  // 400 Bad Request
  INVALID_PARAMETER: 'INVALID_PARAMETER',
  MISSING_PARAMETER: 'MISSING_PARAMETER',
  INVALID_FORMAT: 'INVALID_FORMAT',
  
  // 404 Not Found
  CONTENT_NOT_FOUND: 'CONTENT_NOT_FOUND',
  RESOURCE_NOT_FOUND: 'RESOURCE_NOT_FOUND',
  
  // 500 Internal Server Error
  INTERNAL_ERROR: 'INTERNAL_ERROR',
  DATABASE_ERROR: 'DATABASE_ERROR',
  TIMEOUT_ERROR: 'TIMEOUT_ERROR',
  CONNECTION_ERROR: 'CONNECTION_ERROR',
} as const;

/**
 * Content-Type定義
 */
export const CONTENT_TYPES = {
  JSON: 'application/json',
  TEXT: 'text/plain',
} as const;

// =============================================================================
// ユーティリティ型定義
// =============================================================================

/**
 * 必須プロパティを持つ型
 */
export type RequiredFields<T, K extends keyof T> = T & Required<Pick<T, K>>;

/**
 * オプションプロパティを持つ型
 */
export type OptionalFields<T, K extends keyof T> = Omit<T, K> & Partial<Pick<T, K>>;

/**
 * 作成時用のContent型（IDと日時を除く）
 */
export type CreateContentData = Omit<Content, 'id' | 'createdAt' | 'updatedAt'>;

/**
 * 更新時用のContent型（IDと作成日時を除く）
 */
export type UpdateContentData = Omit<Content, 'id' | 'createdAt'>;

/**
 * コンテンツサマリー型（一覧表示用の軽量版）
 */
export type ContentSummary = Pick<Content, 
  | 'id' 
  | 'title' 
  | 'status' 
  | 'createdAt' 
  | 'updatedAt' 
  | 'publishedAt'
> & {
  excerpt: string;
  category: string;
  tags: string[];
};