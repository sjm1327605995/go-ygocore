package ygocore

/*
#cgo CFLAGS: -Iinclude
#cgo windows,amd64 LDFLAGS: -L${SRCDIR}/libs -locgcore
#include "ocgcore.h"
*/
import "C"
import (
	"fmt"
	"math/rand"
	"os"
	"unsafe"
)

//export goScriptReader
func goScriptReader(scriptName *C.char, slen *C.int) *C.uchar {
	// 将C字符串转换为Go字符串

	*slen = 0
	s := C.GoString(scriptName)
	fmt.Println(s)
	// 调用适当的函数读取脚本内容
	data, _ := os.ReadFile(C.GoString(scriptName))
	if len(data) == 0 {
		// 处理错误
		return (*C.uchar)(nil)
	}

	// 将数据长度设置到slen指针
	*slen = C.int(len(data))
	// 创建C字节数组并将数据复制到其中
	return (*C.uchar)(C.CBytes(data))

}

//export goMessageHandler
func goMessageHandler(data unsafe.Pointer, size C.uint32_t) {
	fmt.Println("goMessageHandler")
	// 处理消息
}

//export goCardReader
func goCardReader(cardID C.uint32_t, data *C.card_data) C.uint32_t {

	//TODO 这里进行了内存拷贝需要重新操作下

	return 0

}

func CreateGame() int64 {
	seed := rand.Int31()
	fmt.Println("seed", seed)
	pDuel := C.create_duel(C.int(seed))
	fmt.Println(pDuel)
	return int64(pDuel)
}
