package day4

import (
	"io"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const size = 5

type board struct {
	grid           [size][size]string
	matchedLines   [size]int
	matchedColumns [size]int
}

func readFile(fn string) ([]string, []*board, error) {
	f, err := os.OpenFile(fn, os.O_RDONLY, 0)
	if err != nil {
		return nil, nil, err
	}
	content, err := io.ReadAll(f)
	if err != nil {
		return nil, nil, err
	}

	lines := strings.Split(string(content), "\n")
	draw := strings.Split(strings.TrimSpace(lines[0]), ",")
	boards := make([]*board, 0)
	for i := 2; i < len(lines); i += 6 {
		var brd board
		for r := range 5 {
			values := strings.Split(lines[i+r], " ")
			var c int
			for _, v := range values {
				v := strings.TrimSpace(v)
				if v == "" {
					continue
				}
				brd.grid[r][c] = v
				c++
			}
		}
		boards = append(boards, &brd)
	}
	return draw, boards, nil
}

func findBingo(draw []string, boards []*board, firstWinner bool) ([]string, *board) {
	winners := make(map[int]int)
	for ndx, d := range draw {
		for ndxBoard, b := range boards {
			if winners[ndxBoard] > 0 {
				continue
			}
			for r := range b.grid {
				for c := range b.grid[r] {
					if d == b.grid[r][c] {
						b.matchedLines[r]++
						b.matchedColumns[c]++
						if b.matchedLines[r] == 5 || b.matchedColumns[c] == 5 {
							winners[ndxBoard] = ndx
							if firstWinner {
								return draw[:ndx+1], b
							}
						}
					}
				}
			}
		}
	}
	winnerBoard := 0
	winnerDraw := 0
	for ndxBoard, winner := range winners {
		if winner > winnerDraw {
			winnerBoard = ndxBoard
			winnerDraw = winner
		}
	}
	return draw[:winnerDraw+1], boards[winnerBoard]
}

func scoreCard(draw []string, b *board) int {
	setBoard := make(map[string]bool)
	for r := range b.grid {
		for c := range b.grid[r] {
			setBoard[b.grid[r][c]] = false
		}
	}
	for _, d := range draw {
		setBoard[d] = true
	}
	var res int
	for k, v := range setBoard {
		if !v {
			i, _ := strconv.Atoi(k)
			res += i
		}
	}
	d, _ := strconv.Atoi(draw[len(draw)-1])
	res *= d
	return res
}

func TestDay4Phase1(t *testing.T) {
	draw, boards, err := readFile("light.txt")
	require.NoError(t, err)

	v, brd := findBingo(draw, boards, true)
	res := scoreCard(v, brd)
	assert.Equal(t, 4512, res)

	draw, boards, err = readFile("large.txt")
	require.NoError(t, err)
	v, brd = findBingo(draw, boards, true)
	res = scoreCard(v, brd)
	assert.Equal(t, 28082, res)
}

func TestDay4Phase2(t *testing.T) {
	draw, boards, err := readFile("light.txt")
	require.NoError(t, err)

	v, brd := findBingo(draw, boards, false)
	res := scoreCard(v, brd)
	assert.Equal(t, 1924, res)

	draw, boards, err = readFile("large.txt")
	require.NoError(t, err)
	v, brd = findBingo(draw, boards, false)
	res = scoreCard(v, brd)
	assert.Equal(t, 8224, res)
}
