package qdrant

import (
	"context"
	"time"

	"github.com/qdrant/go-client/qdrant"
)

func InsertNoteEmbedding(noteID, userID string, vector []float32) error {
	ctx := context.Background()

	_, err := GetQdrantClient().Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: "notes",
		Points: []*qdrant.PointStruct{
			{
				Id: &qdrant.PointId{
					PointIdOptions: &qdrant.PointId_Uuid{Uuid: noteID},
				},
				Vectors: &qdrant.Vectors{
					VectorsOptions: &qdrant.Vectors_Vector{Vector: &qdrant.Vector{Data: vector}},
				},
				Payload: map[string]*qdrant.Value{
					"user_id":    {Kind: &qdrant.Value_StringValue{StringValue: userID}},
					"note_id":    {Kind: &qdrant.Value_StringValue{StringValue: noteID}},
					"created_at": {Kind: &qdrant.Value_StringValue{StringValue: time.Now().Format(time.RFC3339)}},
				},
			},
		},
	})
	return err
}
