package authentication

import (
	"dominos_cmdb/models"
	"dominos_cmdb/tools"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type LoginController struct {
	beego.Controller
}

func (l *LoginController) Login() {
	if l.Ctx.Input.IsGet() {
		//id, base64, err := tools.GetCaptcha()
		//if err != nil {
		//	return
		//}
		//l.Data["captcha"] = tools.Captcha{Id: id, BS64: base64}
		l.TplName = "login/login.html"
	} else {
		username := l.GetString("username")
		password := l.GetString("password")
		//captcha := l.GetString("captcha")
		//captchaId := l.GetString("captcha_id")
		passMD5 := tools.GetMd5Sum(password)
		userInfo := models.User{}
		o := orm.NewOrm()
		isExist := o.QueryTable(new(models.User)).Filter("username", username).Filter("password", passMD5).Exist()
		//isEq := tools.VerifyCaptcha(captchaId, captcha)
		_ = o.QueryTable(new(models.User)).Filter("username", username).Filter("password", passMD5).One(&userInfo)
		resMsg := map[string]interface{}{}
		if !isExist {
			resMsg["code"] = 123
			resMsg["msg"] = "用户名或密码错误"
			l.Data["json"] = resMsg
			//} else if !isEq {
			//	resMsg["code"] = 124
			//	resMsg["msg"] = "验证码错误"
			//	l.Data["json"] = resMsg
		} else if userInfo.IsActive == 0 {
			resMsg["code"] = 125
			resMsg["msg"] = "该用户已停用，请联系管理员"
			l.Data["json"] = resMsg
		} else {
			l.SetSession("username", username)
			l.SetSession("id", userInfo.Id)
			l.SetSession("loginType", "sys")
			UpdateLastLogin(username)
			resMsg["code"] = 0
			resMsg["msg"] = "登陆成功"
			l.Data["json"] = resMsg
			// 更新Last Login时间
		}
		l.ServeJSON()
	}
}

func (l *LoginController) RefreshCaptcha() {
	message := map[string]interface{}{}
	id, base64, err := tools.GetCaptcha()
	if err != nil {
		// 说明有报错
		message["msg"] = "生成验证码失败"
		message["Code"] = 123
		l.Data["json"] = message
	} else {
		l.Data["json"] = tools.Captcha{Id: id, BS64: base64, Code: 0}
	}
	l.ServeJSON()

}

func (l *LoginController) Logout() {
	l.DelSession("username")
	l.Redirect(beego.URLFor("LoginController.Login"), 302)
}
