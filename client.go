package main

import (
	"encoding/binary"
	"fmt"
)

type ClientInterface interface {
	Write([]byte) error
}

type Packet interface {
	Marshal() ([]byte, error)
}
type BytesPacket []byte

func (b BytesPacket) Marshal() ([]byte, error) {
	return b, nil
}
func SendPacketToPlayer(c ClientInterface, MSG uint8, packet Packet) error {
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
	err := c.Write(buff)
	if err != nil {
		return err
	}
	for i := range resend {
		err = resend[i].Write(buff)
		if err != nil {
			return err
		}
	}
	return nil
}

type ConsoleClient struct {
	id int
}

func (c *ConsoleClient) Write(arr []byte) error {
	fmt.Println("console ", c.id, ":", arr)
	return nil
}
