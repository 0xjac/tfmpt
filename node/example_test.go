package node_test

import (
	"fmt"

	"go.0xjac.com/tfmpt/node"
)

func ExampleDecode() {
	enc := []byte{0xc8, 0x83, 0x6b, 0x65, 0x79, 0x83, 0x76, 0x61, 0x6c}
	n, _ := node.Decode(enc, nil)
	fmt.Printf("%s\n", n)
	// Output:
	// [06 05 07 09 10, val]
}
