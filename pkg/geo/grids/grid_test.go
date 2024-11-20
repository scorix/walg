package grids_test

import (
	"testing"

	"github.com/scorix/walg/pkg/geo/distance"
	"github.com/scorix/walg/pkg/geo/grids"
	"github.com/scorix/walg/pkg/geo/grids/gaussian"
	"github.com/scorix/walg/pkg/geo/grids/latlon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 添加 gridTestCase 类型定义
type gridTestCase struct {
	name        string
	lat         float64
	lon         float64
	expectedIdx int
	gridLat     float64
	gridLon     float64
}

// 添加 runGridTests 函数
func runGridTests(t *testing.T, grid grids.Grid, tests []gridTestCase, mode grids.ScanMode) {
	const iterations = 5

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx := grids.GridIndex(grid, tt.lat, tt.lon, mode)
			recoveredLat, recoveredLon := grids.GridPoint(grid, idx, mode)

			actualDist := distance.VincentyIterations(tt.lat, tt.lon, recoveredLat, recoveredLon, iterations)
			t.Logf("Actual #%d index: (%.3f, %.3f), dist: %fkm from (%.3f, %.3f)",
				idx, recoveredLat, recoveredLon, actualDist, tt.lat, tt.lon)
			assert.GreaterOrEqual(t, actualDist, 0.0)

			nearestIdxs := grids.NewNearestGrids(grid).NearestGrids(tt.lat, tt.lon, mode)
			for i, nearIdx := range nearestIdxs {
				nearLat, nearLon := grids.GridPoint(grid, nearIdx, mode)
				dist := distance.VincentyIterations(tt.lat, tt.lon, nearLat, nearLon, iterations)
				t.Logf("Nearest %d index #%d: (%.3f, %.3f), dist: %fkm from (%.3f, %.3f)",
					i, nearIdx, nearLat, nearLon, dist, tt.lat, tt.lon)
				assert.GreaterOrEqual(t, dist, actualDist)
			}

			expectedLat, expectedLon := grids.GridPoint(grid, tt.expectedIdx, mode)
			expectedDist := distance.VincentyIterations(tt.lat, tt.lon, expectedLat, expectedLon, iterations)
			t.Logf("Expected #%d index: (%.3f, %.3f), dist: %f from (%.3f, %.3f)",
				tt.expectedIdx, expectedLat, expectedLon, expectedDist, tt.lat, tt.lon)

			if idx != tt.expectedIdx {
				require.GreaterOrEqual(t, expectedDist, actualDist)
				assert.InDelta(t, expectedDist, actualDist, 1e-3)
			}

			if t.Failed() {
				assert.Equal(t, tt.expectedIdx, idx, "GridIndex mismatch for point (%.3f, %.3f)", tt.lat, tt.lon)
				assert.Equal(t, tt.gridLat, recoveredLat, "Latitude mismatch")
				assert.Equal(t, tt.gridLon, recoveredLon, "Longitude mismatch")

				t.Logf("Expected grid index %d: (%.3f, %.3f) -> (%.3f, %.3f) dist: %f",
					tt.expectedIdx, tt.lat, tt.lon, expectedLat, expectedLon, expectedDist)
			}

			// guess
			guessIdx := grids.GuessGridIndex(grid, tt.lat, tt.lon, mode)
			guessLat, guessLon := grids.GridPoint(grid, guessIdx, mode)
			guessDist := distance.VincentyIterations(tt.lat, tt.lon, guessLat, guessLon, iterations)
			t.Logf("Guess index #%d: (%.3f, %.3f), dist: %f from (%.3f, %.3f)",
				guessIdx, guessLat, guessLon, guessDist, tt.lat, tt.lon)

			assert.InDelta(t, expectedDist, guessDist, 1e-3)
		})
	}
}

var defaultRegularLatLonGridTests = []gridTestCase{
	// 网格点精确匹配
	{name: "North Pole", lat: 90.0, lon: 0.0, expectedIdx: 0, gridLat: 90.0, gridLon: 0.0},
	{name: "Second Row Start", lat: 89.75, lon: 0.0, expectedIdx: 1440, gridLat: 89.75, gridLon: 0.0},

	// 纬度边界值测试（经度固定为0.0）
	{name: "Near 90 degrees", lat: 89.99, lon: 0.0, expectedIdx: 0, gridLat: 90.0, gridLon: 0.0},
	{name: "Close to 90", lat: 89.88, lon: 0.0, expectedIdx: 0, gridLat: 90.0, gridLon: 0.0},
	{name: "Middle Point", lat: 89.875, lon: 0.0, expectedIdx: 1440, gridLat: 89.75, gridLon: 0.0},
	{name: "Near 89.75", lat: 89.87, lon: 0.0, expectedIdx: 1440, gridLat: 89.75, gridLon: 0.0},
	{name: "Very Close to 89.75", lat: 89.76, lon: 0.0, expectedIdx: 1440, gridLat: 89.75, gridLon: 0.0},
	{name: "Exact 89.75", lat: 89.75, lon: 0.0, expectedIdx: 1440, gridLat: 89.75, gridLon: 0.0},

	// 经度边界值测试（纬度固定为90.0）
	{name: "0 Longitude", lat: 90.0, lon: 0.0, expectedIdx: 0, gridLat: 90.0, gridLon: 0.0},
	{name: "Very Near 0", lat: 90.0, lon: 0.12, expectedIdx: 0, gridLat: 90.0, gridLon: 0.0},
	{name: "Near 0.25", lat: 90.0, lon: 0.13, expectedIdx: 1, gridLat: 90.0, gridLon: 0.25},
	{name: "Middle Longitude", lat: 90.0, lon: 0.125, expectedIdx: 0, gridLat: 90.0, gridLon: 0.0},
	{name: "Close to 0.25", lat: 90.0, lon: 0.24, expectedIdx: 1, gridLat: 90.0, gridLon: 0.25},
	{name: "Exact 0.25", lat: 90.0, lon: 0.25, expectedIdx: 1, gridLat: 90.0, gridLon: 0.25},
	{name: "0 Longitude", lat: 90.0, lon: 360.0, expectedIdx: 0, gridLat: 90.0, gridLon: 0.0},

	// 组合测试
	{name: "Combined Near 90", lat: 89.88, lon: 0.13, expectedIdx: 1, gridLat: 90.0, gridLon: 0.25},
	{name: "Combined Near 89.75", lat: 89.87, lon: 0.13, expectedIdx: 1441, gridLat: 89.75, gridLon: 0.25},
	{name: "Combined Middle", lat: 89.875, lon: 0.125, expectedIdx: 0, gridLat: 90.0, gridLon: 0.0},

	// 特殊地理位置
	{name: "Last Point", lat: -90.0, lon: 359.75, expectedIdx: (90+90)/0.25*1440 + (1440 - 1), gridLat: -90.0, gridLon: 359.75},
	{name: "Equator 180", lat: 0.0, lon: 180.0, expectedIdx: (90-0)/0.25*1440 + 720, gridLat: 0.0, gridLon: 180.0},
	{name: "Equator Greenwich", lat: 0.0, lon: 0.0, expectedIdx: (90-0)/0.25*1440 + 0, gridLat: 0.0, gridLon: 0.0},
	{name: "Date Line", lat: 0.0, lon: 180.0, expectedIdx: (90-0)/0.25*1440 + 720, gridLat: 0.0, gridLon: 180.0},
	{name: "South Pole", lat: -90.0, lon: 0.0, expectedIdx: (90 + 90) / 0.25 * 1440, gridLat: -90.0, gridLon: 0.0},

	// 重要纬线
	{name: "Arctic Circle", lat: 66.5, lon: 0.0, expectedIdx: (90 - 66.5) / 0.25 * 1440, gridLat: 66.5, gridLon: 0.0},
	{name: "Tropic of Cancer", lat: 23.5, lon: 0.0, expectedIdx: (90 - 23.5) / 0.25 * 1440, gridLat: 23.5, gridLon: 0.0},
	{name: "Equator", lat: 0.0, lon: 0.0, expectedIdx: (90 - 0) / 0.25 * 1440, gridLat: 0.0, gridLon: 0.0},
	{name: "Tropic of Capricorn", lat: -23.5, lon: 0.0, expectedIdx: (90 + 23.5) / 0.25 * 1440, gridLat: -23.5, gridLon: 0.0},
	{name: "Antarctic Circle", lat: -66.5, lon: 0.0, expectedIdx: (90 + 66.5) / 0.25 * 1440, gridLat: -66.5, gridLon: 0.0},

	// 重要经线
	{name: "Prime Meridian", lat: 0.0, lon: 0.0, expectedIdx: (90-0)/0.25*1440 + 0, gridLat: 0.0, gridLon: 0.0},
	{name: "180 Meridian", lat: 0.0, lon: 180.0, expectedIdx: (90-0)/0.25*1440 + 720, gridLat: 0.0, gridLon: 180.0},
	{name: "90E Meridian", lat: 0.0, lon: 90.0, expectedIdx: (90-0)/0.25*1440 + 360, gridLat: 0.0, gridLon: 90.0},
	{name: "90W Meridian", lat: 0.0, lon: -90.0, expectedIdx: (90-0)/0.25*1440 + 1080, gridLat: 0.0, gridLon: 270.0},

	// 著名地点
	{name: "London", lat: 51.5, lon: -0.13, expectedIdx: 223199, gridLat: 51.5, gridLon: 359.75},
	{name: "New York", lat: 40.75, lon: -74.0, expectedIdx: 284824, gridLat: 40.75, gridLon: 286.0},
	{name: "Beijing", lat: 39.9, lon: 116.4, expectedIdx: 288466, gridLat: 40.0, gridLon: 116.5},
	{name: "Tokyo", lat: 35.7, lon: 139.7, expectedIdx: 313039, gridLat: 35.75, gridLon: 139.75},
	{name: "Sydney", lat: -33.9, lon: 151.2, expectedIdx: 714845, gridLat: -34.0, gridLon: 151.25},

	// 边界情况
	{name: "Near North Pole Low Lon", lat: 89.99, lon: 0.12, expectedIdx: 0, gridLat: 90.0, gridLon: 0.0},
	{name: "Near North Pole High Lon", lat: 89.99, lon: 0.13, expectedIdx: 1, gridLat: 90.0, gridLon: 0.25},
	{name: "Near South Pole High Lon", lat: -89.99, lon: 359.88, expectedIdx: 1036800, gridLat: -90.0, gridLon: 0.0},
	{name: "Near South Pole Low Lon", lat: -89.99, lon: 359.87, expectedIdx: 1038239, gridLat: -90.0, gridLon: 359.75},
	{name: "Near Date Line", lat: 0.13, lon: 179.87, expectedIdx: 517679, gridLat: 0.25, gridLon: 179.75},
	{name: "Near Prime Meridian", lat: 0.13, lon: -0.13, expectedIdx: 518399, gridLat: 0.25, gridLon: 359.75},
}

