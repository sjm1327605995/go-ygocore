package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/sjm1327605995/go-ygocore/msg/host"
	"github.com/sjm1327605995/go-ygocore/msg/stoc"
	"math/rand"
	"time"
)

func init() {
	rand.NewSource(time.Now().UnixNano())
}

// func main() {
//
//		NewDataManager()
//		RegisterDo()
//		n := rand.Int31n(10000)
//
//		pduel := CreateDuel(n)
//		duel := &SingleDuel{
//			players: []ClientInterface{&ConsoleClient{id: 0}, &ConsoleClient{id: 1}},
//			pduel:   pduel,
//		}
//		SetPlayerInfo(pduel, 0, 8000, 5, 1)
//		SetPlayerInfo(pduel, 1, 8000, 5, 1)
//		var (
//			mainCards = []uint32{14124483, 9411399, 9411399, 18094166, 18094166, 18094166, 40044918, 40044918, 59392529, 50720316, 50720316, 27780618, 27780618, 16605586, 16605586, 22865492, 22865492, 23434538, 23434538, 14558127, 14558127,
//				13650422, 83965310, 81439173, 8949584, 8949584, 32807846, 52947044, 45906428, 24094653, 21143940, 21143940, 21143940, 48130397, 24224830, 24224830, 12071500, 24299458, 24299458, 10045474}
//			exidCards = []uint32{73580471, 79606837, 79606837, 79606837, 21521304, 27552504, 1174075, 1174075, 1174075, 73898890, 73898890, 72336818, 41999284, 94259633, 94259633}
//		)
//		for i := len(mainCards) - 1; i >= 0; i-- {
//			NewCard(pduel, mainCards[i], 0, 0, LOCATION_DECK, 0, POS_FACEDOWN_DEFENSE)
//		}
//		for i := len(exidCards) - 1; i >= 0; i-- {
//			NewCard(pduel, exidCards[i], 0, 0, LOCATION_EXTRA, 0, POS_FACEDOWN_DEFENSE)
//		}
//		for i := len(mainCards) - 1; i >= 0; i-- {
//			NewCard(pduel, mainCards[i], 1, 1, LOCATION_DECK, 0, POS_FACEDOWN_DEFENSE)
//		}
//		for i := len(exidCards) - 1; i >= 0; i-- {
//			NewCard(pduel, exidCards[i], 1, 1, LOCATION_EXTRA, 0, POS_FACEDOWN_DEFENSE)
//		}
//		count1 := QueryFieldCount(pduel, 0, 0x1)
//		count2 := QueryFieldCount(pduel, 0, 0x40)
//		count3 := QueryFieldCount(pduel, 1, 0x1)
//		count4 := QueryFieldCount(pduel, 1, 0x40)
//		fmt.Println(count1, count2, count3, count4)
//		duel.RefreshExtraDef(0)
//		duel.RefreshExtraDef(1)
//		opt := 5 << 16
//		StartDuel(pduel, int32(opt))
//
//		duel.Process()
//	}
func (d *SingleDuel) RefreshExtraDef(player uint8) {
	d.RefreshExtra(d.pDuel, player, 0xe81fff, 1)
}
func (d *SingleDuel) RefreshMzoneDef(player uint8) {
	d.RefreshExtra(d.pDuel, player, 0x881fff, 1)
}

func (d *SingleDuel) RefreshSzoneDef(player uint8) {
	d.RefreshSzone(d.pDuel, player, 0x681fff, 1)
}
func (d *SingleDuel) RefreshHandDef(player uint8) {
	d.RefreshHand(d.pDuel, player, 0x681fff, 1)
}
func (d *SingleDuel) RefreshGraveDef(player uint8) {
	d.RefreshGrave(d.pDuel, player, 0x81fff, 1)
}
func (d *SingleDuel) RefreshSingleDef(player uint8, location uint8, seq uint8) {
	d.RefreshSingle(player, location, seq, 0xf81fff)
}

// void RefreshSingle(int player, int location, int sequence, int flag = 0xf81fff);
// void RefreshGrave(int player, int flag = 0x81fff, int use_cache = 1);
func (d *SingleDuel) RefreshMzone(pdule uintptr, player uint8, flag, use_cache int32) {
	fmt.Println("RefreshMzone")
	var (
		originBuf = make([]byte, 0x2000)
		qbuf      = originBuf[3:]
	)

	qbuf[0] = MSG_UPDATE_DATA
	qbuf[1] = player
	qbuf[2] = LOCATION_MZONE
	length := QueryFieldCard(pdule, player, LOCATION_MZONE, flag, qbuf, use_cache)
	SendBufferToPlayer(d.players[player], STOC_GAME_MSG, qbuf[:length])
	var (
		qlen   int32
		clen   int32
		buffer = NewBuffer(qbuf[:length])
	)

	for qlen < length {
		_ = binary.Read(buffer, binary.LittleEndian, &clen)
		qlen += clen
		if clen == 4 {
			continue
		}
		if qbuf[11]&POS_FACEDOWN == 1 {
			var i int32
			for ; i < clen-4; i++ {
				qbuf[i] = 0
			}
		}
		qbuf = qbuf[clen-4:]

	}
	var list []ClientInterface
	for i := range d.observers {
		list = append(list, d.observers[i])
	}
	SendBufferToPlayer(d.players[1-player], STOC_GAME_MSG, originBuf[:length+6], list...)

}

func (d *SingleDuel) RefreshSzone(pdule uintptr, player uint8, flag, use_cache int32) {
	fmt.Println("RefreshSzone")
	var (
		originBuf = make([]byte, 0x2000)
		qbuf      = originBuf[3:]
	)

	qbuf[0] = MSG_UPDATE_DATA
	qbuf[1] = player
	qbuf[2] = LOCATION_SZONE
	length := QueryFieldCard(pdule, player, LOCATION_SZONE, flag, qbuf, use_cache)
	SendBufferToPlayer(d.players[player], STOC_GAME_MSG, originBuf[:length+6])
	var (
		qlen   int32
		clen   int32
		buffer = NewBuffer(qbuf[:length])
	)

	for qlen < length {
		_ = binary.Read(buffer, binary.LittleEndian, &clen)
		qlen += clen
		if clen == 4 {
			continue
		}
		if qbuf[11]&POS_FACEDOWN == 1 {
			var i int32
			for ; i < clen-4; i++ {
				qbuf[i] = 0
			}
		}
		qbuf = qbuf[clen-4:]

	}
	var list []ClientInterface
	for i := range d.observers {
		list = append(list, d.observers[i])
	}
	SendBufferToPlayer(d.players[1-player], STOC_GAME_MSG, qbuf[:length], list...)

}
func (d *SingleDuel) RefreshHand(pdule uintptr, player uint8, flag, use_cache int32) {
	fmt.Println("RefreshHand")
	var (
		originBuf = make([]byte, 0x2000)
		qbuf      = originBuf[3:]
	)

	qbuf[0] = MSG_UPDATE_DATA
	qbuf[1] = player
	qbuf[2] = LOCATION_HAND
	length := QueryFieldCard(pdule, player, LOCATION_HAND, flag|QUERY_POSITION, qbuf[3:], use_cache)
	SendBufferToPlayer(d.players[player], STOC_GAME_MSG, originBuf[:length+6])
	var (
		qlen   int32
		buffer = NewBuffer(qbuf[:length])
	)

	for qlen < length {
		var (
			slen   int32
			qflag  uint32
			offset = 8
		)
		_ = binary.Read(buffer, binary.LittleEndian, &slen)
		qflag = binary.LittleEndian.Uint32(buffer.Bytes())
		if qflag&QUERY_CODE == 0 {
			offset -= 4
		}
		buffer.Next(offset)

		position := (int32(binary.LittleEndian.Uint32(buffer.Bytes())) >> 24) & 0xff
		if position&POS_FACEUP == 0 {
			for i := 0; i < 4; i++ {
				qbuf[i] = 0
			}
		}

		buffer.Next(int(slen - 4))
		qlen += slen
	}
	var list []ClientInterface
	for i := range d.observers {
		list = append(list, d.observers[i])
	}
	SendBufferToPlayer(d.players[1-player], STOC_GAME_MSG, qbuf[:length], list...)

}
func (d *SingleDuel) RefreshGrave(pdule uintptr, player uint8, flag, use_cache int32) {
	fmt.Println("RefreshGrave")
	var (
		originBuf = make([]byte, 0x2000)
		qbuf      = originBuf[3:]
	)

	qbuf[0] = MSG_UPDATE_DATA
	qbuf[1] = player
	qbuf[2] = LOCATION_GRAVE
	length := QueryFieldCard(pdule, player, LOCATION_GRAVE, flag, qbuf[3:], use_cache)
	var list []ClientInterface
	list = append(list, d.players[1])
	for i := range d.observers {
		list = append(list, d.observers[i])
	}
	SendBufferToPlayer(d.players[0], STOC_GAME_MSG, originBuf[:length+6], list...)

}
func (d *SingleDuel) RefreshExtra(pdule uintptr, player uint8, flag, use_cache int32) {
	fmt.Println("RefreshExtra")
	var (
		originBuf = make([]byte, 0x2000)
		qbuf      = originBuf[3:]
	)
	qbuf[0] = MSG_UPDATE_DATA
	qbuf[1] = player
	qbuf[2] = LOCATION_EXTRA
	length := QueryFieldCard(pdule, player, LOCATION_EXTRA, flag, qbuf[3:], use_cache)

	_ = SendBufferToPlayer(d.players[player], STOC_GAME_MSG, originBuf[:length+6])

}

