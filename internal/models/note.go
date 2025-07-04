package models

import (
	"time"
)

type Note struct {
	ID         string    `bson:"_id" json:"id"`
	Title      string    `bson:"title" json:"title"`
	Content    string    `bson:"content" json:"content"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
	ModifiedAt time.Time `bson:"modified_at" json:"modified_at"`
}

type CreateNoteRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
