package repository

import (
	"context"
	"encoding/json"
	"io"

	"fmt"
	"fold/internal/startup"
	"strings"

	"github.com/elastic/go-elasticsearch/v8/esapi"
)

func PerformSearch(index string, searchQuery string) (*esapi.Response, error) {

	cl := startup.EsClient
	// Run the match query
	return cl.Search(
		cl.Search.WithContext(context.Background()),
		cl.Search.WithIndex(index),
		cl.Search.WithBody(strings.NewReader(searchQuery)),
		cl.Search.WithTrackTotalHits(true),
	)
}

// decodeResponse decodes the Elasticsearch response into a map.
func DecodeResponse(res *esapi.Response, v interface{}) error {
	defer res.Body.Close()
	if err := decode(res.Body, v); err != nil {
		return fmt.Errorf("error decoding response body: %w", err)
	}
	return nil
}

func decode(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}