func TestLatLon_DefaultScanMode(t *testing.T) {
	// Convert GRIB units (millionths of a degree) to degrees
	firstLat := 90.0  // 90000000 millionths -> 90.0 degrees
	lastLat := -90.0  // -90000000 millionths -> -90.0 degrees
	firstLon := 0.0   // 0 millionths -> 0.0 degrees
	lastLon := 359.75 // 359750000 millionths -> 359.75 degrees
	latStep := 0.25   // 250000 millionths -> 0.25 degrees
	lonStep := 0.25   // 250000 millionths -> 0.25 degrees

	grid := latlon.NewLatLonGrid(
		lastLat,  // minLat
		firstLat, // maxLat
		firstLon, // minLon
		lastLon,  // maxLon
		latStep,
		lonStep,
	)

	runGridTests(t, grid, defaultRegularLatLonGridTests, 0)
}

var negativeIRegularLatLonGridTests = []gridTestCase{
	// 网格点精确匹配
	{name: "First Row Start", lat: 90.0, lon: 359.75, expectedIdx: 0, gridLat: 90.0, gridLon: 359.75},
	{name: "First Row End", lat: 90.0, lon: 0.0, expectedIdx: 1439, gridLat: 90.0, gridLon: 0.0},
	{name: "Second Row Start", lat: 89.75, lon: 359.75, expectedIdx: 1440, gridLat: 89.75, gridLon: 359.75},
	{name: "Second Row End", lat: 89.75, lon: 0.0, expectedIdx: 2879, gridLat: 89.75, gridLon: 0.0},

	// 纬度边界值测试（经度固定为359.75）
	{name: "Near 90 degrees", lat: 89.99, lon: 359.75, expectedIdx: 0, gridLat: 90.0, gridLon: 359.75},
	{name: "Close to 90", lat: 89.88, lon: 359.75, expectedIdx: 0, gridLat: 90.0, gridLon: 359.75},
	{name: "Middle Point", lat: 89.875, lon: 359.75, expectedIdx: 1440, gridLat: 89.75, gridLon: 359.75},
	{name: "Near 89.75", lat: 89.87, lon: 359.75, expectedIdx: 1440, gridLat: 89.75, gridLon: 359.75},
	{name: "Very Close to 89.75", lat: 89.76, lon: 359.75, expectedIdx: 1440, gridLat: 89.75, gridLon: 359.75},
	{name: "Exact 89.75", lat: 89.75, lon: 359.75, expectedIdx: 1440, gridLat: 89.75, gridLon: 359.75},

	// 经度边界值测试（纬度固定为90.0）
	{name: "359.75 Longitude", lat: 90.0, lon: 359.75, expectedIdx: 0, gridLat: 90.0, gridLon: 359.75},
	{name: "Very Near 359.75", lat: 90.0, lon: 359.87, expectedIdx: 0, gridLat: 90.0, gridLon: 359.75},
	{name: "Near 359.5", lat: 90.0, lon: 359.62, expectedIdx: 1, gridLat: 90.0, gridLon: 359.5},
	{name: "Middle Longitude", lat: 90.0, lon: 359.625, expectedIdx: 1, gridLat: 90.0, gridLon: 359.50},
	{name: "Close to 359.5", lat: 90.0, lon: 359.51, expectedIdx: 1, gridLat: 90.0, gridLon: 359.5},
	{name: "Exact 359.5", lat: 90.0, lon: 359.5, expectedIdx: 1, gridLat: 90.0, gridLon: 359.5},

	// 组合测试
	{name: "Combined Near 90", lat: 89.88, lon: 359.62, expectedIdx: 1, gridLat: 90.0, gridLon: 359.5},
	{name: "Combined Near 89.75", lat: 89.87, lon: 359.62, expectedIdx: 1441, gridLat: 89.75, gridLon: 359.5},
	{name: "Combined Middle", lat: 89.875, lon: 359.625, expectedIdx: 1, gridLat: 90.0, gridLon: 359.50},

	// 特殊地理位置
	{name: "Last Point", lat: -90.0, lon: 0.0, expectedIdx: 721*1440 - 1, gridLat: -90.0, gridLon: 0.0},
	{name: "Equator 180", lat: 0.0, lon: 180.0, expectedIdx: 360*1440 + 719, gridLat: 0.0, gridLon: 180.0},
	{name: "Equator Greenwich", lat: 0.0, lon: 359.75, expectedIdx: 360*1440 + 0, gridLat: 0.0, gridLon: 359.75},
	{name: "Date Line", lat: 0.0, lon: 180.0, expectedIdx: 360*1440 + 719, gridLat: 0.0, gridLon: 180.0},
	{name: "South Pole", lat: -90.0, lon: 359.75, expectedIdx: 720 * 1440, gridLat: -90.0, gridLon: 359.75},

	// 重要纬线
	{name: "Arctic Circle", lat: 66.5, lon: 359.75, expectedIdx: 94*1440 + 0, gridLat: 66.5, gridLon: 359.75},
	{name: "Tropic of Cancer", lat: 23.5, lon: 359.75, expectedIdx: 266*1440 + 0, gridLat: 23.5, gridLon: 359.75},
	{name: "Equator", lat: 0.0, lon: 359.75, expectedIdx: 360*1440 + 0, gridLat: 0.0, gridLon: 359.75},
	{name: "Tropic of Capricorn", lat: -23.5, lon: 359.75, expectedIdx: 454*1440 + 0, gridLat: -23.5, gridLon: 359.75},
	{name: "Antarctic Circle", lat: -66.5, lon: 359.75, expectedIdx: 626*1440 + 0, gridLat: -66.5, gridLon: 359.75},

	// 重要经线
	{name: "Prime Meridian", lat: 0.0, lon: 0.0, expectedIdx: 360*1440 + 1439, gridLat: 0.0, gridLon: 0.0},
	{name: "180 Meridian", lat: 0.0, lon: 180.0, expectedIdx: 360*1440 + 719, gridLat: 0.0, gridLon: 180.0},
	{name: "90E Meridian", lat: 0.0, lon: 90.0, expectedIdx: 360*1440 + 1079, gridLat: 0.0, gridLon: 90.0},
	{name: "90W Meridian", lat: 0.0, lon: 270.0, expectedIdx: 360*1440 + 359, gridLat: 0.0, gridLon: 270.0},

	// 著名地点
	{name: "London", lat: 51.5, lon: -0.13, expectedIdx: 154 * 1440, gridLat: 51.5, gridLon: 359.75},
	{name: "New York", lat: 40.75, lon: -74.0, expectedIdx: 197*1440 + 295, gridLat: 40.75, gridLon: 286.0},
	{name: "Beijing", lat: 39.9, lon: 116.4, expectedIdx: 200*1440 + 973, gridLat: 40.0, gridLon: 116.5},
	{name: "Tokyo", lat: 35.7, lon: 139.7, expectedIdx: 217*1440 + 880, gridLat: 35.75, gridLon: 139.75},
	{name: "Sydney", lat: -33.9, lon: 151.2, expectedIdx: 496*1440 + 834, gridLat: -34.0, gridLon: 151.25},

	// Negative I scan mode 特有测试点
	{name: "Row Start High Lon", lat: 45.0, lon: 359.75, expectedIdx: 180*1440 + 0, gridLat: 45.0, gridLon: 359.75},
	{name: "Row End Low Lon", lat: 45.0, lon: 0.25, expectedIdx: 180*1440 + 1438, gridLat: 45.0, gridLon: 0.25},
	{name: "Cross 0 Longitude High", lat: 30.0, lon: 359.99, expectedIdx: 240*1440 + 1439, gridLat: 30.0, gridLon: 0.00},
	{name: "Cross 0 Longitude Low", lat: 30.0, lon: 0.01, expectedIdx: 240*1440 + 1439, gridLat: 30.0, gridLon: 0.0},
	{name: "Last Point in Row", lat: 15.0, lon: 0.0, expectedIdx: 300*1440 + 1439, gridLat: 15.0, gridLon: 0.0},
	{name: "Near Date Line", lat: 0.13, lon: 179.87, expectedIdx: 359*1440 + 720, gridLat: 0.25, gridLon: 179.75},
}

