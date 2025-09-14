package antssnippet

import (
	"testing"
	"time"

	"github.com/panjf2000/ants/v2"
)

func TestMultiPool(t *testing.T) {
	pool, _ := ants.NewMultiPool(10, 20, ants.RoundRobin)

	for i := 0; i < 10; i++ {
		i := i
		pool.Submit(func() {
			time.Sleep(1 * time.Second)
			t.Log("i: ", i)
		})
	}

	err := pool.ReleaseTimeout(2 * time.Second)
	if err != nil {
		// ReleaseTimeout error: pool 9: operation timed out | pool 4: operation timed out | pool 0: operation timed out | pool 1: operation timed out | pool 5: operation timed out | pool 3: operation timed out | pool 7: operation timed out | pool 2: operation timed out | pool 6: operation timed out | pool 8: operation timed out
		t.Error("ReleaseTimeout error:", err)
	}
}
