package middleware

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

// LoggingMiddleware logs the details of each HTTP request.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Start a timer to measure request duration.
		start := time.Now()

		// Pass the request to the next handler in the chain.
		next.ServeHTTP(w, r)

		// Log the request details after the handler has executed.
		log.Info().
			Str("method", r.Method).
			Str("url", r.RequestURI).
			Str("duration", time.Since(start).String()).
			Msg("Request handled")
	})
}
