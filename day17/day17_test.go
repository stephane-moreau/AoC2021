package day17

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

type point struct {
	x, y int
}

type target struct {
	topleft, bottomright point
}

var (
	testInput = target{
		point{20, -5},
		point{30, -10},
	}

	input = target{
		point{70, -121},
		point{125, -159},
	}
)

func xPos(v, t int) int {
	if t <= v {
		return t*v - t*(t-1)/2
	}
	return v * (v + 1) / 2
}

func yPos(v, t int) int {
	return t*v - t*(t-1)/2
}

// x speed
// x(t) = vx * t - t * (t-1) / 2  (vx>0)
// 20 < t^2/2 + t (vx-1/2) < 30
// y(t) = vy * t - t * (t-1) / 2
// y'(t) = vy - t -> max obtained a t = vy
func findInitialVelocity(start point, tgt target) (point, int, int) {
	var yMax = -1
	var velocity point
	var hit = 0
	for vx := int(math.Sqrt(float64(2 * tgt.topleft.x))); vx <= tgt.bottomright.x; vx++ {
		if xPos(vx, vx) < tgt.topleft.x {
			continue
		}
		for vy := tgt.bottomright.y; vy <= -tgt.bottomright.y; vy++ {
			for t := 1; xPos(vx, t) <= tgt.bottomright.x && t < 1000; t++ {
				x := xPos(vx, t)
				y := yPos(vy, t)
				if x < tgt.topleft.x {
					continue
				}
				if y <= tgt.topleft.y && y >= tgt.bottomright.y {
					var yMaxVy int
					if vy > 0 {
						// Max is obtained at t = vy
						yMaxVy = vy*vy - vy*(vy-1)/2
					}
					hit++
					if yMaxVy > yMax {
						yMax = yMaxVy
						velocity = point{vx, vy}
					}
					break
				}
				if y < tgt.bottomright.y {
					break
				}
			}
		}
	}
	return velocity, yMax, hit
}

func TestDay17(t *testing.T) {
	v, yMax, hits := findInitialVelocity(point{}, testInput)
	assert.Equal(t, 45, yMax)
	assert.Equal(t, 112, hits)
	assert.Equal(t, point{6, 9}, v)

	v, yMax, hits = findInitialVelocity(point{}, input)
	assert.Equal(t, 12561, yMax)
	assert.Equal(t, 3785, hits)
	assert.Equal(t, point{12, 158}, v)
}
