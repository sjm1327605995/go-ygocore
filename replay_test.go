package main

import (
	"bufio"
	"encoding/binary"
	"encoding/hex"
	"os"
	"testing"
)

func TestReplay(t *testing.T) {
	h, _ := hex.DecodeString("140001040005401f0000401f000028000f0028000f00")
	t.Log(h)
	f, err := os.Open("hero.yrp3d")
	if err != nil {
		t.Error(err)
		return
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	var tp uint8
	err = binary.Read(reader, binary.LittleEndian, &tp)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(tp)
	var length uint32
	err = binary.Read(reader, binary.LittleEndian, &length)
	if err != nil {

	}
	t.Log(length)
	var arr = make([]byte, length)
	n, err := reader.Read(arr)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(arr[:n])
}
