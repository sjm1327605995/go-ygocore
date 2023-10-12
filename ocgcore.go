package ygocore

import "C"

/*
#include "ocgcore.h"
*/
import "C"
import (
	"log"
	"unsafe"
)

var (
	scriptReader   func(name string) []byte
	messageHandler func(data unsafe.Pointer, tp int32)
	cardReader     func(cardID uint32, card *CardDataC) bool
)

func RegisterScriptReader(f func(name string) []byte) {
	scriptReader = f
}
func RegisterMessageHandler(f func(data unsafe.Pointer, tp int32)) {
	messageHandler = f
}
func RegisterCardReader(f func(cardID uint32, card *CardDataC) bool) {
	cardReader = f
}
func init() {
	scriptReader = func(name string) []byte {
		log.Println("script", name)
		return nil
	}
	messageHandler = func(data unsafe.Pointer, tp int32) {
		// 将 uintptr 转换为 int64 这里暂时不需要打印日志，所以不写该方法
		//value := int64(uintptr(data))
		//return nil
	}
	cardReader = func(cardID uint32, card *CardDataC) bool {
		return false
	}
}

type CardDataC struct {
	code        uint32
	alias       uint32
	setcode     uint64
	typ         uint32
	level       uint32
	attribute   uint32
	race        uint32
	attack      int32
	defense     int32
	lscale      uint32
	rscale      uint32
	link_marker uint32
	ot          uint32
	category    int
}

//export goScriptReader
func goScriptReader(scriptName *C.char, slen *C.int) *C.uchar {
	// 将C字符串转换为Go字符串

	*slen = 0
	// 调用适当的函数读取脚本内容
	data := scriptReader(C.GoString(scriptName))
	if len(data) == 0 {
		// 处理错误
		return (*C.uchar)(nil)
	}
	// 将数据长度设置到slen指针
	*slen = C.int(len(data))
	// 创建C字节数组并将数据复制到其中
	return (*C.uchar)(C.CBytes(data))

}

//export goMessageHandler
func goMessageHandler(data unsafe.Pointer, size C.uint32_t) {
	messageHandler(data, int32(size))
	// 处理消息
}

//export goCardReader
func goCardReader(cardID C.uint32_t, data *C.card_data) C.uint32_t {

	//TODO 这里进行了内存拷贝需要重新操作下
	var (
		dataC CardDataC
	)
	if cardReader(uint32(cardID), &dataC) {
		data.code = C.uint32_t(dataC.code)
		data.alias = C.uint32_t(dataC.alias)
		data.setcode = C.uint64_t(dataC.setcode)
		data._type = C.uint32_t(dataC.typ)
		data.level = C.uint32_t(dataC.level)
		data.attribute = C.uint32_t(dataC.attribute)
		data.race = C.uint32_t(dataC.race)
		data.attack = C.int32_t(dataC.attack)
		data.defense = C.int32_t(dataC.defense)
		data.lscale = C.uint32_t(dataC.lscale)
		data.rscale = C.uint32_t(dataC.rscale)
		data.link_marker = C.uint32_t(dataC.link_marker)

	} else {
		data.code = C.uint32_t(0)
		data.alias = C.uint32_t(0)
		data.setcode = C.uint64_t(0)
		data._type = C.uint32_t(0)
		data.level = C.uint32_t(0)
		data.attribute = C.uint32_t(0)
		data.race = C.uint32_t(0)
		data.attack = C.int32_t(0)
		data.defense = C.int32_t(0)
		data.lscale = C.uint32_t(0)
		data.rscale = C.uint32_t(0)
		data.link_marker = C.uint32_t(0)
	}
	return 0
}

func CreateGame(seed int32) uintptr {
	pDuel := C.create_duel(C.int(seed))
	return uintptr(pDuel)
}
func StartDuel(pduel uintptr, options int32) {
	C.start_duel(C.longlong(pduel), C.int32_t(options))
}
func EndDuel(pduel uintptr) {
	C.end_duel(C.longlong(pduel))
}
func SetPlayerInfo(pduel uintptr, playerId, lp, startCount, drawCount int32) {
	C.set_player_info(C.longlong(pduel), C.int32_t(playerId), C.int32_t(lp), C.int32_t(startCount), C.int32_t(drawCount))
}

const (
	LogMessageBufLen = 1024
	MessageBufLen    = 0x1000
	QueryCardBufLen  = 0x2000
	ResponsebBufLen  = 64
)

// GetLogMessage 返回[]byte 长度固定为1024
func GetLogMessage(pduel uintptr) []byte {
	var buf = make([]byte, LogMessageBufLen)
	C.get_log_message(C.longlong(pduel), (*C.uchar)(unsafe.Pointer(&buf[0])))
	return buf
}

func GetMessage(pduel uintptr) ([]byte, int32) {
	var buff = make([]byte, MessageBufLen)
	return buff, int32(C.get_message(C.longlong(pduel), (*C.uchar)(unsafe.Pointer(&buff[0]))))
}
func Process(pduel uintptr) int32 {
	return int32(C.process(C.longlong(pduel)))
}
func NewCard(pduel uintptr, code uint32, owner, playerid, location, sequence, position uint8) {
	C.new_card(C.longlong(pduel), C.uint32_t(code), C.uint8_t(owner), C.uint8_t(playerid), C.uint8_t(location), C.uint8_t(sequence), C.uint8_t(position))
}

// QueryCard  buf 长度要大于 0x2000
func QueryCard(pduel uintptr, playerid, location, sequence uint8, queryFlag int32, buf []byte, useCache int32) int32 {
	return int32(C.query_card(C.longlong(pduel), C.uint8_t(playerid), C.uint8_t(location), C.uint8_t(sequence), C.int32_t(queryFlag), (*C.uchar)(unsafe.Pointer(&buf[0])), C.int32_t(useCache)))
}

func QueryFieldCount(pduel uintptr, playerid, location uint8) int32 {
	return int32(C.query_field_count(C.longlong(pduel), C.uint8_t(playerid), C.uint8_t(location)))
}
func QueryFieldCard(pduel uintptr, playerid, location uint8, queryFlag int32, buf []byte, useCache int32) int32 {
	return int32(C.query_field_card(C.longlong(pduel), C.uint8_t(playerid), C.uint8_t(location), C.int32_t(queryFlag), (*C.uchar)(unsafe.Pointer(&buf[0])), C.int32_t(useCache)))
}

func QueryFieldInfo(pduel uintptr, buf []byte) int32 {
	return int32(C.query_field_info(C.longlong(pduel), (*C.uchar)(unsafe.Pointer(&buf[0]))))
}
func SetResponsei(pduel uintptr, value int32) {
	C.set_responsei(C.longlong(pduel), C.int32_t(value))
}
func SetResponseb(pduel uintptr, buf []byte) {
	C.set_responseb(C.longlong(pduel), (*C.uchar)(unsafe.Pointer(&buf[0])))
}
func PreloadScript(pduel uintptr, script []byte) int32 {
	return int32(C.preload_script(C.longlong(pduel), (*C.char)(unsafe.Pointer(&script[0])), C.int32_t(len(script))))
}
