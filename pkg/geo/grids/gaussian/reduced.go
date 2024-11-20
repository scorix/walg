package gaussian

import (
	"fmt"
	"math"

	"github.com/scorix/walg/pkg/geo/grids"
	"golang.org/x/sync/singleflight"
)

var reducedCache = make(map[int]*reduced)
var reducedCacheGroup singleflight.Group

// reduced represents a reduced Gaussian grid
type reduced struct {
	n              int // Number of latitude lines between pole and equator
	TotalLatPoints int // Total number of latitude points (2*N)
	latitudes      []float64
	lonPoints      []int // Number of longitude points for each latitude
}

func NewReduced(n int) *reduced {
	cacheKey := fmt.Sprintf("N%d", n)

	res, _, _ := reducedCacheGroup.Do(cacheKey, func() (any, error) {
		if grid, ok := reducedCache[n]; ok {
			return grid, nil
		}

		grid := newReduced(n)
		reducedCache[n] = grid
		return grid, nil
	})

	return res.(*reduced)
}

// newReduced creates a new reduced Gaussian grid
func newReduced(n int) *reduced {
	if n <= 0 {
		return nil
	}

	grid := &reduced{
		n:              n,
		TotalLatPoints: 2 * n,
		latitudes:      make([]float64, 2*n),
		lonPoints:      make([]int, 2*n),
	}

	// Calculate Gaussian latitudes (simplified version)
	grid.calculateLatitudes()
	// Calculate number of longitude points for each latitude
	grid.calculateLonPoints()

	return grid
}

// calculateLonPoints calculates number of longitude points for each latitude
func (g *reduced) calculateLonPoints() {
	for i := 0; i < g.TotalLatPoints; i++ {
		lat := g.latitudes[i]
		// Convert latitude to radians
		latRad := lat * math.Pi / 180.0

		// Basic formula: nlon = 4 * N * cos(lat)
		// Round to nearest multiple of 4
		nlon := 4 * math.Round(float64(g.n)*math.Cos(latRad))

		// Ensure minimum number of points
		if nlon < 4 {
			nlon = 4
		}

		g.lonPoints[i] = int(nlon)
	}
}

// Size returns total number of grid points
func (g *reduced) Size() int {
	total := 0
	for _, nlon := range g.lonPoints {
		total += nlon
	}
	return total
}

// Latitudes returns latitudes
func (g *reduced) Latitudes() []float64 {
	return g.latitudes
}

// Longitudes returns longitude points for given latitude index
func (g *reduced) Longitudes() []float64 {
	// Return maximum longitude points (at equator)
	maxLon := 4 * g.n
	lons := make([]float64, maxLon)
	dlon := 360.0 / float64(maxLon)

	for i := 0; i < maxLon; i++ {
		lons[i] = float64(i) * dlon
	}

	return lons
}

func (g *reduced) LonPoints() []int {
	return g.lonPoints
}

// calculateLatitudes calculates Gaussian latitudes
// Note: This is a simplified version. For production use,
// you should implement proper Gaussian quadrature calculation
func (g *reduced) calculateLatitudes() {
	// Simplified version using equal spacing (not true Gaussian)
	dlat := 180.0 / float64(g.TotalLatPoints-1)
	for i := 0; i < g.TotalLatPoints; i++ {
		g.latitudes[i] = 90.0 - float64(i)*dlat
	}
}

func (g *reduced) LongitudesOnLat(lat float64) []float64 {
	// indicesLat := grids.FindNearestIndices(lat, g.latitudes)
	// lat0, lat1 := g.latitudes[indicesLat[0]], g.latitudes[indicesLat[1]]

	return nil
}

func (g *reduced) GuessNearestIndex(lat, lon float64) (int, int) {
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
