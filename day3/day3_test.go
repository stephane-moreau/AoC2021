package day3

import (
	"io"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func readLines(fn string) ([]string, error) {
	f, err := os.OpenFile(fn, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	content, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return strings.Split(string(content), "\n"), nil
}

func compile(fn string) (int, error) {
	var g, e int
	lines, err := readLines(fn)
	if err != nil {
		return 0, err
	}
	length := len(strings.TrimSpace(lines[0]))
	for i := range length {
		g = g << 1
		e = e << 1
		var zeros, ones int
		for _, l := range lines {
			if l[i] == '0' {
				zeros++
			} else {
				ones++
			}
		}
		if zeros > ones {
			g += 1
		} else {
			e += 1
		}
	}
	return g * e, nil
}

func split(bit int, lines []string) ([]string, []string) {
	zeroLines := make([]string, 0, len(lines)/2)
	oneLines := make([]string, 0, len(lines)/2)
	for _, l := range lines {
		if l[bit] == '0' {
			zeroLines = append(zeroLines, l)
		} else {
			oneLines = append(oneLines, l)
		}
	}
	return zeroLines, oneLines
}

func computeRatings(fn string) (int64, error) {
	lines, err := readLines(fn)
	if err != nil {
		return 0, err
	}
	length := len(strings.TrimSpace(lines[0]))
	var o2, co2 int64
	o2Lines := lines
	co2Lines := lines
	for i := range length {
		if len(o2Lines) > 1 {
			zeroLines, oneLines := split(i, o2Lines)
			if len(zeroLines) > len(oneLines) {
				o2Lines = zeroLines
			} else {
				o2Lines = oneLines
			}
		}
		if len(co2Lines) > 1 {
			zeroLines, oneLines := split(i, co2Lines)
			if len(zeroLines) > len(oneLines) {
				co2Lines = oneLines
			} else {
				co2Lines = zeroLines
			}
		}
	}
	o2, err = strconv.ParseInt(strings.TrimSpace(o2Lines[0]), 2, 64)
	if err != nil {
		return 0, err
	}
	co2, err = strconv.ParseInt(strings.TrimSpace(co2Lines[0]), 2, 64)
	if err != nil {
		return 0, err
	}
	return o2 * co2, nil
}

func TestDay3Phase1(t *testing.T) {
	res, err := compile("light.txt")
	require.NoError(t, err)
	require.Equal(t, 198, res)

	res, err = compile("large.txt")
	assert.NoError(t, err)
	assert.Equal(t, 2261546, res)
}

func TestDay3Phase2(t *testing.T) {
	res, err := computeRatings("light.txt")
	require.NoError(t, err)
	require.Equal(t, 230, res)

	res, err = computeRatings("large.txt")
	assert.NoError(t, err)
	assert.Equal(t, 6775520, res)
}