// 然后修改测试函数
func TestLatLon_NegativeIScanMode(t *testing.T) {
	// Convert GRIB units (millionths of a degree) to degrees
	firstLat := 90.0  // 90000000 millionths -> 90.0 degrees
	lastLat := -90.0  // -90000000 millionths -> -90.0 degrees
	firstLon := 0.0   // 0 millionths -> 0.0 degrees
	lastLon := 359.75 // 359750000 millionths -> 359.75 degrees
	latStep := 0.25   // 250000 millionths -> 0.25 degrees
	lonStep := 0.25   // 250000 millionths -> 0.25 degrees

	grid := latlon.NewLatLonGrid(
		lastLat,  // minLat
		firstLat, // maxLat
		firstLon, // minLon
		lastLon,  // maxLon
		latStep,
		lonStep,
	)

	runGridTests(t, grid, negativeIRegularLatLonGridTests, grids.ScanModeNegativeI)
}

var positiveJRegularLatLonGridTests = []gridTestCase{
	// 网格点精确匹配
	{name: "First Row Start", lat: -90.0, lon: 0.0, expectedIdx: 0, gridLat: -90.0, gridLon: 0.0},
	{name: "First Row End", lat: -90.0, lon: 359.75, expectedIdx: 1439, gridLat: -90.0, gridLon: 359.75},
	{name: "Second Row Start", lat: -89.75, lon: 0.0, expectedIdx: 1440, gridLat: -89.75, gridLon: 0.0},
	{name: "Second Row End", lat: -89.75, lon: 359.75, expectedIdx: 2879, gridLat: -89.75, gridLon: 359.75},

	// 纬度边界值测试（经度固定为0.0）
	{name: "Near -90 degrees", lat: -89.99, lon: 0.0, expectedIdx: 0, gridLat: -90.0, gridLon: 0.0},
	{name: "Close to -90", lat: -89.88, lon: 0.0, expectedIdx: 0, gridLat: -90.0, gridLon: 0.0},
	{name: "Middle Point", lat: -89.875, lon: 0.0, expectedIdx: 0, gridLat: -90.0, gridLon: 0.0},
	{name: "Near -89.75", lat: -89.87, lon: 0.0, expectedIdx: 1440, gridLat: -89.75, gridLon: 0.0},
	{name: "Very Close to -89.75", lat: -89.76, lon: 0.0, expectedIdx: 1440, gridLat: -89.75, gridLon: 0.0},
	{name: "Exact -89.75", lat: -89.75, lon: 0.0, expectedIdx: 1440, gridLat: -89.75, gridLon: 0.0},

	// 经度边界值测试（纬度固定为-90.0）
	{name: "0 Longitude", lat: -90.0, lon: 0.0, expectedIdx: 0, gridLat: -90.0, gridLon: 0.0},
	{name: "Very Near 0", lat: -90.0, lon: 0.12, expectedIdx: 0, gridLat: -90.0, gridLon: 0.0},
	{name: "Near 0.25", lat: -90.0, lon: 0.13, expectedIdx: 1, gridLat: -90.0, gridLon: 0.25},
	{name: "Middle Longitude", lat: -90.0, lon: 0.125, expectedIdx: 0, gridLat: -90.0, gridLon: 0.0},
	{name: "Close to 0.25", lat: -90.0, lon: 0.24, expectedIdx: 1, gridLat: -90.0, gridLon: 0.25},
	{name: "Exact 0.25", lat: -90.0, lon: 0.25, expectedIdx: 1, gridLat: -90.0, gridLon: 0.25},

	// 组合测试
	{name: "Combined Near -90", lat: -89.88, lon: 0.13, expectedIdx: 1, gridLat: -90.0, gridLon: 0.25},
	{name: "Combined Near -89.75", lat: -89.87, lon: 0.13, expectedIdx: 1441, gridLat: -89.75, gridLon: 0.25},
	{name: "Combined Middle", lat: -89.875, lon: 0.125, expectedIdx: 0, gridLat: -90.0, gridLon: 0.0},

	// 特殊地理位置
	{name: "First Point", lat: -90.0, lon: 0.0, expectedIdx: 0, gridLat: -90.0, gridLon: 0.0},
	{name: "Equator 180", lat: 0.0, lon: 180.0, expectedIdx: 360*1440 + 720, gridLat: 0.0, gridLon: 180.0},
	{name: "Equator Greenwich", lat: 0.0, lon: 0.0, expectedIdx: 360 * 1440, gridLat: 0.0, gridLon: 0.0},
	{name: "Date Line", lat: 0.0, lon: 180.0, expectedIdx: 360*1440 + 720, gridLat: 0.0, gridLon: 180.0},
	{name: "North Pole", lat: 90.0, lon: 0.0, expectedIdx: 720 * 1440, gridLat: 90.0, gridLon: 0.0},

	// 重要纬线
	{name: "Antarctic Circle", lat: -66.5, lon: 0.0, expectedIdx: 94 * 1440, gridLat: -66.5, gridLon: 0.0},
	{name: "Tropic of Capricorn", lat: -23.5, lon: 0.0, expectedIdx: 266 * 1440, gridLat: -23.5, gridLon: 0.0},
	{name: "Equator", lat: 0.0, lon: 0.0, expectedIdx: 360 * 1440, gridLat: 0.0, gridLon: 0.0},
	{name: "Tropic of Cancer", lat: 23.5, lon: 0.0, expectedIdx: 454 * 1440, gridLat: 23.5, gridLon: 0.0},
	{name: "Arctic Circle", lat: 66.5, lon: 0.0, expectedIdx: 626 * 1440, gridLat: 66.5, gridLon: 0.0},

	// 重要经线
	{name: "Prime Meridian", lat: 0.0, lon: 0.0, expectedIdx: 360 * 1440, gridLat: 0.0, gridLon: 0.0},
	{name: "180 Meridian", lat: 0.0, lon: 180.0, expectedIdx: 360*1440 + 720, gridLat: 0.0, gridLon: 180.0},
	{name: "90E Meridian", lat: 0.0, lon: 90.0, expectedIdx: 360*1440 + 360, gridLat: 0.0, gridLon: 90.0},
	{name: "90W Meridian", lat: 0.0, lon: -90.0, expectedIdx: 360*1440 + 1080, gridLat: 0.0, gridLon: 270.0},

	// 著名地点
	{name: "Sydney", lat: -33.9, lon: 151.2, expectedIdx: 224*1440 + 605, gridLat: -34.0, gridLon: 151.25},
	{name: "Tokyo", lat: 35.7, lon: 139.7, expectedIdx: 503*1440 + 559, gridLat: 35.75, gridLon: 139.75},
	{name: "Beijing", lat: 39.9, lon: 116.4, expectedIdx: 520*1440 + 466, gridLat: 40.0, gridLon: 116.5},
	{name: "London", lat: 51.5, lon: -0.13, expectedIdx: 566*1440 + 1439, gridLat: 51.5, gridLon: 359.75},
	{name: "New York", lat: 40.75, lon: -74.0, expectedIdx: 523*1440 + 1144, gridLat: 40.75, gridLon: 286.0},

	// Positive J scan mode 特有测试点
	{name: "First Row South", lat: -90.0, lon: 0.25, expectedIdx: 1, gridLat: -90.0, gridLon: 0.25},
	{name: "Cross Equator South", lat: -0.01, lon: 0.0, expectedIdx: 360 * 1440, gridLat: 0.0, gridLon: 0.0},
	{name: "Cross Equator North", lat: 0.01, lon: 0.0, expectedIdx: 360 * 1440, gridLat: 0.0, gridLon: 0.0},
	{name: "Near Last Row", lat: 89.87, lon: 0.0, expectedIdx: 719 * 1440, gridLat: 90.0, gridLon: 0.0},
	{name: "Last Row North", lat: 90.0, lon: 0.25, expectedIdx: 720*1440 + 1, gridLat: 90.0, gridLon: 0.25},
	{name: "Near Date Line", lat: 0.13, lon: 179.87, expectedIdx: 361*1440 + 719, gridLat: 0.25, gridLon: 179.75},
}

