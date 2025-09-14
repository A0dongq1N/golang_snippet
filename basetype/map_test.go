package basetypesnippet

import (
	"testing"
)

func TestMap(t *testing.T) {
	a := map[int64]int64{
		1: 1,
		2: 2,
		3: 3,
	}
	t.Logf("len(a): %v", len(a))

	for k, v := range a {
		t.Logf("k=%v, v=%v", k, v)
	}

}
