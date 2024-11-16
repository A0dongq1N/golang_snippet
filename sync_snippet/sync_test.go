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
