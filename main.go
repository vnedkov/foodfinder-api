package main

import (
	"context"
	"foodfinder-api/elasticsearch"
	"foodfinder-api/handlers"
	"net/http"

	"github.com/rs/zerolog/log"
)

func main() {
	log.Info().Msg("Starting server")
	es := elasticsearch.NewTypedClient()
	res, err := es.Info().Do(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("Error getting Elasticsearch info")
	}
	log.Info().Interface("result", res).Msgf("Elasticsearch is accessible")

	http.HandleFunc("/search", handlers.SearchHandler)
	http.HandleFunc("/health", handlers.HealthHandler)

	log.Info().Msg("Server is running on port 8080")
	log.Fatal().Err(http.ListenAndServe(":8080", nil))
}
