# CMS API 実装TODO

## 概要

- 全タスク数: 19
- 推定作業時間: 42時間
- クリティカルパス: TASK-001 → TASK-002 → TASK-003 → TASK-090 → TASK-101 → TASK-102 → TASK-103 → TASK-201 → TASK-301
- 参照要件: REQ-001〜REQ-007, REQ-101〜REQ-105, REQ-201〜REQ-202, REQ-401〜REQ-404, NFR-001〜NFR-403
- 設計文書: Clean Architecture + サーバーレス構成、Aurora PostgreSQL 15、RESTful API、構造化ログ

## todo

### フェーズ1: 基盤構築

- [ ] **TASK-001 [DIRECT]**: データベース環境設定 (REQ-405, REQ-001, REQ-002対応)
  - [ ] Aurora PostgreSQL 15のスキーマ作成（database-schema.sqlから抽出）
  - [ ] 基本テーブル作成（contents, content_types, content_blocks, content_block_data）
  - [ ] 基本インデックス・ビューの作成（content_details, published_contents）
  - [ ] 初期テストデータの投入（サンプルコンテンツタイプとデータ）
  - [ ] データベース接続テスト（Aurora Serverless v2対応）
  - [ ] スキーマ検証テスト（整合性チェック）
  - [ ] 初期データ確認テスト（データ投入確認）
  - [ ] 受け入れ基準：基本テーブルが作成され、テストデータが適切に投入されている（requirements.mdから抽出）
  - [ ] 完了条件：基本ビューが正常に動作し、データベース接続が確立されている

- [ ] **TASK-002 [DIRECT]**: Go アプリケーション基本設定 (REQ-401, REQ-003, REQ-004対応)
  - [ ] Go 1.24.1の環境設定（architecture.mdから抽出）
  - [ ] Echo v4.13.4の導入とWebフレームワーク設定
  - [ ] Clean Architecture構造の作成（internal/domain, usecase, infrastructure, di）
  - [ ] 依存関係注入の実装（DI層の設計）
  - [ ] 環境変数設定（DB_SECRET_ARN, DB_CLUSTER_ARN, LOG_LEVEL）
  - [ ] Go アプリケーション起動テスト（:8080でのローカル開発サーバー）
  - [ ] 依存関係注入テスト（DI設定の検証）
  - [ ] 環境変数読み込みテスト（設定値の正常取得確認）
  - [ ] 受け入れ基準：Clean Architecture構造が実装され、環境変数が正しく読み込まれる（requirements.mdから抽出）
  - [ ] 完了条件：アプリケーションが正常に起動し、依存関係が適切に注入されている

- [ ] **TASK-003 [DIRECT]**: AWS Lambda 設定 (REQ-402, REQ-004対応)
  - [ ] Lambda Handler実装（cmd/lambda/main.go - サーバーレス実行環境）
  - [ ] スタンドアロン版実装（cmd/main.go - ローカル開発環境:8080）
  - [ ] Lambda デプロイ設定（メモリ512MB-1GB、タイムアウト30秒、同時実行数1000）
  - [ ] 環境変数設定（Secrets Manager連携、Aurora接続設定）
  - [ ] Lambda Handler動作テスト（AWS環境での実行確認）
  - [ ] ローカル開発サーバー起動テスト（:8080での動作確認）
  - [ ] 環境変数取得テスト（Secrets Manager認証情報取得）
  - [ ] 受け入れ基準：Lambda環境とローカル環境の両方で正常動作する（requirements.mdから抽出）
  - [ ] 完了条件：デプロイ設定が完了し、ハイブリッド環境対応が実現されている

### フェーズ2: ドメインモデル実装（MVP版）

