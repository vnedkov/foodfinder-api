package handlers

import (
	"encoding/json"
	"foodfinder-api/elasticsearch"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

// SearchHandler is an HTTP handler that performs a search query against the Elasticsearch index
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the search query from the request
	query := r.URL.Query().Get("q")

	// Create a new Elasticsearch client
	es := elasticsearch.FromContext(r.Context())

	// Perform the search query
	res, err := es.Search().
		Index(elasticsearch.Index).
		Query(&types.Query{
			Match: map[string]types.MatchQuery{
				"keywords": {
					Query: query,
				},
			},
		}).
		From(0).
		Size(10).
		Pretty(true).
		Do(r.Context())

	if err != nil {
		http.Error(w, "Error performing search query", http.StatusInternalServerError)
		return
	}

	// Decode the search response
	var searchResponse []interface{}
	for _, hit := range res.Hits.Hits {
		var food map[string]interface{}
		if err := json.Unmarshal(hit.Source_, &food); err != nil {
			http.Error(w, "Error decoding search response", http.StatusInternalServerError)
			return
		}
		searchResponse = append(searchResponse, food)
	}

	// Return the search results as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(searchResponse)
}
