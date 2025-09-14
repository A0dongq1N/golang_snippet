package leetcodesnippet

import (
	"fmt"
	"sort"
	"testing"
)

func TestLargestNumber(t *testing.T) {
	fmt.Println(largestNumber([]int{3, 30, 34, 5, 9}))
}

func largestNumber(nums []int) string {
	//还是不明白为什么22比较可以得到最大的值
	sort.Slice(nums, func(i, j int) bool {
		a, b := nums[i], nums[j]
		sa, sb := fmt.Sprintf("%d%d", a, b), fmt.Sprintf("%d%d", b, a)
		return sa > sb
	})
	if nums[0] == 0 {
		return "0"
	}
	res := ""
	for _, v := range nums {
		res += fmt.Sprintf("%d", v)
	}
	return res
}
