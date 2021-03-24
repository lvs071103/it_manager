package tools

import (
	"strconv"
	"unsafe"
)

func StrToInt(str string) int {
	valueInt64, _ := strconv.ParseInt(str, 10, 64)
	valueInt := *(*int)(unsafe.Pointer(&valueInt64))
	return valueInt
}
