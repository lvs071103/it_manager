package auth

import (
	"dominos_cmdb/models"
	"dominos_cmdb/tools"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type GroupController struct {
	beego.Controller
}

func (g *GroupController) List() {
	// 根据当前uid检查是否拥有view_group的权限, 如果没有则重定向到403页面
	uName := g.GetSession("username").(string)
	fmt.Println(uName)
	if !tools.HasPermission(uName, "view_group") {
		g.Redirect(beego.URLFor("ErrorController.Error403"), 302)
	}

	o := orm.NewOrm()
	qs := o.QueryTable(new(models.Group)).Filter("is_delete", 0)
	var groups []models.Group
	var count int64
	limitNum := 10
	currentPage, err := g.GetInt("page")
	if err != nil { // 说明没有获取到当前页
		currentPage = 1
	}
	offsetNum := limitNum * (currentPage - 1)
	searchKey := g.GetString("ss")
	if searchKey == "" { // 搜索有值
		count, _ = qs.Count()
		_, _ = qs.Limit(limitNum).Offset(offsetNum).All(&groups)
	} else {
		count, _ = qs.Filter("group_name__contains", searchKey).Count()
		_, _ = qs.Filter("group_name__contains", searchKey).Limit(limitNum).Offset(offsetNum).All(&groups)
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

	g.Data["count"] = count
	g.Data["groups"] = groups
	g.Data["previousPage"] = previousPage
	g.Data["nextPage"] = nextPage
	g.Data["firstPageNum"] = firstPageNum
	g.Data["leftPages"] = leftPages
	g.Data["rightPages"] = rightPages
	g.Data["leftHasMore"] = leftHasMore
	g.Data["rightHasMore"] = rightHasMore
	g.Data["hasPrev"] = hasPrev
	g.Data["hasNext"] = hasNext
	g.Data["currentPage"] = currentPage
	g.Data["ss"] = searchKey
	g.Data["endPageNum"] = pageNum
	g.TplName = "auth/group_list.html"
}

func (g *GroupController) Add() {
	if g.Ctx.Input.IsGet() {
		uName := g.GetSession("username").(string)
		if !tools.HasPermission(uName, "add_group") {
			g.Redirect(beego.URLFor("ErrorController.Error403"), 302)
		}
		g.Data["action"] = "Add"
		g.TplName = "auth/group_form.html"
	} else if g.Ctx.Input.IsPost() {
		uName := g.GetSession("username").(string)
		if !tools.HasPermission(uName, "add_group") {
			g.Redirect(beego.URLFor("ErrorController.Error403"), 302)
		}
		groupName := g.GetString("group_name")
		desc := g.GetString("desc")
		isActive, _ := g.GetInt("is_active")
		o := orm.NewOrm()
		groups := models.Group{GroupName: groupName, Desc: desc, IsActive: isActive}
		_, err := o.Insert(&groups)
		resMsg := map[string]interface{}{}
		if err == nil {
			resMsg["code"] = 0
			resMsg["msg"] = "添加成功"
		} else {
			resMsg["code"] = 1001
			resMsg["msg"] = err.Error()
		}
		g.Data["json"] = resMsg
		g.ServeJSON()
	}
}

func (g *GroupController) Edit() {
	uName := g.GetSession("username").(string)
	if !tools.HasPermission(uName, "change_group") {
		g.Redirect(beego.URLFor("ErrorController.Error403"), 302)
	}
	if g.Ctx.Input.IsGet() {
		gid := g.GetString("gid")
		o := orm.NewOrm()
		group := models.Group{}
		err := o.QueryTable(new(models.Group)).Filter("id", gid).One(&group)

		if err == nil {
			g.Data["group"] = group
			g.Data["action"] = "Edit"
		} else {
			fmt.Println(err.Error())
		}
		g.TplName = "auth/group_form.html"
	} else if g.Ctx.Input.IsPost() {
		uName := g.GetSession("username").(string)
		if !tools.HasPermission(uName, "change_group") {
			g.Redirect(beego.URLFor("ErrorController.Error403"), 302)
		}
		gid, _ := g.GetInt("gid")
		groupName := g.GetString("group_name")
		isActive, _ := g.GetInt("is_active")
		desc := g.GetString("desc")
		o := orm.NewOrm()
		_, err := o.QueryTable(new(models.Group)).Filter("id", gid).Update(orm.Params{
			"group_name": groupName,
			"desc":       desc,
			"is_active":  isActive,
		})
		resMsg := map[string]interface{}{}
		if err == nil {
			resMsg["code"] = 0
			resMsg["msg"] = "更新成功"
		} else {
			resMsg["code"] = 1001
			resMsg["msg"] = err.Error()
		}
		g.Data["json"] = resMsg
		g.ServeJSON()
	}
}

func (g *GroupController) IsActive() {
	uName := g.GetSession("username").(string)
	if !tools.HasPermission(uName, "change_group") {
		g.Redirect(beego.URLFor("ErrorController.Error403"), 302)
	}
	isActive, _ := g.GetInt("is_active")
	id, _ := g.GetInt("id")
	g.Data["json"] = tools.CommonUpdateIsActive(new(models.Group), id, isActive)
	g.ServeJSON()
}

func (g *GroupController) Delete() {
	uName := g.GetSession("username").(string)
	if !tools.HasPermission(uName, "change_group") {
		g.Redirect(beego.URLFor("ErrorController.Error403"), 302)
	}
	id, _ := g.GetInt("id")
	g.Data["json"] = tools.CommonLogicDelete(id, new(models.Group))
	g.ServeJSON()
}

func (g *GroupController) BatchDelete() {
	uName := g.GetSession("username").(string)
	if !tools.HasPermission(uName, "change_group") {
		g.Redirect(beego.URLFor("ErrorController.Error403"), 302)
	}
	ids := g.GetString("ids")
	g.Data["json"] = tools.CommonBatchLogicDelete(ids, new(models.Group))
	g.ServeJSON()
}
