# CMS API プロジェクト概要

## プロジェクトの目的
このプロジェクトは、ブログ管理システム（Asunaro Blog Admin）のCMS APIサーバーです。記事（Article）を管理するためのRESTful APIを提供します。

## 技術スタック
- **言語**: Go 1.24.1
- **Webフレームワーク**: Echo v4.13.4
- **データベース**: AWS DynamoDB
- **AWS**: AWS Lambda, AWS SDK for Go v2
- **テスト**: testify, testcontainers-go
- **モック生成**: mockery v3.3.2
- **リンター**: golangci-lint v2.2.1

## アーキテクチャ
Clean Architecture パターンを採用:
- `internal/domain/entity`: エンティティ層（Article構造体）
- `internal/usecase`: ユースケース層（ビジネスロジック）
- `internal/infrastructure`: インフラストラクチャ層（リポジトリ、コントローラー）
- `internal/di`: 依存性注入とルーティング設定
- `cmd`: エントリーポイント

## エントリーポイント
- `cmd/main.go`: APIサーバーのメインエントリーポイント（:8080でListen）
- `cmd/lambda/main.go`: AWS Lambda用エントリーポイント

## 主要エンティティ
Article構造体:
- ID, CreatedAt, UpdatedAt
- Title, Description, Body
- CoverImage, PublishedAt, Status
- CategoryID, Tags