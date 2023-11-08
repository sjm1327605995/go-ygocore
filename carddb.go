package main

import (
	_ "github.com/mattn/go-sqlite3"
)

type CardDataC struct {
	Code       int32  `gorm:"column:id"`
	Ot         uint32 `gorm:"column:ot"`
	Alias      uint32 `gorm:"column:alias"`
	SetCode    uint64 `gorm:"column:setcode"`
	Type       uint32 `gorm:"column:type"`
	Attack     int32  `gorm:"column:atk"`
	Defense    int32  `gorm:"column:def"`
	Level      uint32 `gorm:"column:level"`
	Attribute  uint32 `gorm:"column:attribute"`
	Race       uint32 `gorm:"column:race"`
	Category   int    `gorm:"column:category"`
	LScale     uint32 `gorm:"-"`
	RScale     uint32 `gorm:"-"`
	LinkMarker uint32 `gorm:"-"`
}

//	func (c *CardDataC)TableName()string  {
//		return "datas"
//	}
type CardData struct {
	Code       int32  `gorm:"column:id"`
	Ot         uint32 `gorm:"column:ot"`
	Alias      uint32 `gorm:"column:alias"`
	SetCode    uint64 `gorm:"column:setcode"`
	Type       uint32 `gorm:"column:type"`
	Attack     int32  `gorm:"column:atk"`
	Defense    int32  `gorm:"column:def"`
	Level      uint32 `gorm:"column:level"`
	Attribute  uint32 `gorm:"column:attribute"`
	Race       uint32 `gorm:"column:race"`
	LScale     uint32 `gorm:"-"`
	RScale     uint32 `gorm:"-"`
	LinkMarker uint32 `gorm:"-"`
}

type CardString struct {
	Name string
	Text string
	Desc string
}

var (
	UNKNOWN_STRING string = "???"
	datas          map[uint32]CardDataC
	cardStrings    map[uint32]CardString
	sysStrings     map[uint32]string
	victoryStrings map[uint32]string
	counterStrings map[uint32]string
	dataBasePath   string = "cards.cdb"
	stringsPath    string = "strings.conf"
)