func (d *SingleDuel) RefreshSingle(player uint8, location uint8, sequence uint8, flag int32) {
	fmt.Println("RefreshSingle")
	var (
		originBuf = make([]byte, 0x2000)
		qbuf      = originBuf[3:]
	)
	qbuf[0] = MSG_UPDATE_DATA
	qbuf[1] = player
	qbuf[2] = location
	qbuf[3] = sequence
	length := QueryFieldCard(d.pDuel, player, LOCATION_GRAVE, flag|QUERY_POSITION, qbuf[3:], 0)
	SendBufferToPlayer(d.players[player], STOC_GAME_MSG, originBuf[:length+7])
	if location == LOCATION_REMOVED && (qbuf[15]&POS_FACEDOWN) != 0 {
		return
	}
	if location&0x90 != 0 || (location&0x2c != 0 && qbuf[15]&POS_FACEUP != 0) {
		SendBufferToPlayer(d.players[1-player], STOC_GAME_MSG, originBuf[:length+7])
		for i := range d.observers {
			SendBufferToPlayer(d.observers[i], STOC_GAME_MSG, originBuf[:length+7])
		}
	}
}

type SingleDuel struct {
	DuelModeBase

	players      [2]*DuelPlayer
	pplayer      [2]*DuelPlayer
	observers    []*DuelPlayer
	engineBuffer []byte
	hostInfo     host.HostInfo
	Ready        [2]bool
	matchResult  [3]uint8
	duelCount    int
	deckError    [2]int32
	pdeck        [2]Deck
	timeLimit    [2]uint16
	tpPlayer     uint8
}

func (d *SingleDuel) Chat(dp *DuelPlayer, buff []byte) {
	//TODO implement me
	panic("implement me")
}

const PRO_VERSION uint16 = 0x1360

func (d *SingleDuel) JoinGame(dp *DuelPlayer, buff *bytes.Buffer, isCreator bool) {

	if isCreator {
		if dp.game != nil || dp.Type != 0xff {
			var scem stoc.ErrorMsg
			scem.Msg = ERRMSG_JOINERROR
			scem.Code = 0
			SendPacketToPlayer(dp, STOC_ERROR_MSG, scem)
			DisconnectPlayer(dp)

		}
	}
	dp.game = d
	if d.players[0] == nil && d.players[1] == nil && len(d.observers) == 0 {
		d.HostPlayer = dp
	}
	var (
		scjg stoc.JoinGame
		sctc stoc.TypeChange
	)
	scjg.Info = d.hostInfo
	if d.HostPlayer == dp {
		sctc.Type = 0x10
	}
	if d.players[0] == nil || d.players[1] == nil {
		var scpe stoc.HSPlayerEnter
		scpe.Name = dp.Name
		if d.players[0] == nil {
			scpe.Pos = 0
		} else {
			scpe.Pos = 1
		}
		if d.players[0] != nil {
			SendPacketToPlayer(d.players[0], STOC_HS_PLAYER_ENTER, scpe)
		}
		if d.players[1] != nil {
			SendPacketToPlayer(d.players[1], STOC_HS_PLAYER_ENTER, scpe)
		}

		for i := range d.observers {
			SendPacketToPlayer(d.observers[i], STOC_HS_PLAYER_ENTER, scpe)
		}
		if d.players[0] == nil {
			d.players[0] = dp
			dp.Type = NETPLAYER_TYPE_PLAYER1
			sctc.Type |= NETPLAYER_TYPE_PLAYER1
		} else {
			d.players[1] = dp
			dp.Type = NETPLAYER_TYPE_PLAYER2
			sctc.Type |= NETPLAYER_TYPE_PLAYER2
		}
	} else {
		d.observers = append(d.observers, dp)
		dp.Type = NETPLAYER_TYPE_OBSERVER
		sctc.Type |= NETPLAYER_TYPE_OBSERVER
		var scwc stoc.HSWatchChange
		scwc.WatchCount = uint16(len(d.observers))
		if d.players[0] != nil {
			SendPacketToPlayer(d.players[0], STOC_HS_WATCH_CHANGE, scwc)
		}
		if d.players[1] != nil {
			SendPacketToPlayer(d.players[1], STOC_HS_WATCH_CHANGE, scwc)
		}
		for i := range d.observers {
			SendPacketToPlayer(d.observers[i], STOC_HS_WATCH_CHANGE, scwc)
		}
	}
	SendPacketToPlayer(dp, STOC_JOIN_GAME, scjg)
	SendPacketToPlayer(dp, STOC_TYPE_CHANGE, sctc)
	if d.players[0] != nil {
		var scpe stoc.HSPlayerEnter
		scpe.Name = d.players[0].Name
		scpe.Pos = 0
		SendPacketToPlayer(dp, STOC_HS_PLAYER_ENTER, scpe)
		if d.Ready[0] {
			var scpc = stoc.HSPlayerChange{
				Status: PLAYERCHANGE_READY,
			}
			SendPacketToPlayer(dp, STOC_HS_PLAYER_CHANGE, scpc)
		}
	}
	if d.players[1] != nil {
		var scpe stoc.HSPlayerEnter
		scpe.Name = d.players[1].Name
		scpe.Pos = 1
		SendPacketToPlayer(dp, STOC_HS_PLAYER_ENTER, scpe)

		if d.Ready[1] {
			var scpc = stoc.HSPlayerChange{
				Status: 0x10 | PLAYERCHANGE_READY,
			}
			SendPacketToPlayer(dp, STOC_HS_PLAYER_CHANGE, scpc)
		}

	}
	if len(d.observers) > 0 {
		var scwc stoc.HSWatchChange
		scwc.WatchCount = uint16(len(d.observers))
		SendPacketToPlayer(dp, STOC_HS_WATCH_CHANGE, scwc)
	}
}

func (d *SingleDuel) LeaveGame(dp *DuelPlayer) {
	if dp == d.HostPlayer {
		EndDuel(d.pDuel)
		//  TODO
		//	NetServer::StopServer();
	} else if dp.Type == NETPLAYER_TYPE_OBSERVER {
		//delete(d.observers,dp.Name)
		if d.DuelStage == DUEL_STAGE_BEGIN {
			var scwc stoc.HSWatchChange
			scwc.WatchCount = uint16(len(d.observers))
			if d.players[0] != nil {
				SendPacketToPlayer(d.players[0], STOC_HS_WATCH_CHANGE, scwc)
			}
			if d.players[1] != nil {
				SendPacketToPlayer(d.players[1], STOC_HS_WATCH_CHANGE, scwc)
			}
			for i := range d.observers {
				SendPacketToPlayer(d.observers[i], STOC_HS_WATCH_CHANGE, scwc)
			}
		}
		DisconnectPlayer(dp)
	} else {
		if d.DuelStage == DUEL_STAGE_BEGIN {
			var scpc stoc.HSPlayerChange
			d.players[dp.Type] = nil
			d.Ready[dp.Type] = false
			scpc.Status = uint8(dp.Type<<4 | PLAYERCHANGE_LEAVE)
			if d.players[0] != nil && dp.Type != 0 {
				SendPacketToPlayer(d.players[0], STOC_HS_PLAYER_CHANGE, scpc)
			}
			if d.players[1] != nil && dp.Type != 1 {
				SendPacketToPlayer(d.players[1], STOC_HS_PLAYER_CHANGE, scpc)
			}
			DisconnectPlayer(dp)
		} else {
			if d.DuelStage == DUEL_STAGE_SIDING {
				if !d.Ready[0] {
					SendPacketToPlayer(d.players[0], STOC_DUEL_START, nil)
				}
				if !d.Ready[1] {
					SendPacketToPlayer(d.players[1], STOC_DUEL_START, nil)
				}
			}
			if d.DuelStage != DUEL_STAGE_END {
				wbuf := make([]byte, 6)
				wbuf[3] = MSG_WIN
				wbuf[4] = byte(1 - dp.Type)
				wbuf[5] = 0x24
				var resendList []ClientInterface
				resendList = append(resendList, d.players[1])
				for i := range d.observers {
					resendList = append(resendList, d.observers[i])
				}
				SendBufferToPlayer(d.players[0], STOC_GAME_MSG, wbuf, resendList...)
				EndDuel(d.pDuel)
			}
			DisconnectPlayer(dp)
		}
	}
}

func (d *SingleDuel) ToObserver(dp *DuelPlayer) {
	//TODO implement me
	panic("implement me")
}

