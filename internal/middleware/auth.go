package middleware

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

// AuthMiddleware checks for a valid API key in the request header.
// It's a simple example of authentication middleware.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Hardcoded API key for demonstration purposes
		const apiSecretKey = "my-secret-key"

		// Get the API key from the "X-API-Key" header
		apiKey := r.Header.Get("X-API-Key")

		// Check if the API key matches
		if apiKey != apiSecretKey {
			log.Warn().
				Str("url", r.RequestURI).
				Msg("Unauthorized access attempt - Invalid API Key")

			// Respond with a 401 Unauthorized status
			http.Error(w, "Unauthorized - Invalid API Key", http.StatusUnauthorized)
			return
		}

		// Pass the request to the next handler in the chain
		next.ServeHTTP(w, r)
	})
}
