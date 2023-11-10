package main

import (
	"fmt"
	"github.com/sjm1327605995/go-ygocore/config"
	"math/rand"
)

//func main() {
//
//	config.InitConf()
//	InitDB()
//	NewDeckManger()
//	DkManager.LoadLFList()
//	DataCache.LoadDB()
//	addr := fmt.Sprintf("tcp://127.0.0.1:8080")
//
//	//TCP 和UDP 都支持。对TCP分装的。可以通过TCP添加一层协议解析获取内容
//	var srv = NewServer()
//
//	log.Println("server exits:", gnet.Run(srv, addr, gnet.WithMulticore(true), gnet.WithReusePort(true), gnet.WithTicker(false)))
//}

func main() {
	config.InitConf()
	InitDB()
	NewDeckManger()
	DkManager.LoadLFList()
	DataCache.LoadDB()
	RegisterDo()
	n := rand.Int31n(10000)

	pduel := CreateDuel(n)
	duel := &SingleDuel{
		players: [2]*DuelPlayer{&DuelPlayer{}, &DuelPlayer{}},
		DuelModeBase: DuelModeBase{
			pDuel: pduel,
		},
	}
	SetPlayerInfo(pduel, 0, 8000, 5, 1)
	SetPlayerInfo(pduel, 1, 8000, 5, 1)
	var (
		mainCards = []uint32{14124483, 9411399, 9411399, 18094166, 18094166, 18094166, 40044918, 40044918, 59392529, 50720316, 50720316, 27780618, 27780618, 16605586, 16605586, 22865492, 22865492, 23434538, 23434538, 14558127, 14558127,
			13650422, 83965310, 81439173, 8949584, 8949584, 32807846, 52947044, 45906428, 24094653, 21143940, 21143940, 21143940, 48130397, 24224830, 24224830, 12071500, 24299458, 24299458, 10045474}
		exidCards = []uint32{73580471, 79606837, 79606837, 79606837, 21521304, 27552504, 1174075, 1174075, 1174075, 73898890, 73898890, 72336818, 41999284, 94259633, 94259633}
	)
	for i := len(mainCards) - 1; i >= 0; i-- {
		NewCard(pduel, mainCards[i], 0, 0, LOCATION_DECK, 0, POS_FACEDOWN_DEFENSE)
	}
	for i := len(exidCards) - 1; i >= 0; i-- {
		NewCard(pduel, exidCards[i], 0, 0, LOCATION_EXTRA, 0, POS_FACEDOWN_DEFENSE)
	}
	for i := len(mainCards) - 1; i >= 0; i-- {
		NewCard(pduel, mainCards[i], 1, 1, LOCATION_DECK, 0, POS_FACEDOWN_DEFENSE)
	}
	for i := len(exidCards) - 1; i >= 0; i-- {
		NewCard(pduel, exidCards[i], 1, 1, LOCATION_EXTRA, 0, POS_FACEDOWN_DEFENSE)
	}
	count1 := QueryFieldCount(pduel, 0, 0x1)
	count2 := QueryFieldCount(pduel, 0, 0x40)
	count3 := QueryFieldCount(pduel, 1, 0x1)
	count4 := QueryFieldCount(pduel, 1, 0x40)
	fmt.Println(count1, count2, count3, count4)
	duel.RefreshExtraDef(0)
	duel.RefreshExtraDef(1)
	opt := 5 << 16
	StartDuel(pduel, int32(opt))

	duel.Process()
}
