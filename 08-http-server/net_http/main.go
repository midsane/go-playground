package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const userContextKey contextKey = "user"

var jwtSecret = []byte("super-secret-key")

// =========================
// Models
// =========================

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// =========================
// Utility Helpers
// =========================

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func parseJSON(r *http.Request, dst interface{}) error {
	if r.Header.Get("Content-Type") != "application/json" {
		return errors.New("content type must be application/json")
	}
	return json.NewDecoder(r.Body).Decode(dst)
}

// =========================
// Middleware
// =========================

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s %v",
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
			time.Since(start),
		)
	})
}

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("PANIC:", err)
				writeJSON(w, http.StatusInternalServerError, ErrorResponse{
					Error: "internal server error",
				})
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")

		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			writeJSON(w, http.StatusUnauthorized, ErrorResponse{
				Error: "missing or invalid token",
			})
			return
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			writeJSON(w, http.StatusUnauthorized, ErrorResponse{
				Error: "invalid token",
			})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			writeJSON(w, http.StatusUnauthorized, ErrorResponse{
				Error: "invalid claims",
			})
			return
		}

		userID := claims["user_id"].(string)

		ctx := context.WithValue(r.Context(), userContextKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// =========================
// Handlers
// =========================

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{
			Error: "method not allowed",
		})
		return
	}

	type LoginRequest struct {
		UserID string `json:"user_id"`
	}

	var req LoginRequest

	if err := parseJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	if req.UserID == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: "user_id required",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": req.UserID,
		"exp":     time.Now().Add(time.Hour).Unix(),
	})

	tokenStr, err := token.SignedString(jwtSecret)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error: "failed to sign token",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"token": tokenStr,
	})
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userContextKey).(string)

	writeJSON(w, http.StatusOK, map[string]string{
		"user_id": userID,
		"message": "protected profile data",
	})
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{
			Error: "method not allowed",
		})
		return
	}

	var user User

	if err := parseJSON(r, &user); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	if user.ID == "" || user.Email == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: "id and email required",
		})
		return
	}

	writeJSON(w, http.StatusCreated, user)
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: "id query param required",
		})
		return
	}

	writeJSON(w, http.StatusOK, User{
		ID:    id,
		Name:  "Demo User",
		Email: "demo@example.com",
	})
}

// =========================
// Middleware Chain Builder
// =========================

func Chain(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

// =========================
// Main
// =========================

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	mux := http.NewServeMux()

	mux.Handle("/login", http.HandlerFunc(LoginHandler))
	mux.Handle("/users", http.HandlerFunc(CreateUserHandler))
	mux.Handle("/user", http.HandlerFunc(GetUserHandler))

	protectedProfile := Chain(
		http.HandlerFunc(ProfileHandler),
		JWTMiddleware,
	)

	mux.Handle("/profile", protectedProfile)

	finalHandler := Chain(
		mux,
		RecoveryMiddleware,
		LoggingMiddleware,
	)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      finalHandler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	logger.Println("Server running on :8080")

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
