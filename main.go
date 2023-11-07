package main

import (
	"fmt"
	"github.com/panjf2000/gnet/v2"
	"github.com/sjm1327605995/go-ygocore/config"
	"log"
)

func main() {

	config.InitConf()
	InitDB()
	NewDeckManger()
	DkManager.LoadLFList()
	DataCache.LoadDB()
	addr := fmt.Sprintf("tcp://127.0.0.1:8080")

	//TCP 和UDP 都支持。对TCP分装的。可以通过TCP添加一层协议解析获取内容
	var srv = new(Server)

	log.Println("server exits:", gnet.Run(srv, addr, gnet.WithMulticore(true), gnet.WithReusePort(true), gnet.WithTicker(false)))
}
