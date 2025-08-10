package repository

import (
	"cms_api/internal/domain/entity"
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ContentRepository はコンテンツリポジトリのインターフェース
type ContentRepository interface {
	// コンテンツ操作
	GetContentByID(ctx context.Context, id uuid.UUID) (*entity.Content, error)
	GetContents(ctx context.Context, limit, offset int, filters ContentFilters) ([]*entity.Content, int64, error)
	CreateContent(ctx context.Context, content *entity.Content) error
	UpdateContent(ctx context.Context, content *entity.Content) error
	DeleteContent(ctx context.Context, id uuid.UUID) error
	
	// コンテンツタイプ操作
	GetContentTypes(ctx context.Context) ([]*entity.ContentType, error)
	GetContentTypeByID(ctx context.Context, id uuid.UUID) (*entity.ContentType, error)
	CreateContentType(ctx context.Context, contentType *entity.ContentType) error
}

// ContentFilters はコンテンツ検索時のフィルター条件
type ContentFilters struct {
	Status     *entity.ContentStatus
	Category   string
	Tags       []string
	Search     string
	Sort       string
	Order      string
	AuthorID   string
}

type contentRepository struct {
	db *gorm.DB
}

// NewContentRepository は新しいContentRepositoryインスタンスを作成します
func NewContentRepository(db *gorm.DB) ContentRepository {
	return &contentRepository{
		db: db,
	}
}

// GetContentByID はIDでコンテンツを取得します
func (r *contentRepository) GetContentByID(ctx context.Context, id uuid.UUID) (*entity.Content, error) {
	var contentModel ContentModel
	
	err := r.db.WithContext(ctx).
		Preload("ContentType").
		Preload("Blocks.Data").
		Where("id = ?", id).
		First(&contentModel).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("コンテンツが見つかりません: %s", id.String())
		}
		return nil, fmt.Errorf("コンテンツの取得に失敗しました: %w", err)
	}
	
	return contentModel.ToContentEntity(), nil
}

// GetContents はコンテンツ一覧を取得します
func (r *contentRepository) GetContents(ctx context.Context, limit, offset int, filters ContentFilters) ([]*entity.Content, int64, error) {
	query := r.db.WithContext(ctx).Model(&ContentModel{}).
		Preload("ContentType").
		Preload("Blocks.Data")
	
	// フィルター条件の適用
	if filters.Status != nil {
		query = query.Where("status = ?", string(*filters.Status))
	}
	
	if filters.AuthorID != "" {
		query = query.Where("author_id = ?", filters.AuthorID)
	}
	
	if filters.Search != "" {
		query = query.Where("title ILIKE ? OR slug ILIKE ?", 
			"%"+filters.Search+"%", "%"+filters.Search+"%")
	}
	
	// ソート条件の適用
	orderBy := "created_at DESC" // デフォルト
	if filters.Sort != "" && filters.Order != "" {
		orderBy = fmt.Sprintf("%s %s", filters.Sort, filters.Order)
	}
	query = query.Order(orderBy)
	
	// 総数を取得
	var total int64
	countQuery := query.Session(&gorm.Session{})
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("コンテンツ総数の取得に失敗しました: %w", err)
	}
	
	// ページネーション適用してモデルを取得
	var contentModels []ContentModel
	err := query.Limit(limit).Offset(offset).Find(&contentModels).Error
	if err != nil {
		return nil, 0, fmt.Errorf("コンテンツ一覧の取得に失敗しました: %w", err)
	}
	
	// ドメインエンティティに変換
	contents := make([]*entity.Content, len(contentModels))
	for i, model := range contentModels {
		contents[i] = model.ToContentEntity()
	}
	
	return contents, total, nil
}

// CreateContent は新しいコンテンツを作成します
func (r *contentRepository) CreateContent(ctx context.Context, content *entity.Content) error {
	// バリデーション
	if err := content.Validate(); err != nil {
		return fmt.Errorf("コンテンツのバリデーションエラー: %w", err)
	}
	
	// トランザクション内で作成
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// ドメインエンティティからGormモデルに変換
		var contentModel ContentModel
		contentModel.FromContentEntity(content)
		
		// コンテンツの作成
		if err := tx.Create(&contentModel).Error; err != nil {
			return fmt.Errorf("コンテンツの作成に失敗しました: %w", err)
		}
		
		// IDを更新（DB生成の場合）
		content.ID = contentModel.ID
		
		// ブロックがある場合は作成
		for i := range content.Blocks {
			content.Blocks[i].ContentID = content.ID
			
			var blockModel ContentBlockModel
			blockModel.FromContentBlockEntity(&content.Blocks[i])
			
			if err := tx.Create(&blockModel).Error; err != nil {
				return fmt.Errorf("コンテンツブロックの作成に失敗しました: %w", err)
			}
			
			content.Blocks[i].ID = blockModel.ID
			
			// ブロックデータがある場合は作成
			if content.Blocks[i].Data != nil {
				content.Blocks[i].Data.BlockID = content.Blocks[i].ID
				
				var dataModel ContentBlockDataModel
				dataModel.FromContentBlockDataEntity(content.Blocks[i].Data)
				
				if err := tx.Create(&dataModel).Error; err != nil {
					return fmt.Errorf("ブロックデータの作成に失敗しました: %w", err)
				}
				
				content.Blocks[i].Data.ID = dataModel.ID
			}
		}
		
		return nil
	})
}

