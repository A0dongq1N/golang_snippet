package goroutinuesnippet

import (
	"fmt"
	"sync"
	"testing"
)

// 1. 通过sync.WaitGroup等待goroutine结束
func TestGoroutinue01(t *testing.T) {
	var wg = new(sync.WaitGroup)
	wg.Add(1)
	go f(wg, 10)
	// var input string
	// fmt.Scanln(&input)
	wg.Wait()
}

func f(wg *sync.WaitGroup, n int) {

	for i := 0; i < 10; i++ {
		n++
		fmt.Println(n)
	}
	defer wg.Done()

}
