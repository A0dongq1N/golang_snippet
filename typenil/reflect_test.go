package typenilsnippet

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
)

func TestReflect(t *testing.T) {
	var p *int = nil
	v := reflect.ValueOf(p)
	fmt.Println(v)                        //nil
	fmt.Println(reflect.TypeOf(v))        //reflect.value
	fmt.Println(v)                        //nil
	fmt.Println(reflect.TypeOf(p))        //*int
	fmt.Println(reflect.TypeOf(p).Kind()) //ptr
	fmt.Println(v.IsNil())                // true，正确判断

	var i interface{} = p
	v2 := reflect.ValueOf(i)
	fmt.Println(v2)         // nil
	fmt.Println(v2.IsNil()) //true

}

func TestTypeAssert(t *testing.T) {
	var p *int = nil
	var i interface{} = p
	if v, ok := i.(*int); ok {
		fmt.Println(v) // 这里v是*int类型的nil，但ok为true
	} else {
		fmt.Println("类型断言失败")
	}

}

func worker(m *sync.Mutex) {
	m.Lock()
	defer m.Unlock()
	// ...
}

func TestSync(t *testing.T) {
	var wg sync.WaitGroup
	var m *sync.Mutex = nil
	wg.Add(1)
	go worker(m)
	wg.Wait()
}
