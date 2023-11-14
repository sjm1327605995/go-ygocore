package main

import (
	"bufio"
	"os"
	"testing"
)

func TestReplay(t *testing.T) {
	f, err := os.Open("hero.yrp3d")
	if err != nil {
		t.Error(err)
		return
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	arr := make([]byte, 20)
	n, err := reader.Read(arr)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(arr[:n])
}
