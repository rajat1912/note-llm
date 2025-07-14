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

	userId := r.Context().Value(UserIDKey).(string)

	note := models.Note{
		ID:         uuid.New().String(),
		Title:      req.Title,
		Content:    req.Content,
		UserID:     userId,
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}

	collection := db.GetMongoDatabase().Collection("notes")
	_, err := collection.InsertOne(ctx, note)
	if err != nil {
		http.Error(w, "Failed to save note", http.StatusInternalServerError)
		fmt.Printf("Insert error: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(note)
}

func GetNoteHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	noteID := chi.URLParam(r, "id")
	if noteID == "" {
		http.Error(w, "Missing note ID", http.StatusBadRequest)
		return
	}

	var note models.Note
	collection := db.GetMongoDatabase().Collection("notes")

	userId := r.Context().Value(UserIDKey).(string)

	err := collection.FindOne(ctx, bson.M{"_id": noteID, "user_id": userId}).Decode(&note)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			http.Error(w, "Note not found", http.StatusNotFound)
		} else {
			fmt.Println(err.Error())
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(note)
}

func GetAllNotesHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	userId := r.Context().Value(UserIDKey).(string)

	collection := db.GetMongoDatabase().Collection("notes")
	cursor, err := collection.Find(ctx, bson.M{"user_id": userId})
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		fmt.Printf("Find error: %v\n", err)
		return
	}
	defer cursor.Close(ctx)

	var notes []models.Note
	if err := cursor.All(ctx, &notes); err != nil {
		http.Error(w, "Failed to parse notes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

func UpdateNoteHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	noteID := chi.URLParam(r, "id")
	if noteID == "" {
		http.Error(w, "Missing note ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	userId := r.Context().Value(UserIDKey).(string)

	collection := db.GetMongoDatabase().Collection("notes")
	filter := bson.M{"_id": noteID, "user_id": userId}
	update := bson.M{
		"$set": bson.M{
			"title":       req.Title,
			"content":     req.Content,
			"modified_at": time.Now(),
		},
	}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		http.Error(w, "Failed to update note", http.StatusInternalServerError)
		return
	}

	if result.MatchedCount == 0 {
		http.Error(w, "Note not found or unauthorized", http.StatusNotFound)
		return
	}

	var updatedNote models.Note
	err = collection.FindOne(ctx, filter).Decode(&updatedNote)
	if err != nil {
		http.Error(w, "Failed to retrieve updated note", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedNote)
}

func DeleteNoteHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	noteID := chi.URLParam(r, "id")
	if noteID == "" {
		http.Error(w, "Missing note ID", http.StatusBadRequest)
		return
	}

	userId := r.Context().Value(UserIDKey).(string)

	collection := db.GetMongoDatabase().Collection("notes")
	result, err := collection.DeleteOne(ctx, bson.M{"_id": noteID, "user_id": userId})
	if err != nil {
		http.Error(w, "Failed to delete note", http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		http.Error(w, "Note not found or unauthorized", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content
}
