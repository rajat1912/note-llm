package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"note-llm/internal/db"
	"note-llm/internal/models"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
)

func generateJWT(email string) (string, error) {
	jwtSecret := viper.GetString("JWT_SECRET")
	claims := jwt.MapClaims{
		"sub": email,
		"exp": jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		"iat": jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func SetupAuthProviders() {
	key := viper.GetString("GOOGLE_KEY")
	secret := viper.GetString("GOOGLE_SECRET")

	goth.UseProviders(
		google.New(key, secret, "http://localhost:8080/auth/google/callback", "openid", "email", "profile"),
	)
}

func createUserIfNotExists(ctx context.Context, user goth.User) int {

	collection := db.GetMongoDatabase().Collection("users")
	var existing models.User
	err := collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&existing)

	if err != nil {

		fmt.Println(err.Error())
		// Not found → create new user
		newUser := models.User{
			ID:        uuid.New().String(),
			Name:      user.Name,
			Email:     user.Email,
			Provider:  user.Provider,
			CreatedAt: time.Now(),
		}

		_, err := collection.InsertOne(ctx, newUser)
		if err != nil {
			fmt.Println(err.Error())
			return http.StatusInternalServerError
		}
	}

	return http.StatusOK
}

func Provider(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 50*time.Second)
	defer cancel()

	provider := chi.URLParam(r, "provider")
	q := r.URL.Query()
	q.Add("provider", provider)
	r.URL.RawQuery = q.Encode()
	if user, err := gothic.CompleteUserAuth(w, r); err == nil {
		status := createUserIfNotExists(ctx, user)
		if status != http.StatusOK {
			http.Error(w, "Failed to save user", status)
			return
		}
	} else {
		gothic.BeginAuthHandler(w, r)
	}
}

func Callback(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	provider := chi.URLParam(r, "provider")
	q := r.URL.Query()
	q.Add("provider", provider)
	r.URL.RawQuery = q.Encode()

	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Typically you'd check if user exists in DB here
	status := createUserIfNotExists(ctx, user)
	if status != http.StatusOK {
		http.Error(w, "Failed to save user", status)
		return
	}

	tokenString, err := generateJWT(user.Email)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	redirectURL := fmt.Sprintf("http://localhost:5173/auth/callback#token=%s&email=%s", tokenString, user.Email)
	http.Redirect(w, r, redirectURL, http.StatusFound)
}
