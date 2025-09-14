package testingsnippet

import "testing"

// 被测函数示例
func double(input int) int { return input * 2 }

// 一个 TestXxx 中包含多个 t.Run 子测试（表驱动）
func TestDouble(t *testing.T) {
	testCases := []struct {
		name                string
		input               int
		expectedDoubleValue int
	}{
		{name: "base", input: 1, expectedDoubleValue: 2},
		{name: "zero", input: 0, expectedDoubleValue: 0},
		{name: "negative", input: -3, expectedDoubleValue: -6},
	}

	for _, testCase := range testCases {
		testCase := testCase // 避免并发/闭包捕获问题
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel() // 可选：子测试并行

			actual := double(testCase.input)
			if actual != testCase.expectedDoubleValue {
				t.Fatalf("double(%d) got=%d want=%d", testCase.input, actual, testCase.expectedDoubleValue)
			}
		})
	}
}
