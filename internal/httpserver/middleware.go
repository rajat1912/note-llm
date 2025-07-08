package httpserver

import (
	"context"
	"net/http"
	"strings"
	"time"

	"note-llm/internal/db"
	"note-llm/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
)

type contextKey string

const UserEmailKey contextKey = "userEmail"
const UserIDKey contextKey = "userId"

func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		tokenStr := strings.TrimPrefix(auth, "Bearer ")

		jwtSecret := viper.GetString("JWT_SECRET")
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		emailClaim, ok := claims["sub"].(string)
		if !ok || emailClaim == "" {
			http.Error(w, "Invalid email claim", http.StatusUnauthorized)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var user models.User
		err = db.GetMongoDatabase().Collection("users").FindOne(ctx, bson.M{"email": emailClaim}).Decode(&user)
		if err != nil {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}

		ctxWithUser := context.WithValue(r.Context(), UserEmailKey, user.Email)
		ctxWithUser = context.WithValue(ctxWithUser, UserIDKey, user.ID)

		next.ServeHTTP(w, r.WithContext(ctxWithUser))
	})
}