func (d *SingleDuel) PlayerReady(dp *DuelPlayer, isReady bool) {
	if dp.Type > 1 {
		return
	}
	if d.Ready[dp.Type] == isReady {
		return
	}
	if isReady {
		var deckError int32
		if !d.hostInfo.NoCheckDeck {
			if d.deckError[dp.Type] != 0 {
				deckError = (DECKERROR_UNKNOWNCARD << 28) + d.deckError[dp.Type]
			} else {
				deckError = DkManager.CheckDeck(&d.pdeck[dp.Type], d.hostInfo.Lflist, d.hostInfo.Rule)
			}
		}
		if deckError != 0 {
			var scpc stoc.HSPlayerChange
			scpc.Status = uint8(dp.Type<<4 | PLAYERCHANGE_NOTREADY)
			SendPacketToPlayer(dp, STOC_HS_PLAYER_CHANGE, scpc)
			var scem stoc.ErrorMsg
			scem.Msg = ERRMSG_DECKERROR
			scem.Code = uint32(deckError)
			SendPacketToPlayer(dp, STOC_ERROR_MSG, scem)
			return
		}
	}
	d.Ready[dp.Type] = isReady
	var scpc stoc.HSPlayerChange
	scpc.Status = uint8(dp.Type<<4) | IFELSE[uint8](isReady, PLAYERCHANGE_READY, PLAYERCHANGE_NOTREADY)
	SendPacketToPlayer(d.players[dp.Type], STOC_HS_PLAYER_CHANGE, scpc)
	if d.players[1-dp.Type] != nil {
		SendPacketToPlayer(d.players[1-dp.Type], STOC_HS_PLAYER_CHANGE, scpc)
	}
	for i := range d.observers {
		SendPacketToPlayer(d.observers[i], STOC_HS_PLAYER_CHANGE, scpc)
	}
}
func IFELSE[T any](condition bool, trueBlock, falseBlock T) T {
	if condition {
		return trueBlock
	}
	return falseBlock
}
func (d *SingleDuel) PlayerKick(dp *DuelPlayer, pos uint8) {
	//TODO implement me
	panic("implement me")
}

func (d *SingleDuel) UpdateDeck(dp *DuelPlayer, buff []byte) error {
	if dp.Type > 1 || d.Ready[dp.Type] {
		return nil
	}
	reader := bytes.NewReader(buff)
	var (
		mainc int32
		sidec int32
	)
	binary.Read(reader, binary.LittleEndian, &mainc)
	binary.Read(reader, binary.LittleEndian, &sidec)
	possibleMaxLength := int32((len(buff) - 8) / 4)
	if mainc > possibleMaxLength || sidec > possibleMaxLength || mainc+sidec > possibleMaxLength {
		var scem stoc.ErrorMsg
		scem.Msg = ERRMSG_DECKERROR
		scem.Code = 0
		SendPacketToPlayer(dp, STOC_ERROR_MSG, scem)
	}
	if d.duelCount == 0 {
		d.deckError[dp.Type] = DkManager.LoadDeck(&d.pdeck[dp.Type], buff[8:], int(mainc), int(sidec), false)
	} else {
		if DkManager.LoadSide(&d.pdeck[dp.Type], buff[8:], int(mainc), int(sidec)) {
			d.Ready[dp.Type] = true
			SendPacketToPlayer(dp, STOC_DUEL_START, nil)
			if d.Ready[0] && d.Ready[1] {
				SendPacketToPlayer(d.players[d.tpPlayer], STOC_SELECT_TP, nil)
				d.players[1-d.tpPlayer].Status = 0xff
				d.players[d.tpPlayer].Status = CTOS_TP_RESULT
				d.DuelStage = DUEL_STAGE_FIRSTGO
			} else {
				var scem stoc.ErrorMsg
				scem.Msg = ERRMSG_SIDEERROR
				scem.Code = 0
				SendPacketToPlayer(dp, STOC_ERROR_MSG, scem)
			}
		}
	}
	return nil
}

func (d *SingleDuel) StartDuel(dp *DuelPlayer) {
	if dp != d.HostPlayer {
		return
	}
	if !d.Ready[0] || !d.Ready[1] {
		return
	}
	//TODO
	//NetServer::StopListen(); 貌似是停止广播
	var sendPlayer []ClientInterface
	sendPlayer = append(sendPlayer, d.players[1])
	for i := range d.observers {
		d.observers[i].Status = CTOS_LEAVE_GAME
		sendPlayer = append(sendPlayer, d.observers[i])
	}
	SendPacketToPlayer(d.players[0], STOC_DUEL_START, nil, sendPlayer...)
	deckBuff := make([]byte, 3, 15)
	pbuf := deckBuff
	binary.LittleEndian.AppendUint16(pbuf, uint16(len(d.pdeck[0].main)))
	binary.LittleEndian.AppendUint16(pbuf, uint16(len(d.pdeck[0].extra)))
	binary.LittleEndian.AppendUint16(pbuf, uint16(len(d.pdeck[0].side)))
	binary.LittleEndian.AppendUint16(pbuf, uint16(len(d.pdeck[1].main)))
	binary.LittleEndian.AppendUint16(pbuf, uint16(len(d.pdeck[1].extra)))
	binary.LittleEndian.AppendUint16(pbuf, uint16(len(d.pdeck[1].side)))
	SendBufferToPlayer(d.players[0], STOC_DECK_COUNT, deckBuff)
	move(deckBuff[3:])
	SendBufferToPlayer(d.players[1], STOC_DECK_COUNT, deckBuff)
	SendPacketToPlayer(d.players[0], STOC_SELECT_HAND, nil, d.players[1])
	d.handResult[0] = 0
	d.handResult[1] = 0
	d.players[0].Status = CTOS_HAND_RESULT
	d.players[1].Status = CTOS_HAND_RESULT
	d.DuelStage = DUEL_STAGE_FINGER
}
func move(arr []byte) {
	mid := len(arr) / 2
	for i := 0; i < mid; i++ {
		arr[i], arr[i+mid] = arr[i+mid], arr[i]
	}
}
func (d *SingleDuel) HandResult(dp *DuelPlayer, res uint8) {
	if res > 3 {
		return
	}
	if dp.Status != CTOS_HAND_RESULT {
		return
	}
	d.handResult[dp.Type] = res
	if d.handResult[0] != 0 && d.handResult[1] != 0 {
		var schr = stoc.HandResult{
			Res1: d.handResult[0],
			Res2: d.handResult[1],
		}
		var list []ClientInterface
		for i := range d.observers {
			list = append(list, d.observers[i])
		}
		SendPacketToPlayer(d.players[0], STOC_HAND_RESULT, schr, list...)
		schr.Res1 = d.handResult[1]
		schr.Res2 = d.handResult[0]
		SendPacketToPlayer(d.players[1], STOC_HAND_RESULT, schr)
		if d.handResult[0] == d.handResult[1] {
			SendPacketToPlayer(d.players[0], STOC_SELECT_HAND, nil, d.players[1])
			d.handResult[0] = 0
			d.handResult[1] = 1
			d.players[0].Status = CTOS_HAND_RESULT
			d.players[1].Status = CTOS_HAND_RESULT
		} else if (d.handResult[0] == 1 && d.handResult[1] == 2) ||
			(d.handResult[0] == 2 && d.handResult[1] == 3) ||
			(d.handResult[0] == 3 && d.handResult[1] == 1) {
			SendPacketToPlayer(d.players[1], STOC_SELECT_TP, nil)
			d.tpPlayer = 1
			d.players[0].Status = 0xff
			d.players[1].Status = CTOS_TP_RESULT
			d.DuelStage = DUEL_STAGE_FIRSTGO
		} else {
			SendPacketToPlayer(d.players[0], STOC_SELECT_TP, nil)
			d.players[1].Status = 0xff
			d.players[0].Status = CTOS_TP_RESULT
			d.tpPlayer = 0
			d.DuelStage = DUEL_STAGE_FIRSTGO
		}
	}
}

