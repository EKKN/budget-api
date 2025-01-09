package main

import (
	"context"
	"net/http"
	"strings"
)

// Authenticate middleware to check JWT in the header
func (s *APIServer) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the token from the Authorization header
		_, requestLog, _ := s.prepareRequest(r)
		jobID := r.Header.Get("jobID")
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {

			AppLog(LogRequestResponse(requestLog, map[string]interface{}{"status": "error", "message": "Authorization required"}))
			WriteJSON(w, http.StatusBadRequest, APIError{
				Status:  "error",
				JobID:   jobID,
				Message: "Authorization required",
			})
			return
		}

		// Remove the "Bearer " prefix if it exists
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Validate the token
		token, err := validateJWT(tokenString)
		if err != nil {

			AppLog(LogRequestResponse(requestLog, map[string]interface{}{"status": "error", "message": "invalid token" + err.Error()}))
			WriteJSON(w, http.StatusBadRequest, APIError{
				Status:  "error",
				JobID:   jobID,
				Message: "Invalid token",
			})
			return
		}

		// Check if token is valid
		if !token.Valid {
			AppLog(LogRequestResponse(requestLog, map[string]interface{}{"status": "error", "message": "invalid token"}))
			WriteJSON(w, http.StatusBadRequest, APIError{
				Status:  "error",
				JobID:   jobID,
				Message: "Invalid token",
			})
			return
		}

		// Set the token claims in the request context
		claims, ok := token.Claims.(*UserClaims)
		if !ok || !token.Valid {
			AppLog(LogRequestResponse(requestLog, map[string]interface{}{"status": "error", "message": "Invalid token claims"}))
			WriteJSON(w, http.StatusBadRequest, APIError{
				Status:  "error",
				JobID:   jobID,
				Message: "Invalid token claims",
			})

			return
		}
		ctx := context.WithValue(r.Context(), userContextKey, claims)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
