package repository

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// ContentModel はGorm用のコンテンツモデル
type ContentModel struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ContentTypeID uuid.UUID `gorm:"type:uuid;not null"`
	Title         string    `gorm:"size:255;not null"`
	Slug          string    `gorm:"size:255;not null;unique"`
	Status        string    `gorm:"type:varchar(20);not null;default:'draft'"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
	PublishedAt   *time.Time
	AuthorID      string `gorm:"size:255;not null"`
	Version       int    `gorm:"default:1"`
	
	// リレーション
	ContentType *ContentTypeModel   `gorm:"foreignKey:ContentTypeID"`
	Blocks      []ContentBlockModel `gorm:"foreignKey:ContentID"`
}

// TableName はテーブル名を指定
func (ContentModel) TableName() string {
	return "contents"
}

// BeforeCreate はレコード作成前のフック
func (c *ContentModel) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// ContentTypeModel はGorm用のコンテンツタイプモデル
type ContentTypeModel struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string    `gorm:"size:100;not null;unique"`
	DisplayName string    `gorm:"size:255;not null"`
	Description string    `gorm:"type:text"`
	Icon        string    `gorm:"size:255"`
	IsActive    bool      `gorm:"default:true"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	CreatedBy   string    `gorm:"size:255;not null"`
}

// TableName はテーブル名を指定
func (ContentTypeModel) TableName() string {
	return "content_types"
}

// BeforeCreate はレコード作成前のフック
func (ct *ContentTypeModel) BeforeCreate(tx *gorm.DB) error {
	if ct.ID == uuid.Nil {
		ct.ID = uuid.New()
	}
	return nil
}

// ContentBlockModel はGorm用のコンテンツブロックモデル
type ContentBlockModel struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ContentID  uuid.UUID `gorm:"type:uuid;not null"`
	BlockType  string    `gorm:"type:varchar(50);not null"`
	BlockOrder int       `gorm:"not null"`
	IsVisible  bool      `gorm:"default:true"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
	
	// リレーション
	Content *ContentModel          `gorm:"foreignKey:ContentID"`
	Data    *ContentBlockDataModel `gorm:"foreignKey:BlockID"`
}

// TableName はテーブル名を指定
func (ContentBlockModel) TableName() string {
	return "content_blocks"
}

// BeforeCreate はレコード作成前のフック
func (cb *ContentBlockModel) BeforeCreate(tx *gorm.DB) error {
	if cb.ID == uuid.Nil {
		cb.ID = uuid.New()
	}
	return nil
}

// ContentBlockDataModel はGorm用のブロックデータモデル
type ContentBlockDataModel struct {
	ID                  uuid.UUID        `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	BlockID             uuid.UUID        `gorm:"type:uuid;not null;unique"`
	DataType            string           `gorm:"type:varchar(50);not null"`
	ContentText         string           `gorm:"type:text"`
	ContentRichtext     json.RawMessage  `gorm:"type:jsonb"`
	ContentNumber       *decimal.Decimal `gorm:"type:decimal(20,6)"`
	ContentURL          string           `gorm:"size:2048"`
	ContentJSON         json.RawMessage  `gorm:"type:jsonb"`
	ReferencedContentID *uuid.UUID       `gorm:"type:uuid"`
	Settings            json.RawMessage  `gorm:"type:jsonb"`
	CreatedAt           time.Time        `gorm:"autoCreateTime"`
	UpdatedAt           time.Time        `gorm:"autoUpdateTime"`
	
	// リレーション
	Block             *ContentBlockModel `gorm:"foreignKey:BlockID"`
	ReferencedContent *ContentModel      `gorm:"foreignKey:ReferencedContentID"`
}

// TableName はテーブル名を指定
func (ContentBlockDataModel) TableName() string {
	return "content_block_data"
}

// BeforeCreate はレコード作成前のフック
func (cbd *ContentBlockDataModel) BeforeCreate(tx *gorm.DB) error {
	if cbd.ID == uuid.Nil {
		cbd.ID = uuid.New()
	}
	return nil
}