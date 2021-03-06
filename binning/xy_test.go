package binning_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/unchartedsoftware/veldt/binning"
	"github.com/unchartedsoftware/veldt/geometry"
)

var _ = Describe("bounds", func() {

	const (
		epsilon = 0.000001
	)

	var (
		t000 = binning.TileCoord{
			X: 0,
			Y: 0,
			Z: 0,
		}
		t100 = binning.TileCoord{
			X: 0,
			Y: 0,
			Z: 1,
		}
		t111 = binning.TileCoord{
			X: 1,
			Y: 1,
			Z: 1,
		}
		extent = geometry.NewBounds(-1, 1, -1, 1)

		bottomLeft = geometry.Coord{
			X: -1,
			Y: -1,
		}
		center = geometry.Coord{
			X: 0,
			Y: 0,
		}
		topRight = geometry.Coord{
			X: 1,
			Y: 1,
		}
	)

	Describe("GetTileBounds", func() {
		It("should return bounds from a tile coordinate", func() {
			bounds := binning.GetTileBounds(&t000, extent)
			Expect(bounds.Left).To(BeNumerically("~", -1.0, epsilon))
			Expect(bounds.Bottom).To(BeNumerically("~", -1.0, epsilon))
			Expect(bounds.Right).To(BeNumerically("~", 1.0, epsilon))
			Expect(bounds.Top).To(BeNumerically("~", 1.0, epsilon))

			bounds = binning.GetTileBounds(&t100, extent)
			Expect(bounds.Left).To(BeNumerically("~", -1.0, epsilon))
			Expect(bounds.Bottom).To(BeNumerically("~", -1.0, epsilon))
			Expect(bounds.Right).To(BeNumerically("~", 0.0, epsilon))
			Expect(bounds.Top).To(BeNumerically("~", 0.0, epsilon))

			bounds = binning.GetTileBounds(&t111, extent)
			Expect(bounds.Left).To(BeNumerically("~", 0.0, epsilon))
			Expect(bounds.Bottom).To(BeNumerically("~", 0.0, epsilon))
			Expect(bounds.Right).To(BeNumerically("~", 1.0, epsilon))
			Expect(bounds.Top).To(BeNumerically("~", 1.0, epsilon))
		})
		It("should support bounds where right > left", func() {
			bounds := binning.GetTileBounds(&t000, extent)
			Expect(bounds.Left).To(BeNumerically("~", -1.0, epsilon))
			Expect(bounds.Bottom).To(BeNumerically("~", -1.0, epsilon))
			Expect(bounds.Right).To(BeNumerically("~", 1.0, epsilon))
			Expect(bounds.Top).To(BeNumerically("~", 1.0, epsilon))
		})
		It("should support bounds where bottom > top", func() {
			bounds := binning.GetTileBounds(&t000, extent)
			Expect(bounds.Left).To(BeNumerically("~", -1.0, epsilon))
			Expect(bounds.Bottom).To(BeNumerically("~", -1.0, epsilon))
			Expect(bounds.Right).To(BeNumerically("~", 1.0, epsilon))
			Expect(bounds.Top).To(BeNumerically("~", 1.0, epsilon))
		})
		It("should support bounds where left > right", func() {
			bounds := binning.GetTileBounds(&t000, geometry.NewBounds(1, -1, -1, 1))

			Expect(bounds.Left).To(BeNumerically("~", 1.0, epsilon))
			Expect(bounds.Bottom).To(BeNumerically("~", -1.0, epsilon))
			Expect(bounds.Right).To(BeNumerically("~", -1.0, epsilon))
			Expect(bounds.Top).To(BeNumerically("~", 1.0, epsilon))
		})
		It("should support bounds where top > bottom", func() {
			bounds := binning.GetTileBounds(&t000, geometry.NewBounds(-1, 1, 1, -1))

			Expect(bounds.Left).To(BeNumerically("~", -1.0, epsilon))
			Expect(bounds.Bottom).To(BeNumerically("~", 1.0, epsilon))
			Expect(bounds.Right).To(BeNumerically("~", 1.0, epsilon))
			Expect(bounds.Top).To(BeNumerically("~", -1.0, epsilon))
		})
	})

	Describe("CoordToFractionalTile", func() {
		It("should return a fractional tile coordinate", func() {
			tile := binning.CoordToFractionalTile(&bottomLeft, 0, extent)
			Expect(tile.X).To(BeNumerically("~", 0.0, epsilon))
			Expect(tile.Y).To(BeNumerically("~", 0.0, epsilon))

			tile = binning.CoordToFractionalTile(&center, 1, extent)
			Expect(tile.X).To(BeNumerically("~", 1.0, epsilon))
			Expect(tile.Y).To(BeNumerically("~", 1.0, epsilon))

			tile = binning.CoordToFractionalTile(&topRight, 1, extent)
			Expect(tile.X).To(BeNumerically("~", 2.0, epsilon))
			Expect(tile.Y).To(BeNumerically("~", 2.0, epsilon))
		})
	})

})
