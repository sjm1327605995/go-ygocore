package main

import (
	"fmt"
	"github.com/sjm1327605995/go-ygocore/config"
	"math/rand"
)

//
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
	fmt.Println("pduel", pduel)
	duel := &SingleDuel{
		players: [2]*DuelPlayer{&DuelPlayer{}, &DuelPlayer{}},
		DuelModeBase: DuelModeBase{
			pDuel: pduel,
		},
	}

	SetPlayerInfo(pduel, 0, 8000, 5, 1)
	SetPlayerInfo(pduel, 1, 8000, 5, 1)
	var (
		mainCards = []uint32{
			37351133,
			37351133,
			37351133,
			26077387,
			26077387,
			26077387,
			23434538,
			23434538,
			23434538,
			73642296,
			14558127,
			14558127,
			14558127,
			97268402,
			97268402,
			97268402,
			18144506,
			25955749,
			99550630,
			35261759,
			35261759,
			73628505,
			63166095,
			63166095,
			67169062,
			32807846,
			8267140,
			8267140,
			51227866,
			51227866,
			52340444,
			25733157,
			24224830,
			65681983,
			98338152,
			98338152,
			24010609,
			24010609,
			50005218,
			50005218,
		}
		exidCards = []uint32{
			41209827,
			86066372,
			30194529,
			45819647,
			2857636,
			50588353,
			75147529,
			75147529,
			12421694,
			90673288,
			90673288,
			90673288,
			63288573,
			8491308,
			8491308,
		}
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
