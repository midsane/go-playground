package server

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

var jwtSecret = []byte("super-secret-key")

type ErrorResponse struct {
	Error string `json:"error"`
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

		// token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		// 	return jwtSecret, nil
		// })
		claims := jwt.StandardClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (any, error) {
			return jwtSecret, nil
		})
		/*

		 */
		if err != nil || !token.Valid {
			writeJSON(w, http.StatusUnauthorized, ErrorResponse{
				Error: "invalid token",
			})
			return
		}
		userID := claims.Id

		// claims, ok := token.Claims.(jwt.MapClaims)
		// if !ok {
		// 	writeJSON(w, http.StatusUnauthorized, ErrorResponse{
		// 		Error: "invalid claims",
		// 	})
		// }

		/*if we want to avoid the jwt.MapClaims assertion here because of which
		we are forces to check !ok as well.

		we simply use ParseWithClaims and pass pointer to already declared Claims which will can have
		custom static type as we want so no assertion needed.
		*/

		ctx := context.WithValue(r.Context(), "user", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

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
	userID := r.Context().Value("user").(string)

	writeJSON(w, http.StatusOK, map[string]string{
		"user_id": userID,
		"message": "protected profile data",
	})
}