func TestLatLon_PositiveJScanMode(t *testing.T) {
	// Convert GRIB units (millionths of a degree) to degrees
	firstLat := 90.0  // 90000000 millionths -> 90.0 degrees
	lastLat := -90.0  // -90000000 millionths -> -90.0 degrees
	firstLon := 0.0   // 0 millionths -> 0.0 degrees
	lastLon := 359.75 // 359750000 millionths -> 359.75 degrees
	latStep := 0.25   // 250000 millionths -> 0.25 degrees
	lonStep := 0.25   // 250000 millionths -> 0.25 degrees

	grid := latlon.NewLatLonGrid(
		lastLat,  // minLat
		firstLat, // maxLat
		firstLon, // minLon
		lastLon,  // maxLon
		latStep,
		lonStep,
	)

	runGridTests(t, grid, positiveJRegularLatLonGridTests, grids.ScanModePositiveJ)
}

var consecutiveJRegularLatLonGridTests = []gridTestCase{
	// 网格点精确匹配
	{name: "First Column Start", lat: 90.0, lon: 0.0, expectedIdx: 0, gridLat: 90.0, gridLon: 0.0},
	{name: "First Column End", lat: -90.0, lon: 0.0, expectedIdx: 720, gridLat: -90.0, gridLon: 0.0},
	{name: "Second Column Start", lat: 90.0, lon: 0.25, expectedIdx: 721, gridLat: 90.0, gridLon: 0.25},
	{name: "Second Column End", lat: -90.0, lon: 0.25, expectedIdx: 721 + 720, gridLat: -90.0, gridLon: 0.25},
	{name: "Third Column Start", lat: 90.0, lon: 0.50, expectedIdx: 721 * 2, gridLat: 90.0, gridLon: 0.50},
	{name: "Third Column End", lat: -90.0, lon: 0.50, expectedIdx: 721*2 + 720, gridLat: -90.0, gridLon: 0.50},

	// 纬度边界值测试（经度固定为0.0）
	{name: "Near 90 degrees", lat: 89.99, lon: 0.0, expectedIdx: 0, gridLat: 90.0, gridLon: 0.0},
	{name: "Close to 90", lat: 89.88, lon: 0.0, expectedIdx: 0, gridLat: 90.0, gridLon: 0.0},
	{name: "Middle Point", lat: 89.875, lon: 0.0, expectedIdx: 0, gridLat: 90.0, gridLon: 0.0},
	{name: "Near 89.75", lat: 89.87, lon: 0.0, expectedIdx: 1, gridLat: 89.75, gridLon: 0.0},
	{name: "Very Close to 89.75", lat: 89.76, lon: 0.0, expectedIdx: 1, gridLat: 89.75, gridLon: 0.0},
	{name: "Exact 89.75", lat: 89.75, lon: 0.0, expectedIdx: 1, gridLat: 89.75, gridLon: 0.0},

	// 经度边界值测试（纬度固定为90.0）
	{name: "0 Longitude", lat: 90.0, lon: 0.0, expectedIdx: 0, gridLat: 90.0, gridLon: 0.0},
	{name: "Very Near 0", lat: 90.0, lon: 0.12, expectedIdx: 0, gridLat: 90.0, gridLon: 0.0},
	{name: "Near 0.25", lat: 90.0, lon: 0.13, expectedIdx: 721, gridLat: 90.0, gridLon: 0.25},
	{name: "Middle Longitude", lat: 90.0, lon: 0.125, expectedIdx: 0, gridLat: 90.0, gridLon: 0.0},
	{name: "Close to 0.25", lat: 90.0, lon: 0.24, expectedIdx: 721, gridLat: 90.0, gridLon: 0.25},
	{name: "Exact 0.25", lat: 90.0, lon: 0.25, expectedIdx: 721, gridLat: 90.0, gridLon: 0.25},

	// 组合测试
	{name: "Combined Near 90", lat: 89.88, lon: 0.13, expectedIdx: 721, gridLat: 90.0, gridLon: 0.25},
	{name: "Combined Near 89.75", lat: 89.87, lon: 0.13, expectedIdx: 722, gridLat: 89.75, gridLon: 0.25},
	{name: "Combined Middle", lat: 89.875, lon: 0.125, expectedIdx: 0, gridLat: 90.0, gridLon: 0.0},

	// 特殊地理位置
	{name: "Last Point", lat: -90.0, lon: 359.75, expectedIdx: 1439*721 + 720, gridLat: -90.0, gridLon: 359.75},
	{name: "Equator 180", lat: 0.0, lon: 180.0, expectedIdx: 720*721 + 360, gridLat: 0.0, gridLon: 180.0},
	{name: "Equator Greenwich", lat: 0.0, lon: 0.0, expectedIdx: 360, gridLat: 0.0, gridLon: 0.0},
	{name: "Date Line", lat: 0.0, lon: 180.0, expectedIdx: 720*721 + 360, gridLat: 0.0, gridLon: 180.0},
	{name: "South Pole", lat: -90.0, lon: 0.0, expectedIdx: 720, gridLat: -90.0, gridLon: 0.0},

	// 重要纬线
	{name: "Arctic Circle", lat: 66.5, lon: 0.0, expectedIdx: 94, gridLat: 66.5, gridLon: 0.0},
	{name: "Tropic of Cancer", lat: 23.5, lon: 0.0, expectedIdx: 266, gridLat: 23.5, gridLon: 0.0},
	{name: "Equator", lat: 0.0, lon: 0.0, expectedIdx: 360, gridLat: 0.0, gridLon: 0.0},
	{name: "Tropic of Capricorn", lat: -23.5, lon: 0.0, expectedIdx: 454, gridLat: -23.5, gridLon: 0.0},
	{name: "Antarctic Circle", lat: -66.5, lon: 0.0, expectedIdx: 626, gridLat: -66.5, gridLon: 0.0},

	// 重要经线
	{name: "Prime Meridian", lat: 0.0, lon: 0.0, expectedIdx: 360, gridLat: 0.0, gridLon: 0.0},
	{name: "180 Meridian", lat: 0.0, lon: 180.0, expectedIdx: 720*721 + 360, gridLat: 0.0, gridLon: 180.0},
	{name: "90E Meridian", lat: 0.0, lon: 90.0, expectedIdx: 360*721 + 360, gridLat: 0.0, gridLon: 90.0},
	{name: "90W Meridian", lat: 0.0, lon: -90.0, expectedIdx: 1080*721 + 360, gridLat: 0.0, gridLon: 270.0},

	// 著名地点
	{name: "London", lat: 51.5, lon: -0.13, expectedIdx: 1439*721 + 154, gridLat: 51.5, gridLon: 359.75},
	{name: "New York", lat: 40.75, lon: -74.0, expectedIdx: 1144*721 + 197, gridLat: 40.75, gridLon: 286.0},
	{name: "Beijing", lat: 39.9, lon: 116.4, expectedIdx: 466*721 + 200, gridLat: 40.0, gridLon: 116.5},
	{name: "Tokyo", lat: 35.7, lon: 139.7, expectedIdx: 559*721 + 217, gridLat: 35.75, gridLon: 139.75},
	{name: "Sydney", lat: -33.9, lon: 151.2, expectedIdx: 605*721 + 496, gridLat: -34.0, gridLon: 151.25},

	// 边界情况
	{name: "Near North Pole", lat: 89.99, lon: 0.13, expectedIdx: 721, gridLat: 90.0, gridLon: 0.25},
	{name: "Near South Pole", lat: -89.99, lon: 359.87, expectedIdx: 1439*721 + 720, gridLat: -90.0, gridLon: 359.75},
	{name: "Near Date Line", lat: 0.13, lon: 179.87, expectedIdx: 719*721 + 359, gridLat: 0.0, gridLon: 179.75},
	{name: "Near Prime Meridian", lat: 0.13, lon: -0.13, expectedIdx: 1439*721 + 359, gridLat: 0.25, gridLon: 359.75},

	// Consecutive J scan mode 特有测试点
	{name: "First Column Second Point", lat: 89.75, lon: 0.0, expectedIdx: 1, gridLat: 89.75, gridLon: 0.0},
	{name: "Second Column Top", lat: 90.0, lon: 0.25, expectedIdx: 721, gridLat: 90.0, gridLon: 0.25},
	{name: "Column Boundary Low", lat: 89.87, lon: 0.24, expectedIdx: 722, gridLat: 89.75, gridLon: 0.25},
	{name: "Column Boundary High", lat: 89.87, lon: 0.26, expectedIdx: 722, gridLat: 89.75, gridLon: 0.25},
	{name: "Last Column First Point", lat: 90.0, lon: 359.75, expectedIdx: 1439 * 721, gridLat: 90.0, gridLon: 359.75},
	{name: "Last Column Last Point", lat: -90.0, lon: 359.75, expectedIdx: 1439*721 + 720, gridLat: -90.0, gridLon: 359.75},
}

