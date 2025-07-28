package search

import (
	"context"
	"fmt"
	"note-llm/internal/llm"
	"note-llm/internal/qdrant"

	qdrantpb "github.com/qdrant/go-client/qdrant"
)

func SearchRelevantNotes(userID string, query string) ([]string, error) {
	// Step 1: Embed the query
	embeddings, err := llm.GetEmbeddings([]string{query})
	if err != nil {
		return nil, fmt.Errorf("failed to embed query: %w", err)
	}

	// Step 2: Search Qdrant for similar notes
	ctx := context.Background()
	resp, err := qdrant.GetQdrantClient().Query(ctx, &qdrantpb.QueryPoints{
		CollectionName: "notes",
		Query:          qdrantpb.NewQuery(embeddings[0]...),
		Filter: &qdrantpb.Filter{
			Must: []*qdrantpb.Condition{
				qdrantpb.NewMatch("user_id", userID),
			},
		},
		WithPayload: qdrantpb.NewWithPayload(true),
	})
	if err != nil {
		return nil, fmt.Errorf("qdrant search error: %w", err)
	}

	// Step 3: Extract note IDs or payloads
	var notes []string
	for _, hit := range resp {
		payload := hit.GetPayload()
		if payload["note_id"] != nil {
			notes = append(notes, payload["note_id"].GetStringValue())
		}
	}

	return notes, nil
}
