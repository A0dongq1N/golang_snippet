package designpatternsnippet

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
)

var lock sync.Mutex
var once sync.Once

type Singleton struct{}

var sig *Singleton

func GetInstance1() *Singleton {
	lock.Lock()
	defer lock.Unlock()
	if sig == nil {
		sig = new(Singleton)
	}
	return sig
}

func GetInstance2() *Singleton {
	once.Do(func() {
		sig = new(Singleton)
	})
	return sig
}

func (s *Singleton) Show() {
	fmt.Println("singleton show")
}

func TestGetInstance(t *testing.T) {
	ins1 := GetInstance1()
	ins1.Show()

	ins2 := GetInstance2()
	ins2.Show()

	fmt.Println(ins1 == ins2)
	fmt.Println(reflect.DeepEqual(ins1, ins2))

}
