package goroutinuesnippet

import (
	"fmt"
	"sync"
	"testing"
)

/*
题目要求：启动2个协程，交替输出0,1,2,3,4,,,,100
*/

func TestGoroutinue02(t *testing.T) {
	ch1 := make(chan bool, 1)
	ch2 := make(chan bool, 1)
	wg := sync.WaitGroup{}
	wg.Add(2)

	// 启动第一个协程，打印奇数
	go func1(&wg, ch1, ch2)

	// 启动第二个协程，打印偶数
	go func2(&wg, ch1, ch2)

	// 启动交替打印
	ch1 <- true
	wg.Wait()
}

func func1(wg *sync.WaitGroup, ch1 chan bool, ch2 chan bool) {
	defer wg.Done()
	for i := 1; i <= 100; i += 2 {
		<-ch1 // 等待信号
		fmt.Printf("协程1: %d\n", i)
		if i < 100 {
			ch2 <- true // 通知协程2
		}
	}
}

func func2(wg *sync.WaitGroup, ch1 chan bool, ch2 chan bool) {
	defer wg.Done()
	for i := 2; i <= 100; i += 2 {
		<-ch2 // 等待信号
		fmt.Printf("协程2: %d\n", i)
		if i < 100 {
			ch1 <- true // 通知协程1
		}
	}
}