- [ ] **TASK-090 [TDD]**: ドメインモデル実装（MVP版） (REQ-001, REQ-002, Clean Architecture原則対応)
  - [ ] Content構造体とビジネスルール実装（internal/domain/entity/content.go）
  - [ ] ContentType構造体実装（内部のコンテンツタイプ管理）
  - [ ] ContentBlock構造体実装（ブロックベースコンテンツ設計）
  - [ ] ContentBlockData構造体実装（柔軟なデータ型対応）
  - [ ] ドメインバリデーションロジック（UUID検証、ステータス検証、title検証）
  - [ ] ドメイン固有の値オブジェクト実装（ContentStatus, BlockType, DataType）
  - [ ] 単体テスト：Content エンティティ（ビジネスルール検証）
  - [ ] 単体テスト：ContentType エンティティ（タイプ定義検証）
  - [ ] 単体テスト：ContentBlock エンティティ（ブロック構造検証）
  - [ ] 単体テスト：ContentBlockData エンティティ（データ型検証）
  - [ ] 単体テスト：ドメインバリデーション（エラーケース処理）
  - [ ] 単体テスト：値オブジェクト（不変性、値の等価性）
  - [ ] エラーハンドリング：無効なUUID、ステータス値の検証（Edgeケースから抽出）
  - [ ] 受け入れ基準：全ドメインエンティティが実装され、ビジネスルールが適切に動作する（requirements.mdから抽出）

### フェーズ3: コアAPI実装

- [ ] **TASK-101 [TDD]**: ヘルスチェックエンドポイント実装 (REQ-201, NFR-403対応)
  - [ ] GET /healthcheck エンドポイント実装（api-endpoints.mdから抽出）
  - [ ] ヘルスチェックユースケース実装（internal/usecase/healthcheck/）
  - [ ] データベース接続確認機能（Aurora Serverless v2対応）
  - [ ] システム状態チェック（services: aurora, secretsManager）
  - [ ] レスポンス形式の実装（JSON形式、status, timestamp, version, environment）
  - [ ] 単体テスト：ヘルスチェックユースケース（正常時/異常時）
  - [ ] 単体テスト：ヘルスチェックコントローラー（HTTPハンドリング）
  - [ ] 統合テスト：エンドポイント動作確認（Echo + Aurora接続）
  - [ ] エラーハンドリングテスト：DB接続失敗時（503 Service Unavailable）
  - [ ] エラーハンドリング：データベース接続失敗（503 Service Unavailable）
  - [ ] エラーハンドリング：部分的サービス障害（200 OK詳細情報付き）
  - [ ] エラーハンドリング：予期しないエラー（500 Internal Server Error）
  - [ ] 受け入れ基準：ヘルスチェックエンドポイントが正常動作し、適切なHTTPステータスコードが返却される（requirements.mdから抽出）

- [ ] **TASK-102 [TDD]**: コンテンツ詳細取得API実装 (REQ-001, REQ-101対応)
  - [ ] GET /contents/{id} エンドポイント実装（api-endpoints.mdから抽出）
  - [ ] コンテンツ詳細取得ユースケース実装（internal/usecase/content/）
  - [ ] コンテンツリポジトリ実装（Aurora PostgreSQL対応、internal/infrastructure/repository/）
  - [ ] パスパラメータ検証（UUID形式、正規表現バリデーション）
  - [ ] ブロックデータ取得・構築（ContentBlock + ContentBlockData統合）
  - [ ] レスポンス形式の実装（JSON、metadata + dynamic blocks構造）
  - [ ] 単体テスト：コンテンツ取得ユースケース（正常系/異常系）
  - [ ] 単体テスト：コンテンツリポジトリ（データアクセス層）
  - [ ] 単体テスト：コンテンツコントローラー（HTTPハンドリング、Echoルーティング）
  - [ ] 統合テスト：エンドポイント動作確認（E2Eテスト）
  - [ ] 境界値テスト：有効/無効なUUID（正規表現チェック）
  - [ ] エラーハンドリング：無効なUUID形式（400 Bad Request）
  - [ ] エラーハンドリング：コンテンツが見つからない（404 Not Found）
  - [ ] エラーハンドリング：データベース接続エラー（500 Internal Server Error）
  - [ ] エラーハンドリング：タイムアウト（504 Gateway Timeout）
  - [ ] 受け入れ基準：正常なコンテンツ詳細取得ができ、エラーケースが適切に処理される（requirements.mdから抽出）

