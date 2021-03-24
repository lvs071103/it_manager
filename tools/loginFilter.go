package tools

import "github.com/astaxie/beego/context"

func LoginFilter(ctx *context.Context) {
	// 获取session
	id := ctx.Input.Session("username")
	if id == nil {
		ctx.Redirect(302, "/login")
	}
}
