package rag

import (
	"context"
	"fmt"

	"note-llm/internal/db"
	"note-llm/internal/llm"
	"note-llm/internal/search"
)

// AnswerFromUserNotes performs the full RAG flow:
// embed → vector search → fetch from MongoDB → pass to LLM
func AnswerFromUserNotes(ctx context.Context, userID, question string) (string, error) {
	// Step 1: Get relevant note IDs from Qdrant
	noteIDs, err := search.SearchRelevantNotes(userID, question)
	if err != nil {
		return "", fmt.Errorf("search failed: %w", err)
	}
	if len(noteIDs) == 0 {
		return "No relevant notes found.", nil
	}

	// Step 2: Fetch full note content from MongoDB
	notes, err := db.FetchNotesByIDs(ctx, noteIDs, userID)
	if err != nil {
		return "", fmt.Errorf("failed to fetch notes from DB: %w", err)
	}

	// Step 3: Extract note content
	var noteTexts []string
	for _, note := range notes {
		noteTexts = append(noteTexts, fmt.Sprintf("%s\n\n%s", note.Title, note.Content))
	}

	// Step 4: Ask the LLM
	answer, err := llm.Summarize(question, noteTexts)
	if err != nil {
		return "", fmt.Errorf("LLM call failed: %w", err)
	}

	return answer, nil
}
