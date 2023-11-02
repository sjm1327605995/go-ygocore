package ctos

import (
	"bytes"
	"github.com/sjm1327605995/go-ygocore/msg/host"
)

type PlayerInfo struct {
	Name     string
	RealName []byte
}

const (
	StrLimit = 40
)

func (p *PlayerInfo) Parse(b *bytes.Buffer) (err error) {
	// 将二进制数组转换为字符串
	return nil
}

type TPResult struct {
	Res uint8
}

func (h *TPResult) Parse(b *bytes.Buffer) (err error) {
	//return utils.GetData(b, &h.Res)
	return nil
}

type CreateGame struct {
	Info host.HostInfo
	Name string
	Pass string
}

func (h *CreateGame) Parse(b *bytes.Buffer) (err error) {
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
	Pos uint16
}

func (h *Kick) Parse(b *bytes.Buffer) (err error) {
	return nil

}

type HandResult struct {
	Res uint8
}

func (h *HandResult) Parse(b *bytes.Buffer) (err error) {
	return nil

}
