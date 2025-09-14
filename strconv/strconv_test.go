package strconvsnippet

import (
	"fmt"
	"strconv"
	"testing"
)

func TestStrconv(t *testing.T) {
	//字符串变为浮点数
	f, err := strconv.ParseFloat("3.1415926", 2)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(f)

}
