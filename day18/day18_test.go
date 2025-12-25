package day18

import (
	"io"
	"math"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	INVALID_VALUE = -1
)

type number struct {
	pair        bool
	value       int
	left, right *number
}

func (n number) Dump() string {
	if !n.pair {
		return strconv.Itoa(n.value)
	}
	return "[" + n.left.Dump() + "," + n.right.Dump() + "]"
}

func parseNumber(l string, pos int) (*number, int) {
	for l[pos] == ' ' {
		pos++
	}
	if l[pos] == '[' {
		var res number
		res.pair = true
		res.left, pos = parseNumber(l, pos+1)
		if l[pos] != ',' {
			panic(" wrong parsing")
		}
		res.right, pos = parseNumber(l, pos+1)
		if l[pos] != ']' {
			panic(" wrong parsing")
		}
		return &res, pos + 1
	}
	v, err := strconv.Atoi(l[pos : pos+1])
	if err != nil {
		panic(err)
	}
	return &number{value: v}, pos + 1
}

func loadNumbers(fn string) ([]string, error) {
	f, err := os.OpenFile(fn, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	content, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\r\n")
	return lines, nil
}

func split(num *number) (*number, bool) {
	if !num.pair {
		if num.value >= 10 {
			return &number{true,
				0,
				&number{value: int(math.Floor(float64(num.value) / 2.))},
				&number{value: int(math.Ceil(float64(num.value) / 2.))}}, true
		}
		return num, false
	}

	var done bool
	num.left, done = split(num.left)
	if !done {
		num.right, done = split(num.right)
	}
	return num, done
}

func explode(num *number, depth int) (*number, bool, *number) {
	if !num.pair {
		return num, false, nil
	}
	if depth == 4 {
		return &number{}, true, num
	}
	var done bool
	var exploded *number
	num.left, done, exploded = explode(num.left, depth+1)
	if done {
		if exploded.right.value != INVALID_VALUE {
			if !num.right.pair {
				num.right.value += exploded.right.value
				exploded.right.value = INVALID_VALUE
			} else {
				n := num.right
				for n.left.pair {
					n = n.left
				}
				n.left.value += exploded.right.value
				exploded.right.value = INVALID_VALUE
			}
		}
		return num, done, exploded
	}
	num.right, done, exploded = explode(num.right, depth+1)
	if done {
		if exploded.left.value != INVALID_VALUE {
			if !num.left.pair {
				num.left.value += exploded.left.value
				exploded.left.value = INVALID_VALUE
			} else {
				n := num.left
				for n.right.pair {
					n = n.right
				}
				n.right.value += exploded.left.value
				exploded.left.value = INVALID_VALUE
			}
		}
		return num, done, exploded
	}
	return num, false, nil
}

func reduce(num *number) *number {
	chExp, chSplit := true, true
	for chExp || chSplit {
		num, chExp, _ = explode(num, 0)
		if !chExp {
			num, chSplit = split(num)
		}
	}
	return num
}

func add(left, right *number) *number {
	res := reduce(&number{true, 0, left, right})
	return res
}

func sum(nums []string) *number {
	res, _ := parseNumber(nums[0], 0)
	for i := 1; i < len(nums); i++ {
		other, _ := parseNumber(nums[i], 0)
		res = add(res, other)
	}
	return res
}

func TestAddition(t *testing.T) {
	testCases := []struct {
		nums []string
		res  string
	}{
		{
			[]string{"[[[[4,3],4],4],[7,[[8,4],9]]]", "[1,1]"},
			"[[[[0,7],4],[[7,8],[6,0]]],[8,1]]",
		},
		{
			[]string{
				"[1,1]",
				"[2,2]",
				"[3,3]",
				"[4,4]",
			},
			"[[[[1,1],[2,2]],[3,3]],[4,4]]",
		},
		{
			[]string{
				"[1,1]",
				"[2,2]",
				"[3,3]",
				"[4,4]",
				"[5,5]",
			},
			"[[[[3,0],[5,3]],[4,4]],[5,5]]",
		},
		{
			[]string{
				"[1,1]",
				"[2,2]",
				"[3,3]",
				"[4,4]",
				"[5,5]",
				"[6,6]",
			},
			"[[[[5,0],[7,4]],[5,5]],[6,6]]",
		},
	}

	for ndx, tc := range testCases {
		var res *number
		for i := range len(tc.nums) - 1 {
			l, _ := parseNumber(tc.nums[i], 0)
			if i > 0 {
				l = res
			}
			r, _ := parseNumber(tc.nums[i+1], 0)
			res = add(l, r)
		}
		require.Equal(t, tc.res, res.Dump(), "test case number %d", ndx)
	}
}

func TestAddition2(t *testing.T) {
	vals := []string{
		"[[[0,[4,5]],[0,0]],[[[4,5],[2,6]],[9,5]]]",
		"[7,[[[3,7],[4,3]],[[6,3],[8,8]]]]",
		"[[2,[[0,8],[3,4]]],[[[6,7],1],[7,[1,6]]]]",
		"[[[[2,4],7],[6,[0,5]]],[[[6,8],[2,8]],[[2,1],[4,5]]]]",
		"[7,[5,[[3,8],[1,4]]]]",
		"[[2,[2,2]],[8,[8,1]]]",
		"[2,9]",
		"[1,[[[9,3],9],[[9,0],[0,7]]]]",
		"[[[5,[7,4]],7],1]",
		"[[[[4,2],2],6],[8,7]]		",
	}
	steps := []string{
		"[[[[4,0],[5,4]],[[7,7],[6,0]]],[[8,[7,7]],[[7,9],[5,0]]]]",
		"[[[[6,7],[6,7]],[[7,7],[0,7]]],[[[8,7],[7,7]],[[8,8],[8,0]]]]",
		"[[[[7,0],[7,7]],[[7,7],[7,8]]],[[[7,7],[8,8]],[[7,7],[8,7]]]]",
		"[[[[7,7],[7,8]],[[9,5],[8,7]]],[[[6,8],[0,8]],[[9,9],[9,0]]]]",
		"[[[[6,6],[6,6]],[[6,0],[6,7]]],[[[7,7],[8,9]],[8,[8,1]]]]",
		"[[[[6,6],[7,7]],[[0,7],[7,7]]],[[[5,5],[5,6]],9]]",
		"[[[[7,8],[6,7]],[[6,8],[0,8]]],[[[7,7],[5,0]],[[5,5],[5,6]]]]",
		"[[[[7,7],[7,7]],[[8,7],[8,7]]],[[[7,0],[7,7]],9]]",
		"[[[[8,7],[7,7]],[[8,6],[7,7]]],[[[0,7],[6,6]],[8,7]]]",
	}
	var res *number
	for i := 0; i < len(vals)-1; i++ {
		l, _ := parseNumber(vals[i], 0)
		if i != 0 {
			l = res
		}
		r, _ := parseNumber(vals[i+1], 0)
		res = add(l, r)
		require.Equal(t, steps[i], res.Dump())
	}
}

func magnitude(num number) int {
	if !num.pair {
		return num.value
	}
	return 3*magnitude(*num.left) + 2*magnitude(*num.right)
}

func TestMagnitude(t *testing.T) {
	tc := []struct {
		num string
		res int
	}{
		{"[[1,2],[[3,4],5]]", 143},
		{"[[[[0,7],4],[[7,8],[6,0]]],[8,1]]", 1384},
		{"[[[[1,1],[2,2]],[3,3]],[4,4]]", 445},
		{"[[[[3,0],[5,3]],[4,4]],[5,5]]", 791},
		{"[[[[5,0],[7,4]],[5,5]],[6,6]]", 1137},
		{"[[[[8,7],[7,7]],[[8,6],[7,7]]],[[[0,7],[6,6]],[8,7]]]", 3488},
	}
	for _, c := range tc {
		n, _ := parseNumber(c.num, 0)
		require.Equal(t, c.res, magnitude(*n))
	}
}

func largestMag(nums []string) int {
	max := 0
	for i := range nums {
		for j := i + 1; j < len(nums); j++ {
			nI, _ := parseNumber(nums[i], 0)
			nJ, _ := parseNumber(nums[j], 0)
			ij := magnitude(*add(nI, nJ))
			nI, _ = parseNumber(nums[i], 0)
			nJ, _ = parseNumber(nums[j], 0)
			ji := magnitude(*add(nJ, nI))
			if ij > max {
				max = ij
			}
			if ji > max {
				max = ji
			}
		}
	}
	return max
}

func TestDay18(t *testing.T) {
	nums, err := loadNumbers("light.txt")
	require.NoError(t, err)
	res := sum(nums)
	assert.Equal(t, "[[[[6,6],[7,6]],[[7,7],[7,0]]],[[[7,7],[7,7]],[[7,8],[9,9]]]]", res.Dump())
	assert.Equal(t, 4140, magnitude(*res))
	assert.Equal(t, 3993, largestMag(nums))

	nums, err = loadNumbers("large.txt")
	require.NoError(t, err)
	res = sum(nums)
	assert.Equal(t, "[[[[7,7],[7,7]],[[7,0],[7,7]]],[[[8,8],[8,7]],[[7,8],[7,8]]]]", res.Dump())
	assert.Equal(t, 4289, magnitude(*res))
	assert.Equal(t, 4807, largestMag(nums))
}
