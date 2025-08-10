package repository

import (
	"cms_api/internal/domain/entity"
)

// ToContentEntity はContentModelをドメインエンティティに変換
func (c *ContentModel) ToContentEntity() *entity.Content {
	content := &entity.Content{
		ID:            c.ID,
		ContentTypeID: c.ContentTypeID,
		Title:         c.Title,
		Slug:          c.Slug,
		Status:        entity.ContentStatus(c.Status),
		CreatedAt:     c.CreatedAt,
		UpdatedAt:     c.UpdatedAt,
		PublishedAt:   c.PublishedAt,
		AuthorID:      c.AuthorID,
		Version:       c.Version,
	}

	// コンテンツタイプの変換
	if c.ContentType != nil {
		content.ContentType = c.ContentType.ToContentTypeEntity()
	}

	// ブロックの変換
	if len(c.Blocks) > 0 {
		content.Blocks = make([]entity.ContentBlock, len(c.Blocks))
		for i, block := range c.Blocks {
			content.Blocks[i] = *block.ToContentBlockEntity()
		}
	}

	return content
}

// FromContentEntity はドメインエンティティからContentModelを作成
func (c *ContentModel) FromContentEntity(content *entity.Content) {
	c.ID = content.ID
	c.ContentTypeID = content.ContentTypeID
	c.Title = content.Title
	c.Slug = content.Slug
	c.Status = string(content.Status)
	c.CreatedAt = content.CreatedAt
	c.UpdatedAt = content.UpdatedAt
	c.PublishedAt = content.PublishedAt
	c.AuthorID = content.AuthorID
	c.Version = content.Version
}

// ToContentTypeEntity はContentTypeModelをドメインエンティティに変換
func (ct *ContentTypeModel) ToContentTypeEntity() *entity.ContentType {
	return &entity.ContentType{
		ID:          ct.ID,
		Name:        ct.Name,
		DisplayName: ct.DisplayName,
		Description: ct.Description,
		Icon:        ct.Icon,
		IsActive:    ct.IsActive,
		CreatedAt:   ct.CreatedAt,
		UpdatedAt:   ct.UpdatedAt,
		CreatedBy:   ct.CreatedBy,
	}
}

// FromContentTypeEntity はドメインエンティティからContentTypeModelを作成
func (ct *ContentTypeModel) FromContentTypeEntity(contentType *entity.ContentType) {
	ct.ID = contentType.ID
	ct.Name = contentType.Name
	ct.DisplayName = contentType.DisplayName
	ct.Description = contentType.Description
	ct.Icon = contentType.Icon
	ct.IsActive = contentType.IsActive
	ct.CreatedAt = contentType.CreatedAt
	ct.UpdatedAt = contentType.UpdatedAt
	ct.CreatedBy = contentType.CreatedBy
}

// ToContentBlockEntity はContentBlockModelをドメインエンティティに変換
func (cb *ContentBlockModel) ToContentBlockEntity() *entity.ContentBlock {
	block := &entity.ContentBlock{
		ID:         cb.ID,
		ContentID:  cb.ContentID,
		BlockType:  entity.BlockType(cb.BlockType),
		BlockOrder: cb.BlockOrder,
		IsVisible:  cb.IsVisible,
		CreatedAt:  cb.CreatedAt,
		UpdatedAt:  cb.UpdatedAt,
	}

	// ブロックデータの変換
	if cb.Data != nil {
		block.Data = cb.Data.ToContentBlockDataEntity()
	}

	return block
}

// FromContentBlockEntity はドメインエンティティからContentBlockModelを作成
func (cb *ContentBlockModel) FromContentBlockEntity(block *entity.ContentBlock) {
	cb.ID = block.ID
	cb.ContentID = block.ContentID
	cb.BlockType = string(block.BlockType)
	cb.BlockOrder = block.BlockOrder
	cb.IsVisible = block.IsVisible
	cb.CreatedAt = block.CreatedAt
	cb.UpdatedAt = block.UpdatedAt
}

// ToContentBlockDataEntity はContentBlockDataModelをドメインエンティティに変換
func (cbd *ContentBlockDataModel) ToContentBlockDataEntity() *entity.ContentBlockData {
	return &entity.ContentBlockData{
		ID:                  cbd.ID,
		BlockID:             cbd.BlockID,
		DataType:            entity.DataType(cbd.DataType),
		ContentText:         cbd.ContentText,
		ContentRichtext:     cbd.ContentRichtext,
		ContentNumber:       cbd.ContentNumber,
		ContentURL:          cbd.ContentURL,
		ContentJSON:         cbd.ContentJSON,
		ReferencedContentID: cbd.ReferencedContentID,
		Settings:            cbd.Settings,
		CreatedAt:           cbd.CreatedAt,
		UpdatedAt:           cbd.UpdatedAt,
	}
}

// FromContentBlockDataEntity はドメインエンティティからContentBlockDataModelを作成
func (cbd *ContentBlockDataModel) FromContentBlockDataEntity(data *entity.ContentBlockData) {
	cbd.ID = data.ID
	cbd.BlockID = data.BlockID
	cbd.DataType = string(data.DataType)
	cbd.ContentText = data.ContentText
	cbd.ContentRichtext = data.ContentRichtext
	cbd.ContentNumber = data.ContentNumber
	cbd.ContentURL = data.ContentURL
	cbd.ContentJSON = data.ContentJSON
	cbd.ReferencedContentID = data.ReferencedContentID
	cbd.Settings = data.Settings
	cbd.CreatedAt = data.CreatedAt
	cbd.UpdatedAt = data.UpdatedAt
}