func TestLatLon_ConsecutiveJScanMode(t *testing.T) {
	// Convert GRIB units (millionths of a degree) to degrees
	firstLat := 90.0  // 90000000 millionths -> 90.0 degrees
	lastLat := -90.0  // -90000000 millionths -> -90.0 degrees
	firstLon := 0.0   // 0 millionths -> 0.0 degrees
	lastLon := 359.75 // 359750000 millionths -> 359.75 degrees
	latStep := 0.25   // 250000 millionths -> 0.25 degrees
	lonStep := 0.25   // 250000 millionths -> 0.25 degrees

	grid := latlon.NewLatLonGrid(
		lastLat,  // minLat
		firstLat, // maxLat
		firstLon, // minLon
		lastLon,  // maxLon
		latStep,
		lonStep,
	)

	runGridTests(t, grid, consecutiveJRegularLatLonGridTests, grids.ScanModeConsecutiveJ)
}

var oppositeRowsRegularLatLonGridTests = []gridTestCase{
	// 网格点精确匹配
	{name: "First Row Start", lat: 90.0, lon: 0.0, expectedIdx: 0, gridLat: 90.0, gridLon: 0.0},
	{name: "First Row End", lat: 90.0, lon: 359.75, expectedIdx: 1439, gridLat: 90.0, gridLon: 359.75},
	{name: "Second Row Start", lat: 89.75, lon: 359.75, expectedIdx: 1440, gridLat: 89.75, gridLon: 359.75},
	{name: "Second Row End", lat: 89.75, lon: 0.0, expectedIdx: 2879, gridLat: 89.75, gridLon: 0.0},
	{name: "Third Row Start", lat: 89.5, lon: 0.0, expectedIdx: 2880, gridLat: 89.5, gridLon: 0.0},
	{name: "Third Row End", lat: 89.5, lon: 359.75, expectedIdx: 4319, gridLat: 89.5, gridLon: 359.75},

	// 纬度边界值测试（经度固定为0.0）
	{name: "Near 90 degrees", lat: 89.99, lon: 0.0, expectedIdx: 0, gridLat: 90.0, gridLon: 0.0},
	{name: "Close to 90", lat: 89.88, lon: 0.0, expectedIdx: 0, gridLat: 90.0, gridLon: 0.0},
	{name: "Middle Point", lat: 89.875, lon: 0.0, expectedIdx: 0, gridLat: 90.0, gridLon: 0.0},
	{name: "Near 89.75", lat: 89.87, lon: 359.75, expectedIdx: 1440, gridLat: 89.75, gridLon: 359.75},
	{name: "Very Close to 89.75", lat: 89.76, lon: 359.75, expectedIdx: 1440, gridLat: 89.75, gridLon: 359.75},
	{name: "Exact 89.75", lat: 89.75, lon: 359.75, expectedIdx: 1440, gridLat: 89.75, gridLon: 359.75},

	// 经度边界值测试（纬度固定为90.0）
	{name: "0 Longitude", lat: 90.0, lon: 0.0, expectedIdx: 0, gridLat: 90.0, gridLon: 0.0},
	{name: "Very Near 0", lat: 90.0, lon: 0.12, expectedIdx: 0, gridLat: 90.0, gridLon: 0.0},
	{name: "Near 0.25", lat: 90.0, lon: 0.13, expectedIdx: 1, gridLat: 90.0, gridLon: 0.25},
	{name: "Middle Longitude", lat: 90.0, lon: 0.125, expectedIdx: 0, gridLat: 90.0, gridLon: 0.0},
	{name: "Close to 0.25", lat: 90.0, lon: 0.24, expectedIdx: 1, gridLat: 90.0, gridLon: 0.25},
	{name: "Exact 0.25", lat: 90.0, lon: 0.25, expectedIdx: 1, gridLat: 90.0, gridLon: 0.25},

	// 组合测试
	{name: "Combined Near 90", lat: 89.88, lon: 0.13, expectedIdx: 1, gridLat: 90.0, gridLon: 0.25},
	{name: "Combined Near 89.75", lat: 89.87, lon: 359.62, expectedIdx: 1441, gridLat: 89.75, gridLon: 359.5},
	{name: "Combined Middle", lat: 89.875, lon: 0.125, expectedIdx: 0, gridLat: 90.0, gridLon: 0.0},

	// 特殊地理位置
	{name: "Last Point", lat: -90.0, lon: 0.0, expectedIdx: (90+90)/0.25*1440 + (1440 - 1), gridLat: -90.0, gridLon: 0.0},
	{name: "Equator 180", lat: 0.0, lon: 180.0, expectedIdx: (90-0)/0.25*1440 + 720, gridLat: 0.0, gridLon: 180.0},
	{name: "Equator Greenwich", lat: 0.0, lon: 0.0, expectedIdx: (90-0)/0.25*1440 + 0, gridLat: 0.0, gridLon: 0.0},
	{name: "Date Line", lat: 0.0, lon: 180.0, expectedIdx: (90-0)/0.25*1440 + 720, gridLat: 0.0, gridLon: 180.0},
	{name: "South Pole", lat: -90.0, lon: 0.0, expectedIdx: (90+90)/0.25*1440 + (1440 - 1), gridLat: -90.0, gridLon: 0.0},

	// 重要纬线
	{name: "Arctic Circle", lat: 66.5, lon: 359.75, expectedIdx: (90-66.5)/0.25*1440 + 1439, gridLat: 66.5, gridLon: 359.75},
	{name: "Tropic of Cancer", lat: 23.5, lon: 0.0, expectedIdx: (90 - 23.5) / 0.25 * 1440, gridLat: 23.5, gridLon: 0.0},
	{name: "Equator", lat: 0.0, lon: 359.75, expectedIdx: (90-0)/0.25*1440 + 1439, gridLat: 0.0, gridLon: 359.75},
	{name: "Tropic of Capricorn", lat: -23.5, lon: 359.75, expectedIdx: (90+23.5)/0.25*1440 + 1439, gridLat: -23.5, gridLon: 359.75},
	{name: "Antarctic Circle", lat: -66.5, lon: 0.0, expectedIdx: (90 + 66.5) / 0.25 * 1440, gridLat: -66.5, gridLon: 0.0},

	// 重要经线
	{name: "Prime Meridian", lat: 0.0, lon: 359.75, expectedIdx: (90-0)/0.25*1440 + 1439, gridLat: 0.0, gridLon: 359.75},
	{name: "180 Meridian", lat: 0.0, lon: 180.0, expectedIdx: (90-0)/0.25*1440 + 720, gridLat: 0.0, gridLon: 180.0},
	{name: "90E Meridian", lat: 0.0, lon: 90.0, expectedIdx: (90-0)/0.25*1440 + 360, gridLat: 0.0, gridLon: 90.0},
	{name: "90W Meridian", lat: 0.0, lon: 270.0, expectedIdx: (90-0)/0.25*1440 + 1080, gridLat: 0.0, gridLon: 270.0},

	// 著名地点
	{name: "London", lat: 51.5, lon: -0.13, expectedIdx: 154*1440 + 1439, gridLat: 51.5, gridLon: 359.75},
	{name: "New York", lat: 40.75, lon: -74.0, expectedIdx: 197*1440 + 295, gridLat: 40.75, gridLon: 286.0},
	{name: "Beijing", lat: 39.9, lon: 116.4, expectedIdx: 200*1440 + 466, gridLat: 40.0, gridLon: 116.5},
	{name: "Tokyo", lat: 35.7, lon: 139.7, expectedIdx: 217*1440 + 880, gridLat: 35.75, gridLon: 139.75},
	{name: "Sydney", lat: -33.9, lon: 151.2, expectedIdx: 496*1440 + 605, gridLat: -34.0, gridLon: 151.25},

	// 边界情况
	{name: "Near North Pole Low Lon", lat: 89.99, lon: 0.12, expectedIdx: 1439, gridLat: 90.0, gridLon: 0.0},
	{name: "Near North Pole High Lon", lat: 89.99, lon: 0.13, expectedIdx: 1438, gridLat: 90.0, gridLon: 0.25},
	{name: "Near South Pole High Lon", lat: -89.99, lon: 359.88, expectedIdx: (90+90)/0.25*1440 + 1439, gridLat: -90.0, gridLon: 0.0},
	{name: "Near South Pole Low Lon", lat: -89.99, lon: 359.87, expectedIdx: (90+90)/0.25*1440 + 0, gridLat: -90.0, gridLon: 359.75},
	{name: "Near Date Line", lat: 0.13, lon: 179.87, expectedIdx: (90-0.25)/0.25*1440 + (1440 - 720), gridLat: 0.25, gridLon: 179.75},
	{name: "Near Prime Meridian", lat: 0.13, lon: -0.13, expectedIdx: (90 - 0.25) / 0.25 * 1440, gridLat: 0.25, gridLon: 359.75},

	// Opposite rows scan mode 特有测试点
	{name: "First Row Forward Start", lat: 90.0, lon: 0.0, expectedIdx: 0, gridLat: 90.0, gridLon: 0.0},
	{name: "First Row Forward End", lat: 90.0, lon: 359.75, expectedIdx: 1439, gridLat: 90.0, gridLon: 359.75},
	{name: "Second Row Reverse Start", lat: 89.75, lon: 359.75, expectedIdx: 1440, gridLat: 89.75, gridLon: 359.75},
	{name: "Second Row Reverse End", lat: 89.75, lon: 0.0, expectedIdx: 2879, gridLat: 89.75, gridLon: 0.0},
	{name: "Third Row Forward Start", lat: 89.5, lon: 0.0, expectedIdx: 2880, gridLat: 89.5, gridLon: 0.0},
	{name: "Third Row Forward End", lat: 89.5, lon: 359.75, expectedIdx: 4319, gridLat: 89.5, gridLon: 359.75},
	{name: "Row Direction Change Point", lat: 89.625, lon: 180.0, expectedIdx: 2160, gridLat: 89.75, gridLon: 180.0},
}