func (d *SingleDuel) TPResult(dp *DuelPlayer, tp uint8) {
	if dp.Status != CTOS_TP_RESULT {
		return
	}
	d.DuelStage = DUEL_STAGE_FINGER
	var (
		swapped bool
	)
	d.pplayer[0] = d.players[0]
	d.pplayer[1] = d.players[1]
	if (tp != 0 && dp.Type == 1) || (tp == 0 && dp.Type == 0) {
		//玩家位置交换
		d.players[0], d.players[1] = d.players[1], d.players[0]
		d.players[0].Type, d.players[1].Type = 0, 1
		d.pdeck[0], d.pdeck[1] = d.pdeck[1], d.pdeck[0]
		swapped = true
	}
	dp.Status = CTOS_RESPONSE
	seed := rand.Uint32()
	//ReplayHeader rh;
	//rh.id = 0x31707279;
	//rh.version = PRO_VERSION;
	//rh.flag = REPLAY_UNIFORM;
	//rh.seed = seed;
	//rh.start_time = (unsigned int)time(nullptr);
	//last_replay.BeginRecord();
	//last_replay.WriteHeader(rh);
	//last_replay.WriteData(players[0]->name, 40, false);
	//last_replay.WriteData(players[1]->name, 40, false);
	if !d.hostInfo.NoShuffleDeck {
		Shuffle[[]*CardDataC](d.pdeck[0].main)
		Shuffle[[]*CardDataC](d.pdeck[1].main)
	}
	for i := range d.timeLimit {
		d.timeLimit[i] = d.hostInfo.TimeLimit
	}
	RegisterDo()
	d.pDuel = CreateDuel(int32(seed))
	SetPlayerInfo(d.pDuel, 0, d.hostInfo.StartLp, int32(d.hostInfo.StartHand), int32(d.hostInfo.DrawCount))
	SetPlayerInfo(d.pDuel, 1, d.hostInfo.StartLp, int32(d.hostInfo.StartHand), int32(d.hostInfo.DrawCount))
	opt := int32(d.hostInfo.DuleRule) << 16
	if d.hostInfo.NoShuffleDeck {
		opt |= DUEL_PSEUDO_SHUFFLE
	}
	//last_replay.WriteInt32(host_info.start_lp, false);
	//last_replay.WriteInt32(host_info.start_hand, false);
	//last_replay.WriteInt32(host_info.draw_count, false);
	//last_replay.WriteInt32(opt, false);
	//last_replay.Flush();
	//last_replay.WriteInt32(pdeck[0].main.size(), false);
	for i := range d.pdeck[0].main {
		NewCard(d.pDuel, uint32(d.pdeck[0].main[i].Code), 0, 0, LOCATION_DECK, 0, POS_FACEDOWN_DEFENSE)
		//last_replay.WriteInt32(pdeck[0].main[i]->first, false);
	}
	//last_replay.WriteInt32(pdeck[0].extra.size(), false);
	for i := range d.pdeck[0].extra {
		NewCard(d.pDuel, uint32(d.pdeck[0].extra[i].Code), 0, 0, LOCATION_EXTRA, 0, POS_FACEDOWN_DEFENSE)
		//last_replay.WriteInt32(pdeck[0].extra[i]->first, false);
	}
	//last_replay.WriteInt32(pdeck[1].main.size(), false);
	for i := range d.pdeck[1].main {
		NewCard(d.pDuel, uint32(d.pdeck[1].main[i].Code), 1, 1, LOCATION_DECK, 0, POS_FACEDOWN_DEFENSE)
		//last_replay.WriteInt32(pdeck[1].main[i]->first, false);
	}
	//last_replay.WriteInt32(pdeck[1].extra.size(), false);
	for i := range d.pdeck[0].extra {
		NewCard(d.pDuel, uint32(d.pdeck[0].extra[i].Code), 1, 1, LOCATION_EXTRA, 0, POS_FACEDOWN_DEFENSE)
		//last_replay.WriteInt32(pdeck[1].extra[i]->first, false);
	}
	//last_replay.Flush();
	var startBuf = make([]byte, 3, 34)
	startBuf = append(startBuf, 0, d.hostInfo.DuleRule)
	binary.LittleEndian.AppendUint32(startBuf, uint32(d.hostInfo.StartLp))
	binary.LittleEndian.AppendUint32(startBuf, uint32(d.hostInfo.StartLp))
	binary.LittleEndian.AppendUint16(startBuf, uint16(QueryFieldCount(d.pDuel, 0, 0x1)))
	binary.LittleEndian.AppendUint16(startBuf, uint16(QueryFieldCount(d.pDuel, 0, 0x40)))
	binary.LittleEndian.AppendUint16(startBuf, uint16(QueryFieldCount(d.pDuel, 1, 0x1)))
	binary.LittleEndian.AppendUint16(startBuf, uint16(QueryFieldCount(d.pDuel, 1, 0x40)))
	SendBufferToPlayer(d.players[0], STOC_GAME_MSG, startBuf[:19])
	startBuf[1] = 0x10
	SendBufferToPlayer(d.players[1], STOC_GAME_MSG, startBuf[:19])
	if swapped {
		startBuf[1] = 0x10
	} else {
		startBuf[1] = 0x11
	}
	for i := range d.observers {
		SendBufferToPlayer(d.observers[i], STOC_GAME_MSG, startBuf[:19])
	}
	d.RefreshExtraDef(0)
	d.RefreshExtraDef(1)
	StartDuel(d.pDuel, opt)
	if d.hostInfo.TimeLimit != 0 {
		//TODO 定时器
	}
	d.Process()
}
func Shuffle[S ~[]E, E any](s S) {
	rand.Shuffle(len(s), func(i, j int) {
		s[i], s[j] = s[j], s[i]
	})
}
func (d *SingleDuel) Surrender(dp *DuelPlayer) {
	//TODO implement me
	panic("implement me")
}

func (d *SingleDuel) GetResponse(dp *DuelPlayer, buff []byte) {
	//TODO implement me
	panic("implement me")
}

func (d *SingleDuel) TimeConfirm(dp *DuelPlayer) {
	//TODO implement me
	panic("implement me")
}

func (d *SingleDuel) EndDuel() {
	//TODO implement me
	panic("implement me")
}

func (d *SingleDuel) Process() {
	var engineBuffer = make([]byte, 0x1000)
	var (
		engFlag int32
		engLen  int32 = 0
		stop          = 0
	)

	for stop == 0 {
		if engFlag == 2 {
			break
		}

		result := Process(d.pDuel)
		engLen = result & 0xffff
		engFlag = result >> 16
		if engLen > 0 {
			_ = GetMessage(d.pDuel, engineBuffer[3:])
			stop = d.Analyze(engineBuffer, engLen)
		}
	}

}

