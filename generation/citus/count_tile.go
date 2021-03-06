package citus

import (
	"fmt"

	"github.com/unchartedsoftware/veldt"
	"github.com/unchartedsoftware/veldt/binning"
)

// Count represents a citus implementation of the count tile.
type Count struct {
	Bivariate
	Tile
}

// NewCountTile instantiates and returns a new tile struct.
func NewCountTile(cfg *Config) veldt.TileCtor {
	return func() (veldt.Tile, error) {
		t := &Count{}
		t.Config = cfg
		return t, nil
	}
}

// Parse parses the provided JSON object and populates the tiles attributes.
func (t *Count) Parse(params map[string]interface{}) error {
	return t.Bivariate.Parse(params)
}

// Create generates a tile from the provided URI, tile coordinate and query
// parameters.
func (t *Count) Create(uri string, coord *binning.TileCoord, query veldt.Query) ([]byte, error) {
	// Initialize the tile processing.
	client, citusQuery, err := t.InitializeTile(uri, query)
	if err != nil {
		return nil, err
	}

	// add tiling query
	citusQuery = t.Bivariate.AddQuery(coord, citusQuery)

	citusQuery.Select("CAST(COUNT(*) AS FLOAT) AS value")
	// send query
	res, err := client.Query(citusQuery.GetQuery(false), citusQuery.QueryArgs...)
	if err != nil {
		return nil, err
	}

	value := float64(0.0)
	for res.Next() {
		err = res.Scan(&value)
		if err != nil {
			return nil, fmt.Errorf("Error parsing count: %v",
				err)
		}
	}

	return []byte(fmt.Sprintf(`{"count":%d}`, uint64(value))), nil
}
