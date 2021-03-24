package main

import (
	real2 "dominos_cmdb/tests/interface/real"
	"fmt"
	"time"
)

type Retriever interface {
	Get(url string) string
}

func download(r Retriever) string {
	return r.Get("http://www.baidu.com")
}

func inspect(r Retriever) {
	fmt.Printf("%T %v\n", r, r)
	switch v := r.(type) {
	case real2.Retriever:
		fmt.Println("UserAgent", v.UserAgent)
	}
}

func main() {
	var r Retriever
	r = real2.Retriever{UserAgent: "Mozilla/5.0", TimeOut: time.Minute}
	fmt.Printf("%T %v\n", r, r)
	//fmt.Println(download(r))
	inspect(r)
}
