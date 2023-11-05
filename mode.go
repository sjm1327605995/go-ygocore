package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/panjf2000/gnet/v2"
)

type DuelPlayer struct {
	Name     [40]byte //40 byte
	Type     uint16
	Status   uint8
	Protocol uint8
	Pass     string
	Conn     gnet.Conn
	game     DuelMode
	Pos      uint8 // 0 玩家1  1 玩家2
}

func (d DuelPlayer) Write(arr []byte) error {
	switch d.Protocol {
	case TCP:
		_, err := d.Conn.Write(arr)
		return err
	case WS:
	}
	return nil
}

type DuelMode interface {
	Chat(dp *DuelPlayer, buff []byte)
	JoinGame(dp *DuelPlayer, buff []byte, isCreator bool)
	LeaveGame(dp *DuelPlayer)
	ToObserver(dp *DuelPlayer)
	PlayerReady(dp *DuelPlayer, isReady bool)
	PlayerKick(dp *DuelPlayer, pos uint8)
	UpdateDeck(dp *DuelPlayer, buff []byte) error
	StartDuel(dp *DuelPlayer)
	HandResult(dp *DuelPlayer, uint82 uint8)
	TPResult(dp *DuelPlayer, uint82 uint8)
	Process()
	Analyze(buff []byte) int
	Surrender(dp *DuelPlayer)
	GetResponse(dp *DuelPlayer, buff []byte)
	TimeConfirm(dp *DuelPlayer)
	EndDuel()
	PDuel() int64
	IsCreator(dp *DuelPlayer) bool
}
type DuelModeBase struct {
	//Etimer
	HostPlayer *DuelPlayer
	DuelStage  int32
	pDuel      int64
	Name       string //40个字节
	Pass       string //40个字节
}

func (d *DuelModeBase) IsCreator(dp *DuelPlayer) bool {
	//TODO implement me
	panic("implement me")
}

func (d *DuelModeBase) Chat(dp *DuelPlayer, buff []byte) {
	//TODO implement me
	panic("implement me")
}

func (d *DuelModeBase) JoinGame(dp *DuelPlayer, buff []byte, flag bool) {
	//TODO implement me
	panic("implement me")
}

func (d *DuelModeBase) LeaveGame(dp *DuelPlayer) {
	//TODO implement me
	panic("implement me")
}

func (d *DuelModeBase) ToObserver(dp *DuelPlayer) {
	//TODO implement me
	panic("implement me")
}

func (d *DuelModeBase) PlayerReady(dp *DuelPlayer, isReady bool) {
	//TODO implement me
	panic("implement me")
}

func (d *DuelModeBase) PlayerKick(dp *DuelPlayer, pos uint8) {
	//TODO implement me
	panic("implement me")
}

func (d *DuelModeBase) UpdateDeck(dp *DuelPlayer, buff []byte) error {
	//TODO implement me
	panic("implement me")
}

func (d *DuelModeBase) StartDuel(dp *DuelPlayer) {
	//TODO implement me
	panic("implement me")
}

func (d *DuelModeBase) HandResult(dp *DuelPlayer, uint82 uint8) {
	//TODO implement me
	panic("implement me")
}

func (d *DuelModeBase) TPResult(dp *DuelPlayer, uint82 uint8) {
	//TODO implement me
	panic("implement me")
}

func (d *DuelModeBase) Process() {
	//TODO implement me
	panic("implement me")
}

func (d *DuelModeBase) Analyze(buff []byte) int {
	//TODO implement me
	panic("implement me")
}

func (d *DuelModeBase) Surrender(dp *DuelPlayer) {
	//TODO implement me
	panic("implement me")
}

func (d *DuelModeBase) GetResponse(dp *DuelPlayer, buff []byte) {
	//TODO implement me
	panic("implement me")
}

func (d *DuelModeBase) TimeConfirm(dp *DuelPlayer) {
	//TODO implement me
	panic("implement me")
}

func (d *DuelModeBase) EndDuel() {
	//TODO implement me
	panic("implement me")
}
func (d *DuelModeBase) PDuel() int64 {
	return 0
}

type ParseMessage interface {
	Parse(*bytes.Buffer) error
}
type BytesMessage interface {
	ToBytes(*bytes.Buffer) error
}

func (d *DuelModeBase) Write(dp *DuelPlayer, proto uint8, msg BytesMessage) error {
	buffer := bytes.NewBuffer(make([]byte, 3, 100))
	err := msg.ToBytes(buffer)
	if err != nil {
		fmt.Println(err)
		return err
	}
	arr := buffer.Bytes()
	binary.LittleEndian.PutUint16(arr, uint16(len(arr)-2))
	arr[2] = proto
	switch dp.Protocol {
	//Websocket
	case WS:
		return wsutil.WriteServerMessage(dp.Conn, ws.OpBinary, arr)
	//TCP
	case TCP:
		_, err := dp.Conn.Write(arr)
		return err
	}
	return nil
}

type BytesMsg []byte

func (c *BytesMsg) ToBytes(buff *bytes.Buffer) error {
	_, err := buff.Write(*c)
	return err

}