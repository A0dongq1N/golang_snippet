package syncsnippet

import (
	"fmt"
	"sync/atomic"
	"testing"
)

func TestSyncAtomic(t *testing.T) {
	var i int32 = 0
	atomic.AddInt32(&i, 1)
	fmt.Println(i)
}

func TestSyncAtomicLoad(t *testing.T) {
	var c atomic.Int64
	c.Add(1)
	n := c.Load()
	c.CompareAndSwap(n, n+10)

	t.Logf("n: %d\n", n)        // 1
	t.Logf("c: %d\n", c.Load()) //11
}

func TestSyncAtomicStore(t *testing.T) {
	var c atomic.Int64
	c.Add(1)
	n := c.Load()
	c.Store(0)
	c.CompareAndSwap(n, n+10) //此时,c=0, n=1, CAS不成立

	t.Logf("n: %d\n", n)        // n=1
	t.Logf("c: %d\n", c.Load()) // c=0
}
