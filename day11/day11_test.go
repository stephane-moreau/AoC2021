package day11

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type grid [][]int

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
	g := make(grid, len(readLines))
	for i, l := range readLines {
		g[i] = make([]int, len(l))
		for j := range l {
			g[i][j] = int(l[j] - '0')
		}
	}
	return g, nil
}

func (g grid) inc(i, j int) int {
	if i < 0 || j < 0 || i >= len(g) || j >= len(g) {
		return 0
	}
	g[i][j]++
	return g[i][j]
}

func (g grid) next() int {
	flashes := make([][]bool, len(g))
	for i := range len(g) {
		flashes[i] = make([]bool, len(g))
	}

	// Increase energy by 1
	for i := range g {
		for j := range g[i] {
			g.inc(i, j)
		}
	}
	// increase by 1 adjacent flashing octopus
	newFlash := true
	for newFlash {
		newFlash = false
		for i := range g {
			for j := range g[i] {
				if g[i][j] > 9 && !flashes[i][j] {
					newFlash = true
					flashes[i][j] = true
					g.inc(i-1, j-1)
					g.inc(i-1, j)
					g.inc(i-1, j+1)
					g.inc(i, j-1)
					g.inc(i, j+1)
					g.inc(i+1, j-1)
					g.inc(i+1, j)
					g.inc(i+1, j+1)
				}
			}
		}
	}
	// reset wave and count flashes
	flash := 0
	for i := range g {
		for j := range g[i] {
			if flashes[i][j] {
				flash++
			}
			if g[i][j] > 9 {
				g[i][j] = 0
			}
		}
	}

	return flash
}

func flashCount(g grid, iter int) int {
	fc := 0
	for range iter {
		fc += g.next()
	}
	return fc
}

func TestLight(t *testing.T) {
	g := grid{
		{1, 1, 1, 1, 1},
		{1, 9, 9, 9, 1},
		{1, 9, 1, 9, 1},
		{1, 9, 9, 9, 1},
		{1, 1, 1, 1, 1},
	}
	g.next()
	g.next()
}

func TestDay11Phase1(t *testing.T) {
	g, err := load("test.txt")
	require.NoError(t, err)
	assert.Equal(t, 204, flashCount(g, 10))

	g, err = load("test.txt")
	require.NoError(t, err)
	assert.Equal(t, 1656, flashCount(g, 100))

	g, err = load("real.txt")
	require.NoError(t, err)
	assert.Equal(t, 1594, flashCount(g, 100))
}

func TestDay11Phase2(t *testing.T) {
	g, err := load("test.txt")
	require.NoError(t, err)
	flashCount := 0
	iter := 0
	for flashCount != 100 {
		iter++
		flashCount = g.next()
	}
	assert.Equal(t, 195, iter)

	g, err = load("real.txt")
	require.NoError(t, err)
	flashCount = 0
	iter = 0
	for flashCount != 100 {
		iter++
		flashCount = g.next()
	}
	assert.Equal(t, 437, iter)
}
