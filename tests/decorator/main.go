package main

import (
	"dominos_cmdb/models"
	"fmt"
	"github.com/astaxie/beego/orm"
)

func decorator(f func(s string)) func(s string) {
	return func(s string) {
		fmt.Println("start")
		f(s)
		fmt.Println("Done")
	}
}

func Hello(s string) {
	fmt.Println(s)
}

func UserPerHandler(f func(uid int)) func(uid int) {
	return func(uid int) {
		o := orm.NewOrm()
		currentUser := models.User{Id: uid}
		_, err := o.LoadRelated(&currentUser, "Permission")
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(currentUser.Permission)
		f(uid)
	}
}

func GetUserId(uid int) {
	fmt.Println(uid)
}

func main() {
	decorator(Hello)("Hello, world")
	hello := decorator(Hello)
	hello("Hello")
	permission := UserPerHandler(GetUserId)
	permission(1)
}
