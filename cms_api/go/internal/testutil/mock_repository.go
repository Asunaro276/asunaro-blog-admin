package testutil

import (
	"cms/internal/repository"
)

// MockContentRepository はContentRepositoryのモック実装です。
type MockContentRepository struct {
	GetContentFunc func(id string) (string, error)
}

// GetContent はモックのGetContent実装です。
func (m *MockContentRepository) GetContent(id string) (string, error) {
	if m.GetContentFunc != nil {
		return m.GetContentFunc(id)
	}
	return "Mocked content", nil
}

// NewMockContentRepository は新しいMockContentRepositoryを作成します。
func NewMockContentRepository() repository.ContentRepository {
	return &MockContentRepository{}
}
