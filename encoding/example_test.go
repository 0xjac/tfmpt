// Copyright (C) 2024 Jacques Dafflon | 0xjac - All Rights Reserved

package encoding_test

import (
	"fmt"

	"go.0xjac.com/tfmpt/encoding"
)

func ExampleToHex() {
	key := []byte("key")
	fmt.Printf("% 0x\n", encoding.ToHex(key))
	fmt.Printf("% 0x\n", encoding.ToHex(nil))
	// Output:
	// 06 0b 06 05 07 09 10
	// 10
}

func ExampleCompact() {
	key := []byte{0x06, 0x0b, 0x06, 0x05, 0x07, 0x09, 0x10}
	fmt.Printf("%s\n", encoding.Compact(key))

	key = []byte{0x1, 0x10}
	fmt.Printf("%s\n", encoding.Compact(key))

	fmt.Printf("%#0x\n", encoding.Compact(nil))
	// Output:
	// key
	// 1
	// 0x00
}

func ExampleExpandToHex() {
	key := []byte("key")
	fmt.Printf("% 0x\n", encoding.ExpandToHex(key))
	fmt.Printf("% 0x\n", encoding.ExpandToHex([]byte{0x01, 0x2}))
	fmt.Printf("% 0x\n", encoding.ExpandToHex([]byte{0x31, 0x2}))
	fmt.Printf("% 0x\n", encoding.ExpandToHex(nil))
	// Output:
	// 06 05 07 09 10
	// 00 02
	// 01 00 02 10
	//
}

func ExampleHexKeyHasTerm() {
	fmt.Println(encoding.HexKeyHasTerm([]byte{0x1, 0x10}))
	fmt.Println(encoding.HexKeyHasTerm([]byte{0x10}))
	fmt.Println(encoding.HexKeyHasTerm([]byte{0x1}))
	fmt.Println(encoding.HexKeyHasTerm([]byte{0x1, 0x2}))
	fmt.Println(encoding.HexKeyHasTerm([]byte{}))
	// Output:
	// true
	// true
	// false
	// false
	// false
}

func ExampleCommonPrefixLen() {
	fmt.Println(encoding.CommonPrefixLen([]int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4, 5}))
	fmt.Println(encoding.CommonPrefixLen([]int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4}))
	fmt.Println(encoding.CommonPrefixLen([]int{1, 2, 3, 4, 5}, []int{1, 2, 3}))
	fmt.Println(encoding.CommonPrefixLen([]int{1, 2, 3, 4, 5}, []int{1, 2}))
	fmt.Println(encoding.CommonPrefixLen([]int{1, 2, 3, 4, 5}, []int{1}))
	fmt.Println(encoding.CommonPrefixLen([]int{1, 2, 3, 4, 5}, []int{}))
	fmt.Println(encoding.CommonPrefixLen([]int{1, 2, 3, 4, 5}, []int{6}))
	// Output:
	// 5
	// 4
	// 3
	// 2
	// 1
	// 0
	// 0
}
