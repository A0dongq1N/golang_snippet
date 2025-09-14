package leetcodesnippet

import (
	"fmt"
	"strings"
	"testing"
)

func TestHammingWeight(t *testing.T) {
	t.Log(hammingWeight(132) == hammingWeight2(132))
	t.Log(hammingWeight(11) == hammingWeight3(11))
	t.Log(hammingWeight(11) == hammingWeight2(11))
	t.Log(hammingWeight(14) == hammingWeight3(14))
	t.Log(hammingWeight(6) == hammingWeight2(6))
	t.Log(hammingWeight(14) == hammingWeight2(14))
	t.Log(hammingWeight(6) == hammingWeight3(6))
}

func hammingWeight(num uint32) int {
	count := 0
	for num != 0 {
		count++
		num &= num - 1
	}
	return count
}

func hammingWeight2(num uint32) int {
	count := 0
	for num != 0 {
		if num&1 == 1 {
			count++
		}
		num >>= 1
	}
	return count
}

func hammingWeight3(num uint32) int {
	count := strings.Count(fmt.Sprintf("%b", num), "1")

	return count
}
