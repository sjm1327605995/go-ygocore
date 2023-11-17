package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/itchio/lzma"
	"io"
	"os"
	"testing"
	"time"
)

func TestReplay(t *testing.T) {
	t.Log(time.Unix(100667136, 0))
	f, err := os.Open("test.yrp")
	if err != nil {
		t.Error(err)
		return
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	//arr := make([]byte, 50)
	//n, err := reader.Read(arr)
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	t.Log(0x31707279)
	var rh ReplayHeader
	err = binary.Read(reader, binary.LittleEndian, &rh)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(rh)
	b := bytes.NewBuffer(make([]byte, 0, 1024*10))
	rh.GetLzmaHeader(b)
	io.Copy(b, reader)
	t.Log(b.Len())
	r := lzma.NewReader(b)
	defer r.Close()
	decodeReader := bytes.NewBuffer(make([]byte, 0, 1024*10))
	io.Copy(decodeReader, r)
	t.Log(decodeReader.Len())
	decodeReader.Next(40)
	decodeReader.Next(40)
	var (
		startLp   int32
		startHand int32
		drawCount int32
		opt       int32
	)
	binary.Read(decodeReader, binary.LittleEndian, &startLp)
	binary.Read(decodeReader, binary.LittleEndian, &startHand)
	binary.Read(decodeReader, binary.LittleEndian, &drawCount)
	binary.Read(decodeReader, binary.LittleEndian, &opt)
	var replay Replay
	replay.Header = &rh
	replay.readDeck(decodeReader)
}

type Replay struct {
	Header      *ReplayHeader
	HostName    string
	ClientName  string
	StartLp     int32
	StartHand   int32
	DrawCount   int32
	Opt         int32
	HostDeck    []int32
	ClientDeck  []int32
	TagHostName string
}

func (r *Replay) readDeck(reader *bytes.Buffer) {
	m := r.readDeckPack(reader)
	ex := r.readDeckPack(reader)
	fmt.Println(m, ex)
}
func (r *Replay) readDeckPack(reader *bytes.Buffer) []int32 {
	var length int32
	binary.Read(reader, binary.LittleEndian, &length)
	var (
		list []int32
		i, j int32
		ref  = length
	)
	i = 1
	j = 1
	for {
		if 1 <= ref {
			if j > ref {
				break
			}
		} else {
			if j < ref {
				break
			}
		}

		var code int32
		binary.Read(reader, binary.LittleEndian, &code)
		list = append(list, code)
		i = 1
		if i <= ref {
			j++
		} else {
			j--
		}
	}
	return list
}