- [ ] **TASK-103 [TDD]**: コンテンツ一覧取得API実装 (REQ-002, REQ-102対応)
  - [ ] GET /contents エンドポイント実装（api-endpoints.mdから抽出）
  - [ ] コンテンツ一覧取得ユースケース実装（ページネーション対応）
  - [ ] クエリパラメータ処理（limit, offset, status, category, tags, search, sort, order）
  - [ ] ページネーション実装（currentPage, totalCount, totalPages, hasPrev, hasNext）
  - [ ] フィルタリング機能実装（status, category, tagsによる絞り込み）
  - [ ] ソート機能実装（createdAt, updatedAt, publishedAt, title）
  - [ ] 単体テスト：コンテンツ一覧取得ユースケース（パラメータ組み合わせ）
  - [ ] 単体テスト：ページネーションロジック（境界値、計算精度）
  - [ ] 単体テスト：フィルタリングロジック（複数条件組み合わせ）
  - [ ] 単体テスト：ソートロジック（昇順/降順、複数ソートキー）
  - [ ] 統合テスト：エンドポイント動作確認（全パラメータ組み合わせ）
  - [ ] 境界値テスト：limit/offsetの境界値（0, 1, 100, 負の値）
  - [ ] エラーハンドリング：無効なクエリパラメータ（400 Bad Request）
  - [ ] エラーハンドリング：limit範囲外の値（デフォルト値使用、1-100制限）
  - [ ] エラーハンドリング：offset負の値（0として扱う）
  - [ ] エラーハンドリング：データベースエラー（500 Internal Server Error）
  - [ ] 受け入れ基準：コンテンツ一覧取得、ページネーション、フィルタリング・ソート機能が正常動作する（requirements.mdから抽出）

### フェーズ4: インフラストラクチャ実装

- [ ] **TASK-201 [DIRECT]**: API Gateway統合 (REQ-003, REQ-403, REQ-404対応)
  - [ ] API Gateway REST APIの設定（リージョナルエンドポイント）
  - [ ] Lambda関数との統合（プロキシ統合、レスポンス変換）
  - [ ] CORS設定（Access-Control-Allow-Origin: *, Methods: GET, OPTIONS）
  - [ ] カスタムドメイン設定（HTTPS通信、TLS 1.2以上）
  - [ ] レスポンス変換設定（エラーマッピング、ステータスコード）
  - [ ] API Gateway経由でのエンドポイントアクセステスト
  - [ ] CORS設定確認テスト（ブラウザからのアクセス検証）
  - [ ] エラーレスポンス形式テスト（4xx/5xxエラーハンドリング）
  - [ ] 受け入れ基準：API Gateway経由でAPIにアクセスでき、CORS設定が適切に動作する（requirements.mdから抽出）
  - [ ] 完了条件：HTTPSでのアクセスが可能で、エラーレスポンスが統一形式である

- [ ] **TASK-202 [DIRECT]**: AWS WAF セキュリティ設定 (NFR-102, NFR-104対応)
  - [ ] AWS WAF Web ACLの作成（DDoS攻撃防御）
  - [ ] マネージドルールの設定（SQLインジェクション、XSS攻撃防御）
  - [ ] レート制限ルールの設定（100 req/min基本制限）
  - [ ] SQLインジェクション防御（AWS Managed Rules適用）
  - [ ] XSS攻撃防御（入力値サニタイズ）
  - [ ] 正常リクエストの通過確認（許可リクエストのテスト）
  - [ ] 悪意のあるリクエストのブロック確認（攻撃パターンテスト）
  - [ ] レート制限の動作確認（スループット制限テスト）
  - [ ] 受け入れ基準：WAFルールが正常に動作し、セキュリティ要件を満たしている（requirements.mdから抽出）
  - [ ] 完了条件：レート制限が適切に機能し、攻撃パターンがブロックされている

