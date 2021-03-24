package tools

import (
	"crypto/md5"
	"fmt"
)

func GetMD5(str string) string {
	strToByte := []byte(str)
	byteRs := md5.Sum(strToByte)
	result := fmt.Sprintf("%x", byteRs)
	return result
}
