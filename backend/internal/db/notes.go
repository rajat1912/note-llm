package db

import (
	"context"
	"note-llm/internal/models"

	"go.mongodb.org/mongo-driver/bson"
)

func FetchNotesByIDs(ctx context.Context, noteIDs []string, userID string) ([]models.Note, error) {
	collection := GetMongoDatabase().Collection("notes")

	filter := bson.M{
		"user_id": userID,
		"_id": bson.M{
			"$in": noteIDs,
		},
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var notes []models.Note
	for cursor.Next(ctx) {
		var note models.Note
		if err := cursor.Decode(&note); err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}

	return notes, nil
}
