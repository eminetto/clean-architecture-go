package entity

import (
	"time"
)

//Bookmark data
type Bookmark struct {
	ID          ID        `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Link        string    `json:"link"`
	Tags        []string  `json:"tags"`
	Favorite    bool      `json:"favorite"`
	CreatedAt   time.Time `json:"created_at"`
}