- [ ] **TASK-203 [DIRECT]**: Secrets Manager統合 (NFR-103対応)
  - [ ] データベース認証情報のSecrets Manager管理
  - [ ] Lambda関数からの認証情報取得（GetSecretValue API）
  - [ ] 自動ローテーション設定（30日間隔）
  - [ ] アクセス権限設定（Lambda実行ロールからのみアクセス可能）
  - [ ] 認証情報取得テスト（Lambda関数からの正常取得）
  - [ ] 権限設定確認テスト（不正アクセス防止）
  - [ ] ローテーション動作テスト（自動更新機能）
  - [ ] 受け入れ基準：データベース認証情報が安全に管理され、自動ローテーションが設定されている（requirements.mdから抽出）
  - [ ] 完了条件：Lambda関数から認証情報を取得でき、セキュリティが確保されている

### フェーズ5: 監視・ログ実装

- [ ] **TASK-301 [DIRECT]**: CloudWatch 監視設定 (NFR-401, NFR-402対応)
  - [ ] CloudWatch Logs設定（ログストリーム作成、保持期間設定）
  - [ ] 構造化ログ実装（JSON形式、requestId, timestamp, level, message）
  - [ ] カスタムメトリクス実装（レスポンス時間、エラー率、リクエスト数）
  - [ ] ダッシュボード作成（Lambda, API Gateway, Aurora メトリクス可視化）
  - [ ] アラート設定（エラー率5%超過、レスポンス時間1秒超過）
  - [ ] ログ出力確認テスト（ERROR, WARN, INFO, DEBUG各レベル）
  - [ ] メトリクス収集確認テスト（カスタムメトリクス動作）
  - [ ] アラート動作確認テスト（閾値超過時の通知）
  - [ ] 受け入れ基準：適切なログが出力され、メトリクスが収集され、アラートが正常に動作する（requirements.mdから抽出）
  - [ ] 完了条件：監視体制が構築され、システム状態が可視化されている

- [ ] **TASK-302 [DIRECT]**: X-Ray トレーシング設定 (NFR-402対応)
  - [ ] X-Ray SDKの統合（Go用AWS X-Ray SDK導入）
  - [ ] トレーシングの実装（API Gateway → Lambda → Aurora フロー）
  - [ ] セグメント・サブセグメント設定（詳細な処理時間計測）
  - [ ] サービスマップの構築（依存関係の可視化）
  - [ ] トレース情報の収集確認（リクエストフローの追跡）
  - [ ] サービスマップの表示確認（コンポーネント間の関係）
  - [ ] パフォーマンス分析機能確認（ボトルネック特定）
  - [ ] 受け入れ基準：トレーシングが正常に動作し、サービスマップが表示され、パフォーマンス分析ができる（requirements.mdから抽出）
  - [ ] 完了条件：分散トレーシングが機能し、システムの可観測性が向上している

### フェーズ6: パフォーマンス最適化・品質保証

- [ ] **TASK-401 [TDD]**: パフォーマンステスト実装 (NFR-001, NFR-002, NFR-003対応)
  - [ ] 負荷テストツールの設定（Apache JMeter、K6等）
  - [ ] レスポンス時間測定（95%のリクエストで1秒以内）
  - [ ] 同時接続数テスト（1000同時接続処理）
  - [ ] データベースクエリ最適化（インデックス追加、クエリチューニング）
  - [ ] コネクションプール最適化（Lambda関数内プール設定）
  - [ ] 1秒以内レスポンス確認テスト（NFR-001達成検証）
  - [ ] 同時接続数1000テスト（NFR-002達成検証）
  - [ ] データベースクエリパフォーマンステスト（500ms以内完了）
  - [ ] 受け入れ基準：NFR-001〜003の性能要件を満たし、パフォーマンステストが自動実行できる（requirements.mdから抽出）
  - [ ] 完了条件：最適化による改善が確認でき、性能基準をクリアしている

