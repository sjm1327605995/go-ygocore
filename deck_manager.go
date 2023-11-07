package main

import (
	"bufio"
	"encoding/binary"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"os"
)

const unknownString = "???"

type DeckManager struct {
	lfList []*LFList
}
type LFList struct {
	hash     uint32
	listName string
	content  map[uint32]uint8
}

var DkManager *DeckManager

func NewDeckManger() {
	DkManager = new(DeckManager)
}
func (d *DeckManager) LoadLFList() {

	d.LoadLFListSingle("expansions/lflist.conf")
	d.LoadLFListSingle("lflist.conf")
	var noLimit = &LFList{listName: "N/A"}
	d.lfList = append(d.lfList, noLimit)
}
func (d *DeckManager) LoadLFListSingle(path string) {

	file, err := os.Open(path)
	if err != nil {
		zap.S().Warn("path err:", err.Error())
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	cur := new(LFList)

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}
		if line[0] == '#' {
			continue
		}

		if line[0] == '!' {
			strBuffer := decodeUTF8(line[1:])
			cur = &LFList{
				hash:     0x7dfcee6a,
				listName: strBuffer,
				content:  make(map[uint32]uint8),
			}
			d.lfList = append(d.lfList, cur)
			continue
		}

		p := 0
		for line[p] != ' ' && line[p] != '\t' && line[p] != 0 {
			p++
		}

		if line[p] == 0 {
			continue
		}

		linebuf := line[:p]
		p++
		sa := p
		code, err := cast.ToUint32E(linebuf)
		if err != nil || code == 0 {
			continue
		}

		for line[p] == ' ' || line[p] == '\t' {
			p++
		}

		for line[p] != ' ' && line[p] != '\t' && line[p] != 0 {
			p++
		}

		linebuf = line[sa:p]
		count, err := cast.ToUint8E(linebuf)
		if err != nil {
			continue
		}

		if cur != nil {
			cur.content[code] = count
			cur.hash = cur.hash ^ ((code << 18) | (code >> 14)) ^ ((code << (27 + count)) | (code >> (5 - count)))
		}
	}

	if err := scanner.Err(); err != nil {
		zap.S().Error("path err:", err.Error())
		return
	}

}
func decodeUTF8(line string) string {
	strBuffer := []rune(line)
	sa := len(strBuffer)
	for strBuffer[sa-1] == '\r' || strBuffer[sa-1] == '\n' {
		sa--
	}
	return string(strBuffer[:sa])
}
func (d *DeckManager) GetLFListName(lfhash uint32) string {
	for _, list := range d.lfList {
		if list.hash == lfhash {
			return list.listName
		}
	}
	return unknownString
}

func (d *DeckManager) GetLFListContent(lfhash uint32) map[uint32]uint8 {
	for _, list := range d.lfList {
		if list.hash == lfhash {
			return list.content
		}
	}
	return nil
}
func (d *DeckManager) LoadDeck(deck *Deck, buff []byte, mainc, sidec int, isPackList bool) int32 {
	if deck == nil {
		deck = &Deck{}
	}
	var (
		code      int32
		errorCode int32
		cards     []int32
	)
	for i := 0; i < len(buff); i += 4 {
		v := binary.LittleEndian.Uint32(buff[i:])
		cards = append(cards, int32(v))
	}
	deck.Clear()

	for i := 0; i < mainc; i++ {
		var cd CardData

		code = cards[i]
		if !DataCache.GetData(code, &cd) {
			errorCode = code
		}
		if cd.Type&TYPE_TOKEN != 0 {
			continue
		} else if isPackList {
			deck.main = append(deck.main, DataCache.GetCodePointer(code))
			continue
		} else if cd.Type&(TYPE_FUSION|TYPE_SYNCHRO|TYPE_XYZ|TYPE_LINK) != 0 {
			if len(deck.extra) >= 15 {
				continue
			}
			deck.extra = append(deck.extra)
		} else if len(deck.main) < 60 {
			deck.main = append(deck.main, DataCache.GetCodePointer(code))
		}
	}
	for i := 0; i < sidec; i++ {
		var cd CardData
		code = cards[i+mainc]
		if !DataCache.GetData(code, &cd) {
			errorCode = code
		}
		if cd.Type&TYPE_TOKEN != 0 {
			continue
		}
		if len(deck.side) < 15 {
			deck.side = append(deck.side, DataCache.GetCodePointer(code))
		}
	}
	return errorCode
}
func (d *DeckManager) LoadSide(deck *Deck, buff []byte, mianc, sidec int) bool {
	var (
		pcount = make(map[int32]int32, 0)
		ncount = make(map[int32]int32, 0)
	)
	for i := range deck.main {
		pcount[deck.main[i].Code]++
	}
	for i := range deck.extra {
		pcount[deck.extra[i].Code]++
	}
	for i := range deck.side {
		pcount[deck.side[i].Code]++
	}
	var newDeck Deck
	d.LoadDeck(&newDeck, buff, mianc, sidec, false)
	if len(newDeck.main) != len(deck.main) || len(newDeck.extra) != len(deck.extra) {
		return false
	}
	for i := range newDeck.main {
		ncount[newDeck.main[i].Code]++
	}
	for i := range newDeck.extra {
		ncount[newDeck.extra[i].Code]++
	}
	for i := range newDeck.side {
		ncount[newDeck.side[i].Code]++
	}
	for i := range pcount {
		if pcount[i] != ncount[i] {
			return false
		}
	}
	deck = &newDeck
	return true
}

