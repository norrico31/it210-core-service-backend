package utils

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/norrico31/it210-core-service-backend/entities"
)

func GenerateJWT(u entities.User) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", fmt.Errorf("jwt secret is not set")
	}

	claims := jwt.MapClaims{
		"user_id":    u.ID,
		"first_name": u.FirstName,
		"last_name":  u.LastName,
		"email":      u.Email,
		"exp":        time.Now().Add(time.Hour * 720).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

const BEARER = "Bearer "

func ValidateJWT(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for key, values := range r.Header {
			fmt.Printf("Headers: %s, Value: %v\n", key, values)
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, BEARER) {
			http.Error(w, "Authorization header missing or invalid", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, BEARER)
		jwtSecret, exists := os.LookupEnv("JWT_SECRET")
		if !exists {
			http.Error(w, " env variable JWT_SECRET not set", http.StatusInternalServerError)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the algorithm
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Extract claims and validate expiration
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if exp, ok := claims["exp"].(float64); ok {
			if int64(exp) < time.Now().Unix() {
				http.Error(w, "Token has expired", http.StatusUnauthorized)
				return
			}
		}

		userID, ok := claims["user_id"].(float64)
		if ok {
			userIDStr := fmt.Sprintf("%.0f", userID) // Convert to string
			r.Header.Set("X-User-ID", userIDStr)
		}

		next(w, r)
	}
}
