package model

import (
	"admin/model"
	"time"
)

func RandomContent() *model.Content {
	return &model.Content{
		ID:        uuid.New().String(),
		Title:     "test title",
		Body:      "test body",
		Author:    "test author",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
