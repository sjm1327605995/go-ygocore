package main

import (
	"encoding/hex"
	"testing"
)

func TestHex(t *testing.T) {
	//36 0 1 4 0 0 0 4 0 0 0 4 0 0 0 4 0 0 0 4 0 0 0 4 0 0 0 4 0 0 0 4 0 0 0 64 1 8

	//36 0 1 6 0 8 4 0 0 0 4 0 0 0 4 0 0 0 4 0 0 0 4 0 0 0 4 0 0 0 4 0 0 0 4 0 0 0
	arr, _ := hex.DecodeString("2b0020de56c65f847661534c720000227f00003300000000000000800d0018227f000000000000000000000000")
	t.Log(arr)
	t.Log(0x20)
}
