package config

var Conf *Config

type Config struct {
	DBType   string
	DBDsn    string
	HostInfo HostInfo
}

func InitConf() {
	Conf = &Config{
		DBType: "sqlite",
		DBDsn:  "cards.cdb",
	}
}

type HostInfo struct {
	Lflist        uint32
	Rule          uint8
	Mode          uint8
	DuleRule      uint8
	NoCheckDeck   bool
	NoShuffleDeck bool
	Unknown       uint16
	Unknown1      uint8
	StartLp       uint32
	StartHand     uint8
	DrawCount     uint8
	TimeLimit     uint16
}
