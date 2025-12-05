package day15

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type grid [][]byte

func (g grid) value(x, y int) int {
	if x < 0 || y < 0 || y >= len(g) || x >= len(g[y]) {
		return -1
	}
	if x == 0 && y == 0 {
		return -1
	}
	return int(g[y][x] - '0')
}

func (g grid) size() int {
	return len(g)
}

func copy(fg, g grid, fx, fy int) {
	s := g.size()
	for y := range g {
		for x := range g {
			v := int(g[y][x]-'0') + fy + fx
			if v > 9 {
				v -= 9
			}
			fg[y+fy*s][x+fx*s] = '0' + byte(v)
		}
	}
}

func (g grid) fullGrid() grid {
	fg := make([][]byte, 5*g.size())
	for i := range fg {
		fg[i] = make([]byte, 5*g.size())
	}
	for fy := range 5 {
		for fx := range 5 {
			copy(fg, g, fx, fy)
		}
	}
	return fg
}

type point struct {
	x, y int
}

func loadGrid(fn string) (grid, error) {
	f, err := os.OpenFile(fn, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	content, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\r\n")
	g := make(grid, len(lines))
	for i := range lines {
		g[i] = []byte(lines[i])
	}
	return g, nil
}

func move(costs [][]int, g grid, pos point) {
	cur := costs[pos.y][pos.x]
	if pos.y == g.size()-1 && pos.x == g.size()-1 {
		return
	}
	for _, d := range []point{{1, 0}, {0, 1}, {0, -1}, {-1, 0}} {
		next := point{pos.x + d.x, pos.y + d.y}
		v := g.value(next.x, next.y)
		if v < 0 {
			// invalid position
			continue
		}
		if costs[g.size()-1][g.size()-1] != 0 && cur+v >= costs[g.size()-1][g.size()-1] {
			// Cost already would be higher than known minial path
			continue
		}
		if costs[next.y][next.x] == 0 ||
			cur+v < costs[next.y][next.x] {
			costs[next.y][next.x] = cur + v
			move(costs, g, point{pos.x + d.x, pos.y + d.y})
		}
	}
}

func findCheapestPath(g grid) int {
	start := point{}
	costs := make([][]int, len(g))
	for i := range costs {
		costs[i] = make([]int, len(g[0]))
	}
	move(costs, g, start)
	return costs[g.size()-1][g.size()-1]
}

func TestDay15Phase1(t *testing.T) {
	g, err := loadGrid("light.txt")
	require.NoError(t, err)
	require.NotNil(t, g)
	now := time.Now()
	assert.Equal(t, 40, findCheapestPath(g))
	fmt.Printf("%f s\n", time.Since(now).Seconds())
	now = time.Now()
	assert.Equal(t, 315, findCheapestPath(g.fullGrid()))
	fmt.Printf("%f s\n", time.Since(now).Seconds())

	g, err = loadGrid("large.txt")
	require.NoError(t, err)
	require.NotNil(t, g)
	now = time.Now()
	assert.Equal(t, 441, findCheapestPath(g))
	fmt.Printf("%f s\n", time.Since(now).Seconds())
	now = time.Now()
	assert.Equal(t, 2849, findCheapestPath(g.fullGrid()))
	fmt.Printf("%f s\n", time.Since(now).Seconds())
}
