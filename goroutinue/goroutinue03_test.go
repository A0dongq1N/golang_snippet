package goroutinuesnippet

import (
	"sync"
	"testing"
)

/*
题目要求：启动2个协程交替打印0,1,2,3,4....100

额外要求：打印的值要求通过chan在2个协程之间传递

*/

func TestGoroutine03(t *testing.T) {
	wg := sync.WaitGroup{}
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)

	wg.Add(1)
	ch1 <- 0
	go func() {
		defer wg.Done()
		for {
			tmp := <-ch1
			t.Logf("i=%d", tmp)

			if tmp > 99 {
				ch2 <- tmp + 1
				break
			}
			ch2 <- tmp + 1

		}

	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			tmp := <-ch2
			t.Logf("i=%d", tmp)

			if tmp > 98 {
				ch1 <- tmp + 1
				break
			}
			ch1 <- tmp + 1

		}

	}()

	wg.Wait()

}
