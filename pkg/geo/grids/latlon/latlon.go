package latlon

import (
	"cmp"
	"math"

	"github.com/scorix/walg/pkg/geo/distance"
	"github.com/scorix/walg/pkg/geo/grids"
)

// latLon 表示一个基于经纬度的等间距网格系统
type latLon struct {
	minLat   float64        // 最小纬度
	maxLat   float64        // 最大纬度
	minLon   float64        // 最小经度
	maxLon   float64        // 最大经度
	latStep  float64        // 纬度步长
	lonStep  float64        // 经度步长
	latCount int            // 纬度方向的网格数量
	lonCount int            // 经度方向的网格数量
	lats     []float64      // 缓存的纬度值
	lons     []float64      // 缓存的经度值
	scanMode grids.ScanMode // 新增: 扫描模式
}

type latlonOption func(*latLon)

func WithScanMode(scanMode grids.ScanMode) latlonOption {
	return func(ll *latLon) {
		ll.scanMode = scanMode
	}
}

// NewLatLonGrid 创建一个新的经纬度网格
func NewLatLonGrid(minLat, maxLat, minLon, maxLon, latStep, lonStep float64, opts ...latlonOption) *latLon {
	ll := &latLon{
		minLat:  math.Min(minLat, maxLat),
		maxLat:  math.Max(minLat, maxLat),
		minLon:  minLon,
		maxLon:  maxLon,
		latStep: latStep,
		lonStep: lonStep,
	}

	// 先应用选项，这样我们可以获取 scanMode
	for _, o := range opts {
		o(ll)
	}

	// 计算点数
	ll.latCount = int(math.Abs(maxLat-minLat)/latStep) + 1
	ll.lonCount = int((maxLon-minLon)/lonStep) + 1

	// 纬度数组总是从北到南排列
	ll.lats = make([]float64, ll.latCount)
	for i := 0; i < ll.latCount; i++ {
		ll.lats[i] = ll.maxLat - float64(i)*latStep
	}

	// 经度数组总是从西到东排列
	ll.lons = make([]float64, ll.lonCount)
	for i := 0; i < ll.lonCount; i++ {
		ll.lons[i] = ll.minLon + float64(i)*lonStep
	}

	return ll
}

// Size 返回网格总数
func (g *latLon) Size() int {
	return g.latCount * g.lonCount
}

// Latitudes 返回所有纬度值
func (g *latLon) Latitudes() []float64 {
	return g.lats
}

// Longitudes 返回所有经度值
func (g *latLon) Longitudes() []float64 {
	return g.lons
}

func (g *latLon) ScanMode() grids.ScanMode {
	return g.scanMode
}

// normalizeLat 将一个纬度值映射到指定区间 [minVal, maxVal]
func (g *latLon) normalizeLat(minVal, maxVal, val float64) float64 {
	intMinVal := int(minVal * 1e6)
	intMaxVal := int(maxVal * 1e6)
	intVal := int(val * 1e6)

	// Handle values outside the range
	interval := intMaxVal - intMinVal
	normalized := intVal

	// Wrap around until we're in range
	for normalized < intMinVal {
		normalized += interval
	}
	for normalized > intMaxVal {
		normalized -= interval
	}

	return float64(normalized) / 1e6
}

// normalizeLon 将一个经度值映射到指定区间 [minVal, maxVal]
func (g *latLon) normalizeLon(minVal, maxVal, val float64, sphere bool) float64 {
	intMinVal := int(minVal * 1e6)
	intMaxVal := int(maxVal * 1e6)
	intVal := int(val * 1e6)
	interval := intMaxVal - intMinVal

	// 处理球状的经度
	if sphere {
		// 先将经度标准化到 [0, 360) 范围内
		normalized := intVal
		for normalized < 0 {
			normalized += 360 * 1e6
		}
		for normalized >= 360*1e6 {
			normalized -= 360 * 1e6
		}

		// 如果标准化后的值超出了 maxVal
		if normalized > intMaxVal {
			// 计算与区间两端的距离
			distToMax := math.Abs(float64(normalized - intMaxVal))
			distToMin := math.Abs(float64(normalized - 360*1e6 - intMinVal))

			// 如果更接近 minVal，则减去 360
			if distToMin < distToMax {
				normalized -= 360 * 1e6
			}
		}

		return float64(normalized) / 1e6
	}

	// 非球状处理保持不变
	normalized := intVal

	for normalized < intMinVal {
		normalized += interval
	}
	for normalized > intMaxVal {
		normalized -= interval
	}

	return float64(normalized) / 1e6
}

func (g *latLon) GetNearestIndex(lat, lon float64) (int, int) {
	isSphere := func() bool {
		maxLon := int(g.maxLon * 1e6)
		minLon := int(g.minLon * 1e6)
		return maxLon+int(g.lonStep*1e6) == minLon+360*1e6
	}()

	normalizedLat := g.normalizeLat(g.minLat, g.maxLat, lat)
	normalizedLon := g.normalizeLon(g.minLon, g.maxLon, lon, isSphere)

	indicesLat := grids.FindNearestIndices(normalizedLat, g.lats)
	indicesLon := grids.FindNearestIndices(normalizedLon, g.lons)

	latIdx := indicesLat[0]
	lonIdx := indicesLon[0]

	// iterations 3 is enough for comparing the minimum distance
	const iterations = 3

	dist := distance.VincentyIterations(lat, lon, g.lats[latIdx], g.lons[lonIdx], iterations)
	for _, i := range indicesLat {
		for _, j := range indicesLon {
			d := distance.VincentyIterations(lat, lon, g.lats[i], g.lons[j], iterations)
			if cmp.Compare(d, dist) < 0 {
				dist = d
				latIdx = i
				lonIdx = j
			}
		}
	}

	return latIdx, lonIdx
}

func (g *latLon) GuessNearestIndex(lat, lon float64) (int, int) {
	isSphere := func() bool {
		maxLon := int(g.maxLon * 1e6)
		minLon := int(g.minLon * 1e6)
		return maxLon+int(g.lonStep*1e6) == minLon+360*1e6
	}()

	normalizedLat := g.normalizeLat(g.minLat, g.maxLat, lat)
	normalizedLon := g.normalizeLon(g.minLon, g.maxLon, lon, isSphere)

	indicesLat := grids.FindNearestIndices(normalizedLat, g.lats)
	indicesLon := grids.FindNearestIndices(normalizedLon, g.lons)

	latIdx := indicesLat[0]
	lonIdx := indicesLon[0]

	if math.Abs(g.lats[latIdx]-lat) > math.Abs(g.lats[indicesLat[1]]-lat) {
		latIdx = indicesLat[1]
	}

	if math.Abs(g.lons[lonIdx]-lon) > math.Abs(g.lons[indicesLon[1]]-lon) {
		lonIdx = indicesLon[1]
	}

	return latIdx, lonIdx
}
