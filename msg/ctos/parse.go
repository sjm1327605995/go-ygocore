package ctos

import (
	"bytes"
	"encoding/binary"
	"github.com/sjm1327605995/go-ygocore/msg/host"
)

type PlayerInfo struct {
	Name []byte
}

const (
	StrLimit = 40
)

func (p *PlayerInfo) Parse(buff []byte) (err error) {
	if len(p.Name) == 0 {
		p.Name = make([]byte, StrLimit)
	}
	reader := bytes.NewReader(buff)
	// 将二进制数组转换为字符串
	return binary.Read(reader, binary.LittleEndian, p)
}

type TPResult struct {
	Res uint8
}

func (h *TPResult) Parse(buff []byte) (err error) {
	reader := bytes.NewReader(buff)
	return binary.Read(reader, binary.LittleEndian, &h)
}

type CreateGame struct {
	Info host.HostInfo
	Name string
	Pass string
}

func (h *CreateGame) Parse(buff []byte) (err error) {

	return nil
}

type JoinGame struct {
	Version uint16
	Align   uint16
	GameId  uint32
	Pass    string
}

// Pass: [40] - 房间密码
func (h *JoinGame) Parse(b *bytes.Buffer) (err error) {
	return nil
}

type Kick struct {
	Pos uint8
}

func (h *Kick) Parse(buff []byte) (err error) {
	reader := bytes.NewReader(buff)
	return binary.Read(reader, binary.LittleEndian, h)

}

type HandResult struct {
	Res uint8
}

func (h *HandResult) Parse(buff []byte) (err error) {
	reader := bytes.NewReader(buff)
	return binary.Read(reader, binary.LittleEndian, h)

}
