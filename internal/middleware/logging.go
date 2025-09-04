package middleware

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

// LoggingMiddleware logs the details of each HTTP request, including the status code.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Wrap the original ResponseWriter to capture the status code.
		writer := &statusWriter{ResponseWriter: w}
		start := time.Now()

		// Pass the wrapped writer and the request to the next handler in the chain.
		next.ServeHTTP(writer, r)

		// Retrieve the status code from the wrapper.
		status := writer.status
		if status == 0 {
			// If WriteHeader was not called, assume a default of 200 OK.
			status = http.StatusOK
		}

		// Log the request details, including the status code.
		log.Info().
			Int("status", status).
			Str("method", r.Method).
			Str("url", r.RequestURI).
			Str("duration", time.Since(start).String()).
			Msg("Request handled")
	})
}

// statusWriter is a wrapper for http.ResponseWriter that captures the status code.
type statusWriter struct {
	http.ResponseWriter
	status int
}

// WriteHeader captures the status code before calling the underlying WriteHeader.
func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}
