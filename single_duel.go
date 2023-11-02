package main

import (
	"encoding/binary"
	"fmt"
)

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
	d.RefreshExtra(d.pduel, player, 0xe81fff, 1)
}
func (d *SingleDuel) RefreshMzoneDef(player uint8) {
	d.RefreshExtra(d.pduel, player, 0x881fff, 1)
}

func (d *SingleDuel) RefreshSzoneDef(player uint8) {
	d.RefreshSzone(d.pduel, player, 0x681fff, 1)
}
func (d *SingleDuel) RefreshHandDef(player uint8) {
	d.RefreshHand(d.pduel, player, 0x681fff, 1)
}

// void RefreshSingle(int player, int location, int sequence, int flag = 0xf81fff);
//void RefreshGrave(int player, int flag = 0x81fff, int use_cache = 1);

func (d *SingleDuel) RefreshExtra(pdule uintptr, player uint8, flag, use_cache int32) {
	fmt.Println("RefreshExtra")
	var (
		originBuf = ocgPool.Get().([]byte)
		qbuf      = originBuf[3:]
	)
	defer ocgPool.Put(originBuf)
	qbuf[0] = MSG_UPDATE_DATA
	qbuf[1] = player
	qbuf[2] = LOCATION_EXTRA
	length := QueryFieldCard(pdule, player, LOCATION_EXTRA, flag, qbuf[3:], use_cache)

	_ = SendBufferToPlayer(d.players[player], STOC_GAME_MSG, originBuf[:length+6])

}
func (d *SingleDuel) RefreshMzone(pdule uintptr, player uint8, flag, use_cache int32) {
	fmt.Println("RefreshMzone")
	var (
		originBuf = ocgPool.Get().([]byte)
		qbuf      = originBuf[3:]
	)
	defer ocgPool.Put(originBuf)
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
	SendBufferToPlayer(d.players[1-player], STOC_GAME_MSG, originBuf[:length+6])
	//	for(auto pit = observers.begin(); pit != observers.end(); ++pit)
	//NetServer::ReSendToPlayer(*pit);
}
func (d *SingleDuel) RefreshSzone(pdule uintptr, player uint8, flag, use_cache int32) {
	fmt.Println("RefreshSzone")
	var (
		originBuf = ocgPool.Get().([]byte)
		qbuf      = originBuf[3:]
	)
	defer ocgPool.Put(originBuf)
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
	SendBufferToPlayer(d.players[1-player], STOC_GAME_MSG, qbuf[:length])
	//	for(auto pit = observers.begin(); pit != observers.end(); ++pit)
	//NetServer::ReSendToPlayer(*pit);
}
func (d *SingleDuel) RefreshHand(pdule uintptr, player uint8, flag, use_cache int32) {
	fmt.Println("RefreshHand")
	var (
		originBuf = ocgPool.Get().([]byte)
		qbuf      = originBuf[3:]
	)
	defer ocgPool.Put(originBuf)
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
	SendBufferToPlayer(d.players[1-player], STOC_GAME_MSG, qbuf[:length])
	//	for(auto pit = observers.begin(); pit != observers.end(); ++pit)
	//NetServer::ReSendToPlayer(*pit);
}

type SingleDuel struct {
	pduel        uintptr
	players      []ClientInterface
	engineBuffer []byte
}

func (d *SingleDuel) Chat(dp *DuelPlayer, buff []byte) {
	//TODO implement me
	panic("implement me")
}

func (d *SingleDuel) JoinGame(dp *DuelPlayer, buff []byte, flag bool) {
	//TODO implement me
	panic("implement me")
}

func (d *SingleDuel) LeaveGame(dp *DuelPlayer) {
	//TODO implement me
	panic("implement me")
}

func (d *SingleDuel) ToObserver(dp *DuelPlayer) {
	//TODO implement me
	panic("implement me")
}

func (d *SingleDuel) PlayerReady(dp *DuelPlayer, isReady bool) {
	//TODO implement me
	panic("implement me")
}

func (d *SingleDuel) PlayerKick(dp *DuelPlayer, pos uint8) {
	//TODO implement me
	panic("implement me")
}

func (d *SingleDuel) UpdateDeck(dp *DuelPlayer, buff []byte) error {
	//TODO implement me
	panic("implement me")
}

func (d *SingleDuel) StartDuel(dp *DuelPlayer) {
	//TODO implement me
	panic("implement me")
}

func (d *SingleDuel) HandResult(dp *DuelPlayer, uint82 uint8) {
	//TODO implement me
	panic("implement me")
}

