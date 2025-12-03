package day13

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type point struct {
	x, y int64
}

type grid []point

func (g grid) apply(fld fold) {
	for _, f := range fld {
		for i, p := range g {
			if f.y != 0 && p.y > f.y {
				p.y = 2*f.y - p.y
			}
			if f.x != 0 && p.x > f.x {
				p.x = 2*f.x - p.x
			}
			g[i] = p
		}
	}
}

func (g grid) size() int {
	dedup := make(map[point]bool)
	for _, p := range g {
		dedup[p] = true
	}
	return len(dedup)
}

func (g grid) dump() {
	sort.SliceStable(g, func(i, j int) bool {
		return g[i].x < g[j].x || g[i].x == g[j].x && g[i].y < g[j].y
	})

	grd := make([][]byte, g[len(g)-1].y+1)
	for i := range len(grd) {
		grd[i] = []byte(strings.Repeat(" ", int(g[len(g)-1].x+1)))
	}
	for _, p := range g {
		grd[p.y][p.x] = '#'
	}
	var s string
	for _, l := range grd {
		s += fmt.Sprintf("%s\n", string(l))
	}
	fmt.Printf("%s\n", s)
}

type fold []point

func load(fn string) (grid, fold, error) {
	f, err := os.OpenFile(fn, os.O_RDONLY, 0)
	if err != nil {
		return nil, nil, err
	}
	content, err := io.ReadAll(f)
	if err != nil {
		return nil, nil, err
	}

	readLines := strings.Split(string(content), "\r\n")
	grd := make(grid, 0, len(readLines))
	fld := make(fold, 0, len(readLines))
	for _, l := range readLines {
		if strings.HasPrefix(l, "fold along x=") {
			x, err := strconv.ParseInt(strings.TrimPrefix(l, "fold along x="), 10, 64)
			if err != nil {
				return nil, nil, err
			}
			fld = append(fld, point{x, 0})
			continue
		}
		if strings.HasPrefix(l, "fold along y=") {
			y, err := strconv.ParseInt(strings.TrimPrefix(l, "fold along y="), 10, 64)
			if err != nil {
				return nil, nil, err
			}
			fld = append(fld, point{0, y})
			continue
		}

		coords := strings.Split(l, ",")
		if len(coords) < 2 {
			continue
		}
		x, err := strconv.ParseInt(coords[0], 10, 64)
		if err != nil {
			return nil, nil, err
		}
		y, err := strconv.ParseInt(coords[1], 10, 64)
		if err != nil {
			return nil, nil, err
		}
		grd = append(grd, point{x, y})
	}
	return grd, fld, nil
}

func TestDay13Phase1(t *testing.T) {
	grd, fld, err := load("light.txt")
	require.NoError(t, err)
	grd.apply(fld[:1])
	assert.Equal(t, 17, grd.size())
	grd.apply(fld[1:])
	assert.Equal(t, 16, grd.size())
	grd.dump()

	grd, fld, err = load("large.txt")
	require.NoError(t, err)
	grd.apply(fld[:1])
	assert.Equal(t, 720, grd.size())
	grd.apply(fld[1:])
	grd.dump()
}
