package models

import (
	"time"
)

type User struct {
	ID        string    `bson:"_id" json:"id"`
	Name      string    `bson:"name" json:"name"`
	Email     string    `bson:"email" json:"email"`
	Provider  string    `bson:"provider" json:"provider"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}
