package antssnippet

import (
	"fmt"
	"testing"

	"github.com/panjf2000/ants/v2"
)

func TestPoolWithFunction(t *testing.T) {
	pool, _ := ants.NewPoolWithFunc(10, func(i interface{}) { fmt.Println(i.(int)) })
	defer pool.Release()

	for i := 0; i < 100; i++ {
		err := pool.Invoke(i)
		if err != nil {
			fmt.Println("Error invoking pool:", err)
		}
	}
}
