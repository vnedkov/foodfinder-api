package middleware

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type Middleware func(http.Handler) http.Handler

// Apply applies the given middlewares to the given handler.
func Apply(h http.Handler, mws ...Middleware) http.Handler {
	for _, mw := range mws {
		h = mw(h)
	}
	return h
}

// LoggingMiddleware is a middleware that logs the request URI.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ts := time.Now().UnixMilli()
		next.ServeHTTP(w, r)
		log.Info().
			Str("method", r.Method).
			Str("uri", r.RequestURI).
			Int64("duration", time.Now().UnixMilli()-ts).Send()
	})
}

// CorsMiddleware is a middleware that adds CORS headers to the response.
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
