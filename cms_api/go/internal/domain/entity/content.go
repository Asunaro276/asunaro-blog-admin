package model

import "time"

type Article struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Body        string    `json:"body"`
	CoverImage  string    `json:"cover_image"`
	PublishedAt time.Time `json:"published_at"`
	Status      string    `json:"status"`
	CategoryID  string    `json:"category_id"`
	Tags        []string  `json:"tags"`
}
