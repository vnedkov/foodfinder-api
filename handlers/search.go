package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"foodfinder-api/elasticsearch"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

const (
	defaultPageSize = 10
	minSearchLength = 3
	maxSearchLength = 100
)

// SearchHandler is an HTTP handler that performs a search query against the Elasticsearch index
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	// Parse the URL query parameters
	queryParams := r.URL.Query()
	q := queryParams.Get("q")
	fromStr := queryParams.Get("from")
	toStr := queryParams.Get("to")
	pageStr := queryParams.Get("page")
	sizeStr := queryParams.Get("size")
	// Get all selected halls. If no hall is selected, get all halls
	hallStr := queryParams["hall"]

	// Validate the query parameters
	if err := validateSearchString(q); handleError(w, r, err, http.StatusBadRequest) {
		return
	}
	var fromDate time.Time
	// Validate the date parameters
	if fromDate, err = validateDateString(fromStr, time.Now()); handleError(w, r, err, http.StatusBadRequest) {
		return
	}

	var toDate time.Time
	if toDate, err = validateDateString(toStr, time.Now().AddDate(0, 0, 14)); handleError(w, r, err, http.StatusBadRequest) {
		return
	}

	// Validate the hall parameters
	var halls []int
	if halls, err = validateHalls(hallStr); handleError(w, r, err, http.StatusBadRequest) {
		return
	}

	// Parse the page and size parameters
	page := 0
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); handleError(w, r, err, http.StatusBadRequest, "Invalid page number ", pageStr) {
			return
		} else {
			page = p
		}
	}

	pageSize := defaultPageSize
	if sizeStr != "" {
		if s, err := strconv.Atoi(sizeStr); handleError(w, r, err, http.StatusBadRequest, "Invalid page size ", sizeStr) {
			return
		} else {
			pageSize = s
		}
	}

	// Create a new Elasticsearch client
	es, err := elasticsearch.FromContext(r.Context())
	if handleError(w, r, err, http.StatusInternalServerError, "Error getting Elasticsearch client") {
		return
	}

	// Elasticsearch query needs *string for date comparisons
	from := fromDate.Format("2006-01-02")
	to := toDate.Format("2006-01-02")

	// Perform the search query
	res, err := es.Search().
		Index(elasticsearch.Index).
		Query(getQuery(&from, &to, &q, halls)).
		From(page * pageSize).
		Size(pageSize).
		Pretty(true).
		Do(r.Context())

	if handleError(w, r, err, http.StatusInternalServerError, "Error performing search query") {
		return
	}

	// Decode the search response
	var searchResponse []interface{}
	for _, hit := range res.Hits.Hits {
		var food map[string]interface{}
		if err := json.Unmarshal(hit.Source_, &food); err != nil {
			_ = handleError(w, r, err, http.StatusInternalServerError, "Error decoding search response")
			return
		}
		searchResponse = append(searchResponse, food)
	}

	// Return the search results as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(searchResponse)
}

// getQuery creates a new Elasticsearch query based on the given parameters
func getQuery(from, to, q *string, halls []int) *types.Query {
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
	// Add the halls filter if provided
	if len(halls) > 0 {
		query.Bool.Filter = append(query.Bool.Filter, types.Query{
			Terms: &types.TermsQuery{
				TermsQuery: map[string]types.TermsQueryField{
					"hall_id": halls,
				},
			},
		})
	}
	return &query
}

// validateDateString validates the given date string based on specific criteria
func validateDateString(date string, defaultIfEmpty time.Time) (time.Time, error) {
	// If the date string is empty, return the default date
	if date == "" {
		return defaultIfEmpty, nil
	}
	if t, err := time.Parse("2006-01-02", date); err != nil {
		return time.Time{}, err
	} else {
		return t, nil
	}
}

// validateSearchString validates the given search string based on specific criteria
func validateSearchString(query string) error {
	// Trim whitespace
	query = strings.TrimSpace(query)

	// Check for empty query
	if query == "" {
		return errors.New("search query cannot be empty")
	}

	// Length constraints
	if len(query) < minSearchLength {
		return fmt.Errorf("search query must be at least %d characters long", minSearchLength)
	}
	if len(query) > maxSearchLength {
		return fmt.Errorf("search query cannot exceed %d characters", maxSearchLength)
	}

	// Allowed characters: alphanumeric and spaces
	validPattern := `^[a-zA-Z0-9\s]+$`
	matched, err := regexp.MatchString(validPattern, query)
	if err != nil {
		return fmt.Errorf("error validating query: %v", err)
	}
	if !matched {
		return errors.New("search query can only contain alphanumeric characters and spaces")
	}

	// All validations passed
	return nil
}

func validateHalls(halls []string) ([]int, error) {
	var validHalls []int
	for _, hall := range halls {
		if hall == "" {
			continue
		}
		hallID, err := strconv.Atoi(hall)
		if err != nil {
			return nil, fmt.Errorf("invalid hall ID: %s", hall)
		}
		validHalls = append(validHalls, hallID)
	}
	return validHalls, nil
}
