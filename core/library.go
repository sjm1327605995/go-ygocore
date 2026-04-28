package core

import (
	"fmt"
	"unsafe"

	"github.com/ebitengine/purego"
)

const (
	SEED_COUNT = 6
)

var (
	currentScriptReader   ScriptReader
	currentCardReader     CardReader
	currentMessageHandler MessageHandler
)

// CardData 卡片数据库信息
type CardData struct {
	Code       uint32
	Alias      uint32
	SetCode    uint64
	Type       uint32
	Level      uint32
	Attribute  uint32
	Race       uint64
	Attack     int32
	Defense    int32
	LScale     uint32
	RScale     uint32
	LinkMarker uint32
}

// ScriptReader 脚本读取器，根据脚本名返回脚本内容
type ScriptReader func(scriptName string) []byte

// CardReader 卡片数据读取器，根据卡片代码返回卡片数据
type CardReader func(code uint32) *CardData

// MessageHandler 消息处理器，处理错误消息
type MessageHandler func(msg string)

// LibraryOptions 动态库加载选项
type LibraryOptions struct {
	LibPath        string
	ScriptReader   ScriptReader
	CardReader     CardReader
	MessageHandler MessageHandler
}

// Library 封装了 ygopro-core 动态库的调用
type Library struct {
	handle uintptr

	setScriptReader   func(f uintptr)
	setCardReader     func(f uintptr)
	setMessageHandler func(f uintptr)

	createDuel      func(seed int32) uintptr
	createDuelV2    func(seedSequence unsafe.Pointer) uintptr
	startDuel       func(pduel uintptr, options uint32)
	endDuel         func(pduel uintptr)
	setPlayerInfo   func(pduel uintptr, playerid int32, lp int32, startcount int32, drawcount int32)
	getLogMessage   func(pduel uintptr, buf unsafe.Pointer)
	getMessage      func(pduel uintptr, buf unsafe.Pointer) int32
	process         func(pduel uintptr) int32
	newCard         func(pduel uintptr, code uint32, owner uint8, playerid uint8, location uint8, sequence uint8, position uint8)
	newTagCard      func(pduel uintptr, code uint32, owner uint8, location uint8)
	queryCard       func(pduel uintptr, playerid uint8, location uint8, sequence uint8, queryFlag uint32, buf unsafe.Pointer, useCache int32) int32
	queryFieldCount func(pduel uintptr, playerid uint8, location uint8) int32
	queryFieldCard  func(pduel uintptr, playerid uint8, location uint8, queryFlag uint32, buf unsafe.Pointer, useCache int32) int32
	queryFieldInfo  func(pduel uintptr, buf unsafe.Pointer) int32
	setResponseI    func(pduel uintptr, value int32)
	setResponseB    func(pduel uintptr, buf unsafe.Pointer)
	preloadScript   func(pduel uintptr, script unsafe.Pointer, length int32) int32
}

func init() {
}

func scriptReaderCallback(scriptName *byte, length *int32) uintptr {
	name := cStringToGoString(scriptName)
	data := currentScriptReader(name)
	if len(data) == 0 {
		*length = 0
		return 0
	}
	*length = int32(len(data))
	return uintptr(unsafe.Pointer(&data[0]))
}

func cardReaderCallback(code uint32, data *CardData) uintptr {
	card := currentCardReader(code)
	if card == nil {
		return 0
	}
	*data = *card
	return 1
}

func messageHandlerCallback(msg *byte) uintptr {
	currentMessageHandler(cStringToGoString(msg))
	return 0
}

// LoadLibrary 加载 ygopro-core 动态库并封装函数
func LoadLibrary(opts LibraryOptions) (*Library, error) {
	handle, err := openLibrary(opts.LibPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load library %s: %w", opts.LibPath, err)
	}

	currentScriptReader = opts.ScriptReader
	currentCardReader = opts.CardReader
	currentMessageHandler = opts.MessageHandler

	lib := &Library{
		handle: handle,
	}

	if opts.ScriptReader != nil {
		purego.RegisterLibFunc(&lib.setScriptReader, handle, "set_script_reader")
		lib.setScriptReader(purego.NewCallback(scriptReaderCallback))
	}
	if opts.CardReader != nil {
		purego.RegisterLibFunc(&lib.setCardReader, handle, "set_card_reader")
		lib.setCardReader(purego.NewCallback(cardReaderCallback))
	}
	if opts.MessageHandler != nil {
		purego.RegisterLibFunc(&lib.setMessageHandler, handle, "set_message_handler")
		lib.setMessageHandler(purego.NewCallback(messageHandlerCallback))
	}

	purego.RegisterLibFunc(&lib.createDuel, handle, "create_duel")
	purego.RegisterLibFunc(&lib.createDuelV2, handle, "create_duel_v2")
	purego.RegisterLibFunc(&lib.startDuel, handle, "start_duel")
	purego.RegisterLibFunc(&lib.endDuel, handle, "end_duel")
	purego.RegisterLibFunc(&lib.setPlayerInfo, handle, "set_player_info")
	purego.RegisterLibFunc(&lib.getLogMessage, handle, "get_log_message")
	purego.RegisterLibFunc(&lib.getMessage, handle, "get_message")
	purego.RegisterLibFunc(&lib.process, handle, "process")
	purego.RegisterLibFunc(&lib.newCard, handle, "new_card")
	purego.RegisterLibFunc(&lib.newTagCard, handle, "new_tag_card")
	purego.RegisterLibFunc(&lib.queryCard, handle, "query_card")
	purego.RegisterLibFunc(&lib.queryFieldCount, handle, "query_field_count")
	purego.RegisterLibFunc(&lib.queryFieldCard, handle, "query_field_card")
	purego.RegisterLibFunc(&lib.queryFieldInfo, handle, "query_field_info")
	purego.RegisterLibFunc(&lib.setResponseI, handle, "set_responsei")
	purego.RegisterLibFunc(&lib.setResponseB, handle, "set_responseb")
	purego.RegisterLibFunc(&lib.preloadScript, handle, "preload_script")

	return lib, nil
}

