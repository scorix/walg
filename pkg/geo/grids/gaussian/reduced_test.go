package gaussian_test

import (
	"fmt"
	"testing"

	"github.com/scorix/walg/pkg/geo/grids/gaussian"
	"github.com/stretchr/testify/assert"
)

func TestReduced_N32(t *testing.T) {
	g := gaussian.NewReduced(32)

	t.Logf("latitudes: %v", g.Latitudes())
	t.Logf("longitudes: %v", g.Longitudes())

	// Test basic grid properties
	assert.Equal(t, 64, len(g.Latitudes()))
	assert.Equal(t, 128, len(g.Longitudes())) // 4*N at equator
	assert.Equal(t, 5144, g.Size())

	assert.Equal(t, 90.0, g.Latitudes()[0])
	assert.Equal(t, 87.14285714285714, g.Latitudes()[1])
	assert.Equal(t, 84.28571428571429, g.Latitudes()[2])
	assert.Equal(t, 81.42857142857143, g.Latitudes()[3])
	assert.Equal(t, 78.57142857142857, g.Latitudes()[4])
	assert.Equal(t, 75.71428571428571, g.Latitudes()[5])
	assert.Equal(t, 72.85714285714286, g.Latitudes()[6])
	assert.Equal(t, 70.0, g.Latitudes()[7])
	assert.Equal(t, 67.14285714285714, g.Latitudes()[8])
	assert.Equal(t, 64.28571428571428, g.Latitudes()[9])
	assert.Equal(t, 61.42857142857143, g.Latitudes()[10])
	assert.Equal(t, 58.57142857142857, g.Latitudes()[11])
	assert.Equal(t, 55.714285714285715, g.Latitudes()[12])
	assert.Equal(t, 52.857142857142854, g.Latitudes()[13])
	assert.Equal(t, 50.0, g.Latitudes()[14])
	assert.Equal(t, 47.14285714285714, g.Latitudes()[15])
	assert.Equal(t, 44.285714285714285, g.Latitudes()[16])
	assert.Equal(t, 41.42857142857143, g.Latitudes()[17])
	assert.Equal(t, 38.57142857142857, g.Latitudes()[18])
	assert.Equal(t, 35.714285714285715, g.Latitudes()[19])
	assert.Equal(t, 32.857142857142854, g.Latitudes()[20])
	assert.Equal(t, 30.0, g.Latitudes()[21])
	assert.Equal(t, 27.14285714285714, g.Latitudes()[22])
	assert.Equal(t, 24.285714285714278, g.Latitudes()[23])
	assert.Equal(t, 21.42857142857143, g.Latitudes()[24])
	assert.Equal(t, 18.57142857142857, g.Latitudes()[25])
	assert.Equal(t, 15.714285714285708, g.Latitudes()[26])
	assert.Equal(t, 12.857142857142861, g.Latitudes()[27])
	assert.Equal(t, 10.0, g.Latitudes()[28])
	assert.Equal(t, 7.142857142857139, g.Latitudes()[29])
	assert.Equal(t, 4.285714285714278, g.Latitudes()[30])
	assert.Equal(t, 1.4285714285714306, g.Latitudes()[31])
	assert.Equal(t, -1.4285714285714306, g.Latitudes()[32])
	assert.Equal(t, -4.285714285714292, g.Latitudes()[33])
	assert.Equal(t, -7.142857142857139, g.Latitudes()[34])
	assert.Equal(t, -10.0, g.Latitudes()[35])
	assert.Equal(t, -12.857142857142861, g.Latitudes()[36])
	assert.Equal(t, -15.714285714285722, g.Latitudes()[37])
	assert.Equal(t, -18.57142857142857, g.Latitudes()[38])
	assert.Equal(t, -21.42857142857143, g.Latitudes()[39])
	assert.Equal(t, -24.285714285714292, g.Latitudes()[40])
	assert.Equal(t, -27.14285714285714, g.Latitudes()[41])
	assert.Equal(t, -30.0, g.Latitudes()[42])
	assert.Equal(t, -32.85714285714286, g.Latitudes()[43])
	assert.Equal(t, -35.71428571428572, g.Latitudes()[44])
	assert.Equal(t, -38.571428571428584, g.Latitudes()[45])
	assert.Equal(t, -41.428571428571445, g.Latitudes()[46])
	assert.Equal(t, -44.28571428571428, g.Latitudes()[47])
	assert.Equal(t, -47.14285714285714, g.Latitudes()[48])
	assert.Equal(t, -50.0, g.Latitudes()[49])
	assert.Equal(t, -52.85714285714286, g.Latitudes()[50])
	assert.Equal(t, -55.71428571428572, g.Latitudes()[51])
	assert.Equal(t, -58.571428571428584, g.Latitudes()[52])
	assert.Equal(t, -61.428571428571445, g.Latitudes()[53])
	assert.Equal(t, -64.28571428571428, g.Latitudes()[54])
	assert.Equal(t, -67.14285714285714, g.Latitudes()[55])
	assert.Equal(t, -70.0, g.Latitudes()[56])
	assert.Equal(t, -72.85714285714286, g.Latitudes()[57])
	assert.Equal(t, -75.71428571428572, g.Latitudes()[58])
	assert.Equal(t, -78.57142857142858, g.Latitudes()[59])
	assert.Equal(t, -81.42857142857144, g.Latitudes()[60])
	assert.Equal(t, -84.28571428571428, g.Latitudes()[61])
	assert.Equal(t, -87.14285714285714, g.Latitudes()[62])
	assert.Equal(t, -90.0, g.Latitudes()[63])

	// 测试经度点数分布
	points := g.LonPoints()
	t.Logf("lon points: %v", points)
	assert.Equal(t, 4, points[0])    // 北极点
	assert.Equal(t, 8, points[1])    // 高纬度区域
	assert.Equal(t, 12, points[2])   // 高纬度区域
	assert.Equal(t, 20, points[3])   // 高纬度区域
	assert.Equal(t, 24, points[4])   // 高纬度区域
	assert.Equal(t, 32, points[5])   // 北半球高纬度
	assert.Equal(t, 36, points[6])   // 北半球高纬度
	assert.Equal(t, 44, points[7])   // 北半球高纬度
	assert.Equal(t, 48, points[8])   // 北半球中纬度
	assert.Equal(t, 56, points[9])   // 北半球中纬度
	assert.Equal(t, 60, points[10])  // 北半球中纬度
	assert.Equal(t, 68, points[11])  // 北半球中纬度
	assert.Equal(t, 72, points[12])  // 北半球中纬度
	assert.Equal(t, 76, points[13])  // 北半球中纬度
	assert.Equal(t, 84, points[14])  // 北半球中纬度
	assert.Equal(t, 88, points[15])  // 北半球中纬度
	assert.Equal(t, 92, points[16])  // 北半球中纬度
	assert.Equal(t, 96, points[17])  // 北半球中纬度
	assert.Equal(t, 100, points[18]) // 北半球中纬度
	assert.Equal(t, 104, points[19]) // 北半球中纬度
	assert.Equal(t, 108, points[20]) // 北半球中纬度
	assert.Equal(t, 112, points[21]) // 北半球中纬度
	assert.Equal(t, 112, points[22]) // 北半球中纬度
	assert.Equal(t, 116, points[23]) // 北半球中纬度
	assert.Equal(t, 120, points[24]) // 北半球中纬度
	assert.Equal(t, 120, points[25]) // 北半球中纬度
	assert.Equal(t, 124, points[26]) // 近赤道区域
	assert.Equal(t, 124, points[27]) // 近赤道区域
	assert.Equal(t, 128, points[28]) // 赤道区域
	assert.Equal(t, 128, points[29]) // 赤道区域
	assert.Equal(t, 128, points[30]) // 赤道区域
	assert.Equal(t, 128, points[31]) // 赤道区域
	assert.Equal(t, 128, points[32]) // 赤道区域
	assert.Equal(t, 128, points[33]) // 赤道区域
	assert.Equal(t, 128, points[34]) // 赤道区域
	assert.Equal(t, 128, points[35]) // 赤道区域
	assert.Equal(t, 124, points[36]) // 近赤道区域
	assert.Equal(t, 124, points[37]) // 近赤道区域
	assert.Equal(t, 120, points[38]) // 南半球中纬度
	assert.Equal(t, 120, points[39]) // 南半球中纬度
	assert.Equal(t, 116, points[40]) // 南半球中纬度
	assert.Equal(t, 112, points[41]) // 南半球中纬度
	assert.Equal(t, 112, points[42]) // 南半球中纬度
	assert.Equal(t, 108, points[43]) // 南半球中纬度
	assert.Equal(t, 104, points[44]) // 南半球中纬度
	assert.Equal(t, 100, points[45]) // 南半球中纬度
	assert.Equal(t, 96, points[46])  // 南半球中纬度
	assert.Equal(t, 92, points[47])  // 南半球中纬度
	assert.Equal(t, 88, points[48])  // 南半球中纬度
	assert.Equal(t, 84, points[49])  // 南半球中纬度
	assert.Equal(t, 76, points[50])  // 南半球中纬度
	assert.Equal(t, 72, points[51])  // 南半球中纬度
	assert.Equal(t, 68, points[52])  // 南半球中纬度
	assert.Equal(t, 60, points[53])  // 南半球中纬度
	assert.Equal(t, 56, points[54])  // 南半球高纬度
	assert.Equal(t, 48, points[55])  // 南半球高纬度
	assert.Equal(t, 44, points[56])  // 南半球高纬度
	assert.Equal(t, 36, points[57])  // 南半球高纬度
	assert.Equal(t, 32, points[58])  // 南半球高纬度
	assert.Equal(t, 24, points[59])  // 高纬度区域
	assert.Equal(t, 20, points[60])  // 高纬度区域
	assert.Equal(t, 12, points[61])  // 高纬度区域
	assert.Equal(t, 8, points[62])   // 高纬度区域
	assert.Equal(t, 4, points[63])   // 南极点
}

func BenchmarkNewReduced(b *testing.B) {
	for _, n := range []int{16, 32, 64, 128} {
		b.Run(fmt.Sprintf("N%d", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				gaussian.NewReduced(n)
			}
		})
	}
}
