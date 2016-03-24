package elastic

import (
	"encoding/json"
	"fmt"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/generation/elastic/param"
	"github.com/unchartedsoftware/prism/generation/elastic/throttle"
	"github.com/unchartedsoftware/prism/generation/tile"
)

// TopicCountTile represents a tiling generator that produces terms counts.
type TopicCountTile struct {
	TileGenerator
	Tiling    *param.Tiling
	Terms     *param.TermsAgg
	Range     *param.Range
	Histogram *param.Histogram
}

// NewTopicCountTile instantiates and returns a pointer to a new generator.
func NewTopicCountTile(host, port string) tile.GeneratorConstructor {
	return func(tileReq *tile.Request) (tile.Generator, error) {
		client, err := NewClient(host, port)
		if err != nil {
			return nil, err
		}
		tiling, err := param.NewTiling(tileReq)
		if err != nil {
			return nil, err
		}
		terms, err := param.NewTermsAgg(tileReq)
		if err != nil {
			return nil, err
		}
		rang, _ := param.NewRange(tileReq)
		histogram, _ := param.NewHistogram(tileReq)
		t := &TopicCountTile{}
		t.Tiling = tiling
		t.Terms = terms
		t.Range = rang
		t.Histogram = histogram
		t.req = tileReq
		t.host = host
		t.port = port
		t.client = client
		return t, nil
	}
}

// GetParams returns a slice of tiling parameters.
func (g *TopicCountTile) GetParams() []tile.Param {
	return []tile.Param{
		g.Tiling,
		g.Terms,
		g.Range,
		g.Histogram,
	}
}

func (g *TopicCountTile) getQuery() elastic.Query {
	// optional filters
	filters := elastic.NewBoolQuery()
	// if range param is provided, add range queries
	if g.Range != nil {
		for _, query := range g.Range.GetQueries() {
			filters.Must(query)
		}
	}
	// if terms param is provided, add terms query
	if g.Terms != nil {
		filters.Must(g.Terms.GetQuery())
	}
	return elastic.NewBoolQuery().
		Must(g.Tiling.GetXQuery()).
		Must(g.Tiling.GetYQuery()).
		Must(filters)
}

func (g *TopicCountTile) addAggs(query *elastic.SearchService) *elastic.SearchService {
	// add all filter aggregations
	for term, agg := range g.Terms.GetAggregations() {
		// if histogram param is provided, add histogram agg
		if g.Histogram != nil {
			agg.SubAggregation(histogramAggName, g.Histogram.GetAggregation())
		}
		query.Aggregation(term, agg)
	}
	return query
}

func (g *TopicCountTile) parseResult(res *elastic.SearchResult) ([]byte, error) {
	// build map of topics and counts
	counts := make(map[string]interface{})
	for _, term := range g.Terms.Terms {
		filter, ok := res.Aggregations.Filter(term)
		if !ok {
			return nil, fmt.Errorf("Filter aggregation '%s' was not found in response for request %s",
				term,
				g.req.String())
		}
		if filter.DocCount > 0 {
			if g.Histogram != nil {
				histogramAgg, ok := filter.Aggregations.Histogram(histogramAggName)
				if !ok {
					return nil, fmt.Errorf("Histogram aggregation '%s' was not found in response for request %s",
						histogramAggName,
						g.req.String())
				}
				counts[term] = g.Histogram.GetBucketMap(histogramAgg)
			} else {
				counts[term] = filter.DocCount
			}
		}
	}
	// marshal results map
	return json.Marshal(counts)
}

// GetTile returns the marshalled tile data.
func (g *TopicCountTile) GetTile() ([]byte, error) {
	// build query
	query := g.client.
		Search(g.req.Index).
		Size(0).
		Query(g.getQuery())
	// add all filter aggregations
	query = g.addAggs(query)
	// send query through equalizer
	res, err := throttle.Send(query)
	if err != nil {
		return nil, err
	}
	// parse and return results
	return g.parseResult(res)
}
