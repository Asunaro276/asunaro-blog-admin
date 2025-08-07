# API エンドポイント仕様

## 概要

CMS APIは RESTful な設計原則に従い、JSON形式でのデータ交換を行います。全てのエンドポイントはHTTPS通信を使用し、適切なHTTPステータスコードとエラーハンドリングを提供します。

### 基本情報

- **ベースURL**: `https://api.cms.example.com/v1`
- **Content-Type**: `application/json`
- **文字エンコーディング**: UTF-8
- **認証**: 現在は不要（将来的にAPI Key認証を予定）

### 共通レスポンス形式

```json
{
  "success": boolean,
  "data": object | array | null,
  "error": {
    "code": "string",
    "message": "string",
    "details": object,
    "timestamp": "string (ISO 8601)",
    "requestId": "string"
  },
  "meta": {
    "requestId": "string",
    "timestamp": "string (ISO 8601)",
    "processingTimeMs": number
  }
}
```

## エンドポイント一覧

### 1. コンテンツ詳細取得

指定されたIDのコンテンツの詳細情報を取得します。

#### リクエスト

```
GET /contents/{id}
```

**パスパラメータ**
- `id` (required): コンテンツID (UUID形式)

**クエリパラメータ**
なし

#### レスポンス

**成功時 (200 OK)**

```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "metadata": {
      "createdAt": "2024-01-15T10:30:00Z",
      "updatedAt": "2024-01-15T14:20:00Z",
      "publishedAt": "2024-01-15T12:00:00Z",
    },
    "block1": ~~~,
    "block2": ~~~,
    "block3": ~~~,
  }
}
```

**エラー時**

```json
// 404 Not Found
{
  "code": "CONTENT_NOT_FOUND",
  "message": "指定されたコンテンツが見つかりません",
}

// 400 Bad Request
{
  "code": "INVALID_PARAMETER",
  "message": "不正なパラメータです",
}
```

### 2. コンテンツ一覧取得

コンテンツの一覧を取得します。ページネーション、フィルタリング、ソートに対応しています。

#### リクエスト

```
GET /contents
```

**クエリパラメータ**

| パラメータ | 型 | 必須 | デフォルト | 説明 |
|-----------|-----|-----|-----------|------|
| `limit` | integer | No | 20 | 取得件数 (1-100) |
| `offset` | integer | No | 0 | オフセット (0以上) |
| `status` | string | No | - | ステータスフィルタ (`draft`, `published`, `archived`) |
| `category` | string | No | - | カテゴリフィルタ |
| `tags` | string | No | - | タグフィルタ (カンマ区切り) |
| `search` | string | No | - | 検索キーワード (タイトル・本文を対象) |
| `sort` | string | No | createdAt | ソート対象 (`createdAt`, `updatedAt`, `publishedAt`, `title`) |
| `order` | string | No | desc | ソート順 (`asc`, `desc`) |

#### レスポンス

**成功時 (200 OK)**

```json
{
  "success": true,
  "data": {
    "contents": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "metadata": {
          "createdAt": "2024-01-15T10:30:00Z",
          "updatedAt": "2024-01-15T14:20:00Z",
          "publishedAt": "2024-01-15T12:00:00Z",
        },
        "block1": ~~~,
        "block2": ~~~,
        "block3": ~~~,
      }
    ],
    "pagination": {
      "currentPage": 1,
      "perPage": 20,
      "totalCount": 25,
      "totalPages": 2,
      "hasPrev": false,
      "hasNext": true,
      "prevPage": null,
      "nextPage": 2
    }
  },
}
```

**エラー時**

```json
// 400 Bad Request (不正なパラメータ)
{
  "code": "INVALID_PARAMETER",
  "message": "不正なパラメータです",
}
```

### 3. ヘルスチェック

システムの動作状態を確認します。

#### リクエスト

```
GET /healthcheck
```

**パスパラメータ**
なし

**クエリパラメータ**
なし

#### レスポンス

**成功時 (200 OK)**

