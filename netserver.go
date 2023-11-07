package main

import (
	"github.com/sjm1327605995/go-ygocore/msg/ctos"
	"github.com/sjm1327605995/go-ygocore/msg/host"
	"unicode/utf16"
)

// HandleCTOSPacket 重构dp结构体优化调用链
func HandleCTOSPacket(dp *DuelPlayer, data []byte, length int) {

	pktType := data[0]
	if (pktType != ctos.CTOS_SURRENDER) && (pktType != ctos.CTOS_CHAT) && (dp.Status == 0xff || (dp.Status == 1 && dp.Status != pktType)) {
		return
	}
	data = data[1:]
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
		dp.game.GetResponse(dp, data[:n])
	case ctos.CTOS_TIME_CONFIRM:
		if dp.game == nil || dp.game.PDuel() == 0 {
			return
		}
		dp.game.TimeConfirm(dp)
	case ctos.CTOS_CHAT:
		if dp.game == nil {
			return
		}
		dp.game.Chat(dp, data)
	case ctos.CTOS_UPDATE_DECK:
		if dp.game == nil {
			return
		}

		dp.game.UpdateDeck(dp, data)
	case ctos.CTOS_HAND_RESULT:
		if dp.game == nil {
			return
		}
		var res ctos.HandResult
		res.Parse(data)
		dp.game.HandResult(dp, res.Res)
	case ctos.CTOS_TP_RESULT:
		if dp.game == nil {
			return
		}
		var res ctos.TPResult
		err := res.Parse(data)
		if err != nil {
			return
		}
		dp.game.TPResult(dp, res.Res)
	case ctos.CTOS_PLAYER_INFO:
		var pkt ctos.PlayerInfo
		_ = pkt.Parse(data)
		dp.Name = pkt.Name
	case ctos.CTOS_CREATE_GAME: //TODO 暂时请求未使用到 比较疑惑
		if dp.game != nil || dp.game == nil {
			return
		}
		var pkt ctos.CreateGame
		pkt.Parse(data)

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
			defaultRoom = &SingleDuel{DuelModeBase: DuelModeBase{Pass: pkt.Pass, Name: pkt.Name, RealName: name, RealPassword: password}}

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
		mode.JoinGame(dp, data, isCreator)
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
		arr := make([]byte, len(dp.Pass))
		for i := range dp.Pass {
			arr[i] = dp.Pass[i]
		}
		password := WSStr(arr)
		mode, isCreator := JoinOrCreateDuelRoom(password, defaultRoom)
		mode.JoinGame(dp, data, isCreator)
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
		pkt.Parse(data)
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
