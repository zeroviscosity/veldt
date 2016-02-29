package elastic

import (
	"fmt"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/generation/tile"
)

// TileGenerator represents a base generator that uses elasticsearch for its
// backend.
type TileGenerator struct {
	host   string
	port   string
	client *elastic.Client
	req    *tile.Request
}

// GetHash returns the hash for this generator.
func (g *TileGenerator) GetHash() string {
	return fmt.Sprintf("%s:%s", g.host, g.port)
}