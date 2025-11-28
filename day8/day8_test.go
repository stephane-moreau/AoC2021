package day8

import (
	"fmt"
	"io"
	"math"
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

type digits struct {
	value  int
	codes  []string
	digits []string
}

func convert(d *digits) {
	sort.SliceStable(d.codes, func(i, j int) bool {
		return len(d.codes[i]) < len(d.codes[j])
	})
	for i, c := range d.codes {
		bytes := []byte(c)
		sort.SliceStable(bytes, func(i, j int) bool {
			return bytes[i] < bytes[j]
		})
		d.codes[i] = string(bytes)
	}
	for i, c := range d.digits {
		bytes := []byte(c)
		sort.SliceStable(bytes, func(i, j int) bool {
			return bytes[i] < bytes[j]
		})
		d.digits[i] = string(bytes)
	}
	values := [10]int{1, 7, 4, 0, 0, 0, 0, 0, 0, 8}
	one := d.codes[0]
	var fiveIdentifier byte
	var twoIdentifier byte
	var five string
	// 6 is the only schema where a led from one is missing
	for i, c := range d.codes {
		if len(c) == 6 {
			for ndx := range one {
				if strings.IndexByte(c, one[ndx]) == -1 {
					values[i] = 6
					fiveIdentifier = one[ndx]
					if ndx == 0 {
						twoIdentifier = one[1]
					} else {
						twoIdentifier = one[0]
					}
					break
				}
			}
		}
	}
	for i, c := range d.codes {
		if len(c) == 5 {
			if strings.IndexByte(c, fiveIdentifier) == -1 {
				values[i] = 5
				five = c
			} else if strings.IndexByte(c, twoIdentifier) == -1 {
				values[i] = 2
			} else {
				values[i] = 3
			}
		}
	}
	for i, c := range d.codes {
		if len(c) == 6 && values[i] == 0 {
			missing := 0
			for d := range c {
				if strings.IndexByte(five, c[d]) == -1 {
					missing++
				}
			}
			if missing == 1 {
				values[i] = 9
			}
		}
	}

	powStart := len(d.digits) - 1
	for i, c := range d.digits {
		for j, n := range d.codes {
			if n == c {
				d.value += values[j] * int(math.Pow10(powStart-i))
			}
		}
	}
}

func TestConvert(t *testing.T) {
	d := digits{
		codes:  []string{"acedgfb", "cdfbe", "gcdfa", "fbcad", "dab", "cefabd", "cdfgeb", "eafb", "cagedb", "ab"},
		digits: []string{"cdfeb", "fcadb", "cdfeb", "cdbaf"},
	}
	convert(&d)
	assert.Equal(t, 5353, d.value)
}

type digitsconnections []digits

func readFile(fn string, useDiagonals bool) (digitsconnections, error) {
	f, err := os.OpenFile(fn, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	content, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	readLines := strings.Split(string(content), "\r\n")
	values := make(digitsconnections, len(readLines))

	for i, l := range readLines {
		rp := strings.Split(l, " | ")
		values[i].codes = strings.Split(rp[0], " ")
		values[i].digits = strings.Split(rp[1], " ")
	}

	return values, nil
}

func countKnownValues(values []string) int {
	c := 0
	for _, v := range values {
		if len(v) == 2 || len(v) == 3 || len(v) == 4 || len(v) == 7 {
			c++
		}
	}
	return c
}

func countKnowns(dc digitsconnections) int {
	c := 0
	for _, l := range dc {
		c += countKnownValues(l.digits)
	}
	return c
}

func TestDay8Phase1(t *testing.T) {
	dc, err := readFile("light.txt", false)
	require.NoError(t, err)
	assert.Equal(t, 26, countKnowns(dc))

	dc, err = readFile("large.txt", false)
	require.NoError(t, err)
	assert.Equal(t, 381, countKnowns(dc))
}

func sumValues(dc digitsconnections) int {
	for i, c := range dc {
		convert(&c)
		dc[i].value = c.value
	}
	s := 0
	for _, c := range dc {
		s += c.value
	}
	return s
}

func TestDay8Phase2(t *testing.T) {
	dc, err := readFile("light.txt", false)
	require.NoError(t, err)
	assert.Equal(t, 61229, sumValues(dc))

	dc, err = readFile("large.txt", false)
	require.NoError(t, err)
	assert.Equal(t, 1023686, sumValues(dc))
}