- [ ] **TASK-402 [TDD]**: セキュリティテスト実装 (NFR-101, NFR-102, NFR-103, NFR-104対応)
  - [ ] セキュリティテストスイート作成（自動化テスト環境）
  - [ ] SQLインジェクションテスト（攻撃パターンでの検証）
  - [ ] XSS攻撃テスト（スクリプトインジェクション防御）
  - [ ] 認証・認可テスト（将来の認証機能準備）
  - [ ] 機密情報漏洩チェック（ログ出力内容検証）
  - [ ] SQLインジェクション防御確認（パラメータ化クエリ検証）
  - [ ] XSS攻撃防御確認（入力値サニタイズ検証）
  - [ ] ログ出力セキュリティチェック（機密情報非出力確認）
  - [ ] CORS設定確認（適切なオリジン制御）
  - [ ] 受け入れ基準：セキュリティ要件を満たし、脆弱性が検出・修正され、セキュリティテストが自動実行できる（requirements.mdから抽出）
  - [ ] 完了条件：セキュリティ基準をクリアし、攻撃に対する防御が機能している

- [ ] **TASK-403 [TDD]**: E2Eテストスイート (全要件対応)
  - [ ] E2Eテストフレームワーク導入（Go用テストツール）
  - [ ] 主要ユーザーフローテスト（コンテンツ取得シナリオ）
  - [ ] エラーシナリオテスト（404, 400, 500エラーハンドリング）
  - [ ] 障害復旧テスト（データベース再接続、Lambda再起動）
  - [ ] CI/CD統合（GitHub Actions、自動テスト実行）
  - [ ] 正常フローのE2Eテスト（全エンドポイント連携）
  - [ ] エラーハンドリングE2Eテスト（異常系フロー）
  - [ ] 障害シナリオテスト（システム復旧確認）
  - [ ] 可用性テスト（99.9%可用性検証）
  - [ ] 受け入れ基準：E2Eテストスイートが完成し、CI/CDパイプラインに統合され、全ての受け入れ基準を満たしている（requirements.mdから抽出）
  - [ ] 完了条件：自動テストが機能し、システムの品質が保証されている

### フェーズ7: 運用準備・文書化

- [ ] **TASK-501 [DIRECT]**: IaC実装（Terraform） (全インフラ要件対応)
  - [ ] Terraform設定ファイルの作成・更新（.tf形式、リソース定義）
  - [ ] 環境別設定（dev/staging/prod環境分離）
  - [ ] 既存Terraformコードとの統合（プロジェクト構成統一）
  - [ ] 状態ファイル管理（S3バックエンド、DynamoDBロック）
  - [ ] Terraform plan/apply動作確認（リソース作成確認）
  - [ ] 環境別デプロイ確認（各環境での動作検証）
  - [ ] リソース作成・削除確認（terraform destroy検証）
  - [ ] 受け入れ基準：全インフラがTerraformで管理され、環境別デプロイが可能で、既存のTerraform構成と統合されている（requirements.mdから抽出）
  - [ ] 完了条件：Infrastructure as Codeが完成し、運用効率が向上している

- [ ] **TASK-502 [DIRECT]**: CI/CD パイプライン構築 (運用要件対応)
  - [ ] GitHub Actions ワークフロー作成（.github/workflows/）
  - [ ] ビルド・テスト自動化（golangci-lint, go test, go build）
  - [ ] 段階的デプロイメント（dev → staging → prod）
  - [ ] ロールバック機能（前バージョンへの自動復旧）
  - [ ] ビルド・テスト自動実行確認（プルリクエスト連携）
  - [ ] デプロイメント動作確認（環境別自動デプロイ）
  - [ ] ロールバック動作確認（障害時の自動復旧）
  - [ ] 受け入れ基準：CI/CDパイプラインが正常動作し、自動テスト・デプロイが機能し、ロールバック機能が利用可能である（requirements.mdから抽出）
  - [ ] 完了条件：開発効率が向上し、デプロイメントの安全性が確保されている

