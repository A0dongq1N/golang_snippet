package antssnippet

import (
	"fmt"
	"testing"
	"time"

	"github.com/panjf2000/ants/v2"
)

func TestMultiPoolWithFunc(t *testing.T) {
	pool, _ := ants.NewMultiPoolWithFunc(10, 20, func(i interface{}) {
		time.Sleep(2 * time.Second)
		fmt.Println(i.(int))
	}, ants.LeastTasks)
	for i := 0; i < 100; i++ {
		pool.Invoke(i)
	}

	err := pool.ReleaseTimeout(2 * time.Second)
	if err != nil {
		// /data/workspace/golang_snippet/ants_snippet/multipoolwithfunc_test.go:22: ReleaseTimeout error: pool 3: operation timed out | pool 6: operation timed out | pool 0: operation timed out | pool 1: operation timed out | pool 7: operation timed out | pool 9: operation timed out | pool 4: operation timed out | pool 5: operation timed out | pool 8: operation timed out | pool 2: operation timed out
		t.Error("ReleaseTimeout error:", err)
	}
}
