package syncsnippet

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

/*
TestRWMutex模拟了读锁的并发行为，同时写锁由于是后开始的，被读锁阻塞着的现象
*/

func TestRWMutex(t *testing.T) {
	var wg sync.WaitGroup
	var lock sync.RWMutex
	var counter int64
	var reading int64
	var i int

	for i = 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			lock.RLock()
			defer lock.RUnlock()
			atomic.AddInt64(&reading, 1)
			t.Logf("[reader %d] %v开始读，当前并发读数量=%d", id, time.Now().UnixMilli(), reading)
			time.Sleep(200 * time.Millisecond)
			_ = counter
		}(i)
	}

	time.Sleep(100 * time.Millisecond)
	wg.Add(1)
	go func(id int) {
		defer wg.Done()
		t.Log("[writer] 等待写锁（此时所有读完成后才能写）...")
		lock.Lock()
		defer lock.Unlock()
		t.Log("[writer] 获得写锁，开始写")
		counter++
		t.Log("[writer] 写完成，释放锁")

	}(i)

	wg.Wait()
	t.Logf("全部结束,counter=%v", counter)

}
