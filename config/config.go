package config

var Conf *Config

type Config struct {
	DBType string
	DBDsn  string
}

func InitConf() {
	Conf = &Config{
		DBType: "sqlite",
		DBDsn:  "cards.cdb",
	}
}
