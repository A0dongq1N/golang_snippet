package antssnippet

import (
	"sync"
	"testing"

	"github.com/panjf2000/ants/v2"
)

func TestCommonPool(t *testing.T) {
	wg := new(sync.WaitGroup)
	pool, _ := ants.NewPool(10)

	for i := 0; i < 100; i++ {
		i := i
		wg.Add(1)
		pool.Submit(func() {
			t.Log("i: ", i)
			wg.Done()
		})
	}
	wg.Wait()

	defer pool.Release()
}
