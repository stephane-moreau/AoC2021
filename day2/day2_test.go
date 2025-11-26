package day2

import (
	"io"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var ()

func compile(fn string) (int, int, error) {
	f, err := os.OpenFile(fn, os.O_RDONLY, 0)
	if err != nil {
		return 0, 0, err
	}
	content, err := io.ReadAll(f)
	if err != nil {
		return 0, 0, err
	}
	var x, d, realD, aim int
	for _, l := range strings.Split(string(content), "\n") {
		if strings.HasPrefix(l, "forward") {
			delta, err := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(l, "forward")))
			if err != nil {
				return 0, 0, err
			}
			x += delta
			realD += aim * delta
		} else if strings.HasPrefix(l, "down") {
			delta, err := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(l, "down")))
			if err != nil {
				return 0, 0, err
			}
			d += delta
			aim += delta
		} else if strings.HasPrefix(l, "up") {
			delta, err := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(l, "up")))
			if err != nil {
				return 0, 0, err
			}
			d -= delta
			aim -= delta
		}
	}
	return x * d, x * realD, nil
}

func TestDay2(t *testing.T) {
	res, realRes, err := compile("light.txt")
	require.NoError(t, err)
	require.Equal(t, 150, res)
	require.Equal(t, 900, realRes)

	res, realRes, err = compile("large.txt")
	assert.NoError(t, err)
	assert.Equal(t, 1427868, res)
	assert.Equal(t, 1568138742, realRes)
}
