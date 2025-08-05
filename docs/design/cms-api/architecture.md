# CMS API アーキテクチャ設計

## システム概要

クライアントサイドから呼び出すCMS用のAPIシステム。コンテンツ管理機能を提供し、高可用性・高拡張性を実現するクラウドネイティブなアーキテクチャを採用する。

## アーキテクチャパターン

- **パターン**: サーバーレス + マイクロサービスアーキテクチャ
- **理由**: 
  - 自動スケーリングによる高可用性の実現
  - 運用コストの最適化（従量課金）
  - 単一障害点の排除
  - 開発・デプロイの独立性確保

## システム構成図

```mermaid
graph TB
    Client[クライアントアプリケーション]
    
    subgraph "AWS Cloud"
        subgraph "外部向けAPI層"
            APIGW[API Gateway<br/>REST API]
            WAF[AWS WAF]
        end
        
        subgraph "アプリケーション層"
            Lambda[Lambda Function<br/>CMS API (Echo)]            
        end
        
        subgraph "データ層"
            Aurora[Aurora Serverless v2<br/>PostgreSQL]
            Secret[Secrets Manager<br/>DB認証情報]
        end
        
        subgraph "監視・ログ"
            CloudWatch[CloudWatch<br/>Logs & Metrics]
        end
    end
    
    Client --> WAF
    WAF --> APIGW
    APIGW --> Lambda
    Lambda --> Aurora
    Lambda --> Secret
    Lambda --> CloudWatch
```

## コンポーネント詳細

### API Gateway
- **役割**: 
  - HTTPSエンドポイントの提供
  - リクエストルーティング
  - レート制限・スロットリング
  - CORS設定
  - リクエスト/レスポンス変換
- **設定**:
  - REST API形式
  - リージョナルエンドポイント
  - カスタムドメイン対応
  - ステージ別デプロイ（dev/staging/prod）

### AWS WAF
- **役割**:
  - DDoS攻撃の防御
  - SQLインジェクション防御
  - XSS攻撃防御
  - 地理的アクセス制御
- **ルール**:
  - AWS Managed Rules適用
  - レート制限ルール
  - カスタムIP許可/拒否リスト

### Lambda Function (CMS API)
- **実装言語**: Go 1.21
- **Webフレームワーク**: Echo v4
- **実行環境**: 
  - メモリ: 512MB-1GB（トラフィックに応じて調整）
  - タイムアウト: 30秒
  - 同時実行数制限: 1000
- **環境変数**:
  - `DB_SECRET_ARN`: データベース認証情報のSecrets Manager ARN
  - `DB_CLUSTER_ARN`: Aurora クラスター ARN
  - `LOG_LEVEL`: ログレベル設定

#### 単一Lambda構成の利点
- **コールドスタート最適化**: 1つのLambda関数による初期化時間の削減
- **リソース効率**: 共通の依存関係とコネクションプールの共有
- **デプロイ簡素化**: 単一のデプロイメントユニット
- **運用コスト削減**: Lambda関数数の削減による管理コスト低減

#### 提供エンドポイント
- **`GET /contents/{id}`**: 特定コンテンツの詳細情報取得
- **`GET /contents`**: コンテンツ一覧の取得（ページネーション対応）
- **`GET /healthcheck`**: システムヘルスチェック

### Aurora Serverless v2
- **エンジン**: PostgreSQL 15
- **設定**:
  - 最小ACU: 0.5
  - 最大ACU: 16
  - 自動一時停止: 5分間非アクティブ後
  - バックアップ保持期間: 7日間
- **セキュリティ**:
  - VPC内プライベートサブネット配置
  - セキュリティグループによるアクセス制御
  - データ暗号化（保存時・転送時）

### Secrets Manager
- **役割**: データベース認証情報の安全な管理
- **ローテーション**: 30日間隔での自動ローテーション
- **アクセス制御**: Lambda実行ロールからのみアクセス可能

## セキュリティ設計

### 認証・認可
- **現フェーズ**: パブリックAPI（認証なし）
- **将来対応**: 
  - API Key認証
  - JWT Bearer Token認証
  - IAM認証（AWS内部サービス向け）

### ネットワークセキュリティ
- **HTTPS通信**: 全通信をTLS 1.2以上で暗号化
- **VPC設計**: 
  - パブリックサブネット: API Gateway（該当なし、マネージドサービス）
  - プライベートサブネット: Aurora Serverless v2
- **セキュリティグループ**: 最小権限の原則でアクセス制御

### データ保護
- **暗号化**: 
  - 保存時: Aurora暗号化（KMS）
  - 転送時: TLS暗号化
- **機密情報管理**: Secrets Managerでの認証情報管理
- **ログセキュリティ**: 機密情報のログ出力防止

## 可用性・拡張性設計

### 高可用性
- **マルチAZ構成**: Aurora Serverless v2の自動マルチAZ
- **自動フェイルオーバー**: Aurora の自動フェイルオーバー機能
- **ヘルスチェック**: Lambda関数レベルでのヘルスチェック実装

### 自動スケーリング
- **Lambda**: リクエスト数に応じた自動スケーリング
- **Aurora**: ACUの自動スケーリング（負荷に応じて0.5-16ACU）
- **API Gateway**: AWSマネージドサービスによる自動スケーリング

### パフォーマンス最適化
- **コネクションプール**: Lambda関数内でのデータベースコネクションプール
- **クエリ最適化**: インデックス設計による高速クエリ実行
- **レスポンスキャッシュ**: 必要に応じてAPI Gatewayキャッシュ機能

## 監視・運用

### ログ戦略
- **アプリケーションログ**: CloudWatch Logs
- **アクセスログ**: API Gateway アクセスログ
- **ログレベル**: ERROR, WARN, INFO, DEBUG
- **構造化ログ**: JSON形式での出力

### メトリクス監視
- **Lambda メトリクス**: 実行時間、エラー率、スロットル数
- **API Gateway メトリクス**: リクエスト数、レイテンシ、エラー率
- **Aurora メトリクス**: CPU使用率、コネクション数、ACU使用率
- **カスタムメトリクス**: ビジネスロジック固有のメトリクス

### アラート設定
- **エラー率**: 5%を超えた場合
- **レスポンス時間**: 1秒を超えた場合
- **Aurora ACU**: 上限に近づいた場合
- **Lambda スロットル**: 発生した場合

### 分散トレーシング
- **X-Ray**: リクエストフローの可視化
- **トレース対象**: API Gateway → Lambda → Aurora

## デプロイ戦略

### インフラストラクチャ as Code
- **IaC ツール**: AWS CDK (TypeScript)
- **環境分離**: dev/staging/prod環境の分離
- **CI/CD**: GitHub Actions による自動デプロイ

### デプロイフロー
1. **開発環境**: feature ブランチでの開発・テスト
2. **ステージング環境**: main ブランチマージ時の自動デプロイ
3. **本番環境**: タグ作成時の手動承認後デプロイ

### ブルーグリーンデプロイ
- **API Gateway**: ステージを利用したブルーグリーンデプロイ
- **Lambda**: エイリアスを利用した段階的トラフィック移行
- **データベース**: マイグレーション戦略による後方互換性確保

## コスト最適化

### 従量課金の活用
- **Lambda**: 実行時間分のみ課金
- **Aurora Serverless v2**: 使用ACU分のみ課金
- **API Gateway**: リクエスト数に応じた課金

### コスト監視
- **AWS Cost Explorer**: 月次コスト分析
- **予算アラート**: 予想コスト超過時のアラート
- **リソース最適化**: 定期的な使用量レビューと最適化