- [ ] **TASK-503 [DIRECT]**: 運用文書作成 (運用要件対応)
  - [ ] 運用マニュアル作成（システム概要、運用手順、トラブルシューティング）
  - [ ] 障害対応手順書作成（エスカレーション、復旧手順）
  - [ ] 監視・アラート手順書作成（メトリクス確認、対応フロー）
  - [ ] デプロイメント手順書作成（リリースプロセス、ロールバック手順）
  - [ ] 手順書の実行可能性確認（実際の手順での動作検証）
  - [ ] 緊急時対応手順確認（障害シナリオでの対応テスト）
  - [ ] 受け入れ基準：運用マニュアルが完成し、障害対応手順が明文化され、手順書通りに操作が実行できる（requirements.mdから抽出）
  - [ ] 完了条件：運用体制が整備され、安定したサービス提供が可能になっている

## 実行順序

1. **基盤構築** (TASK-001, TASK-002, TASK-003) - 理由：データベース環境とGo環境、Lambda環境の基盤が他のタスクの前提条件
2. **ドメイン実装** (TASK-090) - 理由：Clean Architectureの中心となるドメイン層が、上位層の実装の前提
3. **コアAPI実装** (TASK-101, TASK-102, TASK-103) - 理由：ヘルスチェック → コンテンツ詳細 → 一覧の順でAPI機能を段階的に構築
4. **インフラ実装** (TASK-201, TASK-202, TASK-203) - 理由：API Gateway、WAF、Secrets Managerがアプリケーションの本格運用に必要
5. **監視・ログ** (TASK-301, TASK-302) - 理由：CloudWatch、X-Rayがシステムの可観測性向上に必要
6. **品質保証** (TASK-401, TASK-402, TASK-403) - 理由：パフォーマンス、セキュリティ、E2Eテストで品質確保
7. **運用準備** (TASK-501, TASK-502, TASK-503) - 理由：Terraform、CI/CD、文書化でプロダクション運用準備

## 実装プロセス

### TDDタスクの実装プロセス

[TDD]タスクは以下の順序で実装:

1. `/{taskID}/tdd-requirements.md` - 詳細要件定義（要件定義文書REQ-XXXから抽出）
2. `/{taskID}/tdd-testcases.md` - テストケース作成（受け入れ基準とEdgeケース EDGE-XXXから導出）
3. `/{taskID}/tdd-red.md` - テスト実装（失敗テストの作成）
4. `/{taskID}/tdd-green.md` - 最小実装（アーキテクチャ設計Clean Architectureに準拠）
5. `/{taskID}/tdd-refactor.md` - リファクタリング（設計文書architecture.mdとの整合性確認）
6. `/{taskID}/tdd-verify-complete.md` - 品質確認（要件定義の受け入れ基準で検証）

### DIRECTタスクの実装プロセス

[DIRECT]タスクは以下の順序で実装:

1. `/{taskID}/direct-setup.md` - 設定作業の実行（設計文書に基づく直接実装）
2. `/{taskID}/direct-verify.md` - 設定確認（動作確認とテスト実行）

## 文書との連携

- **cms-api-requirements.md**: 機能要件（REQ-001〜REQ-404）、非機能要件（NFR-001〜NFR-403）、受け入れ基準、Edgeケース（EDGE-001〜EDGE-302）
- **architecture.md**: Clean Architecture + サーバーレス構成の全体的な実装方針、Go 1.24.1 + Echo v4.13.4 + Aurora Serverless v2
- **database-schema.sql**: データベース関連タスクの詳細実装（Aurora PostgreSQL 15、UUID、content_types, contents, content_blocks, content_block_data）
- **api-endpoints.md**: API実装タスクの仕様（RESTful API、JSON形式、HTTPSステータスコード、CORS設定）
- **dataflow.md**: データ処理フローと統合テストシナリオ（API Gateway → Lambda → Aurora、エラーハンドリング）

