package elastic

import (
	"fmt"
	"math"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/binning"
	"github.com/unchartedsoftware/prism/tile"
)

// Bivariate represents a bivariate tile generator.
type Bivariate struct {
	tile.Bivariate
	// tiling
	tiling bool
	bounds *binning.Bounds
	minX   int64
	maxX   int64
	minY   int64
	maxY   int64
	// binning
	binning   bool
	binSizeX  float64
	binSizeY  float64
	intervalX int64
	intervalY int64
}

func (b *Bivariate) computeTilingProps(coord *binning.TileCoord) {
	if b.tiling {
		return
	}
	// tiling params
	extents := &binning.Bounds{
		TopLeft: &binning.Coord{
			X: b.Left,
			Y: b.Top,
		},
		BottomRight: &binning.Coord{
			X: b.Right,
			Y: b.Bottom,
		},
	}
	b.bounds = binning.GetTileBounds(coord, extents)
	b.minX = int64(math.Min(b.bounds.TopLeft.X, b.bounds.BottomRight.X))
	b.maxX = int64(math.Max(b.bounds.TopLeft.X, b.bounds.BottomRight.X))
	b.minY = int64(math.Min(b.bounds.TopLeft.Y, b.bounds.BottomRight.Y))
	b.maxY = int64(math.Max(b.bounds.TopLeft.Y, b.bounds.BottomRight.Y))
	// flag as computed
	b.tiling = true
}

func (b *Bivariate) computeBinningProps(coord *binning.TileCoord) {
	if b.binning {
		return
	}
	// ensure we have tiling props
	b.computeTilingProps(coord)
	// binning params
	xRange := math.Abs(b.bounds.BottomRight.X - b.bounds.TopLeft.X)
	yRange := math.Abs(b.bounds.BottomRight.Y - b.bounds.TopLeft.Y)
	b.intervalX = int64(math.Max(1, xRange/float64(b.Resolution)))
	b.intervalY = int64(math.Max(1, yRange/float64(b.Resolution)))
	b.binSizeX = xRange / float64(b.Resolution)
	b.binSizeY = yRange / float64(b.Resolution)
	// flag as computed
	b.binning = true
}

func (b *Bivariate) GetQuery(coord *binning.TileCoord) elastic.Query {
	// compute the tiling properties
	b.computeTilingProps(coord)
	// create the range queries
	query := elastic.NewBoolQuery()
	query.Must(elastic.NewRangeQuery(b.XField).
		Gte(b.minX).
		Lt(b.maxX))
	query.Must(elastic.NewRangeQuery(b.YField).
		Gte(b.minY).
		Lt(b.maxY))
	return query
}

func (b *Bivariate) GetAggs(coord *binning.TileCoord) map[string]elastic.Aggregation {
	// compute the binning properties
	b.computeBinningProps(coord)
	// create the binning aggregations
	x := elastic.NewHistogramAggregation().
		Field(b.XField).
		Offset(b.minX).
		Interval(b.intervalX).
		MinDocCount(1)
	y := elastic.NewHistogramAggregation().
		Field(b.YField).
		Offset(b.minY).
		Interval(b.intervalY).
		MinDocCount(1)
	x.SubAggregation("y", y)
	return map[string]elastic.Aggregation{
		"x": x,
		"y": y,
	}
}

// GetXBin given an x value, returns the corresponding bin.
func (b *Bivariate) getXBin(x int) int {
	fx := float64(x)
	var bin int
	if b.bounds.TopLeft.X > b.bounds.BottomRight.X {
		bin = int(float64(b.Resolution) - ((fx - b.bounds.BottomRight.X) / b.binSizeX))
	} else {
		bin = int((fx - b.bounds.TopLeft.X) / b.binSizeX)
	}
	return b.clampBin(bin)
}

// GetYBin given an y value, returns the corresponding bin.
func (b *Bivariate) getYBin(y int) int {
	fy := float64(y)
	var bin int
	if b.bounds.TopLeft.Y > b.bounds.BottomRight.Y {
		bin = int(float64(b.Resolution) - ((fy - b.bounds.BottomRight.Y) / b.binSizeY))
	} else {
		bin = int((fy - b.bounds.TopLeft.Y) / b.binSizeY)
	}
	return b.clampBin(bin)
}

// GetBins parses the resulting histograms into bins.
func (b *Bivariate) GetBins(res *elastic.SearchResult) ([]*elastic.AggregationBucketHistogramItem, error) {
	if !b.binning {
		return nil, fmt.Errorf("binning properties have not been computed, ensure `GetAggs` is called")
	}
	// parse aggregations
	xAgg, ok := res.Aggregations.Histogram("x")
	if !ok {
		return nil, fmt.Errorf("histogram aggregation `x` was not found")
	}
	// allocate bins
	bins := make([]*elastic.AggregationBucketHistogramItem, b.Resolution*b.Resolution)
	// fill bins
	for _, xBucket := range xAgg.Buckets {
		x := xBucket.Key
		xBin := b.getXBin(int(x))
		yAgg, ok := xBucket.Histogram("y")
		if !ok {
			return nil, fmt.Errorf("histogram aggregation `y` was not found")
		}
		for _, yBucket := range yAgg.Buckets {
			y := yBucket.Key
			yBin := b.getYBin(int(y))
			index := xBin + b.Resolution*yBin
			bins[index] = yBucket
		}
	}
	return bins, nil
}

func (p *Bivariate) clampBin(bin int) int {
	if bin > p.Resolution-1 {
		return p.Resolution - 1
	}
	if bin < 0 {
		return 0
	}
	return bin
}