// Analyze 这里做了特殊处理， msgbuffer前3个自己都是不做使用位。后面才是内容
func (d *SingleDuel) Analyze(msgbuffer []byte, engLen int32) int {

	var (
		offset, pbufw      int
		pbuf               = &Buffer{buf: msgbuffer[3:engLen]}
		player, count, typ uint8
	)

	for pbuf.Len() > 0 {
		offset = pbuf.Position()
		var engType uint8
		_ = binary.Read(pbuf, binary.LittleEndian, &engType)
		fmt.Println("engType", engType)
		switch engType {
		case MSG_RETRY:
			d.WaitForResponse(d.LastResponse)
			SendBufferToPlayer(d.players[d.LastResponse], STOC_GAME_MSG, msgbuffer[offset:pbuf.Position()+3])
			return 1
		case MSG_HINT:
			fmt.Println("MSG_HINT")
			_ = binary.Read(pbuf, binary.LittleEndian, &typ)
			_ = binary.Read(pbuf, binary.LittleEndian, &player)
			pbuf.Next(4)
			switch typ {
			case 1, 2, 3, 5:
				//发送消息给客户端
				SendBufferToPlayer(d.players[player], STOC_GAME_MSG, msgbuffer[offset:pbuf.Position()+3])

			case 4, 6, 7, 8, 9, 11:
				var resendList []ClientInterface
				for i := range d.observers {
					resendList = append(resendList, d.observers[i])
				}
				//发送消息给客户端
				SendBufferToPlayer(d.players[1-player], STOC_GAME_MSG, msgbuffer[offset:pbuf.Position()+3], resendList...)
			case 10:
				//发送消息给客户端
				var resendList []ClientInterface
				resendList = append(resendList, d.players[1])
				for i := range d.observers {
					resendList = append(resendList, d.observers[i])
				}
				SendBufferToPlayer(d.players[0], STOC_GAME_MSG, msgbuffer[offset:pbuf.Position()+3], resendList...)
			}
		case MSG_WIN:
			fmt.Println("MSG_WIN")
			_ = binary.Read(pbuf, binary.LittleEndian, &player)
			_ = binary.Read(pbuf, binary.LittleEndian, &typ)

			var resendList []ClientInterface
			resendList = append(resendList, d.players[1])
			for i := range d.observers {
				resendList = append(resendList, d.observers[i])
			}
			SendBufferToPlayer(d.players[0], STOC_GAME_MSG, msgbuffer[offset:pbuf.Position()+3], resendList...)
			if player > 1 {
				d.matchResult[d.duelCount] = uint8(1 - player)
				d.duelCount++
				d.tpPlayer = 1 - d.tpPlayer
			}
			d.EndDuel()
			return 2
		case MSG_SELECT_BATTLECMD:
			_ = binary.Read(pbuf, binary.LittleEndian, &player)
			_ = binary.Read(pbuf, binary.LittleEndian, &count)
			pbuf.Next(int(count * 11))
			_ = binary.Read(pbuf, binary.LittleEndian, &count)
			pbuf.Next(int(count*8) + 2)
			d.RefreshMzoneDef(0)
			d.RefreshMzoneDef(1)
			d.RefreshSzoneDef(0)
			d.RefreshSzoneDef(1)
			d.RefreshHandDef(0)
			d.RefreshHandDef(1)
			d.WaitForResponse(player)
			SendBufferToPlayer(d.players[player], STOC_GAME_MSG, msgbuffer[offset:pbuf.Position()+3])
			return 1
		case MSG_SELECT_IDLECMD:
			_ = binary.Read(pbuf, binary.LittleEndian, &player)
			_ = binary.Read(pbuf, binary.LittleEndian, &count)
			pbuf.Next(int(count * 7))
			_ = binary.Read(pbuf, binary.LittleEndian, &count)
			pbuf.Next(int(count * 7))
			_ = binary.Read(pbuf, binary.LittleEndian, &count)
			pbuf.Next(int(count * 7))
			_ = binary.Read(pbuf, binary.LittleEndian, &count)
			pbuf.Next(int(count * 7))
			_ = binary.Read(pbuf, binary.LittleEndian, &count)
			pbuf.Next(int(count * 7))
			_ = binary.Read(pbuf, binary.LittleEndian, &count)
			pbuf.Next(int(count*11 + 3))
			d.RefreshMzoneDef(0)
			d.RefreshMzoneDef(1)
			d.RefreshSzoneDef(0)
			d.RefreshSzoneDef(1)
			d.RefreshHandDef(0)
			d.RefreshHandDef(1)
			d.WaitForResponse(player)
			SendBufferToPlayer(d.players[player], STOC_GAME_MSG, msgbuffer[offset:pbuf.Position()+3])
			return 1
		case MSG_SELECT_EFFECTYN:
			_ = binary.Read(pbuf, binary.LittleEndian, &player)
			pbuf.Next(12)
			d.WaitForResponse(player)
			SendBufferToPlayer(d.players[player], STOC_GAME_MSG, msgbuffer[offset:pbuf.Position()+3])
			return 1
		case MSG_SELECT_YESNO:
			_ = binary.Read(pbuf, binary.LittleEndian, &player)
			pbuf.Next(4)
			d.WaitForResponse(player)
			SendBufferToPlayer(d.players[player], STOC_GAME_MSG, msgbuffer[offset:pbuf.Position()+3])
			return 1
		case MSG_SELECT_OPTION:
			_ = binary.Read(pbuf, binary.LittleEndian, &player)
			_ = binary.Read(pbuf, binary.LittleEndian, &count)
			pbuf.Next(int(count * 4))
			d.WaitForResponse(player)
			SendBufferToPlayer(d.players[player], STOC_GAME_MSG, msgbuffer[offset:pbuf.Position()+3])
			return 1
		case MSG_SELECT_CARD, MSG_SELECT_TRIBUTE:
			_ = binary.Read(pbuf, binary.LittleEndian, &player)
			pbuf.Next(3)
			_ = binary.Read(pbuf, binary.LittleEndian, &count)
			var (
				c uint8
				i uint8
			)
			for i = 0; i < count; i++ {
				pbufw = pbuf.Position()
				pbuf.Next(4) //code
				_ = binary.Read(pbuf, binary.LittleEndian, &c)
				pbuf.Next(1) //l
				pbuf.Next(1) //s
				pbuf.Next(1) //ss
				if c != player {
					binary.LittleEndian.PutUint32(pbuf.OffsetBytes(pbufw), 0)
				}

			}
			d.WaitForResponse(player)
			SendBufferToPlayer(d.players[player], STOC_GAME_MSG, msgbuffer[offset:pbuf.Position()+3])
			return 1
		case MSG_SELECT_UNSELECT_CARD:
			_ = binary.Read(pbuf, binary.LittleEndian, &player)
			pbuf.Next(4)
			_ = binary.Read(pbuf, binary.LittleEndian, &count)
			var (
				c uint8
				i uint8
			)
			for i = 0; i < count; i++ {
				pbufw = pbuf.Position()
				pbuf.Next(4) //code
				_ = binary.Read(pbuf, binary.LittleEndian, &c)
				pbuf.Next(1) //l
				pbuf.Next(1) //s
				pbuf.Next(1) //ss
				if c != player {
					binary.LittleEndian.PutUint32(pbuf.OffsetBytes(pbufw), 0)
				}
			}
			_ = binary.Read(pbuf, binary.LittleEndian, &count)
			for i = 0; i < count; i++ {
				pbufw = pbuf.Position()
				pbuf.Next(4) //code
				_ = binary.Read(pbuf, binary.LittleEndian, &c)
				pbuf.Next(1) //l
				pbuf.Next(1) //s
				pbuf.Next(1) //ss
				if c != player {
					binary.LittleEndian.PutUint32(pbuf.OffsetBytes(pbufw), 0)
				}
			}
			d.WaitForResponse(player)
			SendBufferToPlayer(d.players[player], STOC_GAME_MSG, msgbuffer[offset:pbuf.Position()+3])
			return 1
		case MSG_SELECT_CHAIN:
			_ = binary.Read(pbuf, binary.LittleEndian, &player)
			_ = binary.Read(pbuf, binary.LittleEndian, &count)
			pbuf.Next(10 + int(count)*13)
			d.WaitForResponse(player)
			SendBufferToPlayer(d.players[player], STOC_GAME_MSG, msgbuffer[offset:pbuf.Position()+3])
			return 1
		case MSG_SELECT_PLACE, MSG_SELECT_DISFIELD:
			_ = binary.Read(pbuf, binary.LittleEndian, &player)
			pbuf.Next(5)
			d.WaitForResponse(player)
			SendBufferToPlayer(d.players[player], STOC_GAME_MSG, msgbuffer[offset:pbuf.Position()+3])
			return 1
		case MSG_SELECT_POSITION:
			_ = binary.Read(pbuf, binary.LittleEndian, &player)
			pbuf.Next(5)
			d.WaitForResponse(player)
			SendBufferToPlayer(d.players[player], STOC_GAME_MSG, msgbuffer[offset:pbuf.Position()+3])
			return 1
		case MSG_SELECT_COUNTER:
			_ = binary.Read(pbuf, binary.LittleEndian, &player)
			pbuf.Next(4)
			_ = binary.Read(pbuf, binary.LittleEndian, &count)
			pbuf.Next(int(count) * 9)
			d.WaitForResponse(player)
			SendBufferToPlayer(d.players[player], STOC_GAME_MSG, msgbuffer[offset:pbuf.Position()+3])
			return 1
		case MSG_SELECT_SUM:
			pbuf.Next(1)
			_ = binary.Read(pbuf, binary.LittleEndian, &player)
			pbuf.Next(6)
			_ = binary.Read(pbuf, binary.LittleEndian, &count)
			pbuf.Next(int(count) * 11)
			_ = binary.Read(pbuf, binary.LittleEndian, &count)
			pbuf.Next(int(count) * 11)
			d.WaitForResponse(player)
			SendBufferToPlayer(d.players[player], STOC_GAME_MSG, msgbuffer[offset:pbuf.Position()+3])
			return 1
		case MSG_SORT_CARD:
			_ = binary.Read(pbuf, binary.LittleEndian, &player)
			_ = binary.Read(pbuf, binary.LittleEndian, &count)
			pbuf.Next(int(count) * 7)
			d.WaitForResponse(player)
			SendBufferToPlayer(d.players[player], STOC_GAME_MSG, msgbuffer[offset:pbuf.Position()+3])
			return 1
		case MSG_CONFIRM_DECKTOP:
			_ = binary.Read(pbuf, binary.LittleEndian, &player)
			_ = binary.Read(pbuf, binary.LittleEndian, &count)
			pbuf.Next(int(count) * 7)
			var resendList []ClientInterface
			resendList = append(resendList, d.players[1])
			for i := range d.observers {
				resendList = append(resendList, d.observers[i])
			}
			SendBufferToPlayer(d.players[player], STOC_GAME_MSG, msgbuffer[offset:pbuf.Position()+3], resendList...)
		case MSG_CONFIRM_EXTRATOP:
			_ = binary.Read(pbuf, binary.LittleEndian, &player)
			_ = binary.Read(pbuf, binary.LittleEndian, &count)
			pbuf.Next(int(count) * 7)
			var resendList []ClientInterface
			resendList = append(resendList, d.players[1])
			for i := range d.observers {
				resendList = append(resendList, d.observers[i])
			}
			SendBufferToPlayer(d.players[player], STOC_GAME_MSG, msgbuffer[offset:pbuf.Position()+3], resendList...)
		case MSG_CONFIRM_CARDS:
			_ = binary.Read(pbuf, binary.LittleEndian, &player)
			_ = binary.Read(pbuf, binary.LittleEndian, &count)
			if pbuf.buf[5] != LOCATION_DECK {
				pbuf.Next(int(count) * 7)
				var resendList []ClientInterface
				resendList = append(resendList, d.players[1])
				for i := range d.observers {
					resendList = append(resendList, d.observers[i])
				}
				SendBufferToPlayer(d.players[player], STOC_GAME_MSG, msgbuffer[offset:pbuf.Position()+3], resendList...)
			} else {
				pbuf.Next(int(count) * 7)
				SendBufferToPlayer(d.players[player], STOC_GAME_MSG, msgbuffer[offset:pbuf.Position()+3])
			}
		case MSG_SHUFFLE_DECK:
			_ = binary.Read(pbuf, binary.LittleEndian, &player)
			var resendList []ClientInterface
			resendList = append(resendList, d.players[1])
			for i := range d.observers {
				resendList = append(resendList, d.observers[i])
			}
			SendBufferToPlayer(d.players[0], STOC_GAME_MSG, msgbuffer[offset:pbuf.Position()+3], resendList...)
		case MSG_SHUFFLE_HAND:
			_ = binary.Read(pbuf, binary.LittleEndian, &player)
			_ = binary.Read(pbuf, binary.LittleEndian, &count)
			SendBufferToPlayer(d.players[player], STOC_GAME_MSG, msgbuffer[offset:pbuf.Position()+3+int(count)*4])
			var (
				i uint8
			)
			for ; i < count; i++ {
				pbuf.Write([]byte{0, 0, 0, 0})
			}
			var resendList []ClientInterface
			for j := range d.observers {
				resendList = append(resendList, d.observers[j])
			}
			SendBufferToPlayer(d.players[1-player], STOC_GAME_MSG, msgbuffer[offset:pbuf.Position()+3], resendList...)
			d.RefreshHand(d.pDuel, player, 0x781fff, 0)
		case MSG_SHUFFLE_EXTRA:
			_ = binary.Read(pbuf, binary.LittleEndian, &player)
			_ = binary.Read(pbuf, binary.LittleEndian, &count)
			SendBufferToPlayer(d.players[player], STOC_GAME_MSG, msgbuffer[offset:pbuf.Position()+3+int(count)*4])
			var (
				i uint8
			)
			for ; i < count; i++ {
				pbuf.Write([]byte{0, 0, 0, 0})
			}
			var resendList []ClientInterface
			for j := range d.observers {
				resendList = append(resendList, d.observers[j])
			}
			SendBufferToPlayer(d.players[1-player], STOC_GAME_MSG, msgbuffer[offset:pbuf.Position()+3], resendList...)
			d.RefreshExtraDef(player)
		case MSG_REFRESH_DECK:
			pbuf.Next(1)
			var resendList []ClientInterface
			resendList = append(resendList, d.players[1])
			for j := range d.observers {
				resendList = append(resendList, d.observers[j])
			}
			SendBufferToPlayer(d.players[0], STOC_GAME_MSG, msgbuffer[offset:pbuf.Position()+3], resendList...)
		case MSG_SWAP_GRAVE_DECK:
			_ = binary.Read(pbuf, binary.LittleEndian, &player)
			var resendList []ClientInterface
			resendList = append(resendList, d.players[1])
			for j := range d.observers {
				resendList = append(resendList, d.observers[j])
			}
			SendBufferToPlayer(d.players[0], STOC_GAME_MSG, msgbuffer[offset:pbuf.Position()+3], resendList...)
			d.RefreshGraveDef(player)
		case MSG_REVERSE_DECK:
			var resendList []ClientInterface
			resendList = append(resendList, d.players[1])
			for j := range d.observers {
				resendList = append(resendList, d.observers[j])
			}
			SendBufferToPlayer(d.players[0], STOC_GAME_MSG, msgbuffer[offset:pbuf.Position()+3], resendList...)
		case MSG_DECK_TOP:
			pbuf.Next(6)
			var resendList []ClientInterface
			resendList = append(resendList, d.players[1])
			for j := range d.observers {
				resendList = append(resendList, d.observers[j])
			}
			SendBufferToPlayer(d.players[0], STOC_GAME_MSG, msgbuffer[offset:pbuf.Position()+3], resendList...)
		case MSG_SHUFFLE_SET_CARD:
			var loc uint8
			_ = binary.Read(pbuf, binary.LittleEndian, &loc)
			_ = binary.Read(pbuf, binary.LittleEndian, &count)
			pbuf.Next(int(count) * 8)
			var resendList []ClientInterface
			resendList = append(resendList, d.players[1])
			for j := range d.observers {
				resendList = append(resendList, d.observers[j])
			}
			SendBufferToPlayer(d.players[0], STOC_GAME_MSG, msgbuffer[offset:pbuf.Position()+3], resendList...)
			if loc == LOCATION_MZONE {
				d.RefreshMzone(d.pDuel, 0, 0x181fff, 0)
				d.RefreshMzone(d.pDuel, 1, 0x181fff, 0)
			} else {
				d.RefreshSzone(d.pDuel, 0, 0x181fff, 0)
				d.RefreshSzone(d.pDuel, 1, 0x181fff, 0)
			}
		case MSG_NEW_TURN:
			d.RefreshMzoneDef(0)
			d.RefreshMzoneDef(1)
			d.RefreshSzoneDef(0)
			d.RefreshSzoneDef(1)
			d.RefreshHandDef(0)
			d.RefreshHandDef(1)

			pbuf.Next(1)
			d.timeLimit[0] = d.hostInfo.TimeLimit
			d.timeLimit[1] = d.hostInfo.TimeLimit
			var resendList []ClientInterface
			resendList = append(resendList, d.players[1])
			for j := range d.observers {
				resendList = append(resendList, d.observers[j])
			}
			SendBufferToPlayer(d.players[0], STOC_GAME_MSG, msgbuffer[offset:pbuf.Position()+3], resendList...)

		case MSG_NEW_PHASE:
			pbuf.Next(2)
			var resendList []ClientInterface
			resendList = append(resendList, d.players[1])
			for j := range d.observers {
				resendList = append(resendList, d.observers[j])
			}
			SendBufferToPlayer(d.players[0], STOC_GAME_MSG, msgbuffer[offset:pbuf.Position()+3], resendList...)
			d.RefreshMzoneDef(0)
			d.RefreshMzoneDef(1)
			d.RefreshSzoneDef(0)
			d.RefreshSzoneDef(1)
			d.RefreshHandDef(0)
			d.RefreshHandDef(1)
		case MSG_MOVE:

			pbufw = pbuf.Position()
			pc := pbuf.buf[4]

			pl := pbuf.buf[5]
			//	/*int ps = pbuf[6];*/
			//	/*int pp = pbuf[7];*/
			cc := pbuf.buf[8]
			cl := pbuf.buf[9]
			cs := pbuf.buf[10]
			cp := pbuf.buf[11]
			pbuf.Next(16)
			SendBufferToPlayer(d.players[cc], STOC_GAME_MSG, msgbuffer[offset:pbuf.Position()+3])
			if (cl&(LOCATION_GRAVE+LOCATION_OVERLAY) == 0) && ((cl&(LOCATION_DECK+LOCATION_HAND) != 0) || (cp&POS_FACEDOWN) != 0) {
				for i := pbufw; i < 4; i++ {
					pbuf.buf[i] = 0
				}
			}
			var resendList []ClientInterface
			for j := range d.observers {
				resendList = append(resendList, d.observers[j])
			}
			SendBufferToPlayer(d.players[1-cc], STOC_GAME_MSG, msgbuffer[offset:pbuf.Position()+3], resendList...)

			if cl != 0 && (cl&LOCATION_OVERLAY) == 0 && (cl != pl || pc != cc) {
				d.RefreshSingleDef(cc, cl, cs)
			}

		case MSG_POS_CHANGE:
			cc := pbuf.buf[4]
			cl := pbuf.buf[5]
			cs := pbuf.buf[6]
			pp := pbuf.buf[7]
			cp := pbuf.buf[8]
			pbuf.Next(9)
			var resendList []ClientInterface
			resendList = append(resendList, d.players[1])
			for j := range d.observers {
				resendList = append(resendList, d.observers[j])
			}
			SendBufferToPlayer(d.players[0], STOC_GAME_MSG, msgbuffer[offset:pbuf.Position()+3], resendList...)
			if (pp&POS_FACEDOWN != 0) && (cp&POS_FACEUP != 0) {
				d.RefreshSingleDef(cc, cl, cs)
			}
		case MSG_SET:
			pbuf.Write([]byte{0, 0, 0, 0})
			pbuf.Next(4)
			var resendList []ClientInterface
			resendList = append(resendList, d.players[1])
			for j := range d.observers {
				resendList = append(resendList, d.observers[j])
			}
			SendBufferToPlayer(d.players[0], STOC_GAME_MSG, msgbuffer[offset:pbuf.Position()+3], resendList...)
		case MSG_SWAP:
			c1 := pbuf.buf[4]
			l1 := pbuf.buf[5]
			s1 := pbuf.buf[6]
			c2 := pbuf.buf[12]
			l2 := pbuf.buf[13]
			s2 := pbuf.buf[14]
			pbuf.Next(16)
			var resendList []ClientInterface
			resendList = append(resendList, d.players[1])
			for j := range d.observers {
				resendList = append(resendList, d.observers[j])
			}
			SendBufferToPlayer(d.players[0], STOC_GAME_MSG, msgbuffer[offset:pbuf.Position()+3], resendList...)
			d.RefreshSingleDef(c1, l1, s1)
			d.RefreshSingleDef(c2, l2, s2)
		case MSG_FIELD_DISABLED:
			pbuf.Next(4)
			var resendList []ClientInterface
			resendList = append(resendList, d.players[1])
			for j := range d.observers {
				resendList = append(resendList, d.observers[j])
			}
			SendBufferToPlayer(d.players[0], STOC_GAME_MSG, msgbuffer[offset:pbuf.Position()+3], resendList...)
		case MSG_SUMMONING:
			pbuf.Next(8)
			var resendList []ClientInterface
			resendList = append(resendList, d.players[1])
			for j := range d.observers {
				resendList = append(resendList, d.observers[j])
			}
			SendBufferToPlayer(d.players[0], STOC_GAME_MSG, msgbuffer[offset:pbuf.Position()+3], resendList...)
		case MSG_SUMMONED:
			var resendList []ClientInterface
			resendList = append(resendList, d.players[1])
			for j := range d.observers {
				resendList = append(resendList, d.observers[j])
			}
			SendBufferToPlayer(d.players[0], STOC_GAME_MSG, msgbuffer[offset:pbuf.Position()+3], resendList...)
			d.RefreshMzoneDef(0)
			d.RefreshMzoneDef(1)
			d.RefreshSzoneDef(0)
			d.RefreshSzoneDef(1)
		case MSG_SPSUMMONED:
			pbuf.Next(8)
			var resendList []ClientInterface
			resendList = append(resendList, d.players[1])
			for j := range d.observers {
				resendList = append(resendList, d.observers[j])
			}
			SendBufferToPlayer(d.players[0], STOC_GAME_MSG, msgbuffer[offset:pbuf.Position()+3], resendList...)
		case MSG_SPSUMMONING:
			var resendList []ClientInterface
			resendList = append(resendList, d.players[1])
			for j := range d.observers {
				resendList = append(resendList, d.observers[j])
			}
			SendBufferToPlayer(d.players[0], STOC_GAME_MSG, msgbuffer[offset:pbuf.Position()+3], resendList...)
			d.RefreshMzoneDef(0)
			d.RefreshMzoneDef(1)
			d.RefreshSzoneDef(0)
			d.RefreshSzoneDef(1)
		//NetServer::SendBufferToPlayer(players[0], STOC_GAME_MSG, offset, pbuf - offset);
		//NetServer::ReSendToPlayer(players[1]);
		//	for(auto oit = observers.begin(); oit != observers.end(); ++oit)
		//NetServer::ReSendToPlayer(*oit);
		//	RefreshMzone(0);
		//	RefreshMzone(1);
		//	RefreshSzone(0);
		//	RefreshSzone(1);
		//	break;
		//}
		case MSG_FLIPSUMMONING:
		//	RefreshSingle(pbuf[4], pbuf[5], pbuf[6]);
		//	pbuf += 8;
		//NetServer::SendBufferToPlayer(players[0], STOC_GAME_MSG, offset, pbuf - offset);
		//NetServer::ReSendToPlayer(players[1]);
		//	for(auto oit = observers.begin(); oit != observers.end(); ++oit)
		//NetServer::ReSendToPlayer(*oit);
		//	break;
		//}
		case MSG_FLIPSUMMONED:
		//NetServer::SendBufferToPlayer(players[0], STOC_GAME_MSG, offset, pbuf - offset);
		//NetServer::ReSendToPlayer(players[1]);
		//	for(auto oit = observers.begin(); oit != observers.end(); ++oit)
		//NetServer::ReSendToPlayer(*oit);
		//	RefreshMzone(0);
		//	RefreshMzone(1);
		//	RefreshSzone(0);
		//	RefreshSzone(1);
		//	break;
		//}
		case MSG_CHAINING:
		//	pbuf += 16;
		//NetServer::SendBufferToPlayer(players[0], STOC_GAME_MSG, offset, pbuf - offset);
		//NetServer::ReSendToPlayer(players[1]);
		//	for(auto oit = observers.begin(); oit != observers.end(); ++oit)
		//NetServer::ReSendToPlayer(*oit);
		//	break;
		//}
		case MSG_CHAINED:
		//	pbuf++;
		//NetServer::SendBufferToPlayer(players[0], STOC_GAME_MSG, offset, pbuf - offset);
		//NetServer::ReSendToPlayer(players[1]);
		//	for(auto oit = observers.begin(); oit != observers.end(); ++oit)
		//NetServer::ReSendToPlayer(*oit);
		//	RefreshMzone(0);
		//	RefreshMzone(1);
		//	RefreshSzone(0);
		//	RefreshSzone(1);
		//	RefreshHand(0);
		//	RefreshHand(1);
		//	break;
		//}
		case MSG_CHAIN_SOLVING:
		//	pbuf++;
		//NetServer::SendBufferToPlayer(players[0], STOC_GAME_MSG, offset, pbuf - offset);
		//NetServer::ReSendToPlayer(players[1]);
		//	for(auto oit = observers.begin(); oit != observers.end(); ++oit)
		//NetServer::ReSendToPlayer(*oit);
		//	break;
		//}
		case MSG_CHAIN_SOLVED:
		//	pbuf++;
		//NetServer::SendBufferToPlayer(players[0], STOC_GAME_MSG, offset, pbuf - offset);
		//NetServer::ReSendToPlayer(players[1]);
		//	for(auto oit = observers.begin(); oit != observers.end(); ++oit)
		//NetServer::ReSendToPlayer(*oit);
		//	RefreshMzone(0);
		//	RefreshMzone(1);
		//	RefreshSzone(0);
		//	RefreshSzone(1);
		//	RefreshHand(0);
		//	RefreshHand(1);
		//	break;
		//}
		case MSG_CHAIN_END:
		//NetServer::SendBufferToPlayer(players[0], STOC_GAME_MSG, offset, pbuf - offset);
		//NetServer::ReSendToPlayer(players[1]);
		//	for(auto oit = observers.begin(); oit != observers.end(); ++oit)
		//NetServer::ReSendToPlayer(*oit);
		//	RefreshMzone(0);
		//	RefreshMzone(1);
		//	RefreshSzone(0);
		//	RefreshSzone(1);
		//	RefreshHand(0);
		//	RefreshHand(1);
		//	break;
		//}
		case MSG_CHAIN_NEGATED:
		//	pbuf++;
		//NetServer::SendBufferToPlayer(players[0], STOC_GAME_MSG, offset, pbuf - offset);
		//NetServer::ReSendToPlayer(players[1]);
		//	for(auto oit = observers.begin(); oit != observers.end(); ++oit)
		//NetServer::ReSendToPlayer(*oit);
		//	break;
		//}
		case MSG_CHAIN_DISABLED:
		//	pbuf++;
		//NetServer::SendBufferToPlayer(players[0], STOC_GAME_MSG, offset, pbuf - offset);
		//NetServer::ReSendToPlayer(players[1]);
		//	for(auto oit = observers.begin(); oit != observers.end(); ++oit)
		//NetServer::ReSendToPlayer(*oit);
		//	break;
		//}
		case MSG_CARD_SELECTED:
		//	player = BufferIO::ReadInt8(pbuf);
		//	count = BufferIO::ReadInt8(pbuf);
		//	pbuf += count * 4;
		//	break;
		//}
		case MSG_RANDOM_SELECTED:
		//	player = BufferIO::ReadInt8(pbuf);
		//	count = BufferIO::ReadInt8(pbuf);
		//	pbuf += count * 4;
		//NetServer::SendBufferToPlayer(players[player], STOC_GAME_MSG, offset, pbuf - offset);
		//NetServer::ReSendToPlayer(players[1]);
		//	for(auto oit = observers.begin(); oit != observers.end(); ++oit)
		//NetServer::ReSendToPlayer(*oit);
		//	break;
		//}
		case MSG_BECOME_TARGET:
		//	count = BufferIO::ReadInt8(pbuf);
		//	pbuf += count * 4;
		//NetServer::SendBufferToPlayer(players[0], STOC_GAME_MSG, offset, pbuf - offset);
		//NetServer::ReSendToPlayer(players[1]);
		//	for(auto oit = observers.begin(); oit != observers.end(); ++oit)
		//NetServer::ReSendToPlayer(*oit);
		//	break;
		//}
		case MSG_DRAW:

			_ = binary.Read(pbuf, binary.LittleEndian, &player)
			_ = binary.Read(pbuf, binary.LittleEndian, &count)
			fmt.Println("MSG_DRAW", player, count)
			var cards = make([]uint32, count)
			_ = binary.Read(pbuf, binary.LittleEndian, &cards)

			SendPacketToPlayer(d.players[player], STOC_GAME_MSG, BytesPacket(pbuf.OffsetBytes(offset)))
			index := 3
			for i := 0; i < int(count); i++ {
				if msgbuffer[index+3]&0x80 == 0 {
					for j := 0; j < 4; j++ {
						msgbuffer[index+j] = 0
					}
				}
				index += 4
			}
			SendPacketToPlayer(d.players[1-player], STOC_GAME_MSG, BytesPacket(pbuf.OffsetBytes(offset)))
			for i := range d.observers {
				SendPacketToPlayer(d.observers[i], STOC_GAME_MSG, BytesPacket(pbuf.OffsetBytes(offset)))
			}
			//for(auto oit = observers.begin(); oit != observers.end(); ++oit)
		//NetServer::ReSendToPlayer(*oit);

		case MSG_DAMAGE:
		//	pbuf += 5;
		//NetServer::SendBufferToPlayer(players[0], STOC_GAME_MSG, offset, pbuf - offset);
		//NetServer::ReSendToPlayer(players[1]);
		//	for(auto oit = observers.begin(); oit != observers.end(); ++oit)
		//NetServer::ReSendToPlayer(*oit);
		//	break;
		//}
		case MSG_RECOVER:
		//	pbuf += 5;
		//NetServer::SendBufferToPlayer(players[0], STOC_GAME_MSG, offset, pbuf - offset);
		//NetServer::ReSendToPlayer(players[1]);
		//	for(auto oit = observers.begin(); oit != observers.end(); ++oit)
		//NetServer::ReSendToPlayer(*oit);
		//	break;
		//}
		case MSG_EQUIP:
		//	pbuf += 8;
		//NetServer::SendBufferToPlayer(players[0], STOC_GAME_MSG, offset, pbuf - offset);
		//NetServer::ReSendToPlayer(players[1]);
		//	for(auto oit = observers.begin(); oit != observers.end(); ++oit)
		//NetServer::ReSendToPlayer(*oit);
		//	break;
		//}
		case MSG_LPUPDATE:
		//	pbuf += 5;
		//NetServer::SendBufferToPlayer(players[0], STOC_GAME_MSG, offset, pbuf - offset);
		//NetServer::ReSendToPlayer(players[1]);
		//	for(auto oit = observers.begin(); oit != observers.end(); ++oit)
		//NetServer::ReSendToPlayer(*oit);
		//	break;
		//}
		case MSG_UNEQUIP:
		//	pbuf += 4;
		//NetServer::SendBufferToPlayer(players[0], STOC_GAME_MSG, offset, pbuf - offset);
		//NetServer::ReSendToPlayer(players[1]);
		//	for(auto oit = observers.begin(); oit != observers.end(); ++oit)
		//NetServer::ReSendToPlayer(*oit);
		//	break;
		//}
		case MSG_CARD_TARGET:
		//	pbuf += 8;
		//NetServer::SendBufferToPlayer(players[0], STOC_GAME_MSG, offset, pbuf - offset);
		//NetServer::ReSendToPlayer(players[1]);
		//	for(auto oit = observers.begin(); oit != observers.end(); ++oit)
		//NetServer::ReSendToPlayer(*oit);
		//	break;
		//}
		case MSG_CANCEL_TARGET:
		//	pbuf += 8;
		//NetServer::SendBufferToPlayer(players[0], STOC_GAME_MSG, offset, pbuf - offset);
		//NetServer::ReSendToPlayer(players[1]);
		//	for(auto oit = observers.begin(); oit != observers.end(); ++oit)
		//NetServer::ReSendToPlayer(*oit);
		//	break;
		//}
		case MSG_PAY_LPCOST:
		//	pbuf += 5;
		//NetServer::SendBufferToPlayer(players[0], STOC_GAME_MSG, offset, pbuf - offset);
		//NetServer::ReSendToPlayer(players[1]);
		//	for(auto oit = observers.begin(); oit != observers.end(); ++oit)
		//NetServer::ReSendToPlayer(*oit);
		//	break;
		//}
		case MSG_ADD_COUNTER:
		//	pbuf += 7;
		//NetServer::SendBufferToPlayer(players[0], STOC_GAME_MSG, offset, pbuf - offset);
		//NetServer::ReSendToPlayer(players[1]);
		//	for(auto oit = observers.begin(); oit != observers.end(); ++oit)
		//NetServer::ReSendToPlayer(*oit);
		//	break;
		//}
		case MSG_REMOVE_COUNTER:
		//	pbuf += 7;
		//NetServer::SendBufferToPlayer(players[0], STOC_GAME_MSG, offset, pbuf - offset);
		//NetServer::ReSendToPlayer(players[1]);
		//	for(auto oit = observers.begin(); oit != observers.end(); ++oit)
		//NetServer::ReSendToPlayer(*oit);
		//	break;
		//}
		case MSG_ATTACK:
		//	pbuf += 8;
		//NetServer::SendBufferToPlayer(players[0], STOC_GAME_MSG, offset, pbuf - offset);
		//NetServer::ReSendToPlayer(players[1]);
		//	for(auto oit = observers.begin(); oit != observers.end(); ++oit)
		//NetServer::ReSendToPlayer(*oit);
		//	break;
		//}
		case MSG_BATTLE:
		//	pbuf += 26;
		//NetServer::SendBufferToPlayer(players[0], STOC_GAME_MSG, offset, pbuf - offset);
		//NetServer::ReSendToPlayer(players[1]);
		//	for(auto oit = observers.begin(); oit != observers.end(); ++oit)
		//NetServer::ReSendToPlayer(*oit);
		//	break;
		//}
		case MSG_ATTACK_DISABLED:
		//NetServer::SendBufferToPlayer(players[0], STOC_GAME_MSG, offset, pbuf - offset);
		//NetServer::ReSendToPlayer(players[1]);
		//	for(auto oit = observers.begin(); oit != observers.end(); ++oit)
		//NetServer::ReSendToPlayer(*oit);
		//	break;
		//}
		case MSG_DAMAGE_STEP_START:
		//NetServer::SendBufferToPlayer(players[0], STOC_GAME_MSG, offset, pbuf - offset);
		//NetServer::ReSendToPlayer(players[1]);
		//	for(auto oit = observers.begin(); oit != observers.end(); ++oit)
		//NetServer::ReSendToPlayer(*oit);
		//	RefreshMzone(0);
		//	RefreshMzone(1);
		//	break;
		//}
		case MSG_DAMAGE_STEP_END:
		//NetServer::SendBufferToPlayer(players[0], STOC_GAME_MSG, offset, pbuf - offset);
		//NetServer::ReSendToPlayer(players[1]);
		//	for(auto oit = observers.begin(); oit != observers.end(); ++oit)
		//NetServer::ReSendToPlayer(*oit);
		//	RefreshMzone(0);
		//	RefreshMzone(1);
		//	break;
		//}
		case MSG_MISSED_EFFECT:
		//	player = pbuf[0];
		//	pbuf += 8;
		//NetServer::SendBufferToPlayer(players[player], STOC_GAME_MSG, offset, pbuf - offset);
		//	break;
		//}
		case MSG_TOSS_COIN:
		//	player = BufferIO::ReadInt8(pbuf);
		//	count = BufferIO::ReadInt8(pbuf);
		//	pbuf += count;
		//NetServer::SendBufferToPlayer(players[0], STOC_GAME_MSG, offset, pbuf - offset);
		//NetServer::ReSendToPlayer(players[1]);
		//	for(auto oit = observers.begin(); oit != observers.end(); ++oit)
		//NetServer::ReSendToPlayer(*oit);
		//	break;
		//}
		case MSG_TOSS_DICE:
		//	player = BufferIO::ReadInt8(pbuf);
		//	count = BufferIO::ReadInt8(pbuf);
		//	pbuf += count;
		//NetServer::SendBufferToPlayer(players[0], STOC_GAME_MSG, offset, pbuf - offset);
		//NetServer::ReSendToPlayer(players[1]);
		//	for(auto oit = observers.begin(); oit != observers.end(); ++oit)
		//NetServer::ReSendToPlayer(*oit);
		//	break;
		//}
		case MSG_ROCK_PAPER_SCISSORS:
		//	player = BufferIO::ReadInt8(pbuf);
		//	WaitforResponse(player);
		//NetServer::SendBufferToPlayer(players[player], STOC_GAME_MSG, offset, pbuf - offset);
		//	return 1;
		//}
		case MSG_HAND_RES:
		//	pbuf += 1;
		//NetServer::SendBufferToPlayer(players[0], STOC_GAME_MSG, offset, pbuf - offset);
		//NetServer::ReSendToPlayer(players[1]);
		//	for (auto oit = observers.begin(); oit != observers.end(); ++oit)
		//NetServer::ReSendToPlayer(*oit);
		//	break;
		//}
		case MSG_ANNOUNCE_RACE:
		//	player = BufferIO::ReadInt8(pbuf);
		//	pbuf += 5;
		//	WaitforResponse(player);
		//NetServer::SendBufferToPlayer(players[player], STOC_GAME_MSG, offset, pbuf - offset);
		//	return 1;
		//}
		case MSG_ANNOUNCE_ATTRIB:
		//	player = BufferIO::ReadInt8(pbuf);
		//	pbuf += 5;
		//	WaitforResponse(player);
		//NetServer::SendBufferToPlayer(players[player], STOC_GAME_MSG, offset, pbuf - offset);
		//	return 1;
		//}
		case MSG_ANNOUNCE_CARD, MSG_ANNOUNCE_NUMBER:
		//	player = BufferIO::ReadInt8(pbuf);
		//	count = BufferIO::ReadUInt8(pbuf);
		//	pbuf += 4 * count;
		//	WaitforResponse(player);
		//NetServer::SendBufferToPlayer(players[player], STOC_GAME_MSG, offset, pbuf - offset);
		//	return 1;
		//}
		case MSG_CARD_HINT:
		//	pbuf += 9;
		//NetServer::SendBufferToPlayer(players[0], STOC_GAME_MSG, offset, pbuf - offset);
		//NetServer::ReSendToPlayer(players[1]);
		//	for(auto oit = observers.begin(); oit != observers.end(); ++oit)
		//NetServer::ReSendToPlayer(*oit);
		//	break;
		//}
		case MSG_PLAYER_HINT:
		//	pbuf += 6;
		//NetServer::SendBufferToPlayer(players[0], STOC_GAME_MSG, offset, pbuf - offset);
		//NetServer::ReSendToPlayer(players[1]);
		//	for(auto oit = observers.begin(); oit != observers.end(); ++oit)
		//NetServer::ReSendToPlayer(*oit);
		//	break;
		//}
		case MSG_MATCH_KILL:
			//	int code = BufferIO::ReadInt32(pbuf);
			//	if(match_mode) {
			//		match_kill = code;
			//	NetServer::SendBufferToPlayer(players[0], STOC_GAME_MSG, offset, pbuf - offset);
			//	NetServer::ReSendToPlayer(players[1]);
			//		for(auto oit = observers.begin(); oit != observers.end(); ++oit)
			//	NetServer::ReSendToPlayer(*oit);
			//	}
			//	break;
			//}
		}
	}
	return 0
}
func (d *SingleDuel) IsCreator(dp *DuelPlayer) bool {
	if d.HostPlayer == nil {
		return true
	}
	return d.HostPlayer == dp
}
func (d *SingleDuel) SetHostInfo(info host.HostInfo) {
	d.hostInfo = info
}
func (d *SingleDuel) WaitForResponse(playerId uint8) {
	d.LastResponse = playerId
	var msg = make([]byte, 3, 4)
	msg = append(msg, MSG_WAITING)
	SendBufferToPlayer(d.players[1-playerId], STOC_GAME_MSG, msg)
	if d.hostInfo.TimeLimit != 0 {
		var sctl = stoc.TimeLimit{
			Player:   playerId,
			LeftTime: d.timeLimit[playerId],
		}
		SendPacketToPlayer(d.players[0], STOC_TIME_LIMIT, sctl, d.players[1])
		d.players[playerId].Status = CTOS_TIME_CONFIRM
	} else {
		d.players[playerId].Status = CTOS_RESPONSE
	}
}
