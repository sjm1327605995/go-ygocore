package ygocore

/*
typedef struct {
unsigned int code;
unsigned int alias;
unsigned int setcode;
unsigned int type;
unsigned int level;
unsigned int attribute;
unsigned int race;
long attack;
long defense;
unsigned int lscale;
unsigned int rscale;
unsigned int link_marker;
}card_data;
*/
import "C"

import (
	"bytes"
	"github.com/ebitengine/purego"
	"github.com/sjm1327605995/go-ygocore/core"
)

type CardData struct {
	Code       uint32 `gorm:"column:id"`
	Ot         uint32 `gorm:"column:ot"`
	Alias      uint32 `gorm:"column:alias"`
	SetCode    uint64 `gorm:"column:setcode"`
	Typ        int32  `gorm:"column:type"`
	Level      uint32 `gorm:"column:level"`
	Race       uint32 `gorm:"column:race"`
	Attribute  uint32 `gorm:"column:attribute"`
	Attack     int32  `gorm:"column:atk"`
	Defense    int32  `gorm:"column:def"`
	Lscale     uint32 `gorm:"-"`
	Rscale     uint32 `gorm:"-"`
	LinkMarker int32  `gorm:"-"`
	//
	//	Category int
}

func NewYGOCore(libPath string, scriptReader ScriptReader, CardReader CardReader, msgHandler MessageHandler) *YGOCore {
	libc, err := core.LoadLib(libPath)
	if err != nil {
		panic(err)
	}
	if scriptReader == nil {
		panic("scriptReader is nil")
	}
	if CardReader == nil {
		panic("CardReader is nil")
	}
	if msgHandler == nil {
		panic("msgHandler is nil")
	}
	var ygoCore = &YGOCore{
		ScriptReader:   scriptReader,
		CardReader:     CardReader,
		MessageHandler: msgHandler,
	}

	purego.RegisterLibFunc(&ygoCore.CreateDuel, libc, "create_duel")
	purego.RegisterLibFunc(&ygoCore.StartDuel, libc, "start_duel")
	purego.RegisterLibFunc(&ygoCore.EndDuel, libc, "end_duel")
	purego.RegisterLibFunc(&ygoCore.SetPlayerInfo, libc, "set_player_info")
	purego.RegisterLibFunc(&ygoCore.GetLogMessage, libc, "get_log_message")
	purego.RegisterLibFunc(&ygoCore.GetMessage, libc, "get_message")
	purego.RegisterLibFunc(&ygoCore.Process, libc, "process")
	purego.RegisterLibFunc(&ygoCore.NewCard, libc, "new_card")
	purego.RegisterLibFunc(&ygoCore.QueryCard, libc, "query_card")
	purego.RegisterLibFunc(&ygoCore.QueryFieldCount, libc, "query_field_count")
	purego.RegisterLibFunc(&ygoCore.QueryFieldCard, libc, "query_field_card")
	purego.RegisterLibFunc(&ygoCore.QueryFieldInfo, libc, "query_field_info")
	purego.RegisterLibFunc(&ygoCore.SetResponseI, libc, "set_responsei")
	purego.RegisterLibFunc(&ygoCore.SetResponseB, libc, "set_responseb")
	purego.RegisterLibFunc(&ygoCore.PreloadScript, libc, "preload_script")

	purego.RegisterLibFunc(&ygoCore.setScriptReader, libc, "set_script_reader")
	purego.RegisterLibFunc(&ygoCore.setCardReader, libc, "set_card_reader")
	purego.RegisterLibFunc(&ygoCore.setMessageHandler, libc, "set_message_handler")
	scriptReaderLib := func(scriptName *C.char, slen *C.int) *C.uchar {
		*slen = 0
		scriptNameStr := C.GoString(scriptName)
		// 调用适当的函数读取脚本内容
		data := ygoCore.ScriptReader(scriptNameStr)
		if len(data) == 0 {
			// 处理错误
			return (*C.uchar)(nil)
		}

		// 将数据长度设置到slen指针
		*slen = C.int(len(data))
		// 创建C字节数组并将数据复制到其中
		return (*C.uchar)(C.CBytes(data))
	}

	scriptReaderCb := purego.NewCallback(scriptReaderLib)

	ygoCore.setScriptReader(scriptReaderCb)

	cardReaderLib := func(cardId uint, card *C.card_data) uint {
		data := ygoCore.CardReader(int32(cardId))
		if data != nil {
			card.code = C.uint(data.Code)
			card.alias = C.uint(data.Alias)
			card.setcode = C.uint(data.SetCode)
			card._type = C.uint(data.Typ)
			card.level = C.uint(data.Level)
			card.attribute = C.uint(data.Attribute)
			card.race = C.uint(data.Race)
			card.attack = C.long(data.Attack)
			card.defense = C.long(data.Defense)
			card.lscale = C.uint(data.Lscale)
			card.rscale = C.uint(data.Rscale)
			card.link_marker = C.uint(data.LinkMarker)
		}

		return cardId
	}
	cardReaderCb := purego.NewCallback(cardReaderLib)
	ygoCore.setCardReader(cardReaderCb)

	msgHandlerLib := func(u uintptr, u2 uint32) uint {
		var buf = make([]byte, 256)
		ygoCore.GetLogMessage(u, buf)
		index := bytes.IndexByte(buf, 0)
		if index > 0 {
			buf = buf[:index]
		}
		ygoCore.MessageHandler(string(buf))
		return 0
	}

	msgHandlerCb := purego.NewCallback(msgHandlerLib)
	ygoCore.setMessageHandler(msgHandlerCb)
	return ygoCore

}

type ScriptReader func(scriptName string) (data []byte)
type MessageHandler func(msg string)
type CardReader func(cardId int32) *CardData
type YGOCore struct {
	CreateDuel        func(seed int32) uintptr
	StartDuel         func(pduel uintptr, options int32)
	EndDuel           func(pduel uintptr)
	SetPlayerInfo     func(pduel uintptr, playerId, LP, startCount, drawCount int32)
	GetLogMessage     func(pduel uintptr, buf []byte)
	GetMessage        func(pduel uintptr, buf []byte) int32
	Process           func(pduel uintptr) uint32
	NewCard           func(pduel uintptr, code uint32, owner, playerid, location, sequence, position uint8)
	QueryCard         func(pduel uintptr, playerid, location, sequence uint8, queryFlag int32, buf []byte, useCache int32) int32
	QueryFieldCount   func(pduel uintptr, playerid, location uint8) int32
	QueryFieldCard    func(pduel uintptr, playerId, location uint8, queryFlag int32, buf []byte, useCache int32) int32
	QueryFieldInfo    func(pduel uintptr, buf []byte) int32
	SetResponseI      func(pduel uintptr, value int32)
	SetResponseB      func(pduel uintptr, buf []byte)
	PreloadScript     func(pduel uintptr, script string, len int32) int32
	setScriptReader   func(f uintptr)
	ScriptReader      ScriptReader
	setCardReader     func(f uintptr)
	CardReader        CardReader
	MessageHandler    MessageHandler
	setMessageHandler func(f uintptr)
}
