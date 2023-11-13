package main

import (
	"fmt"
	"github.com/sjm1327605995/go-ygocore/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)
var DataCache = &DataManager{}

type DataManager struct {
	datas   map[int32]*CardDataC
	strings map[int32]*CardString
}

func InitDB() error {
	var (
		err       error
		dialector gorm.Dialector
	)

	switch config.Conf.DBType {
	case "sqlite":
		dialector = sqlite.Open(config.Conf.DBDsn)
	case "mysql":
		dialector = mysql.Open(config.Conf.DBDsn)
	case "postgres", "pg":
		dialector = postgres.Open(config.Conf.DBDsn)
	}
	db, err = gorm.Open(dialector)
	return err
}

type CardsRes struct {
	CardDataC
	CardString
}

func (d *DataManager) LoadDB() error {
	var results []CardsRes
	//这种加载效率太慢。后考虑优化
	d.datas = make(map[int32]*CardDataC, 10000)
	d.strings = make(map[int32]*CardString, 10000)
	err := db.Raw("SELECT * FROM datas INNER JOIN texts ON datas.id = texts.id").Scan(&results).Error
	if err != nil {
		return err
	}
	for i := range results {
		d.datas[results[i].Code] = &results[i].CardDataC
		d.strings[results[i].Code] = &results[i].CardString
	}
	return nil
}
func (d *DataManager) GetCodePointer(code int32) *CardDataC {
	return d.datas[code]
}

// GetData 只能用户查询单卡效率太慢
func (d *DataManager) GetData(code int32, cd *CardData) (has bool) {
	v, has := d.datas[code]
	if !has {
		return false
	}

	cd.Code = v.Code
	cd.Ot = v.Ot
	cd.Alias = v.Alias
	cd.SetCode = v.SetCode
	cd.Type = v.Type
	cd.Attack = v.Attack
	cd.Defense = v.Defense
	cd.Level = v.Level
	cd.Attribute = v.Attribute
	cd.Race = v.Race
	cd.LScale = v.LScale
	cd.RScale = v.RScale
	cd.LinkMarker = v.LinkMarker
	return has
}

func getDataForCore(code uint32, pdata *CardDataC) bool {
	//target,continuation of the code:

	val, ok := DataCache.datas[int32(code)]
	if !ok {
		fmt.Println("getCode", code)
		return false
	}

	pdata = val
	return true
}

//
