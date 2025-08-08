# データフロー図

## ユーザーインタラクションフロー

### システム全体フロー

```mermaid
flowchart TD
    User[ユーザー] --> Frontend[フロントエンドアプリ]
    Frontend --> |HTTPS Request| WAF[AWS WAF]
    WAF --> |Security Check| APIGW[API Gateway]
    
    APIGW --> |"Route /contents/{id}"| DetailLambda[Content Detail Lambda]
    
    DetailLambda --> |SQL Query| Aurora[Aurora Serverless v2]
    
    DetailLambda --> |Get Credentials| Secrets[Secrets Manager]
    
    Aurora --> |Query Result| DetailLambda
    
    DetailLambda --> |JSON Response| APIGW
    
    APIGW --> |HTTPS Response| Frontend
    Frontend --> |Display Content| User
    
    DetailLambda --> |Logs & Metrics| CloudWatch[CloudWatch]
```

## データ処理フロー

### コンテンツ詳細取得フロー

```mermaid
sequenceDiagram
    participant U as ユーザー
    participant F as フロントエンド
    participant W as AWS WAF
    participant A as API Gateway
    participant L as Detail Lambda
    participant S as Secrets Manager
    participant D as Aurora DB
    participant C as CloudWatch
    
    U->>F: コンテンツ詳細要求
    F->>W: GET /contents/{id}
    W->>A: セキュリティチェック通過
    A->>L: Lambda関数呼び出し
    
    L->>S: DB認証情報取得
    S-->>L: 認証情報返却
    
    L->>D: SELECT * FROM contents WHERE id = ?
    
    alt コンテンツが存在する場合
        D-->>L: コンテンツデータ返却
        L->>C: 成功ログ出力
        L-->>A: 200 OK + コンテンツデータ
        A-->>F: JSONレスポンス
        F-->>U: コンテンツ表示
    else コンテンツが存在しない場合
        D-->>L: 空の結果セット
        L->>C: 404エラーログ出力
        L-->>A: 404 Not Found
        A-->>F: エラーレスポンス
        F-->>U: エラーメッセージ表示
    end
```

### コンテンツ一覧取得フロー

```mermaid
sequenceDiagram
    participant U as ユーザー
    participant F as フロントエンド
    participant W as AWS WAF
    participant A as API Gateway
    participant L as List Lambda
    participant S as Secrets Manager
    participant D as Aurora DB
    participant C as CloudWatch
    
    U->>F: コンテンツ一覧要求
    F->>W: GET /contents?limit=20&offset=0
    W->>A: セキュリティチェック通過
    A->>L: Lambda関数呼び出し
    
    L->>L: パラメータ検証・デフォルト値設定
    
    alt パラメータが有効な場合
        L->>S: DB認証情報取得
        S-->>L: 認証情報返却
        
        L->>D: SELECT COUNT(*) FROM contents
        D-->>L: 総件数返却
        
        L->>D: SELECT * FROM contents LIMIT ? OFFSET ?
        D-->>L: コンテンツ一覧返却
        
        L->>L: ページネーション情報計算
        L->>C: 成功ログ出力
        L-->>A: 200 OK + 一覧データ + ページ情報
        A-->>F: JSONレスポンス
        F-->>U: コンテンツ一覧表示
    else パラメータが無効な場合
        L->>C: 400エラーログ出力
        L-->>A: 400 Bad Request
        A-->>F: エラーレスポンス
        F-->>U: エラーメッセージ表示
    end
```

## エラーハンドリングフロー

### データベースエラー処理

```mermaid
flowchart TD
    Request[APIリクエスト] --> Lambda[Lambda関数]
    Lambda --> DBConnect{DB接続試行}
    
    DBConnect -->|成功| Query[クエリ実行]
    DBConnect -->|接続失敗| ConnError[接続エラー処理]
    
    Query -->|成功| Response[正常レスポンス]
    Query -->|タイムアウト| TimeoutError[タイムアウトエラー]
    Query -->|SQL文エラー| SQLError[SQLエラー処理]
    
    ConnError --> Log1[エラーログ出力]
    TimeoutError --> Log2[タイムアウトログ出力]
    SQLError --> Log3[SQLエラーログ出力]
    
    Log1 --> Return500[500 Internal Server Error]
    Log2 --> Return500
    Log3 --> Return500
    
    Return500 --> Client[クライアントへエラー返却]
```

