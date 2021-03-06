package tile

import (
	"sort"

	"github.com/unchartedsoftware/veldt/util/json"
)

// Micro represents a tile that returns individual data points with optional
// included attributes.
type Micro struct {
	LOD       int
	xField    string
	yField    string
	xIncluded bool
	yIncluded bool
}

// Parse parses the provided JSON object and populates the structs attributes.
func (m *Micro) Parse(params map[string]interface{}) error {
	// parse LOD
	m.LOD = json.GetIntDefault(params, 0, "lod")
	return nil
}

// ParseIncludes parses the included attributes to ensure they include the raw
// data coordinates.
func (m *Micro) ParseIncludes(includes []string, xField string, yField string) []string {
	// store x / y field
	m.xField = xField
	m.yField = yField
	// ensure that the x / y field are included
	if !existsIn(xField, includes) {
		includes = append(includes, xField)
	} else {
		m.xIncluded = true
	}
	if !existsIn(yField, includes) {
		includes = append(includes, yField)
	} else {
		m.yIncluded = true
	}
	return includes
}

// Encode will encode the tile results based on the LOD property.
func (m *Micro) Encode(hits []map[string]interface{}, points []float32) ([]byte, error) {
	emptyHits := false

	// remove any non-included fields from hits
	if !m.xIncluded || !m.yIncluded {
		for _, hit := range hits {
			// remove fields if they weren't explicitly included
			if !m.xIncluded {
				delete(hit, m.xField)
			}
			if !m.yIncluded {
				delete(hit, m.yField)
			}
			if !emptyHits && len(hit) == 0 {
				emptyHits = true
			}
		}
	}

	// if no hit contains any data, occlude them from response
	if emptyHits {
		// no point returning an array of empty hits
		hits = nil
	}

	// encode using LOD
	if m.LOD > 0 {
		// NOTE: during LOD points are sorted by morton code, therefore we sort
		// the hits by morton code as well to ensure both arrays align by index.
		sortHitsArray(hits, points)
		// sort points and get offsets
		sorted, offsets := LOD(points, m.LOD)
		return json.Marshal(map[string]interface{}{
			"points":  sorted,
			"offsets": offsets,
			"hits":    hits,
		})
	}
	// encode without LOD
	return json.Marshal(map[string]interface{}{
		"points": points,
		"hits":   hits,
	})
}

func existsIn(val string, arr []string) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}

func sortHitsArray(hits []map[string]interface{}, points []float32) {
	// exit early if no hits
	if hits == nil {
		return
	}
	// sort hits by morton code so they align
	hitsArr := make(hitsArray, len(hits))
	for i, hit := range hits {
		// add to hits array
		hitsArr[i] = &hitWrapper{
			x:    points[i*2],
			y:    points[i*2+1],
			data: hit,
		}
	}
	sort.Sort(hitsArr)
	// copy back into same arr
	for i, hit := range hitsArr {
		hits[i] = hit.data
	}
}

type hitWrapper struct {
	x    float32
	y    float32
	data map[string]interface{}
}

type hitsArray []*hitWrapper

func (h hitsArray) Len() int {
	return len(h)
}
func (h hitsArray) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}
func (h hitsArray) Less(i, j int) bool {
	return Morton(h[i].x, h[i].y) < Morton(h[j].x, h[j].y)
}
