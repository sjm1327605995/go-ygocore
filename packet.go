package main

import (
	"bytes"
	"encoding/binary"
	"github.com/panjf2000/ants/v2"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"sync"
	"time"
)

const (
	TCP uint8 = iota + 1
	WS
)
const poolSize = 10000

type Server struct {
	gnet.BuiltinEventEngine

	addr      string
	multicore bool
	eng       gnet.Engine
	goPool    *ants.Pool
	bytesPool *BytesPool
}

type BytesPool struct {
	pool *sync.Pool
}

func NewBytesPool() *BytesPool {
	return &BytesPool{pool: &sync.Pool{
		New: func() any {
			return bytes.NewBuffer(make([]byte, 0, 1024))
		},
	}}
}
func (b *BytesPool) Get() *bytes.Buffer {
	return b.pool.Get().(*bytes.Buffer)
}
func (b *BytesPool) Put(buffer *bytes.Buffer) {
	buffer.Reset()

	if buffer != nil || buffer.Cap() <= 1024 {
		b.pool.Put(buffer)
	}
	buffer = nil
}
func NewServer() *Server {
	var err error
	goPool, err := ants.NewPool(10000, ants.WithExpiryDuration(time.Second*5))
	if err != nil {
		panic(err)
	}
	return &Server{
		goPool:    goPool,
		bytesPool: NewBytesPool(),
	}
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
		Type:     0xff,
		Protocol: TCP,
		Conn:     c,
	}
	ctx.ticker = time.NewTicker(time.Second * 3)
	ctx.netServer = &NetServer{
		queue: make(chan *bytes.Buffer, 50),
	}
	err := wss.goPool.Submit(func() {
		ctx.netServer.HandleCTOSPacket(ctx.dp)

	})
	if err != nil {
		return nil, gnet.Close
	}
	c.SetContext(ctx)

	return nil, gnet.None
}

func (wss *Server) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	if err != nil {
		logging.Warnf("error occurred on connection=%s, %v\n", c.RemoteAddr().String(), err)
	}
	ctx := c.Context().(*Context)
	if ctx.netServer != nil {
		close(ctx.netServer.queue)
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

		buff := wss.bytesPool.Get()
		buff.Write(arr[2:])
		ctx.netServer.queue <- buff

		if err != nil {
			return gnet.Close
		}
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
	Id        uint64
	nextOp    tcpReadOp
	msgLen    int
	dp        *DuelPlayer
	netServer *NetServer
	ticker    *time.Ticker
}