func TestLatLon_OppositeRowsScanMode(t *testing.T) {
	// Convert GRIB units (millionths of a degree) to degrees
	firstLat := 90.0  // 90000000 millionths -> 90.0 degrees
	lastLat := -90.0  // -90000000 millionths -> -90.0 degrees
	firstLon := 0.0   // 0 millionths -> 0.0 degrees
	lastLon := 359.75 // 359750000 millionths -> 359.75 degrees
	latStep := 0.25   // 250000 millionths -> 0.25 degrees
	lonStep := 0.25   // 250000 millionths -> 0.25 degrees

	grid := latlon.NewLatLonGrid(
		lastLat,  // minLat
		firstLat, // maxLat
		firstLon, // minLon
		lastLon,  // maxLon
		latStep,
		lonStep,
	)

	runGridTests(t, grid, oppositeRowsRegularLatLonGridTests, grids.ScanModeOppositeRows)
}

var defaultRegularGaussianGridTests = []gridTestCase{
	// 网格点精确匹配
	{name: "First Row Start", lat: 88.572169, lon: 0.0, expectedIdx: 0, gridLat: 88.572169, gridLon: 0.0},
	{name: "First Row Second Point", lat: 88.572169, lon: 1.875, expectedIdx: 1, gridLat: 88.572169, gridLon: 1.875},
	{name: "First Row End", lat: 88.572169, lon: 358.125, expectedIdx: 191, gridLat: 88.572169, gridLon: 358.125},
	{name: "Second Row Start", lat: 86.722531, lon: 0.0, expectedIdx: 192, gridLat: 86.722531, gridLon: 0.0},
	{name: "Second Row End", lat: 86.722531, lon: 358.125, expectedIdx: 383, gridLat: 86.722531, gridLon: 358.125},
	{name: "Last Row Start", lat: -88.572169, lon: 0.0, expectedIdx: 95*192 + 0, gridLat: -88.572169, gridLon: 0.0},
	{name: "Last Row End", lat: -88.572169, lon: 358.125, expectedIdx: 95*192 + 191, gridLat: -88.572169, gridLon: 358.125},
}

