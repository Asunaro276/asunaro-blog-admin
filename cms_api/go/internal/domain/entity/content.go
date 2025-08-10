package entity

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ContentStatus はコンテンツの状態を表す列挙型
type ContentStatus string

const (
	ContentStatusDraft     ContentStatus = "draft"
	ContentStatusPublished ContentStatus = "published"
	ContentStatusArchived  ContentStatus = "archived"
)

// BlockType はブロックの種類を表す列挙型
type BlockType string

const (
	BlockTypeText      BlockType = "text"
	BlockTypeRichText  BlockType = "richtext"
	BlockTypeImage     BlockType = "image"
	BlockTypeVideo     BlockType = "video"
	BlockTypeEmbed     BlockType = "embed"
	BlockTypeReference BlockType = "reference"
)

// DataType はデータの種類を表す列挙型
type DataType string

const (
	DataTypeText      DataType = "text"
	DataTypeRichText  DataType = "richtext"
	DataTypeNumber    DataType = "number"
	DataTypeURL       DataType = "url"
	DataTypeJSON      DataType = "json"
	DataTypeReference DataType = "reference"
)

// Content はコンテンツのドメインエンティティ
type Content struct {
	ID            uuid.UUID     `json:"id"`
	ContentTypeID uuid.UUID     `json:"content_type_id"`
	Title         string        `json:"title"`
	Slug          string        `json:"slug"`
	Status        ContentStatus `json:"status"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
	PublishedAt   *time.Time    `json:"published_at"`
	AuthorID      string        `json:"author_id"`
	Version       int           `json:"version"`
	
	// リレーション
	ContentType *ContentType   `json:"content_type,omitempty"`
	Blocks      []ContentBlock `json:"blocks,omitempty"`
}

// ContentType はコンテンツタイプのドメインエンティティ
type ContentType struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	DisplayName string    `json:"display_name"`
	Description string    `json:"description"`
	Icon        string    `json:"icon"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedBy   string    `json:"created_by"`
}

// ContentBlock はコンテンツブロックのドメインエンティティ
type ContentBlock struct {
	ID         uuid.UUID `json:"id"`
	ContentID  uuid.UUID `json:"content_id"`
	BlockType  BlockType `json:"block_type"`
	BlockOrder int       `json:"block_order"`
	IsVisible  bool      `json:"is_visible"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	
	// リレーション
	Content *Content          `json:"content,omitempty"`
	Data    *ContentBlockData `json:"data,omitempty"`
}

// ContentBlockData はブロックデータのドメインエンティティ
type ContentBlockData struct {
	ID                  uuid.UUID        `json:"id"`
	BlockID             uuid.UUID        `json:"block_id"`
	DataType            DataType         `json:"data_type"`
	ContentText         string           `json:"content_text"`
	ContentRichtext     json.RawMessage  `json:"content_richtext"`
	ContentNumber       *decimal.Decimal `json:"content_number"`
	ContentURL          string           `json:"content_url"`
	ContentJSON         json.RawMessage  `json:"content_json"`
	ReferencedContentID *uuid.UUID       `json:"referenced_content_id"`
	Settings            json.RawMessage  `json:"settings"`
	CreatedAt           time.Time        `json:"created_at"`
	UpdatedAt           time.Time        `json:"updated_at"`
	
	// リレーション
	Block             *ContentBlock `json:"block,omitempty"`
	ReferencedContent *Content      `json:"referenced_content,omitempty"`
}

// IsPublished はコンテンツが公開されているかを確認
func (c *Content) IsPublished() bool {
	return c.Status == ContentStatusPublished && c.PublishedAt != nil
}

// IsActive はコンテンツが有効かを確認
func (c *Content) IsActive() bool {
	return c.Status != ContentStatusArchived
}

// Validate はContentの基本的なバリデーション
func (c *Content) Validate() error {
	if c.Title == "" {
		return fmt.Errorf("タイトルは必須です")
	}
	if c.Slug == "" {
		return fmt.Errorf("スラッグは必須です")
	}
	if c.AuthorID == "" {
		return fmt.Errorf("作成者IDは必須です")
	}
	if c.ContentTypeID == uuid.Nil {
		return fmt.Errorf("コンテンツタイプIDは必須です")
	}
	return nil
}

// Validate はContentTypeの基本的なバリデーション
func (ct *ContentType) Validate() error {
	if ct.Name == "" {
		return fmt.Errorf("名前は必須です")
	}
	if ct.DisplayName == "" {
		return fmt.Errorf("表示名は必須です")
	}
	if ct.CreatedBy == "" {
		return fmt.Errorf("作成者は必須です")
	}
	return nil
}