package main

import (
	"github.com/sjm1327605995/go-ygocore/msg/ctos"
	"github.com/sjm1327605995/go-ygocore/msg/host"
)

var duelMode DuelMode

// HandleCTOSPacket 重构dp结构体优化调用链
func HandleCTOSPacket(dp *DuelPlayer, data []byte, length int) {

	pktType := data[0]
	if (pktType != ctos.CTOS_SURRENDER) && (pktType != ctos.CTOS_CHAT) && (dp.Status == 0xff || (dp.Status == 1 && dp.Status != pktType)) {
		return
	}
	data = data[1:]
	switch pktType {
	case ctos.CTOS_RESPONSE:
		if dp.game == nil || duelMode.PDuel() == 0 {
			return
		}
		n := 0
		if length > 64 {
			n = 64
		} else {
			n = length - 1
		}
		duelMode.GetResponse(dp, data[:n])
	case ctos.CTOS_TIME_CONFIRM:
		if dp.game == nil || duelMode.PDuel() == 0 {
			return
		}
		duelMode.TimeConfirm(dp)
	case ctos.CTOS_CHAT:
		if dp.game == nil {
			return
		}
		duelMode.Chat(dp, data)
	case ctos.CTOS_UPDATE_DECK:
		if dp.game == nil {
			return
		}

		duelMode.UpdateDeck(dp, data)
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
		if dp.game != nil || duelMode == nil {
			return
		}
		var pkt ctos.CreateGame
		pkt.Parse(data)
		switch pkt.Info.Mode {
		case MODE_SINGLE, MODE_MATCH:
			duelMode = new(SingleDuel)
		case MODE_TAG:
			panic("tag duel unsupported")
		}
		if pkt.Info.DuleRule > 5 {
			pkt.Info.DuleRule = 5
		}
		if pkt.Info.Mode > 2 {
			pkt.Info.Mode = 2
		}
		//TODO
		//unsigned int hash = 1;
		//		for(auto lfit = deckManager._lfList.begin(); lfit != deckManager._lfList.end(); ++lfit) {
		//			if(pkt->info.lflist == lfit->hash) {
		//				hash = pkt->info.lflist;
		//				break;
		//			}
		//		}
		//		if(hash == 1)
		//			pkt->info.lflist = deckManager._lfList[0].hash;
		//		duel_mode->host_info = pkt->info;
		//		BufferIO::CopyWStr(pkt->name, duel_mode->name, 20);
		//		BufferIO::CopyWStr(pkt->pass, duel_mode->pass, 20);
		duelMode.JoinGame(dp, data, true)
		//StartBroadcast();
	case ctos.CTOS_JOIN_GAME: //TODO 现在如果game为空就进行初始化
		if duelMode == nil {
			s := new(SingleDuel)
			s.hostInfo = host.HostInfo{
				DuleRule:      5,
				NoCheckDeck:   false,
				NoShuffleDeck: false,
				StartLp:       8000,
				StartHand:     5,
				DrawCount:     1,
				TimeLimit:     180,
			}
			duelMode = s

		}
		isCreator := duelMode.IsCreator(dp)
		duelMode.JoinGame(dp, data, isCreator)

	case ctos.CTOS_LEAVE_GAME:
		duelMode.LeaveGame(dp)
	case ctos.CTOS_SURRENDER:
		duelMode.Surrender(dp)
	case ctos.CTOS_HS_TODUELIST:
		if dp.game == nil || duelMode.PDuel() == 0 {
			return
		}
	//TODO
	//duelMode.ToDuelist(dp);
	case ctos.CTOS_HS_TOOBSERVER:
		if dp.game == nil || duelMode.PDuel() == 0 {
			return
		}
		duelMode.ToObserver(dp)
	case ctos.CTOS_HS_READY, ctos.CTOS_HS_NOTREADY:
		if dp.game == nil || duelMode.PDuel() != 0 {
			return
		}
		duelMode.PlayerReady(dp, (CTOS_HS_NOTREADY-pktType) != 0)
	case CTOS_HS_KICK:
		if dp.game == nil || duelMode.PDuel() == 0 {
			return
		}

		var pkt ctos.Kick
		pkt.Parse(data)
		duelMode.PlayerKick(dp, pkt.Pos)
	case CTOS_HS_START:
		if dp.game == nil || duelMode.PDuel() == 0 {
			return
		}
		duelMode.StartDuel(dp)
	}
}
func WSStr(arr []byte) []byte {
	var i int
	for ; i < len(arr)-1; i = i + 2 {
		if arr[i] == 0 && arr[i+1] == 0 {
			break
		}
	}
	for {
		if i > len(arr)-1 {
			break
		}
		arr[i] = 0
		i++
	}
	return arr
}
