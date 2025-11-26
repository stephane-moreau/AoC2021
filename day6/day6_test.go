package day6

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	light = []int{3, 4, 3, 1, 2}
	large = []int{1, 1, 3, 5, 1, 3, 2, 1, 5, 3, 1, 4, 4, 4, 1, 1, 1, 3, 1, 4, 3, 1, 2, 2, 2, 4, 1, 1, 5, 5, 4, 3, 1, 1, 1, 1, 1, 1, 3, 4, 1, 2, 2, 5, 1, 3, 5, 1, 3, 2, 5, 2, 2, 4, 1, 1, 1, 4, 3, 3, 3, 1, 1, 1, 1, 3, 1, 3, 3, 4, 4, 1, 1, 5, 4, 2, 2, 5, 4, 5, 2, 5, 1, 4, 2, 1, 5, 5, 5, 4, 3, 1, 1, 4, 1, 1, 3, 1, 3, 4, 1, 1, 2, 4, 2, 1, 1, 2, 3, 1, 1, 1, 4, 1, 3, 5, 5, 5, 5, 1, 2, 2, 1, 3, 1, 2, 5, 1, 4, 4, 5, 5, 4, 1, 1, 3, 3, 1, 5, 1, 1, 4, 1, 3, 3, 2, 4, 2, 4, 1, 5, 5, 1, 2, 5, 1, 5, 4, 3, 1, 1, 1, 5, 4, 1, 1, 4, 1, 2, 3, 1, 3, 5, 1, 1, 1, 2, 4, 5, 5, 5, 4, 1, 4, 1, 4, 1, 1, 1, 1, 1, 5, 2, 1, 1, 1, 1, 2, 3, 1, 4, 5, 5, 2, 4, 1, 5, 1, 3, 1, 4, 1, 1, 1, 4, 2, 3, 2, 3, 1, 5, 2, 1, 1, 4, 2, 1, 1, 5, 1, 4, 1, 1, 5, 5, 4, 3, 5, 1, 4, 3, 4, 4, 5, 1, 1, 1, 2, 1, 1, 2, 1, 1, 3, 2, 4, 5, 3, 5, 1, 2, 2, 2, 5, 1, 2, 5, 3, 5, 1, 1, 4, 5, 2, 1, 4, 1, 5, 2, 1, 1, 2, 5, 4, 1, 3, 5, 3, 1, 1, 3, 1, 4, 4, 2, 2, 4, 3, 1, 1}

	debug = false
)

func debugLog(s string, vals ...any) {
	if debug {
		fmt.Printf(s, vals...)
	}
}
func nextGen(fishes []int) []int {
	nextgen := fishes
	for i, f := range fishes {
		if f == 0 {
			nextgen[i] = 6
			nextgen = append(nextgen, 8)
		} else {
			nextgen[i]--
		}
	}
	return nextgen
}

func countFishes(gen int, fishes []int) int {
	cache := map[int]int{}
	res := 0
	for _, fish := range fishes {
		if cache[fish] == 0 {
			curFishes := []int{fish}
			for i := range gen {
				curFishes = nextGen(curFishes)
				debugLog("% 3d: %v\n", gen-i, curFishes)
			}
			cache[fish] = len(curFishes)
		}
		res += cache[fish]
	}
	return res
}

func TestDay6Phase1(t *testing.T) {
	require.Equal(t, 26, countFishes(18, light))
	require.Equal(t, 5934, countFishes(80, light))

	assert.Equal(t, 360610, countFishes(80, large))
}

const (
	firstRegenDelta = 2
	regularGen      = 7
)

var cachedSubGenCount = make(map[int]int)

func computeFishCount(gen int, initial int, root bool) int {
	if gen < initial {
		return 0
	}
	count := gen / regularGen
	if gen-count*regularGen > initial {
		count++
	} else if gen-count*regularGen == 0 {
		count--
	}
	if root {
		count++
	}
	debugLog("%d / %d: %d\n", gen, initial, count)
	if !root {
		initial++
	}
	if initial == regularGen {
		if cachedSubGenCount[gen] != 0 {
			return cachedSubGenCount[gen]
		}
	}
	for curGen := range count {
		count += computeFishCount(gen-initial-curGen*regularGen-firstRegenDelta, 6, false)
	}
	if initial == regularGen {
		cachedSubGenCount[gen] = count
	}
	return count
}

func computeFishesCount(gen int, fishes []int) int {
	var res int
	cache := map[int]int{}
	for _, f := range fishes {
		if cache[f] == 0 {
			cache[f] = computeFishCount(gen, f, true)
		}
		res += cache[f]
	}
	return res
}

func TestDay6Phase2(t *testing.T) {

	require.Equal(t, countFishes(18, []int{3}), computeFishCount(18, 3, true))
	require.Equal(t, countFishes(18, []int{4}), computeFishCount(18, 4, true))
	require.Equal(t, countFishes(18, []int{1}), computeFishCount(18, 1, true))
	require.Equal(t, countFishes(18, []int{2}), computeFishCount(18, 2, true))
	require.Equal(t, countFishes(80, []int{3}), computeFishCount(80, 3, true))

	require.Equal(t, 26, computeFishesCount(18, light))
	require.Equal(t, 5934, computeFishesCount(80, light))
	require.Equal(t, 26984457539, computeFishesCount(256, light))

	assert.Equal(t, 360610, computeFishesCount(80, large))
	assert.Equal(t, 1631629590423, computeFishesCount(256, large))
}
