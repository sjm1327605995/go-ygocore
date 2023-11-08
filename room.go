package main

import "sync"

var rooms sync.Map

func JoinOrCreateDuelRoom(password string, mode DuelMode) (DuelMode, bool) {
	val, has := rooms.LoadOrStore(password, mode)
	if has {
		return val.(DuelMode), false
	}

	return mode, true
}
