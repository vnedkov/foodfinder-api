package main

import (
	"context"
	"foodfinder-api/elasticsearch"
	"foodfinder-api/handlers"
	"foodfinder-api/middleware"
	"net/http"

	"github.com/rs/zerolog/log"
)

func main() {
	log.Info().Msg("Starting server")
	es := elasticsearch.TypedClient()
	res, err := es.Info().Do(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("Error getting Elasticsearch info")
	}
	log.Info().Interface("result", res).Msgf("Elasticsearch is accessible")
	healthMiddlewareChain := middleware.Apply(http.HandlerFunc(handlers.HealthHandler), middleware.LoggingMiddleware, middleware.CorsMiddleware)
	searchMiddlewareChain := middleware.Apply(http.HandlerFunc(handlers.SearchHandler), middleware.LoggingMiddleware, middleware.CorsMiddleware)

	http.HandleFunc("/search", searchMiddlewareChain.ServeHTTP)
	http.HandleFunc("/health", healthMiddlewareChain.ServeHTTP)

	log.Info().Msg("Server is running on port 8080")
	log.Fatal().Err(http.ListenAndServe(":8080", nil))
}