### Lambda関数エラー処理

```mermaid
flowchart TD
    Trigger[API Gateway Trigger] --> LambdaStart[Lambda関数開始]
    
    LambdaStart --> Process{処理実行}
    
    Process -->|正常| Success[成功レスポンス]
    Process -->|ビジネスロジックエラー| BusinessError[400/404エラー]
    Process -->|予期しないエラー| UnexpectedError[500エラー]
    Process -->|タイムアウト| LambdaTimeout[Lambda タイムアウト]
    
    BusinessError --> LogWarn[WARNレベルログ]
    UnexpectedError --> LogError[ERRORレベルログ]
    LambdaTimeout --> LogError
    
    LogWarn --> ClientError[4xxエラーレスポンス]
    LogError --> ServerError[5xxエラーレスポンス]
    
    Success --> APIGateway[API Gateway]
    ClientError --> APIGateway
    ServerError --> APIGateway
    
    APIGateway --> Client[クライアント]
```

## リトライ・フォールバック戦略

### データベース接続リトライ

```mermaid
sequenceDiagram
    participant L as Lambda
    participant D as Aurora DB
    participant C as CloudWatch
    
    L->>D: 接続試行 (1回目)
    
    alt 接続成功
        D-->>L: 接続確立
    else 接続失敗
        L->>C: 接続失敗ログ
        L->>L: 1秒待機
        L->>D: 接続試行 (2回目)
        
        alt 接続成功
            D-->>L: 接続確立
        else 接続失敗
            L->>C: 接続失敗ログ
            L->>L: 2秒待機
            L->>D: 接続試行 (3回目)
            
            alt 接続成功
                D-->>L: 接続確立
            else 接続失敗
                L->>C: 最終接続失敗ログ
                L->>L: 500エラー生成
            end
        end
    end
```

## パフォーマンス最適化フロー

### コネクションプール管理

```mermaid
stateDiagram-v2
    [*] --> LambdaColdStart: Lambda初回実行
    LambdaColdStart --> CreatePool: コネクションプール作成
    CreatePool --> PoolReady: プール準備完了
    
    PoolReady --> GetConnection: コネクション取得
    GetConnection --> ExecuteQuery: クエリ実行
    ExecuteQuery --> ReleaseConnection: コネクション返却
    ReleaseConnection --> PoolReady: プールに返却
    
    PoolReady --> LambdaWarmStart: 次回実行（Warm Start）
    LambdaWarmStart --> GetConnection: 既存プール利用
    
    PoolReady --> PoolTimeout: 5分間非アクティブ
    PoolTimeout --> [*]: Lambda終了
```

### キャッシュ戦略（将来拡張）

```mermaid
flowchart TD
    Request[APIリクエスト] --> CacheCheck{キャッシュ確認}
    
    CacheCheck -->|Hit| CacheReturn[キャッシュから返却]
    CacheCheck -->|Miss| DBQuery[データベースクエリ]
    
    DBQuery --> DBResult[DB結果取得]
    DBResult --> CacheStore[キャッシュに保存]
    CacheStore --> Response[レスポンス返却]
    
    CacheReturn --> Client[クライアント]
    Response --> Client
    
    CacheStore --> TTL{TTL設定}
    TTL --> |コンテンツ詳細| TTL5min[5分間キャッシュ]
    TTL --> |コンテンツ一覧| TTL1min[1分間キャッシュ]
```

## 監視・アラートフロー

### メトリクス収集フロー

```mermaid
flowchart LR
    Lambda[Lambda実行] --> Metrics[メトリクス生成]
    
    Metrics --> Duration[実行時間]
    Metrics --> Errors[エラー数]
    Metrics --> Invocations[実行回数]
    Metrics --> Throttles[スロットル数]
    
    Duration --> CloudWatch[CloudWatch Metrics]
    Errors --> CloudWatch
    Invocations --> CloudWatch
    Throttles --> CloudWatch
    
    CloudWatch --> Alarm{閾値チェック}
    Alarm -->|閾値超過| SNS[SNS通知]
    Alarm -->|正常範囲| Monitor[継続監視]
    
    SNS --> Email[メール通知]
    SNS --> Slack[Slack通知]
```
