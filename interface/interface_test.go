package interfacesnippet

import "testing"

type I interface {
	Get() int
	Set(int)
}

type S struct {
	Age int
}

type R struct {
	Age int
}

func (s S) Get() int {
	return s.Age
}

func (s R) Get() int {
	return s.Age
}

func (s *S) Set(age int) {
	s.Age = age
}

func (s *R) Set(age int) {
	s.Age = age
}

func TestInterface(t *testing.T) {
	//如何判断 interface 变量存储的是哪种类型
	s := S{1}
	var i I = &s
	t.Log("i.Get()=", i.Get()) // 1
	if _, ok := i.(*S); ok {
		t.Log("i impl S") // i impl S
	}

}

func TestInterface2(t *testing.T) {
	//如何判断 interface 变量存储的是哪种类型
	s := S{1}
	var i I = &s
	switch i.(type) {
	case *S:
		t.Log("*S impl") // *S impl
	case *R:
		t.Log("*R impl")
	}
}
