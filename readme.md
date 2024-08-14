### 调用ocgcore动态库


### api

  ```

type ScriptReader func(scriptName string) (data []byte)
type MessageHandler func(msg string)
type CardReader func(cardId int32) *CardData

	CreateDuel        func(seed int32) uintptr //Create a the instance of the duel using a random number.
	StartDuel         func(pduel uintptr, options int32) // Starts the duel
	EndDuel           func(pduel uintptr) //ends the duel
	SetPlayerInfo     func(pduel uintptr, playerId, LP, startCount, drawCount int32) // sets the duel up
	GetLogMessage     func(pduel uintptr, buf []byte)
	GetMessage        func(pduel uintptr, buf []byte) int32
	Process           func(pduel uintptr) int32 //do a game tick
	NewCard           func(pduel uintptr, code uint32, owner, playerid, location, sequence, position uint8) // add a card to the duel state.
	QueryCard         func(pduel uintptr, playerid, location, sequence uint8, queryFlag int32, buf []byte, useCache int32) int32  //find out about a card in a specific spot.
	QueryFieldCount   func(pduel uintptr, playerid, location uint8) int32 // Get the number of cards in a specific field/zone.
	QueryFieldCard    func(pduel uintptr, playerId, location uint8, queryFlag int32, buf []byte, useCache int32) int32
	QueryFieldInfo    func(pduel uintptr, buf []byte) int32
	SetResponseI      func(pduel uintptr, value int32)
	SetResponseB      func(pduel uintptr, buf []byte)
	PreloadScript     func(pduel uintptr, script string, len int32) int32
	ScriptReader      ScriptReader //Interface provided returns scripts based on number that corresponds to a lua file, send in a string.
	CardReader        CardReader //Interface provided function that provides database information from the data table of cards.cdb.
	MessageHandler    MessageHandler //Interface provided function that handles errors



  ```


