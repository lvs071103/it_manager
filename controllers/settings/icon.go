package settings

import (
	"dominos_cmdb/models"
	"dominos_cmdb/tools"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strings"
	"time"
)

type IconController struct {
	beego.Controller
}

func (i *IconController) List() {
	uName := i.GetSession("username").(string)
	if !tools.HasPermission(uName, "view_icon") {
		i.Redirect(beego.URLFor("ErrorController.Error403"), 302)
	}
	o := orm.NewOrm()
	qs := o.QueryTable(new(models.Icon))
	var icons []models.Icon
	limitNum := 10
	currentPage, err := i.GetInt("page")
	if err != nil { // 说明没有获取到当前页
		currentPage = 1
	}
	offsetNum := limitNum * (currentPage - 1)
	var count int64 = 0
	searchKey := i.GetString("ss")
	if searchKey == "" {
		// 未输入搜索字段
		count, _ = qs.Filter("is_delete", 0).Count()
		_, _ = qs.Filter("is_delete", 0).Limit(limitNum).Offset(offsetNum).All(&icons)
	} else {
		_, _ = qs.Filter("is_delete", 0).Filter("ClassName__contains", searchKey).Limit(limitNum).
			Offset(offsetNum).All(&icons)
	}
	pageNum := tools.GetPageNum(count, limitNum)
	previousPage := 1
	if currentPage == 1 {
		previousPage = currentPage
	} else if currentPage > 1 {
		previousPage = currentPage - 1
	}
	nextPage := 1
	if currentPage < pageNum {
		nextPage = currentPage + 1
	} else if currentPage >= pageNum {
		nextPage = currentPage
	}
	firstPageNum := 1
	leftPages, rightPages, leftHasMore, rightHasMore := tools.GetPaginationData(pageNum, currentPage, 4)
	hasPrev, hasNext, _, _ := tools.HasNext(currentPage, pageNum)
	i.Data["icons"] = icons
	i.Data["count"] = count
	i.Data["previousPage"] = previousPage
	i.Data["firstPageNum"] = firstPageNum
	i.Data["leftPages"] = leftPages
	i.Data["rightPages"] = rightPages
	i.Data["leftHasMore"] = leftHasMore
	i.Data["rightHasMore"] = rightHasMore
	i.Data["hasPrev"] = hasPrev
	i.Data["hasNext"] = hasNext
	i.Data["nextPage"] = nextPage
	i.Data["ss"] = searchKey
	i.Data["currentPage"] = currentPage
	i.Data["endPageNum"] = pageNum
	i.TplName = "settings/icon_list.html"
}

func (i *IconController) Add() {
	uName := i.GetSession("username").(string)
	if !tools.HasPermission(uName, "add_icon") {
		i.Redirect(beego.URLFor("ErrorController.Error403"), 302)
	}
	if i.Ctx.Input.IsGet() {
		i.Data["action"] = "Add"
		i.TplName = "settings/icon_form.html"
	} else if i.Ctx.Input.IsPost() {
		name := i.GetString("name")
		className := i.GetString("class_name")
		desc := i.GetString("desc")
		isActive, _ := i.GetInt("is_active")
		o := orm.NewOrm()
		icons := models.Icon{Name: name, ClassName: className, Desc: desc, IsActive: isActive, CreateTime: time.Now()}
		_, err := o.Insert(&icons)
		resMsg := map[string]interface{}{}
		if err == nil {
			resMsg["code"] = 0
			resMsg["msg"] = "添加成功"
		} else {
			resMsg["code"] = 1001
			resMsg["msg"] = err.Error()
		}
		i.Data["json"] = resMsg

		i.ServeJSON()
	}
}

func (i *IconController) Delete() {
	id, _ := i.GetInt("id")
	o := orm.NewOrm()
	_, err := o.QueryTable("sys_icon").Filter("id", id).Update(orm.Params{"is_delete": 1})
	returnMp := map[string]interface{}{}
	if err == nil {
		returnMp["code"] = 0
		returnMp["msg"] = "删除成功"
	} else {
		returnMp["code"] = 1001
		returnMp["msg"] = err.Error()
	}
	i.Data["json"] = returnMp
	i.ServeJSON()
}

func (i *IconController) IsActive() {
	id, _ := i.GetInt("id")
	isActive, _ := i.GetInt("is_active")
	o := orm.NewOrm()
	qs := o.QueryTable("sys_icon").Filter("id", id)
	returnMp := map[string]interface{}{}
	if isActive == 0 {
		_, err := qs.Update(orm.Params{"is_active": 1})
		if err == nil {
			returnMp["msg"] = "启用成功"
		} else {
			returnMp["msg"] = err.Error()
		}
	} else if isActive == 1 {
		_, err := qs.Update(orm.Params{"is_active": 0})
		if err == nil {
			returnMp["msg"] = "停用成功"
		} else {
			returnMp["msg"] = err.Error()
		}
	}
	i.Data["json"] = returnMp
	i.ServeJSON()
}

func (i *IconController) Edit() {
	uName := i.GetSession("username").(string)
	if !tools.HasPermission(uName, "change_icon") {
		i.Redirect(beego.URLFor("ErrorController.Error403"), 302)
	}
	if i.Ctx.Input.IsGet() {
		id, _ := i.GetInt("id")
		o := orm.NewOrm()
		icon := models.Icon{}
		qs := o.QueryTable(new(models.Icon))
		_ = qs.Filter("id", id).One(&icon)
		i.Data["icon"] = icon
		i.Data["action"] = "Edit"
		i.TplName = "settings/icon_form.html"
	} else if i.Ctx.Input.IsPost() {
		id, _ := i.GetInt("iconId")
		name := i.GetString("name")
		className := i.GetString("class_name")
		desc := i.GetString("desc")
		isActive, _ := i.GetInt("is_active")
		fmt.Println(id, className)
		o := orm.NewOrm()
		qs := o.QueryTable(new(models.Icon))
		_, err := qs.Filter("id", id).Update(orm.Params{
			"Name":      name,
			"ClassName": className,
			"Desc":      desc,
			"IsActive":  isActive,
		})

		returnMp := map[string]interface{}{}

		if err == nil {
			returnMp["code"] = 0
			returnMp["msg"] = "更新成功"
		} else {
			returnMp["code"] = 1001
			returnMp["msg"] = err.Error()
		}

		i.Data["json"] = returnMp

		i.ServeJSON()
	}
}

func (i *IconController) BatchDelete() {
	uName := i.GetSession("username").(string)
	if !tools.HasPermission(uName, "change_icon") {
		i.Redirect(beego.URLFor("ErrorController.Error403"), 302)
	}
	ids := i.GetString("ids")
	idsArr := strings.Split(ids, "&")
	o := orm.NewOrm()
	qs := o.QueryTable("sys_icon")
	for _, v := range idsArr {
		_, err := qs.Filter("id", tools.StrToInt(v)).Update(orm.Params{"is_delete": 1})
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	i.Data["json"] = map[string]interface{}{
		"code": 0,
		"msg":  "批量删除成功",
	}

	i.ServeJSON()
}
