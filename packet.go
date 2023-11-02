package main

import (
	"encoding/binary"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"

	"time"
)

const (
	WS uint8 = iota + 1
	TCP
)

type Server struct {
	gnet.BuiltinEventEngine

	addr      string
	multicore bool
	eng       gnet.Engine
}

func (wss *Server) OnBoot(eng gnet.Engine) gnet.Action {
	wss.eng = eng
	logging.Infof("echo server with multi-core=%t is listening on %s", wss.multicore, wss.addr)
	return gnet.None
}

func (wss *Server) OnOpen(c gnet.Conn) ([]byte, gnet.Action) {
	ctx := new(Context)
	ctx.Id, _ = Sf.NextID()
	ctx.dp = &DuelPlayer{
		Name:     "",
		Type:     0xff,
		Protocol: TCP,
		Conn:     c,
	}
	c.SetContext(ctx)

	return nil, gnet.None
}

func (wss *Server) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	if err != nil {
		logging.Warnf("error occurred on connection=%s, %v\n", c.RemoteAddr().String(), err)
	}

	logging.Infof("conn[%v] disconnected", c.RemoteAddr().String())
	return gnet.None
}

func (wss *Server) OnTraffic(c gnet.Conn) (action gnet.Action) {
	ctx := c.Context().(*Context)
	n := c.InboundBuffered()
	if n == 0 {
		return gnet.None
	}
	if ctx.nextOp == readLen {
		arr, err := c.Peek(2)
		if err != nil {
			return gnet.Close
		}
		ctx.msgLen = int(binary.LittleEndian.Uint16(arr))
		ctx.nextOp = readMsg
	}
	if n-2 >= ctx.msgLen {
		arr, err := c.Next(ctx.msgLen + 2)
		if err != nil {
			return gnet.Close
		}
		HandleCTOSPacket(ctx.dp, arr[2:], ctx.msgLen)
		ctx.nextOp = readLen
		ctx.msgLen = 0
	}
	return gnet.None
}

func (wss *Server) OnTick() (delay time.Duration, action gnet.Action) {
	return 3 * time.Second, gnet.None
}

type tcpReadOp uint8

const (
	readLen tcpReadOp = iota
	readMsg
)

type Context struct {
	Id     uint64
	nextOp tcpReadOp
	msgLen int
	dp     *DuelPlayer
}
