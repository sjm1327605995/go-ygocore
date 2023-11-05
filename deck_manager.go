package main

import (
	"encoding/binary"
	"errors"
)

type DeckManager struct {
}

func (d *DeckManager) LoadDeck(deck *Deck, buff []byte, mainc, sidec int, isPackList bool) (int, error) {
	if deck == nil {
		deck = &Deck{}
	}
	var (
		code      int32
		errorCode int
		cards     []int32
	)
	for i := 0; i < len(buff); i += 4 {
		v := binary.LittleEndian.Uint32(buff[i:])
		cards = append(cards, int32(v))
	}
	if len(cards) != mainc+sidec {
		return 0, errors.New("卡牌解析长度解析错误")
	}
	deck.Clear()
	for i := 0; i < mainc; i++ {
		code = cards[i]
	}
}
func (d *DeckManager) GetData(code int32, cd *CardDataC) bool {

}
