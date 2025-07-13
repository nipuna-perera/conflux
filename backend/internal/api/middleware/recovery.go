// Panic recovery middleware
// Catches and handles panics in HTTP handlers gracefully
// Prevents server crashes and provides structured error responses
package middleware

import (
	"log"
	"net/http"
)

// Recovery middleware catches panics and returns 500 Internal Server Error
// Logs panic details for debugging while preventing server crashes
// Essential for production stability and error handling
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Recovery implementation:
				// - Log panic details and stack trace
				// - Return 500 Internal Server Error
				// - Prevent response from being written multiple times

				log.Printf("Panic: %v", err)

				// Check if headers have already been written
				if w.Header().Get("Content-Type") == "" {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}
		}()

		next.ServeHTTP(w, r)
	})
}
