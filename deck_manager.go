package main

import (
	"encoding/binary"
)

type DeckManager struct {
}

var DkManager *DeckManager

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
	var (
		cd CardData
	)
	for i := 0; i < mainc; i++ {
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
