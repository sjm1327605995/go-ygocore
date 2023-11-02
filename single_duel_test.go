package main

import (
	"encoding/hex"
	"testing"
)

func TestHex(t *testing.T) {
	//36 0 1 4 0 0 0 4 0 0 0 4 0 0 0 4 0 0 0 4 0 0 0 4 0 0 0 4 0 0 0 4 0 0 0 64 1 8

	//36 0 1 6 0 8 4 0 0 0 4 0 0 0 4 0 0 0 4 0 0 0 4 0 0 0 4 0 0 0 4 0 0 0 4 0 0 0
	arr, _ := hex.DecodeString("2400010600080400000004000000040000000400000004000000040000000400000004000000")
	t.Log(arr)
}