## 品質保証コマンド

### 必須実行コマンド（各タスク完了時）
```bash
# リンター実行
golangci-lint run

# コードフォーマット
go fmt ./...

# 依存関係整理
go mod tidy

# 全テスト実行
make test

# ビルド確認
go build cmd/main.go
go build cmd/lambda/main.go
```

### 推奨実行コマンド
```bash
# モック更新
make mock

# カバレッジ確認
go test ./... -cover

# サーバー起動テスト
make run
```

## MVP版ドメインモデル詳細

本実装では以下のシンプルで拡張可能なドメインモデルを採用：

### Content（コンテンツ）
- ID: UUID（一意識別子）
- ContentTypeID: UUID（コンテンツタイプ参照）
- Title: string（タイトル、必須）
- Slug: string（URL用スラッグ）
- Status: ContentStatus（draft, published, archived, trash）
- 時刻管理: CreatedAt, UpdatedAt, PublishedAt
- AuthorID: string（作成者ID）
- Version: int（バージョン管理）
- Blocks: []ContentBlock（ブロック構造）

### ContentType（コンテンツタイプ）
- 動的なコンテンツタイプ定義
- Name, DisplayName, Description
- Icon, IsActive（アクティブ状態管理）
- 作成者・時刻管理

### ContentBlock（コンテンツブロック）
- ブロックベースCMSの基本単位
- BlockType（text, richtext, image, video等）
- BlockOrder（表示順序）
- IsVisible（表示制御）
- Data: ContentBlockData（実際のデータ）

### ContentBlockData（ブロックデータ）
- 多様なデータ型対応
- ContentText, ContentRichtext, ContentNumber, ContentURL
- ContentJSON（柔軟なJSONデータ）
- ReferencedContentID（他コンテンツ参照）
- Settings（ブロック固有設定）

## 完了確認

このタスクリストの完了により、以下の成果物が得られます：

- [ ] AWS Lambda + Aurora Serverless v2によるCMS API（サーバーレス構成）
- [ ] Clean Architectureに基づく高品質なGoアプリケーション（Go 1.24.1 + Echo v4.13.4）
- [ ] MVPドメインモデルによるシンプルで拡張可能な設計（ブロックベースCMS）
- [ ] セキュアで高可用性なサーバーレスシステム（AWS WAF + Secrets Manager）
- [ ] 包括的なテストスイート（単体・統合・E2Eテスト）
- [ ] 監視・運用体制の構築（CloudWatch + X-Ray）
- [ ] CI/CDによる自動化されたデプロイメント（GitHub Actions）
- [ ] 運用マニュアル・手順書の完備（障害対応・監視手順）

## 将来拡張予定

以下の機能は将来のフェーズで実装予定：

- [ ] ContentRelation（コンテンツ間関係）実装
- [ ] ContentReferenceHierarchy（参照階層）実装
- [ ] 動的フィールド機能の拡張（カスタムフィールドタイプ）
- [ ] 認証・認可機能（API Key, JWT Bearer Token, OAuth 2.0）
- [ ] 全文検索機能（Elasticsearch連携）
- [ ] キャッシュ機能（Redis、API Gatewayキャッシュ）
- [ ] GraphQL対応検討

## 注意事項

- 各タスクは依存関係に従って順次実行してください
- **ドメインモデル実装（TASK-090）を最優先**で実行し、Clean Architectureの依存関係を遵守してください
- TDDプロセスを遵守し、テストファーストで開発を進めてください
- セキュリティ要件（NFR-101〜NFR-104）・非機能要件（NFR-001〜NFR-403）を常に意識して実装してください
- 実装中に設計変更が必要な場合は、設計文書も併せて更新してください
- **各タスク完了時は必ず上記の品質チェックコマンドを実行してください**
- Edgeケース（EDGE-001〜EDGE-302）を考慮した実装とテストを行ってください