//func (dm *DataManager) loadDB(dataBasePath string) bool {
//	if dataBasePath == "" {
//		return false
//	}
//
//	db, err := sql.Open("sqlite3", "file:"+dataBasePath)
//	if err != nil {
//		fmt.Println(err)
//		return false
//	}
//	defer db.Close()
//
//	rows, err := db.Query("SELECT * FROM datas")
//	if err != nil {
//		fmt.Println(err)
//		return false
//	}
//	defer rows.Close()
//
//	for rows.Next() {
//		var card CardDataC
//
//		err = rows.Scan(&card)
//		if err != nil {
//			fmt.Println(err)
//			return false
//		}
//		card.LScale = (card.Level >> 24) & 0xff
//		card.RScale = (card.Level >> 16) & 0xff
//		datas[uint32(card.Code)] = card
//		if card.Type&TYPE_LINK > 0 {
//			card.LinkMarker = cast.ToUint32(card.Defense)
//			card.Defense = 0
//		}
//	}
//	fmt.Println("data load complete")
//
//	rows, err = db.Query("SELECT * FROM texts")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer rows.Close()
//
//	for rows.Next() {
//		var id uint32
//		var name, text string
//		var desc [16]string
//		err = rows.Scan(&id, &name, &text, &desc[0], &desc[1], &desc[2], &desc[3], &desc[4], &desc[5], &desc[6], &desc[7], &desc[8], &desc[9], &desc[10], &desc[11], &desc[12], &desc[13], &desc[14], &desc[15])
//		if err != nil {
//			log.Fatal(err)
//		}
//		str := CardString{
//			Name: name,
//			Text: text,
//			Desc: desc[:],
//		}
//		cardStrings[id] = str
//	}
//	fmt.Println("text load complete")
//
//	return true
//}
//
//func (dm *DataManager) loadString(stringsPath string) bool {
//	file, err := os.Open(stringsPath)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer file.Close()
//
//	scanner := bufio.NewScanner(file)
//	for scanner.Scan() {
//		line := scanner.Text()
//		line = strings.TrimSpace(line)
//
//		if strings.HasPrefix(line, "#") {
//			continue
//		}
//
//		parts := strings.Fields(line)
//		if len(parts) < 3 {
//			continue
//		}
//
//		switch parts[0] {
//		case "!system":
//			key, err := cast.ToUint32E(parts[1])
//			if err != nil {
//				log.Fatal(err)
//			}
//			sysStrings[key] = parts[2]
//		case "!victory":
//
//			key, err := strconv.ParseInt(parts[1][2:], 16, 32)
//			if err != nil {
//				log.Fatal(err)
//			}
//			victoryStrings[uint32(key)] = parts[2]
//		case "!counter":
//			key, err := strconv.ParseInt(parts[1][2:], 16, 32)
//			if err != nil {
//				log.Fatal(err)
//			}
//			counterStrings[uint32(key)] = parts[2]
//		}
//	}
//
//	if err := scanner.Err(); err != nil {
//		log.Fatal(err)
//	}
//
//	return true
//}
//
//func NewDataManager() *DataManager {
//	dm := DataManager{}
//
//	sysStrings = make(map[uint32]string, 0)
//	victoryStrings = make(map[uint32]string, 0)
//	counterStrings = make(map[uint32]string, 0)
//	cardStrings = make(map[uint32]CardString)
//	// read cards.cdb data to datas
//	datas = make(map[uint32]CardDataC)
//	if !dm.loadDB(dataBasePath) {
//		panic("read data failed")
//	}
//
//	if !dm.loadString(stringsPath) {
//		panic("read strings failed")
//	}
//
//	return &dm
//}
//
//func getDataForCore(code uint32, pdata *CardDataC) bool {
//	//target,continuation of the code:
//
//	val, ok := datas[code]
//	if !ok {
//		return false
//	}
//
//	*pdata = val
//	return true
//}
////
//func getData(code uint32) *CardDataC {
//	if len(datas) == 0 {
//		return nil
//	}
//
//	card, ok := datas[code]
//	if !ok {
//		return nil
//	}
//
//	return &card
//}
//
//func getCardDesc(code uint32) *CardString {
//	if len(cardStrings) == 0 {
//		return nil
//	}
//
//	str, ok := cardStrings[code]
//	if !ok {
//		return nil
//	}
//
//	return &str
//}
//
//func getDesc(key int) string {
//	if key < 10000 {
//		return getSysString(key)
//	}
//
//	code := (key >> 4) & 0x0fffffff
//	offset := key & 0xf
//
//	cardString, ok := cardStrings[uint32(code)]
//	if !ok {
//		return UNKNOWN_STRING
//	}
//
//	if offset < len(cardString.desc) {
//		return cardString.desc[offset]
//	}
//
//	return UNKNOWN_STRING
//}
//
//func formatLocation(location, sequence int) string {
//	if location == LOCATION_SZONE {
//		if sequence < 5 {
//			return getSysString(1003)
//		} else if sequence == 5 {
//			return getSysString(1008)
//		} else {
//			return getSysString(1009)
//		}
//	}
//
//	filter := 1
//	i := 1000
//	for filter != 0x100 && filter != location {
//		filter <<= 1
//		i++
//	}
//
//	if filter == location {
//		return getSysString(i)
//	}
//
//	return UNKNOWN_STRING
//}
//
//func getSysString(key int) string {
//	str, ok := sysStrings[uint32(key)]
//	if !ok {
//		return ""
//	}
//
//	return str
//}
//
//func getVictoryString(key int) string {
//	str, ok := victoryStrings[uint32(key)]
//	if !ok {
//		return ""
//	}
//
//	return str
//}
//
//func getCounterString(key int) string {
//	str, ok := counterStrings[uint32(key)]
//	if !ok {
//		return ""
//	}
//
//	return str
//}
