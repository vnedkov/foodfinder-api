package handlers

import (
	"encoding/json"
	"foodfinder-api/elasticsearch"
	"net/http"
	"strconv"
	"time"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

const (
	defaultPageSize = 10
)

// SearchHandler is an HTTP handler that performs a search query against the Elasticsearch index
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the search query from the request
	q := r.URL.Query().Get("q")
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	pageStr := r.URL.Query().Get("page")
	sizeStr := r.URL.Query().Get("size")
	// hallStr := r.URL.Query().Get("hall")

	// Validate the query parameters
	// TODO: Implement validation for the query string
	page := 0
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err != nil {
			_ = handleError(w, r, err, http.StatusBadRequest, "Invalid page number")
		} else {
			page = p
		}
	}

	pageSize := defaultPageSize
	if sizeStr != "" {
		if s, err := strconv.Atoi(sizeStr); err != nil {
			_ = handleError(w, r, err, http.StatusBadRequest, "Invalid page size")
		} else {
			pageSize = s
		}
	}

	if from == "" {
		from = time.Now().Format("2006-01-02")
	} else if _, err := time.Parse("2006-01-02", from); err != nil {
		_ = handleError(w, r, err, http.StatusBadRequest, "Invalid from date format")
	}
	if to == "" {
		to = time.Now().AddDate(0, 0, 15).Format("2006-01-02")
	} else if _, err := time.Parse("2006-01-02", to); err != nil {
		_ = handleError(w, r, err, http.StatusBadRequest, "Invalid to date format")
	}

	// Create a new Elasticsearch client
	es, err := elasticsearch.FromContext(r.Context())
	if handleError(w, r, err, http.StatusInternalServerError, "Error getting Elasticsearch client") {
		return
	}

	query := getQueryFilters(&from, &to, &q)

	// Perform the search query
	res, err := es.Search().
		Index(elasticsearch.Index).
		Query(query).
		From(page * pageSize).
		Size(pageSize).
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

func getQueryFilters(from, to, q *string) *types.Query {
	// Create a new query filter
	query := types.Query{
		Bool: &types.BoolQuery{
			Filter: []types.Query{
				{
					Range: map[string]types.RangeQuery{
						"date": types.DateRangeQuery{
							Gte: from,
							Lte: to,
						},
					},
				},
			},
			Must: []types.Query{
				{
					Match: map[string]types.MatchQuery{
						"keywords": {
							Query: *q,
						},
					},
				},
			},
		},
	}
	return &query
}
