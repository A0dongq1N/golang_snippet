package designpatternsnippet

import (
	"fmt"
	"testing"
)

type ClothesShop struct{}

func (cs *ClothesShop) OnShop() {
	fmt.Println("休闲的装扮")
}

type ClothesWork struct{}

func (cw *ClothesWork) OnWork() {
	fmt.Println("工作的装扮")
}

func TestSRP(t *testing.T) {
	//第一种创建方式
	cs1 := ClothesShop{}
	cs1.OnShop()

	cw1 := ClothesWork{}
	cw1.OnWork()

	//第二种创建方式
	cs2 := &ClothesShop{}
	cs2.OnShop()

	cw2 := &ClothesWork{}
	cw2.OnWork()

	//第三种创建方式
	cs3 := new(ClothesShop)
	cs3.OnShop()

	cw3 := new(ClothesWork)
	cw3.OnWork()

}
