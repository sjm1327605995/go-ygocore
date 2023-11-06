package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"
)

type ClientInterface interface {
	Write([]byte) error
	Close() error
}

type Packet interface {
	Marshal() ([]byte, error)
}
type BytesPacket []byte

func (b BytesPacket) Marshal() ([]byte, error) {
	return b, nil
}
func SendPacketToPlayer(c ClientInterface, MSG uint8, packet Packet) error {
	if packet == nil {
		return SendBufferToPlayer(c, MSG, make([]byte, 3))
	}
	bytes, err := packet.Marshal()
	if err != nil {
		return err
	}
	buff := make([]byte, len(bytes)+3)
	copy(buff[3:], bytes)

	return SendBufferToPlayer(c, MSG, buff)
}
func SendBufferToPlayer(c ClientInterface, MSG uint8, buff []byte, resend ...ClientInterface) error {

	binary.LittleEndian.PutUint16(buff, uint16(len(buff)-2))
	buff[2] = MSG
	fmt.Println(MSG, hex.EncodeToString(buff))
	err := c.Write(buff)
	if err != nil {
		return err
	}
	for i := range resend {
		fmt.Println(MSG, hex.EncodeToString(buff))
		err = resend[i].Write(buff)
		if err != nil {
			return err
		}
	}
	return nil
}
func init() {
	os.Remove("hex")
	os.Create("hex")
}

type ConsoleClient struct {
	id int
}

func (c *ConsoleClient) Write(arr []byte) error {
	fmt.Println("console ", c.id, ":", arr)
	if c.id == 1 {
		return nil
	}
	f, err := os.OpenFile("hex", os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer f.Close()

	f.Write([]byte(hex.EncodeToString(arr) + "\n"))
	return nil
}
func (c *ConsoleClient) Close() error {
	return nil
}
func DisconnectPlayer(c ClientInterface) error {
	return c.Close()
}