func (d *SingleDuel) TPResult(dp *DuelPlayer, uint82 uint8) {
	//TODO implement me
	panic("implement me")
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

func (d *SingleDuel) PDuel() int64 {
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

		result := Process(d.pduel)
		engLen = result & 0xffff
		engFlag = result >> 16
		if engLen > 0 {
			_ = GetMessage(d.pduel, engineBuffer)
			stop = d.Analyze(engineBuffer[:engLen])
		}
	}

}
func (d *SingleDuel) Analyze(msgbuffer []byte) int {

	var (
		offset             int
		pbuf               = &Buffer{buf: msgbuffer}
		player, count, typ int8
	)

	for pbuf.Len() > 0 {
		offset = pbuf.Position()
		var engType uint8
		_ = binary.Read(pbuf, binary.LittleEndian, &engType)
		fmt.Println("engType", engType)
		switch engType {
		case MSG_RETRY:
			WaitforResponse()
			return 1
		case MSG_HINT:
			fmt.Println("MSG_HINT")
			_ = binary.Read(pbuf, binary.LittleEndian, &typ)
			_ = binary.Read(pbuf, binary.LittleEndian, &player)
			pbuf.Next(4)
			switch typ {
			case 1, 2, 3, 5:
				//发送消息给客户端
				fmt.Println(msgbuffer[offset : pbuf.Position()-offset])
			case 4, 6, 7, 8, 9, 11:
				//发送消息给客户端
				fmt.Println(msgbuffer[offset : pbuf.Position()-offset])
			case 10:
				//发送消息给客户端
				fmt.Println(msgbuffer[offset : pbuf.Position()-offset])
			}
		case MSG_WIN:
		case MSG_SELECT_BATTLECMD:
		case MSG_SELECT_IDLECMD:
		case MSG_SELECT_EFFECTYN:
		case MSG_SELECT_YESNO:
		case MSG_SELECT_OPTION:
		case MSG_SELECT_CARD, MSG_SELECT_TRIBUTE:
		case MSG_SELECT_UNSELECT_CARD:
		case MSG_SELECT_CHAIN:
		case MSG_SELECT_PLACE, MSG_SELECT_DISFIELD:
		case MSG_SELECT_POSITION:
		case MSG_SELECT_COUNTER:
		case MSG_SELECT_SUM:
		case MSG_SORT_CARD:
		case MSG_CONFIRM_DECKTOP:
		case MSG_CONFIRM_EXTRATOP:
		case MSG_CONFIRM_CARDS:
		case MSG_SHUFFLE_DECK:
		case MSG_SHUFFLE_HAND:
		case MSG_SHUFFLE_EXTRA:
		case MSG_REFRESH_DECK:
		case MSG_SWAP_GRAVE_DECK:
		case MSG_REVERSE_DECK:
		case MSG_DECK_TOP:
		case MSG_SHUFFLE_SET_CARD:
		case MSG_NEW_TURN:
			d.RefreshMzoneDef(0)
			d.RefreshMzoneDef(1)
			d.RefreshSzoneDef(0)
			d.RefreshSzoneDef(1)
			d.RefreshHandDef(0)
			d.RefreshHandDef(1)

			pbuf.Next(1)
			//			time_limit[0] = host_info.time_limit;
			//			time_limit[1] = host_info.time_limit;
			SendBufferToPlayer(d.players[0], STOC_GAME_MSG, append([]byte{0, 0, 0}, msgbuffer[offset:pbuf.Position()-offset]...), d.players[1])
		//			for(auto oit = observers.begin(); oit != observers.end(); ++oit)
		//				NetServer::ReSendToPlayer(*oit);
		//			break;
		case MSG_NEW_PHASE:
			pbuf.Next(2)
			SendBufferToPlayer(d.players[0], STOC_GAME_MSG, append([]byte{0, 0, 0}, msgbuffer[offset:pbuf.Position()-offset]...), d.players[1])

			//			for(auto oit = observers.begin(); oit != observers.end(); ++oit)
			//				NetServer::ReSendToPlayer(*oit);
			d.RefreshMzoneDef(0)
			d.RefreshMzoneDef(1)
			d.RefreshSzoneDef(0)
			d.RefreshSzoneDef(1)
			d.RefreshHandDef(0)
			d.RefreshHandDef(1)
		case MSG_MOVE:

		//	pbufw := pbuf
		//	 pc := msgbuffer[4];
		//	pl := msgbuffer[5];
		//	/*int ps = pbuf[6];*/
		//	/*int pp = pbuf[7];*/
		//	cc := msgbuffer[8];
		//	cl := msgbuffer[9];
		//	cs := msgbuffer[10];
		//	cp := msgbuffer[11];
		//	pbuf += 16;
		//NetServer::SendBufferToPlayer(players[cc], STOC_GAME_MSG, offset, pbuf - offset);
		//	if (!(cl & (LOCATION_GRAVE + LOCATION_OVERLAY)) && ((cl & (LOCATION_DECK + LOCATION_HAND)) || (cp & POS_FACEDOWN)))
		//		BufferIO::WriteInt32(pbufw, 0);
		//NetServer::SendBufferToPlayer(players[1 - cc], STOC_GAME_MSG, offset, pbuf - offset);
		//	for(auto oit = observers.begin(); oit != observers.end(); ++oit)
		//NetServer::ReSendToPlayer(*oit);
		//	if (cl != 0 && (cl & LOCATION_OVERLAY) == 0 && (cl != pl || pc != cc))
		//		RefreshSingle(cc, cl, cs);
		case MSG_POS_CHANGE:
		case MSG_SET:
		case MSG_SWAP:
		case MSG_FIELD_DISABLED:
		case MSG_SUMMONING:
		case MSG_SUMMONED:
		case MSG_SPSUMMONING:
		case MSG_SPSUMMONED:
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
			fmt.Println("  cards:", cards)

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
		//观众暂不考虑
		//	for(auto oit = observers.begin(); oit != observers.end(); ++oit)
		//NetServer::ReSendToPlayer(*oit);
		//	break;
		//}
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
func WaitforResponse() {
	fmt.Println("等待用户操作")
}
