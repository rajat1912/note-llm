package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"note-llm/internal/db"
	"note-llm/internal/llm"
	"note-llm/internal/models"
	"note-llm/internal/qdrant"
	"note-llm/internal/rag"

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

	userId, ok := r.Context().Value(UserIDKey).(string)
	if !ok || userId == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	stringToEmbed := fmt.Sprintf("%s\n\n%s", req.Title, req.Content)
	embeddings, err := llm.GetEmbeddings([]string{stringToEmbed})
	if err != nil {
		http.Error(w, "Failed to generate embedding", http.StatusInternalServerError)
		fmt.Printf("Embedding error: %v\n", err)
		return
	}

	note := models.Note{
		ID:         uuid.New().String(),
		Title:      req.Title,
		Content:    req.Content,
		Embeddings: embeddings[0],
		UserID:     userId,
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}

	collection := db.GetMongoDatabase().Collection("notes")
	_, err = collection.InsertOne(ctx, note)
	if err != nil {
		http.Error(w, "Failed to save note", http.StatusInternalServerError)
		fmt.Printf("Insert error: %v\n", err)
		return
	}
	err = qdrant.InsertNoteEmbedding(note.ID, userId, embeddings[0])
	if err != nil {
		fmt.Printf("Qdrant insert error: %v\n", err)
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
	w.WriteHeader(http.StatusCreated)
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

func AskQuestionHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Question string `json:"question"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Question == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value(UserIDKey).(string)
	ctx := r.Context()

	answer, err := rag.AnswerFromUserNotes(ctx, userID, req.Question)
	if err != nil {
		fmt.Print(err.Error())
		http.Error(w, "Failed to generate answer", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"answer": answer})
}
