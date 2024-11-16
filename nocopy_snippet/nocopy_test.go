package nocopysnippet

import (
	"fmt"
	"testing"
)

type noCopy struct{}

func (n *noCopy) Lock()   {}
func (n *noCopy) Unlock() {}

type cool struct {
	noCopy
	val int32
}

func TestNoCopy(t *testing.T) {
	var c cool = cool{val: 1}
	fmt.Println("c.val:", c.val)

	var d cool
	d = c
	d.val = 2
	fmt.Println("d.val:", d.val)
}

// [root@jeremyqin-1erdn1fxga nocopy_snippet]# go vet nocopy_test.go
// # command-line-arguments
// ./nocopy_test.go:23:6: assignment copies lock value to d: command-line-arguments.cool