```json
{
  "data": {
    "status": "healthy",
    "timestamp": "2024-01-15T15:30:00Z",
    "version": "1.0.0",
    "environment": "production",
    "database": {
      "status": "connected",
      "responseTime": 12
    },
    "services": {
      "aurora": "healthy",
      "secretsManager": "healthy"
    }
  },
  "meta": {
    "requestId": "req_health_123456",
    "timestamp": "2024-01-15T15:30:00Z",
    "processingTimeMs": 15
  }
}
```

**エラー時 (503 Service Unavailable)**

```json
{
  "data": {
    "status": "unhealthy",
    "timestamp": "2024-01-15T15:30:00Z",
    "version": "1.0.0",
    "environment": "production",
    "database": {
      "status": "disconnected",
      "responseTime": null,
      "error": "Connection timeout"
    },
    "services": {
      "aurora": "unhealthy",
      "secretsManager": "healthy"
    }
  },
  "error": {
    "code": "SERVICE_UNAVAILABLE",
    "message": "サービスが一時的に利用できません",
  },
}
```

## エラーコード一覧

### 4xx クライアントエラー

| コード | HTTPステータス | 説明 |
|--------|---------------|------|
| `INVALID_PARAMETER` | 400 | パラメータの値が不正です |
| `MISSING_PARAMETER` | 400 | 必須パラメータが不足しています |
| `INVALID_FORMAT` | 400 | データ形式が不正です |
| `CONTENT_NOT_FOUND` | 404 | コンテンツが見つかりません |
| `RESOURCE_NOT_FOUND` | 404 | リソースが見つかりません |

### 5xx サーバーエラー

| コード | HTTPステータス | 説明 |
|--------|---------------|------|
| `INTERNAL_ERROR` | 500 | 内部サーバーエラーが発生しました |
| `DATABASE_ERROR` | 500 | データベースエラーが発生しました |
| `TIMEOUT_ERROR` | 504 | タイムアウトが発生しました |
| `CONNECTION_ERROR` | 500 | 接続エラーが発生しました |
| `SERVICE_UNAVAILABLE` | 503 | サービスが一時的に利用できません |

## APIリクエスト例

### cURL

```bash
# コンテンツ詳細取得
curl -X GET "https://api.cms.example.com/v1/contents/550e8400-e29b-41d4-a716-446655440000" \
  -H "Accept: application/json"

# コンテンツ一覧取得（フィルタ・ソート付き）
curl -X GET "https://api.cms.example.com/v1/contents?status=published&limit=10&sort=publishedAt&order=desc" \
  -H "Accept: application/json"

# 検索クエリ付きリクエスト
curl -X GET "https://api.cms.example.com/v1/contents?search=AWS&category=technology&tags=API,Lambda" \
  -H "Accept: application/json"

# ヘルスチェック
curl -X GET "https://api.cms.example.com/v1/healthcheck" \
  -H "Accept: application/json"
```

### JavaScript (Fetch API)

```javascript
// コンテンツ詳細取得
const getContentDetail = async (contentId) => {
  try {
    const response = await fetch(`https://api.cms.example.com/v1/contents/${contentId}`, {
      method: 'GET',
      headers: {
        'Accept': 'application/json',
      },
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const data = await response.json();
    return data.success ? data.data : null;
  } catch (error) {
    console.error('Error fetching content:', error);
    throw error;
  }
};

