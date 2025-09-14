package mathsnippet

import (
	"math"
	"testing"
)

func TestMath(t *testing.T) {
	var r float64
	r = math.Round(4.32)
	t.Logf("r: %v", r) //4

	r = math.Round(4.52)
	t.Logf("r: %v", r) //5

	r = math.Round(4.62)
	t.Logf("r: %v", r) //5
}
