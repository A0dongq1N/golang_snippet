package antssnippet

import (
	"fmt"
	"sync"
	"testing"

	"github.com/panjf2000/ants/v2"
)

/*
这个测试文件中的测试用例为什么都测试结果中都没有输出大量重复的i值？

原因：
Go 1.22 将“循环迭代变量”改为“每轮重新声明（redeclare）并重新绑定（rebind）”，彻底消除闭包捕获同一地址的历史陷阱；
语义直观、并发更安全，性能无显著回退（仅在闭包捕获时才有必要的逃逸）。


*/

func PrintTask(i int, wg *sync.WaitGroup, t *testing.T) func() {
	return func() {
		t.Logf("i=%d", i)
		wg.Done()
	}
}

func TestCommonPool(t *testing.T) {
	wg := new(sync.WaitGroup)
	pool, _ := ants.NewPool(10)

	for i := 0; i < 100; i++ {
		wg.Add(1)
		err := pool.Submit(PrintTask(i, wg, t))
		if err != nil {
			return
		}
	}
	wg.Wait()

	defer pool.Release()
}

//-------------------------------------

func TestCommonPool2(t *testing.T) {
	wg := new(sync.WaitGroup)
	pool, _ := ants.NewPool(10)

	for i := 0; i < 100; i++ {
		wg.Add(1)
		err := pool.Submit(func() {
			t.Logf("i=%d", i)
			wg.Done()
		})
		if err != nil {
			return
		}
	}
	wg.Wait()

	defer pool.Release()
}

func TestCommonPool3(t *testing.T) {
	wg := new(sync.WaitGroup)
	pool, _ := ants.NewPool(10)
	defer pool.Release()

	for i := 0; i < 100; i++ {
		wg.Add(1)
		err := pool.Submit(func() {
			fmt.Printf("i=%d\n", i)
			wg.Done()
		})
		if err != nil {
			return
		}
	}
	wg.Wait()

}

func TestCommonPool4(t *testing.T) {
	wg := new(sync.WaitGroup)

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			fmt.Printf("i=%d\n", i)
			wg.Done()
		}()
	}
	wg.Wait()

}

func BenchmarkCommonPool(b *testing.B) {
	pool, _ := ants.NewPool(10)
	defer pool.Release()

	for n := 0; n < b.N; n++ {
		wg := new(sync.WaitGroup)
		for i := 0; i < 100; i++ {
			wg.Add(1)
			err := pool.Submit(func() {
				fmt.Printf("i=%d\n", i)
				wg.Done()
			})
			if err != nil {
				return
			}
		}
		wg.Wait()
	}

}