// UpdateContent はコンテンツを更新します
func (r *contentRepository) UpdateContent(ctx context.Context, content *entity.Content) error {
	// バリデーション
	if err := content.Validate(); err != nil {
		return fmt.Errorf("コンテンツのバリデーションエラー: %w", err)
	}
	
	// 存在確認
	var existing ContentModel
	if err := r.db.WithContext(ctx).Where("id = ?", content.ID).First(&existing).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("更新対象のコンテンツが見つかりません: %s", content.ID.String())
		}
		return fmt.Errorf("コンテンツの存在確認に失敗しました: %w", err)
	}
	
	// トランザクション内で更新
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// ドメインエンティティからGormモデルに変換
		var contentModel ContentModel
		contentModel.FromContentEntity(content)
		
		// コンテンツの更新
		if err := tx.Save(&contentModel).Error; err != nil {
			return fmt.Errorf("コンテンツの更新に失敗しました: %w", err)
		}
		
		return nil
	})
}

// DeleteContent はコンテンツを削除します
func (r *contentRepository) DeleteContent(ctx context.Context, id uuid.UUID) error {
	// 存在確認
	var contentModel ContentModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&contentModel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("削除対象のコンテンツが見つかりません: %s", id.String())
		}
		return fmt.Errorf("コンテンツの存在確認に失敗しました: %w", err)
	}
	
	// トランザクション内で削除
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// ブロックデータの削除
		if err := tx.Where("block_id IN (SELECT id FROM content_blocks WHERE content_id = ?)", id).Delete(&ContentBlockDataModel{}).Error; err != nil {
			return fmt.Errorf("ブロックデータの削除に失敗しました: %w", err)
		}
		
		// ブロックの削除
		if err := tx.Where("content_id = ?", id).Delete(&ContentBlockModel{}).Error; err != nil {
			return fmt.Errorf("コンテンツブロックの削除に失敗しました: %w", err)
		}
		
		// コンテンツの削除
		if err := tx.Delete(&contentModel).Error; err != nil {
			return fmt.Errorf("コンテンツの削除に失敗しました: %w", err)
		}
		
		return nil
	})
}

// GetContentTypes はコンテンツタイプ一覧を取得します
func (r *contentRepository) GetContentTypes(ctx context.Context) ([]*entity.ContentType, error) {
	var contentTypeModels []ContentTypeModel
	
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Order("display_name ASC").
		Find(&contentTypeModels).Error
	
	if err != nil {
		return nil, fmt.Errorf("コンテンツタイプ一覧の取得に失敗しました: %w", err)
	}
	
	// ドメインエンティティに変換
	contentTypes := make([]*entity.ContentType, len(contentTypeModels))
	for i, model := range contentTypeModels {
		contentTypes[i] = model.ToContentTypeEntity()
	}
	
	return contentTypes, nil
}

// GetContentTypeByID はIDでコンテンツタイプを取得します
func (r *contentRepository) GetContentTypeByID(ctx context.Context, id uuid.UUID) (*entity.ContentType, error) {
	var contentTypeModel ContentTypeModel
	
	err := r.db.WithContext(ctx).
		Where("id = ? AND is_active = ?", id, true).
		First(&contentTypeModel).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("コンテンツタイプが見つかりません: %s", id.String())
		}
		return nil, fmt.Errorf("コンテンツタイプの取得に失敗しました: %w", err)
	}
	
	return contentTypeModel.ToContentTypeEntity(), nil
}

// CreateContentType は新しいコンテンツタイプを作成します
func (r *contentRepository) CreateContentType(ctx context.Context, contentType *entity.ContentType) error {
	// バリデーション
	if err := contentType.Validate(); err != nil {
		return fmt.Errorf("コンテンツタイプのバリデーションエラー: %w", err)
	}
	
	// 名前の重複チェック
	var existing ContentTypeModel
	err := r.db.WithContext(ctx).Where("name = ?", contentType.Name).First(&existing).Error
	if err == nil {
		return fmt.Errorf("同じ名前のコンテンツタイプが既に存在します: %s", contentType.Name)
	} else if err != gorm.ErrRecordNotFound {
		return fmt.Errorf("コンテンツタイプの重複チェックに失敗しました: %w", err)
	}
	
	// ドメインエンティティからGormモデルに変換
	var contentTypeModel ContentTypeModel
	contentTypeModel.FromContentTypeEntity(contentType)
	
	// 作成
	if err := r.db.WithContext(ctx).Create(&contentTypeModel).Error; err != nil {
		return fmt.Errorf("コンテンツタイプの作成に失敗しました: %w", err)
	}
	
	// IDを更新（DB生成の場合）
	contentType.ID = contentTypeModel.ID
	
	return nil
}