func TestRegularGaussian_DefaultScanMode(t *testing.T) {
	grid := gaussian.NewRegular(48)

	runGridTests(t, grid, defaultRegularGaussianGridTests, 0)
}

var negativeIRegularGaussianGridTests = []gridTestCase{
	// 网格点精确匹配
	{name: "First Row Start", lat: 88.572169, lon: 358.125, expectedIdx: 0, gridLat: 88.572169, gridLon: 358.125},
	{name: "First Row Second Point", lat: 88.572169, lon: 356.25, expectedIdx: 1, gridLat: 88.572169, gridLon: 356.25},
	{name: "First Row End", lat: 88.572169, lon: 0, expectedIdx: 191, gridLat: 88.572169, gridLon: 0},
	{name: "Second Row Start", lat: 86.722531, lon: 358.125, expectedIdx: 192, gridLat: 86.722531, gridLon: 358.125},
	{name: "Second Row End", lat: 86.722531, lon: 0, expectedIdx: 383, gridLat: 86.722531, gridLon: 0},
	{name: "Last Row Start", lat: -88.572169, lon: 358.125, expectedIdx: 95*192 + 0, gridLat: -88.572169, gridLon: 358.125},
	{name: "Last Row End", lat: -88.572169, lon: 0, expectedIdx: 95*192 + 191, gridLat: -88.572169, gridLon: 0},
}

func TestRegularGaussian_NegativeIScanMode(t *testing.T) {
	grid := gaussian.NewRegular(48)

	runGridTests(t, grid, negativeIRegularGaussianGridTests, grids.ScanModeNegativeI)
}

var positiveJRegularGaussianGridTests = []gridTestCase{
	// 网格点精确匹配
	{name: "First Row Start", lat: -88.572169, lon: 0.0, expectedIdx: 0, gridLat: -88.572169, gridLon: 0.0},
	{name: "First Row Second Point", lat: -88.572169, lon: 1.875, expectedIdx: 1, gridLat: -88.572169, gridLon: 1.875},
	{name: "First Row End", lat: -88.572169, lon: 358.125, expectedIdx: 191, gridLat: -88.572169, gridLon: 358.125},
	{name: "Second Row Start", lat: -86.722531, lon: 0.0, expectedIdx: 192, gridLat: -86.722531, gridLon: 0.0},
	{name: "Second Row End", lat: -86.722531, lon: 358.125, expectedIdx: 383, gridLat: -86.722531, gridLon: 358.125},
	{name: "Last Row Start", lat: 88.572169, lon: 0.0, expectedIdx: 95*192 + 0, gridLat: 88.572169, gridLon: 0.0},
	{name: "Last Row End", lat: 88.572169, lon: 358.125, expectedIdx: 95*192 + 191, gridLat: 88.572169, gridLon: 358.125},
}

func TestRegularGaussian_PositiveJScanMode(t *testing.T) {
	grid := gaussian.NewRegular(48)

	runGridTests(t, grid, positiveJRegularGaussianGridTests, grids.ScanModePositiveJ)
}

var consecutiveJRegularGaussianGridTests = []gridTestCase{
	// 网格点精确匹配
	{name: "First Column Start", lat: 88.572169, lon: 0.0, expectedIdx: 0, gridLat: 88.572169, gridLon: 0.0},
	{name: "First Column Second Point", lat: 86.722531, lon: 0.0, expectedIdx: 1, gridLat: 86.722531, gridLon: 0.0},
	{name: "First Column End", lat: -88.572169, lon: 0.0, expectedIdx: 95, gridLat: -88.572169, gridLon: 0.0},
	{name: "Second Column Start", lat: 88.572169, lon: 1.875, expectedIdx: 96, gridLat: 88.572169, gridLon: 1.875},
	{name: "Second Column End", lat: -88.572169, lon: 1.875, expectedIdx: 96 + 95, gridLat: -88.572169, gridLon: 1.875},
	{name: "Last Column Start", lat: 88.572169, lon: 358.125, expectedIdx: 96 * 191, gridLat: 88.572169, gridLon: 358.125},
	{name: "Last Column End", lat: -88.572169, lon: 358.125, expectedIdx: 96*192 - 1, gridLat: -88.572169, gridLon: 358.125},
}

func TestRegularGaussian_ConsecutiveJScanMode(t *testing.T) {
	grid := gaussian.NewRegular(48)

	runGridTests(t, grid, consecutiveJRegularGaussianGridTests, grids.ScanModeConsecutiveJ)
}

var oppositeRowsRegularGaussianGridTests = []gridTestCase{
	// 网格点精确匹配
	{name: "First Row Start", lat: 88.572169, lon: 0.0, expectedIdx: 0, gridLat: 88.572169, gridLon: 0.0},
	{name: "First Row Second Point", lat: 88.572169, lon: 1.875, expectedIdx: 1, gridLat: 88.572169, gridLon: 1.875},
	{name: "First Row End", lat: 88.572169, lon: 358.125, expectedIdx: 191, gridLat: 88.572169, gridLon: 358.125},
	{name: "Second Row Start", lat: 86.722531, lon: 358.125, expectedIdx: 192, gridLat: 86.722531, gridLon: 358.125},
	{name: "Second Row End", lat: 86.722531, lon: 0.0, expectedIdx: 383, gridLat: 86.722531, gridLon: 0.0},
	{name: "Last Row Start", lat: -88.572169, lon: 358.125, expectedIdx: 95*192 + 0, gridLat: -88.572169, gridLon: 358.125},
	{name: "Last Row End", lat: -88.572169, lon: 0.0, expectedIdx: 95*192 + 191, gridLat: -88.572169, gridLon: 0.0},
}

func TestRegularGaussian_OppositeRowsScanMode(t *testing.T) {
	grid := gaussian.NewRegular(48)

	runGridTests(t, grid, oppositeRowsRegularGaussianGridTests, grids.ScanModeOppositeRows)
}