// Close 关闭动态库
func (lib *Library) Close() error {
	if lib.handle != 0 {
		return closeLibrary(lib.handle)
	}
	return nil
}

// CreateDuel 创建决斗实例（使用单个 seed，已废弃，仅用于回放模式）
func (lib *Library) CreateDuel(seed int32) uintptr {
	return lib.createDuel(seed)
}

// CreateDuelV2 创建决斗实例（使用 seed 数组）
func (lib *Library) CreateDuelV2(seedSequence [SEED_COUNT]uint32) uintptr {
	return lib.createDuelV2(unsafe.Pointer(&seedSequence[0]))
}

// StartDuel 开始决斗
func (lib *Library) StartDuel(pduel uintptr, options uint32) {
	lib.startDuel(pduel, options)
}

// EndDuel 结束决斗
func (lib *Library) EndDuel(pduel uintptr) {
	lib.endDuel(pduel)
}

// SetPlayerInfo 设置玩家信息
func (lib *Library) SetPlayerInfo(pduel uintptr, playerid int32, lp int32, startcount int32, drawcount int32) {
	lib.setPlayerInfo(pduel, playerid, lp, startcount, drawcount)
}

// GetLogMessage 获取日志消息
func (lib *Library) GetLogMessage(pduel uintptr, buf []byte) {
	if len(buf) > 0 {
		lib.getLogMessage(pduel, unsafe.Pointer(&buf[0]))
	}
}

// GetMessage 获取消息
func (lib *Library) GetMessage(pduel uintptr, buf []byte) int32 {
	if len(buf) > 0 {
		return lib.getMessage(pduel, unsafe.Pointer(&buf[0]))
	}
	return 0
}

// Process 执行一个游戏 tick
func (lib *Library) Process(pduel uintptr) int32 {
	return lib.process(pduel)
}

// NewCard 添加卡片到决斗状态
func (lib *Library) NewCard(pduel uintptr, code uint32, owner uint8, playerid uint8, location uint8, sequence uint8, position uint8) {
	lib.newCard(pduel, code, owner, playerid, location, sequence, position)
}

// NewTagCard 添加卡片到 tag 池
func (lib *Library) NewTagCard(pduel uintptr, code uint32, owner uint8, location uint8) {
	lib.newTagCard(pduel, code, owner, location)
}

// QueryCard 查询特定位置的卡片
func (lib *Library) QueryCard(pduel uintptr, playerid uint8, location uint8, sequence uint8, queryFlag uint32, buf []byte, useCache int32) int32 {
	if len(buf) > 0 {
		return lib.queryCard(pduel, playerid, location, sequence, queryFlag, unsafe.Pointer(&buf[0]), useCache)
	}
	return 0
}

// QueryFieldCount 查询特定区域卡片数量
func (lib *Library) QueryFieldCount(pduel uintptr, playerid uint8, location uint8) int32 {
	return lib.queryFieldCount(pduel, playerid, location)
}

// QueryFieldCard 查询特定区域所有卡片
func (lib *Library) QueryFieldCard(pduel uintptr, playerid uint8, location uint8, queryFlag uint32, buf []byte, useCache int32) int32 {
	if len(buf) > 0 {
		return lib.queryFieldCard(pduel, playerid, location, queryFlag, unsafe.Pointer(&buf[0]), useCache)
	}
	return 0
}

// QueryFieldInfo 查询场地信息
func (lib *Library) QueryFieldInfo(pduel uintptr, buf []byte) int32 {
	if len(buf) > 0 {
		return lib.queryFieldInfo(pduel, unsafe.Pointer(&buf[0]))
	}
	return 0
}

// SetResponseI 设置整数响应
func (lib *Library) SetResponseI(pduel uintptr, value int32) {
	lib.setResponseI(pduel, value)
}

// SetResponseB 设置字节数组响应
func (lib *Library) SetResponseB(pduel uintptr, buf []byte) {
	if len(buf) > 0 {
		lib.setResponseB(pduel, unsafe.Pointer(&buf[0]))
	}
}

// PreloadScript 预加载脚本
func (lib *Library) PreloadScript(pduel uintptr, script string) int32 {
	cs := append([]byte(script), 0)
	return lib.preloadScript(pduel, unsafe.Pointer(&cs[0]), int32(len(script)))
}

func cStringToGoString(p *byte) string {
	if p == nil {
		return ""
	}
	var s []byte
	for {
		if *p == 0 {
			break
		}
		s = append(s, *p)
		p = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(p)) + 1))
	}
	return string(s)
}
