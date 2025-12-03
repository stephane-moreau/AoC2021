package day12

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type graph map[string][]string

func load(fn string) (graph, error) {
	f, err := os.OpenFile(fn, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	content, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	readLines := strings.Split(string(content), "\r\n")
	g := make(graph)
	for _, l := range readLines {
		nodes := strings.Split(l, "-")
		g[nodes[0]] = append(g[nodes[0]], nodes[1])
		g[nodes[1]] = append(g[nodes[1]], nodes[0])
	}
	return g, nil
}

func countOccurences(path []string, val string) int {
	occ := 0
	for _, c := range path {
		if c == val {
			occ++
		}
	}
	return occ
}

var debug bool

func navigate(g graph, curPath []string, limit int) int {
	pathCount := 0
	for _, next := range g[curPath[len(curPath)-1]] {
		if next == "end" {
			pathCount++
			if debug {
				fmt.Printf("%v\n", append(curPath, "end"))
			}

			continue
		}
		if next == "start" {
			continue
		}
		occ := 0
		if next == strings.ToLower(next) {
			occ = countOccurences(curPath, next)
		}
		if occ >= limit {
			continue
		}
		newLimit := limit
		if occ != 0 && occ == limit-1 {
			newLimit = limit - 1
		}
		pathCount += navigate(g, append(curPath, next), newLimit)
	}
	return pathCount
}

func TestDay1(t *testing.T) {
	g, err := load("first.txt")
	require.NoError(t, err)
	debug = true
	assert.Equal(t, 10, navigate(g, []string{"start"}, 1))
	assert.Equal(t, 36, navigate(g, []string{"start"}, 2))

	debug = false
	g, err = load("second.txt")
	require.NoError(t, err)
	assert.Equal(t, 19, navigate(g, []string{"start"}, 1))
	assert.Equal(t, 103, navigate(g, []string{"start"}, 2))

	g, err = load("large.txt")
	require.NoError(t, err)
	assert.Equal(t, 226, navigate(g, []string{"start"}, 1))
	assert.Equal(t, 3509, navigate(g, []string{"start"}, 2))

	g, err = load("real.txt")
	require.NoError(t, err)
	assert.Equal(t, 5178, navigate(g, []string{"start"}, 1))
	assert.Equal(t, 130094, navigate(g, []string{"start"}, 2))
}
