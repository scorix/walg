package gaussian

import (
	"fmt"
	"math"

	"github.com/scorix/walg/pkg/geo/distance"
	"github.com/scorix/walg/pkg/geo/grids"
	"golang.org/x/sync/singleflight"
)

var octahedralCache = make(map[int]*octahedral)
var octahedralCacheGroup singleflight.Group

type octahedral struct {
	n         int
	latitudes []float64
	lonPoints []int
}

func NewOctahedral(n int) *octahedral {
	cacheKey := fmt.Sprintf("O%d", n)

	res, _, _ := octahedralCacheGroup.Do(cacheKey, func() (any, error) {
		if grid, ok := octahedralCache[n]; ok {
			return grid, nil
		}

		grid := newOctahedral(n)
		octahedralCache[n] = grid
		return grid, nil
	})

	return res.(*octahedral)
}

func newOctahedral(n int) *octahedral {
	o := &octahedral{
		n: n,
	}

	o.latitudes = o.calcLatitudes()
	o.lonPoints = o.calcLonPoints()

	return o
}

// Grid interface implementation
func (g *octahedral) Size() int {
	total := 0
	for _, nlon := range g.lonPoints {
		total += nlon * len(g.latitudes)
	}
	return total
}

func (g *octahedral) Latitudes() []float64 {
	return g.latitudes
}

func (g *octahedral) Longitudes() []float64 {
	// Return maximum longitude points (at equator)
	maxLon := 4 * g.n
	lons := make([]float64, maxLon)
	dlon := 360.0 / float64(maxLon)

	for i := 0; i < maxLon; i++ {
		lons[i] = float64(i) * dlon
	}

	return lons
}

func (g *octahedral) LongitudesOnLat(lat float64) []float64 {
	indicesLat := grids.FindNearestIndices(lat, g.latitudes)
	lat0, lat1 := g.latitudes[indicesLat[0]], g.latitudes[indicesLat[1]]

	var nearestLatIdx int
	if math.Abs(lat0-lat) <= math.Abs(lat1-lat) {
		nearestLatIdx = indicesLat[0]
	} else {
		nearestLatIdx = indicesLat[1]
	}

	fmt.Printf("lat0: %v, lat1: %v, nearestLatIdx: %v\n", lat0, lat1, nearestLatIdx)

	nlon := g.lonPoints[nearestLatIdx]
	dlon := 360.0 / float64(nlon)

	lons := make([]float64, nlon)
	for i := 0; i < nlon; i++ {
		lons[i] = float64(i) * dlon
	}

	return lons
}

func (g *octahedral) LonPoints() []int {
	return g.lonPoints
}

func (g *octahedral) GuessNearestIndex(lat, lon float64) (int, int) {
	latitudes := g.Latitudes()
	longitudes := g.LongitudesOnLat(lat)

	indicesLat := grids.FindNearestIndices(lat, latitudes)
	indicesLon := grids.FindNearestIndices(lon, longitudes)

	latIdx := indicesLat[0]
	lonIdx := indicesLon[0]

	if math.Abs(latitudes[latIdx]-lat) > math.Abs(latitudes[indicesLat[1]]-lat) {
		latIdx = indicesLat[1]
	}

	if math.Abs(longitudes[lonIdx]-lon) > math.Abs(longitudes[indicesLon[1]]-lon) {
		lonIdx = indicesLon[1]
	}

	return latIdx, lonIdx
}

func (g *octahedral) GetNearestIndex(lat, lon float64) (int, int) {
	latitudes := g.Latitudes()
	indicesLat := grids.FindNearestIndices(lat, latitudes)

	latIdx := indicesLat[0]
	lonIdx := 0
	minDist := math.MaxFloat64

	const iterations = 3

	// 标准化经度到 [0, 360)
	lon = math.Mod(lon+360.0, 360.0)

	for _, i := range indicesLat {
		nlon := g.lonPoints[i]
		dlon := 360.0 / float64(nlon)

		// 计算最近的经度索引
		centerLonIdx := int(math.Round(lon / dlon))
		if centerLonIdx == nlon {
			centerLonIdx = 0
		}

		// 检查相邻点，包括跨越 0°/360° 的情况
		indices := []int{
			(centerLonIdx - 1 + nlon) % nlon,
			centerLonIdx % nlon,
			(centerLonIdx + 1) % nlon,
		}

		for _, j := range indices {
			gridLon := float64(j) * dlon
			d := distance.VincentyIterations(lat, lon, latitudes[i], gridLon, iterations)
			if d < minDist {
				minDist = d
				latIdx = i
				lonIdx = j
			}
		}
	}

	return latIdx, lonIdx
}

func (g *octahedral) calcLatitudes() []float64 {
	return gaussLegendreZeros(g.n * 2)
}

func (g *octahedral) calcLonPoints() []int {
	points := make([]int, len(g.latitudes))

	for i, lat := range g.latitudes {
		// Convert to colatitude (0 at North Pole)
		colat := 90.0 - lat
		colatRad := colat * math.Pi / 180.0

		// Calculate number of longitude points using sine of colatitude
		nlon := 4 * g.n * int(math.Round(math.Sin(colatRad)))

		// Ensure minimum number of points
		if nlon < 4 {
			nlon = 4
		}

		points[i] = nlon
	}

	return points
}
