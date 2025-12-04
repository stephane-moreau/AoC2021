package day14

import (
	"io"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func load(fn string) (string, map[string]string, error) {
	f, err := os.OpenFile(fn, os.O_RDONLY, 0)
	if err != nil {
		return "", nil, err
	}
	content, err := io.ReadAll(f)
	if err != nil {
		return "", nil, err
	}

	readLines := strings.Split(string(content), "\r\n")
	var base string
	codes := map[string]string{}
	for i, l := range readLines {
		if i == 0 {
			base = strings.TrimSpace(l)
			continue
		}
		if strings.Contains(l, " -> ") {
			code := strings.Split(l, " -> ")
			codes[code[0]] = code[1]
		}
	}
	return base, codes, nil
}

func apply(base string, codes map[string]string, round int) string {
	var sb strings.Builder
	sb.WriteString(base)
	for range round {
		base = sb.String()
		sb.Reset()
		for i := range base {
			sb.WriteByte(base[i])
			if i < len(base)-1 {
				sb.WriteString(codes[string(base[i:i+2])])
			}
		}
	}
	return sb.String()
}

func TestDay14Expansion(t *testing.T) {
	base, codes, err := load("light.txt")
	require.NoError(t, err)
	assert.Equal(t, "NCNBCHB", apply(base, codes, 1))
	assert.Equal(t, "NBCCNBBBCBHCB", apply(base, codes, 2))
	assert.Equal(t, "NBBBCNCCNBBNBNBBCHBHHBCHB", apply(base, codes, 3))
	assert.Equal(t, "NBBNBNBBCCNBCNCCNBBNBBNBBBNBBNBBCBHCBHHNHCBBCBHCB", apply(base, codes, 4))
}

type charCount struct {
	c     byte
	count int
}

func charsCount(value string) []charCount {
	chars := make([]charCount, 0, len(value))
source:
	for i := range value {
		cur := value[i]
		for c := range chars {
			if chars[c].c == cur {
				chars[c].count++
				continue source
			}
		}
		chars = append(chars, charCount{value[i], 1})
	}
	return chars
}

func merge(chars *[]charCount, newChars []charCount) {
source:
	for _, nc := range newChars {
		for c := range *chars {
			if (*chars)[c].c == nc.c {
				(*chars)[c].count += nc.count
				continue source
			}
		}
		*chars = append(*chars, nc)
	}
}

func score(chars []charCount) int {
	sort.SliceStable(chars, func(i, j int) bool {
		return chars[i].count < chars[j].count
	})
	return chars[len(chars)-1].count - chars[0].count
}

type expansion struct {
	base   string
	result string
	chars  []charCount
}

func expansionRes(part string, codes map[string]string, expRound int) *expansion {
	var exp expansion
	exp.base = part
	exp.result = apply(part, codes, expRound)
	exp.chars = charsCount(exp.result[:len(exp.result)-1])
	return &exp
}

type cache map[int]map[string][]charCount

var cacheHit int

func cmpChars(base string, codes map[string]string, cachedResult cache, round, expRound int, chars *[]charCount) []charCount {
	root := chars == nil
	if root {
		chars = &[]charCount{}
	}
	if cachedResult[round] == nil {
		cachedResult[round] = map[string][]charCount{}
	}

	if cc := cachedResult[round][base]; cc != nil {
		cacheHit++
		merge(chars, cc)
		return nil
	}
	if round == 1 {
		var baseChars []charCount
		for i := range len(base) - 1 {
			part := base[i : i+2]
			cp := cachedResult[1][part]
			if len(cp) == 0 {
				xp := expansionRes(part, codes, expRound)
				cachedResult[1][part] = xp.chars
				cp = xp.chars
			} else {
				cacheHit++
			}
			merge(&baseChars, cp)
		}
		merge(chars, baseChars)
		return nil
	}

	for i := range len(base) - 1 {
		part := base[i : i+2]
		if len(cachedResult[round][part]) == 0 {
			xp := expansionRes(part, codes, expRound)
			var baseChars []charCount
			cmpChars(xp.result, codes, cachedResult, round-1, expRound, &baseChars)
			cachedResult[round][part] = baseChars
		} else {
			cacheHit++
		}
		merge(chars, cachedResult[round][part])
	}
	if root {
		merge(chars, []charCount{{base[len(base)-1], 1}})
	}
	return *chars
}

func TestDay14(t *testing.T) {
	base, codes, err := load("light.txt")
	require.NoError(t, err)
	assert.Equal(t, charsCount(apply("NN", codes, 10)), cmpChars("NN", codes, cache{}, 2, 5, nil))
	assert.Equal(t, 1588, score(cmpChars(base, codes, cache{}, 2, 5, nil)))
	assert.Equal(t, 2188189693529, score(cmpChars(base, codes, cache{}, 8, 5, nil)))

	base, codes, err = load("large.txt")
	require.NoError(t, err)
	assert.Equal(t, 3306, score(cmpChars(base, codes, cache{}, 2, 5, nil)))
	assert.Equal(t, 3760312702877, score(cmpChars(base, codes, cache{}, 8, 5, nil)))
}
