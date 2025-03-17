package model

import "time"

// Content はコンテンツのドメインモデルです
type Content struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewContent は新しいContentインスタンスを作成します
func NewContent(id, title, body, author string) *Content {
	now := time.Now()
	return &Content{
		ID:        id,
		Title:     title,
		Body:      body,
		Author:    author,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
