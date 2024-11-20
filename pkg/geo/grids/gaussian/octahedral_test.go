package gaussian_test

import (
	"fmt"
	"testing"

	"github.com/scorix/walg/pkg/geo/grids/gaussian"
	"github.com/stretchr/testify/assert"
)

func TestOctahedral_O32(t *testing.T) {
	g := gaussian.NewOctahedral(32)

	t.Logf("latitudes: %v", g.Latitudes())
	t.Logf("lon points: %v", g.Longitudes())

	// Test basic grid properties
	assert.Equal(t, 64, len(g.Latitudes()))
	assert.Equal(t, 128, len(g.Longitudes())) // 4*N at equator
	assert.Equal(t, 4*10*64+128*(64-20)*64+4*10*64, g.Size())

	assert.Equal(t, 87.86379883923267, g.Latitudes()[0])
	assert.Equal(t, 85.09652698831732, g.Latitudes()[1])
	assert.Equal(t, 82.31291294788628, g.Latitudes()[2])
	assert.Equal(t, 79.52560657265944, g.Latitudes()[3])
	assert.Equal(t, 76.73689968036832, g.Latitudes()[4])
	assert.Equal(t, 73.94751515398967, g.Latitudes()[5])
	assert.Equal(t, 71.15775201158732, g.Latitudes()[6])
	assert.Equal(t, 68.36775610831316, g.Latitudes()[7])
	assert.Equal(t, 65.57760701082782, g.Latitudes()[8])
	assert.Equal(t, 62.787351798963066, g.Latitudes()[9])
	assert.Equal(t, 59.99702010849129, g.Latitudes()[10])
	assert.Equal(t, 57.20663152764325, g.Latitudes()[11])
	assert.Equal(t, 54.4161995260862, g.Latitudes()[12])
	assert.Equal(t, 51.625733674938246, g.Latitudes()[13])
	assert.Equal(t, 48.83524096625058, g.Latitudes()[14])
	assert.Equal(t, 46.044726631101675, g.Latitudes()[15])
	assert.Equal(t, 43.254194665350944, g.Latitudes()[16])
	assert.Equal(t, 40.46364817811504, g.Latitudes()[17])
	assert.Equal(t, 37.67308962904533, g.Latitudes()[18])
	assert.Equal(t, 34.88252099377346, g.Latitudes()[19])
	assert.Equal(t, 32.09194388174401, g.Latitudes()[20])
	assert.Equal(t, 29.301359621762735, g.Latitudes()[21])
	assert.Equal(t, 26.510769325210994, g.Latitudes()[22])
	assert.Equal(t, 23.720173933534745, g.Latitudes()[23])
	assert.Equal(t, 20.929574254489513, g.Latitudes()[24])
	assert.Equal(t, 18.13897099023935, g.Latitudes()[25])
	assert.Equal(t, 15.348364759491496, g.Latitudes()[26])
	assert.Equal(t, 12.55775611523068, g.Latitudes()[27])
	assert.Equal(t, 9.767145559195566, g.Latitudes()[28])
	assert.Equal(t, 6.976533553948636, g.Latitudes()[29])
	assert.Equal(t, 4.1859205331891545, g.Latitudes()[30])
	assert.Equal(t, 1.3953069108194958, g.Latitudes()[31])
	assert.Equal(t, -1.3953069108194958, g.Latitudes()[32])
	assert.Equal(t, -4.1859205331891545, g.Latitudes()[33])
	assert.Equal(t, -6.976533553948636, g.Latitudes()[34])
	assert.Equal(t, -9.767145559195566, g.Latitudes()[35])
	assert.Equal(t, -12.55775611523068, g.Latitudes()[36])
	assert.Equal(t, -15.348364759491496, g.Latitudes()[37])
	assert.Equal(t, -18.13897099023935, g.Latitudes()[38])
	assert.Equal(t, -20.929574254489513, g.Latitudes()[39])
	assert.Equal(t, -23.720173933534745, g.Latitudes()[40])
	assert.Equal(t, -26.510769325210994, g.Latitudes()[41])
	assert.Equal(t, -29.301359621762735, g.Latitudes()[42])
	assert.Equal(t, -32.09194388174401, g.Latitudes()[43])
	assert.Equal(t, -34.88252099377346, g.Latitudes()[44])
	assert.Equal(t, -37.67308962904533, g.Latitudes()[45])
	assert.Equal(t, -40.46364817811504, g.Latitudes()[46])
	assert.Equal(t, -43.254194665350944, g.Latitudes()[47])
	assert.Equal(t, -46.044726631101675, g.Latitudes()[48])
	assert.Equal(t, -48.83524096625058, g.Latitudes()[49])
	assert.Equal(t, -51.625733674938246, g.Latitudes()[50])
	assert.Equal(t, -54.4161995260862, g.Latitudes()[51])
	assert.Equal(t, -57.20663152764325, g.Latitudes()[52])
	assert.Equal(t, -59.99702010849129, g.Latitudes()[53])
	assert.Equal(t, -62.787351798963066, g.Latitudes()[54])
	assert.Equal(t, -65.57760701082782, g.Latitudes()[55])
	assert.Equal(t, -68.36775610831316, g.Latitudes()[56])
	assert.Equal(t, -71.15775201158732, g.Latitudes()[57])
	assert.Equal(t, -73.94751515398967, g.Latitudes()[58])
	assert.Equal(t, -76.73689968036832, g.Latitudes()[59])
	assert.Equal(t, -79.52560657265944, g.Latitudes()[60])
	assert.Equal(t, -82.31291294788628, g.Latitudes()[61])
	assert.Equal(t, -85.09652698831732, g.Latitudes()[62])
	assert.Equal(t, -87.86379883923267, g.Latitudes()[63])

	// 测试经度点数分布
	points := g.LonPoints()
	t.Logf("lon points: %v", points)
	assert.Equal(t, 4, points[0])    // 北极附近最少4个点
	assert.Equal(t, 4, points[9])    // 北极附近最少4个点
	assert.Equal(t, 128, points[10]) // 赤道附近最多点
	assert.Equal(t, 128, points[53]) // 赤道附近最多点
	assert.Equal(t, 4, points[54])   // 南极附近最少4个点
	assert.Equal(t, 4, points[63])   // 南极附近最少4个点
}

func TestOctahedral_O32_LongitudesOnLat(t *testing.T) {
	g := gaussian.NewOctahedral(32)

	assert.Equal(t, g.LonPoints()[31], len(g.LongitudesOnLat(0)))   // 31是赤道附近
	assert.Equal(t, g.LonPoints()[0], len(g.LongitudesOnLat(90)))   // 90是北极附近
	assert.Equal(t, g.LonPoints()[63], len(g.LongitudesOnLat(-90))) // -90是南极附近
}

func TestOctahedral_O32_GetNearestIndex(t *testing.T) {
	g := gaussian.NewOctahedral(32)

	// Test grid point lookup
	latIdx, lonIdx := g.GetNearestIndex(0.0, 0.0)
	assert.Equal(t, 31, latIdx) // 应该在赤道附近，但不是正好在赤道上
	assert.Equal(t, 0, lonIdx)  // 应该在格林威治子午线上
}

func BenchmarkNewOctahedral(b *testing.B) {
	for _, n := range []int{16, 32, 64, 128} {
		b.Run(fmt.Sprintf("O%d", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				gaussian.NewOctahedral(n)
			}
		})
	}
}
