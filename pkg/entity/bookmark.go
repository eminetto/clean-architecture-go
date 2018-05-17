package entity

import (
	"time"
)

//Bookmark data
type Bookmark struct {
	ID          ID        `json:"id" bson:"_id"`
	Name        string    `json:"name" bson:"name"`
	Description string    `json:"description" bson:"description"`
	Link        string    `json:"link" bson:"link"`
	Tags        []string  `json:"tags" bson:"tags"`
	Favorite    bool      `json:"favorite" bson:"favorite"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
}
