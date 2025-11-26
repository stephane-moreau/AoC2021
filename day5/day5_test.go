package day5

import (
	"io"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type point struct {
	x, y int
}

func move(pt, target point) point {
	var deltaX, deltaY int
	if pt.x < target.x {
		deltaX = 1
	} else if pt.x > target.x {
		deltaX = -1
	}
	if pt.y < target.y {
		deltaY = 1
	} else if pt.y > target.y {
		deltaY = -1
	}
	return point{
		pt.x + deltaX,
		pt.y + deltaY,
	}
}

func parsePoint(s string) point {
	coords := strings.Split(s, ",")
	var pt point
	pt.x, _ = strconv.Atoi(strings.TrimSpace(coords[0]))
	pt.y, _ = strconv.Atoi(strings.TrimSpace(coords[1]))
	return pt
}

type segment struct {
	start, end point
}

type board [][]int

func readFile(fn string, useDiagonals bool) (board, error) {
	f, err := os.OpenFile(fn, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	content, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	readLines := strings.Split(string(content), "\n")
	lines := make([]segment, 0, len(readLines))

	var maxX, maxY int
	for _, l := range readLines {
		rp := strings.Split(l, " -> ")
		start := parsePoint(rp[0])
		end := parsePoint(rp[1])
		lines = append(lines, segment{start, end})
		if start.x > maxX {
			maxX = start.x
		}
		if end.x > maxX {
			maxX = end.x
		}
		if start.y > maxY {
			maxY = start.y
		}
		if end.y > maxY {
			maxY = end.y
		}
	}

	grid := make(board, maxY+1)
	for i := range grid {
		grid[i] = make([]int, maxX+1)
	}

	for _, l := range lines {
		if l.start.x == l.end.x || l.start.y == l.end.y || useDiagonals {
			for pt := l.start; ; pt = move(pt, l.end) {
				grid[pt.y][pt.x]++
				if pt == l.end {
					break
				}
			}
		}
	}
	return grid, nil
}

func countOverlaps(brd board) int {
	var res int
	for _, r := range brd {
		for _, c := range r {
			if c > 1 {
				res++
			}
		}
	}
	return res
}

func TestDay5Phase1(t *testing.T) {
	board, err := readFile("light.txt", false)
	require.NoError(t, err)
	assert.Equal(t, 5, countOverlaps(board))

	board, err = readFile("large.txt", false)
	require.NoError(t, err)
	assert.Equal(t, 5280, countOverlaps(board))
}

func TestDay5Phase2(t *testing.T) {
	board, err := readFile("light.txt", true)
	require.NoError(t, err)
	assert.Equal(t, 12, countOverlaps(board))

	board, err = readFile("large.txt", true)
	require.NoError(t, err)
	assert.Equal(t, 16716, countOverlaps(board))
}
