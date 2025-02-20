package elasticsearch

import (
	"context"
	"foodfinder-api/config"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/rs/zerolog/log"
)

var esTypedClient *elasticsearch.TypedClient
var Index string = "foodfinder"

// init initializes the Elasticsearch client when the package is loaded.
func init() {
	var err error
	esTypedClient, err = NewTypedClient()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create Elasticsearch client")
	}
	if config.ES_INDEX.Get(config.Global()) != "" {
		Index = config.ES_INDEX.Get(config.Global())
	}
}

// NewClient creates a new Elasticsearch client.
func NewTypedClient() (*elasticsearch.TypedClient, error) {
	cfg := config.Global()

	esCfg := elasticsearch.Config{
		Addresses: []string{config.ES_URL.Get(cfg)},
		Username:  config.ES_USER.Get(cfg),
		Password:  config.ES_PASS.Get(cfg),
	}
	if config.DEBUG.Get(cfg) == "true" {
		esCfg.EnableDebugLogger = true
		esCfg.Logger = &CustomLogger{log.Logger}
	}
	return elasticsearch.NewTypedClient(esCfg)
}

type esClientKey struct{}

// Client returns the Elasticsearch client.
func TypedClient() *elasticsearch.TypedClient {
	return esTypedClient
}

// FromContext returns the Elasticsearch client attached to the context.
func FromContext(ctx context.Context) (*elasticsearch.TypedClient, error) {
	if es, ok := ctx.Value(esClientKey{}).(*elasticsearch.TypedClient); ok {
		return es, nil
	}
	return NewTypedClient()
}

// WithContext returns a new context with the Elasticsearch client attached.
func WithContext(ctx context.Context, es *elasticsearch.TypedClient) context.Context {
	return context.WithValue(ctx, esClientKey{}, es)
}
