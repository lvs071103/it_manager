package main

import (
	"fmt"
	"github.com/twinj/uuid"
	_ "github.com/twinj/uuid"
	"math/rand"
	"time"
)

//获得定长字符串
//str 填充字符串
//length 获得定长的长度
//char 不够长时填充的字符
func GetFixedLenString(str string, length int, char byte) string {
	if len(str) == 0 {
		return ""
	}

	if len(str) == length {
		return str
	}

	//超出切后面
	if len(str) > length {
		return string(str[:length])
	}

	//缺少添加char
	if len(str) < length {
		slice := make([]byte, length-len(str))
		for k := range slice {
			slice[k] = char
		}
		return string(append(slice, []byte(str)...))
	}

	return ""
}

func RandomStr(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func main() {
	fmt.Println(RandomStr(15))
	fmt.Println(uuid.NewV4().String())
}
