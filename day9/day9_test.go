package day9

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var debug = true

func debugLog(s string, vals ...any) {
	if debug {
		fmt.Printf(s, vals...)
	}
}

type grid [][]byte

type point struct {
	x, y int
}

const invalidBorder = byte(127)

func load(fn string) (grid, error) {
	f, err := os.OpenFile(fn, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	content, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	readLines := strings.Split(string(content), "\r\n")

	values := make(grid, len(readLines))
	for i, l := range readLines {
		values[i] = []byte(l)
	}
	return values, nil
}

func (g grid) value(x, y int) byte {
	if x < 0 || y < 0 || x >= len(g[0]) || y >= len(g) {
		return invalidBorder
	}
	return g[y][x]
}

func isLocalMinimum(g grid, x, y int) bool {
	v := g.value(x, y)
	return v < g.value(x, y-1) &&
		v < g.value(x-1, y) && v < g.value(x+1, y) &&
		v < g.value(x, y+1)
}

func countLocalMinimum(g grid) int {
	c := 0
	for y := 0; y < len(g); y++ {
		for x := 0; x < len(g[y]); x++ {
			if isLocalMinimum(g, x, y) {
				c += int(1 + g.value(x, y) - '0')
			}
		}
	}
	return c
}

func TestDay9Phase1(t *testing.T) {
	g, err := load("light.txt")
	require.NoError(t, err)
	assert.Equal(t, 15, countLocalMinimum(g))

	g, err = load("large.txt")
	require.NoError(t, err)
	// 1681 high
	// 523 low
	assert.Equal(t, 528, countLocalMinimum(g))
}

func extendBassin(g grid, x, y int) int {
	dedup := map[point]bool{point{x, y}: true}
	size := 0
	for size != len(dedup) {
		size = len(dedup)
		for p := range dedup {
			v := g.value(p.x, p.y)
			if nV := g.value(p.x-1, p.y); nV > v && nV < '9' {
				dedup[point{p.x - 1, p.y}] = true
			}
			if nV := g.value(p.x+1, p.y); nV > v && nV < '9' {
				dedup[point{p.x + 1, p.y}] = true
			}
			if nV := g.value(p.x, p.y-1); nV > v && nV < '9' {
				dedup[point{p.x, p.y - 1}] = true
			}
			if nV := g.value(p.x, p.y+1); nV > v && nV < '9' {
				dedup[point{p.x, p.y + 1}] = true
			}
		}
	}
	return size
}

func largestBassins(g grid) int {
	var sizes []int
	for y := 0; y < len(g); y++ {
		for x := 0; x < len(g[y]); x++ {
			if isLocalMinimum(g, x, y) {
				sizes = append(sizes, extendBassin(g, x, y))
			}
		}
	}
	sort.SliceStable(sizes, func(i, j int) bool {
		return sizes[i] > sizes[j]
	})
	return sizes[0] * sizes[1] * sizes[2]
}

func TestDay9Phase2(t *testing.T) {
	g, err := load("light.txt")
	require.NoError(t, err)
	assert.Equal(t, 1134, largestBassins(g))

	g, err = load("large.txt")
	require.NoError(t, err)
	assert.Equal(t, 920448, largestBassins(g))
}
