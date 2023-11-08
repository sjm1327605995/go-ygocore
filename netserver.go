package main

import (
	"bytes"
	"fmt"
	"github.com/sjm1327605995/go-ygocore/msg/ctos"
	"github.com/sjm1327605995/go-ygocore/msg/host"
	"github.com/sjm1327605995/go-ygocore/msg/stoc"
	"unicode/utf16"
)

// HandleCTOSPacket 重构dp结构体优化调用链
func HandleCTOSPacket(dp *DuelPlayer, buff *bytes.Buffer, length int) {

	pktType, _ := buff.ReadByte()
	fmt.Println(pktType)
	if (pktType != ctos.CTOS_SURRENDER) && (pktType != ctos.CTOS_CHAT) && (dp.Status == 0xff || (dp.Status == 1 && dp.Status != pktType)) {
		return
	}

	switch pktType {
	case ctos.CTOS_RESPONSE:
		if dp.game == nil || dp.game.PDuel() == 0 {
			return
		}
		n := 0
		if length > 64 {
			n = 64
		} else {
			n = length - 1
		}
		dp.game.GetResponse(dp, buff.Next(n))
	case ctos.CTOS_TIME_CONFIRM:
		if dp.game == nil || dp.game.PDuel() == 0 {
			return
		}
		dp.game.TimeConfirm(dp)
	case ctos.CTOS_CHAT:
		if dp.game == nil {
			return
		}
		dp.game.Chat(dp, buff.Bytes())
	case ctos.CTOS_UPDATE_DECK:
		if dp.game == nil {
			return
		}

		dp.game.UpdateDeck(dp, buff.Bytes())
	case ctos.CTOS_HAND_RESULT:
		if dp.game == nil {
			return
		}
		var res ctos.HandResult
		res.Parse(buff)
		dp.game.HandResult(dp, res.Res)
	case ctos.CTOS_TP_RESULT:
		if dp.game == nil {
			return
		}
		var res ctos.TPResult
		err := res.Parse(buff)
		if err != nil {
			return
		}
		dp.game.TPResult(dp, res.Res)
	case ctos.CTOS_PLAYER_INFO:
		var pkt ctos.PlayerInfo
		_ = pkt.Parse(buff)
		dp.Name = pkt.Name
	case ctos.CTOS_CREATE_GAME: //TODO 暂时请求未使用到 比较疑惑
		if dp.game != nil || dp.game == nil {
			return
		}
		var pkt ctos.CreateGame
		pkt.Parse(buff)

		arr := make([]byte, len(pkt.Name))
		for i := range dp.Name {
			arr[i] = dp.Name[i]
		}
		var (
			password = WSStr(arr)
		)
		arr = make([]byte, len(pkt.Name))
		for i := range dp.Name {
			arr[i] = dp.Name[i]
		}
		var (
			name = WSStr(arr)
		)
		var defaultRoom DuelMode = &DuelModeBase{}
		switch pkt.Info.Mode {
		case MODE_SINGLE, MODE_MATCH:
			defaultRoom = &SingleDuel{DuelModeBase: DuelModeBase{Pass: pkt.Pass, Name: pkt.Name, RealName: name}}

		case MODE_TAG:
			panic("tag duel unsupported")
		}
		if pkt.Info.DuleRule > 5 {
			pkt.Info.DuleRule = 5
		}
		if pkt.Info.Mode > 2 {
			pkt.Info.Mode = 2
		}

		var hash uint32 = 1
		for _, lfit := range DkManager.lfList {

			if pkt.Info.Lflist == lfit.hash {
				hash = pkt.Info.Lflist
			}
		}
		if hash == 1 {
			pkt.Info.Lflist = DkManager.lfList[0].hash
		}
		defaultRoom.SetHostInfo(pkt.Info)
		mode, isCreator := JoinOrCreateDuelRoom(password, defaultRoom)
		mode.JoinGame(dp, nil, isCreator)
		//StartBroadcast();
	case ctos.CTOS_JOIN_GAME: //TODO 现在如果game为空就进行初始化

		defaultRoom := new(SingleDuel)
		defaultRoom.hostInfo = host.HostInfo{
			Lflist:        DkManager.lfList[0].hash,
			DuleRule:      5,
			NoCheckDeck:   false,
			NoShuffleDeck: false,
			StartLp:       8000,
			StartHand:     5,
			DrawCount:     1,
			TimeLimit:     180,
		}
		var pkt ctos.JoinGame
		err := pkt.Parse(buff)
		if err != nil {
			var scem stoc.ErrorMsg
			scem.Msg = ERRMSG_JOINERROR
			scem.Code = 0
			SendPacketToPlayer(dp, STOC_ERROR_MSG, scem)
			DisconnectPlayer(dp)
			return
		}

		if pkt.Version != PRO_VERSION {
			var scem stoc.ErrorMsg
			scem.Msg = ERRMSG_VERERROR
			scem.Code = uint32(PRO_VERSION)
			SendPacketToPlayer(dp, STOC_ERROR_MSG, scem)
			DisconnectPlayer(dp)
			return
		}
		arr := make([]byte, len(pkt.Pass))
		for i := range pkt.Pass {
			arr[i] = pkt.Pass[i]
		}
		password := WSStr(arr)
		model, isCreator := JoinOrCreateDuelRoom(password, defaultRoom)
		//这里提前解析，不在JoinGame里面进行 buffer读取。buff实际是个空数组
		model.JoinGame(dp, nil, isCreator)
	case ctos.CTOS_LEAVE_GAME:
		dp.game.LeaveGame(dp)
	case ctos.CTOS_SURRENDER:
		dp.game.Surrender(dp)
	case ctos.CTOS_HS_TODUELIST:
		if dp.game == nil || dp.game.PDuel() != 0 {
			return
		}
	//TODO
	//duelMode.ToDuelist(dp);
	case ctos.CTOS_HS_TOOBSERVER:
		if dp.game == nil || dp.game.PDuel() != 0 {
			return
		}
		dp.game.ToObserver(dp)
	case ctos.CTOS_HS_READY, ctos.CTOS_HS_NOTREADY:
		if dp.game == nil || dp.game.PDuel() != 0 {
			return
		}
		dp.game.PlayerReady(dp, (CTOS_HS_NOTREADY-pktType) != 0)
	case CTOS_HS_KICK:
		if dp.game == nil || dp.game.PDuel() != 0 {
			return
		}

		var pkt ctos.Kick
		pkt.Parse(buff)
		dp.game.PlayerKick(dp, pkt.Pos)
	case CTOS_HS_START:
		if dp.game == nil || dp.game.PDuel() != 0 {
			return
		}
		dp.game.StartDuel(dp)
	}
}
func WSStr(bytes []byte) string {

	total := len(bytes) / 2
	// 将字节数组转换为 uint16 数组
	uint16s := make([]uint16, 0, total)
	for i := 0; i < total; i++ {
		n := uint16(bytes[i*2]) | (uint16(bytes[i*2+1]) << 8)
		if n == 0 {
			break
		}
		uint16s = append(uint16s, n)

	}

	// 解码并构建字符串
	runes := utf16.Decode(uint16s)
	str := string(runes)

	return str
}
