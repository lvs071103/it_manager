package my_workspaces

import (
	"dominos_cmdb/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strconv"
)

type ProfileController struct {
	beego.Controller
}

func (p *ProfileController) Edit() {
	if p.Ctx.Input.IsPost() {
		id, _ := p.GetInt("uid")
		gid, _ := p.GetInt("gid")
		username := p.GetString("username")
		age, _ := p.GetInt("age")
		email := p.GetString("email")
		phone := p.GetString("phone")
		address := p.GetString("address")
		isActive, _ := p.GetInt("is_active")
		gender, _ := p.GetInt("gender")
		job := p.GetString("job")
		phoneInt64, _ := strconv.ParseInt(phone, 10, 64)
		o := orm.NewOrm()
		_, err := o.QueryTable(new(models.User)).Filter("id", id).Update(orm.Params{
			"groups_id": gid,
			"username":  username,
			"age":       age,
			"email":     email,
			"address":   address,
			"is_active": isActive,
			"gender":    gender,
			"phone":     phoneInt64,
			"job":       job,
		})
		resMsg := map[string]interface{}{}
		if err != nil {
			resMsg["code"] = 1062
			resMsg["msg"] = err.Error()
			p.Data["json"] = resMsg
		} else {
			resMsg["code"] = 0
			resMsg["msg"] = "编辑成功"
			p.Data["json"] = resMsg
		}
		p.ServeJSON()
	} else if p.Ctx.Input.IsGet() {
		username := p.GetSession("username")
		user := models.User{}
		var groups []models.Group
		o := orm.NewOrm()
		_, _ = o.QueryTable(new(models.Group)).All(&groups)
		err := o.QueryTable(new(models.User)).Filter("UserName", username).RelatedSel().One(&user)
		if err == nil {
			p.Data["user"] = user
			p.Data["successURI"] = p.Ctx.Request.Referer()
		}
		p.Data["groups"] = groups
		p.TplName = "auth/my_profile.html"
	}
}