var RuleMap = []int{AVAIL_OCG, AVAIL_TCG, AVAIL_SC, AVAIL_CUSTOM, AVAIL_OCGTCG, 0}

func (d *DeckManager) CheckDeck(deck *Deck, lfhash uint32, rule uint8) int32 {
	var (
		ccount = make(map[int32]uint8)
	)
	list := d.GetLFListContent(lfhash)
	if list == nil {
		return 0
	}
	if len(deck.main) < 40 || len(deck.main) > 60 {
		return int32((DECKERROR_MAINCOUNT << 28) + len(deck.main))
	}
	if len(deck.extra) > 15 {
		return int32((DECKERROR_EXTRACOUNT << 28) + len(deck.extra))
	}
	if len(deck.side) > 15 {
		return int32((DECKERROR_SIDECOUNT << 28) + len(deck.side))
	}

	avail := RuleMap[cast.ToInt(rule)]
	gameRuleDeckError := checkCards(ccount, deck.main, mainCards, list, avail)
	if gameRuleDeckError != 0 {
		return gameRuleDeckError
	}
	gameRuleDeckError = checkCards(ccount, deck.extra, extraCards, list, avail)
	if gameRuleDeckError != 0 {
		return gameRuleDeckError
	}
	gameRuleDeckError = checkCards(ccount, deck.side, sideCards, list, avail)
	if gameRuleDeckError != 0 {
		return gameRuleDeckError
	}
	return 0
}

const (
	mainCards uint8 = iota
	extraCards
	sideCards
)

func checkCards(ccount map[int32]uint8, cards []*CardDataC, typ uint8, list map[uint32]uint8, avail int) int32 {
	var dc uint8
	var countError int32
	switch typ {
	case mainCards:
		countError = DECKERROR_MAINCOUNT
	case sideCards:
		countError = DECKERROR_SIDECOUNT
	case extraCards:
		countError = DECKERROR_EXTRACOUNT
	default:
		panic("unknown type")
	}
	for i := range cards {
		cit := cards[i]
		gameRuleDeckError := checkAvail(cit.Ot, uint32(avail))
		if gameRuleDeckError != 0 {
			return (gameRuleDeckError << 28) + cit.Code
		}
		if typ == mainCards {
			if cit.Type&(TYPE_FUSION|TYPE_SYNCHRO|TYPE_XYZ|TYPE_TOKEN|TYPE_LINK) != 0 {
				return DECKERROR_EXTRACOUNT << 28
			}
		}

		code := IFELSE(cit.Alias != 0, int32(cit.Alias), cit.Code)
		ccount[code]++
		dc = ccount[code]
		if dc > 3 {
			return countError<<28 + cit.Code
		}
		it := list[uint32(code)]
		it, ok := list[uint32(code)]
		if ok && dc > it {
			return (DECKERROR_LFLIST << 28) + cit.Code
		}
	}
	return 0
}
func checkAvail(ot uint32, avail uint32) int32 {
	if (ot & avail) == avail {
		return 0
	}

	if (ot&AVAIL_OCG) != 0 && !(avail == AVAIL_OCG) {
		return DECKERROR_OCGONLY
	}

	if (ot&AVAIL_TCG) != 0 && !(avail == AVAIL_TCG) {
		return DECKERROR_TCGONLY
	}

	return DECKERROR_NOTAVAIL
}
