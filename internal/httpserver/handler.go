package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"note-llm/internal/db"
	"note-llm/internal/models"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateNoteHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var req models.CreateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Title == "" || req.Content == "" {
		http.Error(w, "Missing title or content", http.StatusBadRequest)
		return
	}

	note := models.Note{
		ID:         uuid.New(),
		Title:      req.Title,
		Content:    req.Content,
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}

	collection := db.GetMongoDatabase().Collection("notes")
	result, err := collection.InsertOne(ctx, note)
	if err != nil {
		http.Error(w, "Failed to save note", http.StatusInternalServerError)
		fmt.Printf("Insert error: %v\n", err)
		return
	}

	resp := map[string]interface{}{
		"id": result.InsertedID,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func GetNoteHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		http.Error(w, "Missing note ID", http.StatusBadRequest)
		return
	}

	// Convert to UUID if your model uses UUIDs (as in POST handler)
	noteID, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(w, "Invalid UUID format", http.StatusBadRequest)
		return
	}

	var note models.Note
	collection := db.GetMongoDatabase().Collection("notes")
	err = collection.FindOne(ctx, bson.M{"_id": noteID}).Decode(&note)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			http.Error(w, "Note not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(note)
}
