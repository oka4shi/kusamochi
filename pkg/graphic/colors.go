package graphic

import (
	"math"
	"math/rand/v2"
)

func hslToRgb(hue int, sat, light float64) (float64, float64, float64) {
	hue = hue % 360

	if hue < 0 {
		hue += 360
	}

	c := [3]float64{}
	for i, n := range []int{0, 8, 4} {
		k := (n + hue/30) % 12
		a := sat * min(light, 1-light)
		c[i] = light - a*float64(max(-1, min(k-3, 9-k, 1)))
	}

	return c[0], c[1], c[2]
}

func GenerateRandomColor() (float64, float64, float64) {
	h := rand.IntN(35) * 10 // from 0 to 350 (step: 10)
	s := 0.5
	l := float64(40+rand.IntN(2)*20) / 100 // from 0.4 to 0.8 (step: 0.2)
	return hslToRgb(h, s, l)
}

func GenerateColors(count int) [][3]float64 {
	step := 360 / max(count, 12)

	rgbs := [][3]float64{}

	for range int(math.Ceil(float64(count) / (12 * 3))) {
		for j := range int(math.Ceil(float64(count)/12)) % 3 {
			for i := range min(12, count-12*j) {
				h := i * step
				s := 0.5
				l := []float64{0.5, 0.4, 0.6}[j]
				r, g, b := hslToRgb(h, s, l)
				rgbs = append(rgbs, [3]float64{r, g, b})
			}
		}
	}

	return rgbs
}