// コンテンツ一覧取得
const getContentList = async (params = {}) => {
  const searchParams = new URLSearchParams(params);
  
  try {
    const response = await fetch(`https://api.cms.example.com/v1/contents?${searchParams}`, {
      method: 'GET',
      headers: {
        'Accept': 'application/json',
      },
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const data = await response.json();
    return data.success ? data.data : null;
  } catch (error) {
    console.error('Error fetching content list:', error);
    throw error;
  }
};

// 使用例
getContentList({
  status: 'published',
  limit: 10,
  sort: 'publishedAt',
  order: 'desc'
}).then(data => {
  console.log('Content list:', data.contents);
  console.log('Pagination:', data.pagination);
});
```

### TypeScript (with 型定義)

```typescript
import { GetContentDetailResponse, GetContentListResponse, GetContentListParams } from './interfaces';

class CMSApiClient {
  private baseUrl: string;

  constructor(baseUrl: string = 'https://api.cms.example.com/v1') {
    this.baseUrl = baseUrl;
  }

  async getContentDetail(contentId: string): Promise<GetContentDetailResponse> {
    const response = await fetch(`${this.baseUrl}/contents/${contentId}`, {
      method: 'GET',
      headers: {
        'Accept': 'application/json',
      },
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    return response.json();
  }

  async getContentList(params: GetContentListParams = {}): Promise<GetContentListResponse> {
    const searchParams = new URLSearchParams();
    
    Object.entries(params).forEach(([key, value]) => {
      if (value !== undefined && value !== null) {
        searchParams.append(key, String(value));
      }
    });

    const response = await fetch(`${this.baseUrl}/contents?${searchParams}`, {
      method: 'GET',
      headers: {
        'Accept': 'application/json',
      },
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    return response.json();
  }
}

// 使用例
const apiClient = new CMSApiClient();

// コンテンツ詳細取得
const content = await apiClient.getContentDetail('550e8400-e29b-41d4-a716-446655440000');
if (content.success) {
  console.log('Content title:', content.data.title);
}

// コンテンツ一覧取得
const contentList = await apiClient.getContentList({
  status: 'published',
  limit: 20,
  offset: 0,
  sort: 'createdAt',
  order: 'desc'
});

if (contentList.success) {
  contentList.data.contents.forEach(content => {
    console.log(`${content.title} - ${content.status}`);
  });
}
```

## パフォーマンス考慮事項

### キャッシュ戦略

- **コンテンツ詳細**: 5分間のブラウザキャッシュ
- **コンテンツ一覧**: 1分間のブラウザキャッシュ
- **API Gateway**: 必要に応じてレスポンスキャッシュを有効化

### レート制限

- **現在**: 制限なし
- **将来予定**: 
  - 認証なし: 100 req/min
  - API Key認証: 1000 req/min

### ページネーション推奨事項

- **デフォルトページサイズ**: 20件
- **最大ページサイズ**: 100件
- **大量データ取得**: 複数リクエストに分割して実行

## セキュリティ考慮事項

### CORS設定

```
Access-Control-Allow-Origin: *
Access-Control-Allow-Methods: GET, OPTIONS
Access-Control-Allow-Headers: Content-Type, Accept
Access-Control-Max-Age: 86400
```

### 入力値検証

- **UUID形式**: 正規表現による検証
- **数値範囲**: 最小・最大値チェック
- **文字列長**: 最大長制限
- **SQLインジェクション**: パラメータ化クエリによる防御

### エラー情報の制限

- **本番環境**: 詳細なエラー情報は非表示
- **開発環境**: デバッグ用の詳細情報を提供

## 監視・ログ

### アクセスログ

```json
{
  "timestamp": "2024-01-15T15:30:00Z",
  "requestId": "req_123456789",
  "method": "GET",
  "path": "/contents/550e8400-e29b-41d4-a716-446655440000",
  "statusCode": 200,
  "responseTime": 45,
  "userAgent": "Mozilla/5.0...",
  "sourceIp": "192.168.1.100"
}
```

### エラーログ

```json
{
  "timestamp": "2024-01-15T15:30:00Z",
  "level": "ERROR",
  "requestId": "req_123456789",
  "errorCode": "DATABASE_ERROR",
  "message": "Database connection failed",
  "details": {
    "endpoint": "/contents/550e8400-e29b-41d4-a716-446655440000",
    "duration": 5000
  }
}
```

## 将来拡張予定

### 認証機能

- API Key認証
- JWT Bearer Token認証
- OAuth 2.0対応

### 新規エンドポイント

- `POST /contents` - コンテンツ作成
- `PUT /contents/{id}` - コンテンツ更新
- `DELETE /contents/{id}` - コンテンツ削除
- `GET /contents/{id}/history` - バージョン履歴取得

### 機能強化

- 全文検索の精度向上
- レスポンス形式の最適化
- GraphQL対応検討