var consecutiveJOppositeRowsRegularGaussianGridTests = []gridTestCase{
	// 网格点精确匹配
	{name: "First Column Start", lat: 88.572169, lon: 0.0, expectedIdx: 0, gridLat: 88.572169, gridLon: 0.0},
	{name: "First Column Second Point", lat: 86.722531, lon: 0.0, expectedIdx: 1, gridLat: 86.722531, gridLon: 0.0},
	{name: "First Column End", lat: -88.572169, lon: 0.0, expectedIdx: 95, gridLat: -88.572169, gridLon: 0.0},
	{name: "Second Column Start", lat: -88.572169, lon: 1.875, expectedIdx: 96, gridLat: -88.572169, gridLon: 1.875},
	{name: "Second Column End", lat: 88.572169, lon: 1.875, expectedIdx: 96 + 95, gridLat: 88.572169, gridLon: 1.875},
	{name: "Last Column Start", lat: -88.572169, lon: 358.125, expectedIdx: 96 * 191, gridLat: -88.572169, gridLon: 358.125},
	{name: "Last Column End", lat: 88.572169, lon: 358.125, expectedIdx: 96*192 - 1, gridLat: 88.572169, gridLon: 358.125},
}

func TestRegularGaussian_ConsecutiveJOppositeRowsScanMode(t *testing.T) {
	grid := gaussian.NewRegular(48)

	runGridTests(t, grid, consecutiveJOppositeRowsRegularGaussianGridTests, grids.ScanModeConsecutiveJ|grids.ScanModeOppositeRows)
}

func BenchmarkGridIndex(b *testing.B) {
	b.Run("LatLon 0p25", func(b *testing.B) {
		grid := latlon.NewLatLonGrid(-90, 90, 0.0, 359.75, 0.25, 0.25)

		for i := 0; i < b.N; i++ {
			grids.GridIndex(grid, 88.572169, 0.0, 0)
		}
	})

	b.Run("LatLon 0p16", func(b *testing.B) {
		grid := latlon.NewLatLonGrid(-90, 90, 0.0, 359.84, 0.16, 0.16)

		for i := 0; i < b.N; i++ {
			grids.GridIndex(grid, 88.572169, 0.0, 0)
		}
	})

	b.Run("Gaussian F48", func(b *testing.B) {
		grid := gaussian.NewRegular(48)

		for i := 0; i < b.N; i++ {
			grids.GridIndex(grid, 88.572169, 0.0, 0)
		}
	})

	b.Run("Gaussian F768", func(b *testing.B) {
		grid := gaussian.NewRegular(768)

		for i := 0; i < b.N; i++ {
			grids.GridIndex(grid, 88.572169, 0.0, 0)
		}
	})
}

func BenchmarkGridGuessIndex(b *testing.B) {
	b.Run("LatLon 0p25", func(b *testing.B) {
		grid := latlon.NewLatLonGrid(-90, 90, 0.0, 359.75, 0.25, 0.25)

		for i := 0; i < b.N; i++ {
			grids.GuessGridIndex(grid, 88.572169, 0.0, 0)
		}
	})

	b.Run("LatLon 0p16", func(b *testing.B) {
		grid := latlon.NewLatLonGrid(-90, 90, 0.0, 359.84, 0.16, 0.16)

		for i := 0; i < b.N; i++ {
			grids.GuessGridIndex(grid, 88.572169, 0.0, 0)
		}
	})

	b.Run("Gaussian F48", func(b *testing.B) {
		grid := gaussian.NewRegular(48)

		for i := 0; i < b.N; i++ {
			grids.GuessGridIndex(grid, 88.572169, 0.0, 0)
		}
	})

	b.Run("Gaussian F768", func(b *testing.B) {
		grid := gaussian.NewRegular(768)

		for i := 0; i < b.N; i++ {
			grids.GuessGridIndex(grid, 88.572169, 0.0, 0)
		}
	})
}

func BenchmarkGridPoint(b *testing.B) {
	b.Run("LatLon 0p25", func(b *testing.B) {
		grid := latlon.NewLatLonGrid(-90, 90, 0.0, 359.75, 0.25, 0.25)

		for i := 0; i < b.N; i++ {
			grids.GridPoint(grid, i%grid.Size(), 0)
		}
	})

	b.Run("LatLon 0p16", func(b *testing.B) {
		grid := latlon.NewLatLonGrid(-90, 90, 0.0, 359.84, 0.16, 0.16)

		for i := 0; i < b.N; i++ {
			grids.GridPoint(grid, i%grid.Size(), 0)
		}
	})

	b.Run("Gaussian F48", func(b *testing.B) {
		grid := gaussian.NewRegular(48)

		for i := 0; i < b.N; i++ {
			grids.GridPoint(grid, i%grid.Size(), 0)
		}
	})

	b.Run("Gaussian F768", func(b *testing.B) {
		grid := gaussian.NewRegular(768)

		for i := 0; i < b.N; i++ {
			grids.GridPoint(grid, i%grid.Size(), 0)
		}
	})
}

var octahedralGaussianGridTests = []gridTestCase{
	// 网格点精确匹配
	{name: "North Pole", lat: 90.0, lon: 0.0, expectedIdx: 0, gridLat: 90.0, gridLon: 0.0},
	{name: "Near North Pole", lat: 89.99, lon: 0.0, expectedIdx: 0, gridLat: 90.0, gridLon: 0.0},

	// 赤道测试点（经度点数最多）
	{name: "Equator 0", lat: 0.0, lon: 0.0, expectedIdx: 9024, gridLat: 0.933, gridLon: 0.0},
	{name: "Equator 90", lat: 0.0, lon: 90.0, expectedIdx: 9072, gridLat: 0.933, gridLon: 90.0},
	{name: "Equator 180", lat: 0.0, lon: 180.0, expectedIdx: 9120, gridLat: 0.933, gridLon: 180.0},
	{name: "Equator 270", lat: 0.0, lon: 270.0, expectedIdx: 9168, gridLat: 0.933, gridLon: 270.0},

	// 极点附近测试（经度点数最少）
	{name: "Near North Pole 90", lat: 89.5, lon: 90.0, expectedIdx: 48, gridLat: 88.572, gridLon: 90.0},
	{name: "Near South Pole 90", lat: -89.5, lon: 90.0, expectedIdx: 18049, gridLat: -88.572, gridLon: 90.0},

	// 中间纬度测试点
	{name: "Mid Lat North", lat: 45.0, lon: 0.0, expectedIdx: 4416, gridLat: 45.699, gridLon: 0.0},
	{name: "Mid Lat South", lat: -45.0, lon: 0.0, expectedIdx: 13824, gridLat: -44.301, gridLon: 0.0},
	{name: "Near South Pole 90", lat: -89.5, lon: 90.0, expectedIdx: 18144, gridLat: -88.572, gridLon: 90.0},

	// 经度边界测试
	{name: "Date Line Positive", lat: 0.0, lon: 180.0, expectedIdx: 9120, gridLat: 0.933, gridLon: 180.0},
	{name: "Date Line Negative", lat: 0.0, lon: -180.0, expectedIdx: 9120, gridLat: 0.933, gridLon: 180.0},

	// 特殊纬线
	{name: "Arctic Circle", lat: 66.5, lon: 0.0, expectedIdx: 2304, gridLat: 66.5, gridLon: 0.0},
	{name: "Tropic of Cancer", lat: 23.5, lon: 0.0, expectedIdx: 6720, gridLat: 23.315, gridLon: 0.0},
	{name: "Tropic of Capricorn", lat: -23.5, lon: 0.0, expectedIdx: 11712, gridLat: -23.315, gridLon: 0.0},
	{name: "Antarctic Circle", lat: -66.5, lon: 0.0, expectedIdx: 15936, gridLat: -66.5, gridLon: 0.0},
}

// func TestOctahedralGaussian(t *testing.T) {
// 	grid := gaussian.NewOctahedral(48)
// 	runGridTests(t, grid, octahedralGaussianGridTests, 0)
// }
