package myError

import "github.com/astaxie/beego"

type ErrorController struct {
	beego.Controller
}

func (c *ErrorController) Error403() {
	c.Data["content"] = "403"
	c.Data["resCode"] = 403
	c.TplName = "error/template.html"
}
