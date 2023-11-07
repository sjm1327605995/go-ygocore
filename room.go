package main

import "sync"

var rooms sync.Map

func JoinOrCreateDuelRoom(password string, mode DuelMode) DuelMode {
	val, has := rooms.LoadOrStore(password, mode)
	if has {
		return val.(DuelMode)
	}
	return mode
}
