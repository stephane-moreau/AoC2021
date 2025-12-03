package day10

import (
	"io"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func load(fn string) ([]string, error) {
	f, err := os.OpenFile(fn, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	content, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	readLines := strings.Split(string(content), "\r\n")
	return readLines, nil
}

var (
	codePairs = map[byte]byte{
		'(': ')',
		'[': ']',
		'{': '}',
		'<': '>',
	}

	codeValues = map[byte]int{
		')': 3,
		']': 57,
		'}': 1197,
		'>': 25137,
		'(': 1,
		'[': 2,
		'{': 3,
		'<': 4,
	}
)

func isValidLine(l string) (int, []byte) {
	stack := make([]byte, 0, len(l))
	for i := range l {
		c := l[i]
		if codePairs[c] != 0 {
			stack = append(stack, c)
		} else if c != codePairs[stack[len(stack)-1]] {
			return codeValues[c], nil
		} else {
			stack = stack[:len(stack)-1]
		}
	}
	return 0, stack
}

func sumInvalidCodes(lines []string) int {
	s := 0
	for _, l := range lines {
		v, _ := isValidLine(l)
		s += v
	}
	return s
}

func completeScore(lines []string) int {
	scores := make([]int, 0, len(lines))
	for _, l := range lines {
		_, toComplete := isValidLine(l)
		if len(toComplete) == 0 {
			continue
		}
		score := 0
		for i := len(toComplete) - 1; i >= 0; i-- {
			score *= 5
			score += codeValues[toComplete[i]]
		}
		scores = append(scores, score)
	}
	sort.SliceStable(scores, func(i, j int) bool {
		return scores[i] < scores[j]
	})
	return scores[len(scores)/2]
}

func TestDay10Phase1(t *testing.T) {
	lines, err := load("light.txt")
	require.NoError(t, err)
	assert.Equal(t, 26397, sumInvalidCodes(lines))

	lines, err = load("large.txt")
	require.NoError(t, err)
	assert.Equal(t, 374061, sumInvalidCodes(lines))
}

func TestDay10Phase2(t *testing.T) {
	lines, err := load("light.txt")
	require.NoError(t, err)
	assert.Equal(t, 288957, completeScore(lines))

	lines, err = load("large.txt")
	require.NoError(t, err)
	assert.Equal(t, 2116639949, completeScore(lines))
}
