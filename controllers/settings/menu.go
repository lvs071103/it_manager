package settings

import (
	"dominos_cmdb/models"
	"dominos_cmdb/tools"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type MenuController struct {
	beego.Controller
}

func (m *MenuController) List() {
	uName := m.GetSession("username").(string)
	if !tools.HasPermission(uName, "view_menu") {
		m.Redirect(beego.URLFor("ErrorController.Error403"), 302)
	}
	var menuList []models.Menu
	o := orm.NewOrm()
	qs := o.QueryTable(new(models.Menu)).Filter("is_delete", 0)
	// 定义每页显示的记录数
	limitNum := 10
	// 获取当前页
	currentPage, err := m.GetInt("page")
	if err != nil {
		// 说明没有获取到当前页
		currentPage = 1
	}
	offsetNum := limitNum * (currentPage - 1)
	var count int64 = 0
	searchKey := m.GetString("ss")
	if searchKey == "" {
		// 未输入搜索字段
		count, _ = qs.Count()
		_, _ = qs.Limit(limitNum).Offset(offsetNum).All(&menuList)
	} else {
		count, _ = qs.Filter("name__contains", searchKey).Count()
		_, _ = qs.Filter("name__contains", searchKey).Limit(limitNum).Offset(offsetNum).All(&menuList)
	}
	// 计算总页数
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
	m.Data["previousPage"] = previousPage
	m.Data["nextPage"] = nextPage
	m.Data["currentPage"] = currentPage
	m.Data["count"] = count
	m.Data["hasPrev"] = hasPrev
	m.Data["hasNext"] = hasNext
	m.Data["leftHasMore"] = leftHasMore
	m.Data["rightHasMore"] = rightHasMore
	m.Data["leftPages"] = leftPages
	m.Data["rightPages"] = rightPages
	m.Data["firstPage"] = firstPageNum
	m.Data["endPageNum"] = pageNum
	m.Data["ss"] = searchKey
	m.Data["menuList"] = menuList
	m.TplName = "settings/menu_list.html"
}

func (m *MenuController) Add() {
	if m.Ctx.Input.IsGet() {
		// get request
		var menuList []models.Menu
		o := orm.NewOrm()
		qs := o.QueryTable(new(models.Menu))
		_, queryAllErr := qs.Filter("is_delete", 0).All(&menuList)
		if queryAllErr != nil {
			fmt.Println(queryAllErr)
		}
		var iconList []models.Icon
		qs = o.QueryTable(new(models.Icon))
		_, _ = qs.All(&iconList)

		m.Data["menuList"] = menuList
		m.Data["iconList"] = iconList
		m.Data["action"] = "Add"
		m.TplName = "settings/menu_form.html"
	} else if m.Ctx.Input.IsPost() {
		// post request
		pid, _ := m.GetInt("pid")
		name := m.GetString("name")
		urlfor := m.GetString("urlfor")
		weight, _ := m.GetInt("weight")
		desc := m.GetString("desc")
		isActive, _ := m.GetInt("is_active")
		iconsId, _ := m.GetInt("icons_id")
		o := orm.NewOrm()
		icons := models.Icon{Id: iconsId}
		menuData := models.Menu{Pid: pid, Name: name, UrlFor: urlfor, Weight: weight, Icons: &icons, Desc: desc,
			IsActive: isActive, CreateTime: time.Now()}
		_, insertErr := o.Insert(&menuData)
		var resMsg = map[string]interface{}{}
		if insertErr != nil {
			resMsg["msg"] = insertErr
			resMsg["code"] = 1001
		} else {
			resMsg["msg"] = "插入成功"
			resMsg["code"] = 0
		}
		m.Data["json"] = resMsg
		m.ServeJSON()
	}
}

func (m *MenuController) IsActive() {
	uName := m.GetSession("username").(string)
	if !tools.HasPermission(uName, "change_menu") {
		m.Redirect(beego.URLFor("ErrorController.Error403"), 302)
	}
	isActive, _ := m.GetInt("is_active")
	id, _ := m.GetInt("id")
	m.Data["json"] = tools.CommonUpdateIsActive(new(models.Menu), id, isActive)
	m.ServeJSON()
}

func (m *MenuController) Delete() {
	uName := m.GetSession("username").(string)
	if !tools.HasPermission(uName, "change_user") {
		m.Redirect(beego.URLFor("ErrorController.Error403"), 302)
	}
	id, _ := m.GetInt("id")
	m.Data["json"] = tools.CommonLogicDelete(id, new(models.Menu))
	m.ServeJSON()
}

func (m *MenuController) Edit() {
	uName := m.GetSession("username").(string)
	if !tools.HasPermission(uName, "change_menu") {
		m.Redirect(beego.URLFor("ErrorController.Error403"), 302)
	}
	if m.Ctx.Input.IsPost() {
		id, _ := m.GetInt("id")
		pid, _ := m.GetInt("pid")
		name := m.GetString("name")
		urlfor := m.GetString("urlfor")
		weight, _ := m.GetInt("weight")
		desc := m.GetString("desc")
		isActive, _ := m.GetInt("is_active")
		icons, _ := m.GetInt("icons_id")
		o := orm.NewOrm()
		_, err := o.QueryTable(new(models.Menu)).Filter("id", id).Update(orm.Params{
			"pid":       pid,
			"name":      name,
			"url_for":   urlfor,
			"weight":    weight,
			"desc":      desc,
			"is_active": isActive,
			"icons":     icons,
		})

		retMap := map[string]interface{}{}
		if err != nil {
			retMap["code"] = 1062
			retMap["msg"] = err.Error()
			m.Data["json"] = retMap
		} else {
			retMap["code"] = 0
			retMap["msg"] = "编辑成功"
			m.Data["json"] = retMap
		}
		m.ServeJSON()
	} else if m.Ctx.Input.IsGet() {
		id, _ := m.GetInt("id")
		// 单个菜单对象
		menu := models.Menu{}
		// 所有菜单slice
		var menuList []models.Menu

		o := orm.NewOrm()
		qs := o.QueryTable(new(models.Menu))

		_ = qs.Filter("id", id).RelatedSel().One(&menu)
		_, _ = qs.All(&menuList)
		// 定义Icon的slice
		var iconList []models.Icon
		// 查询所有icons
		iconQS := o.QueryTable(new(models.Icon))
		_, _ = iconQS.All(&iconList)

		m.Data["menu"] = menu
		m.Data["iconList"] = iconList
		m.Data["menuList"] = menuList
		m.Data["action"] = "Edit"
		m.Data["successURI"] = m.Ctx.Request.Referer()
		m.TplName = "settings/menu_form.html"
	}
}

func (m *MenuController) BatchDelete() {
	uName := m.GetSession("username").(string)
	if !tools.HasPermission(uName, "change_user") {
		m.Redirect(beego.URLFor("ErrorController.Error403"), 302)
	}
	ids := m.GetString("ids")
	m.Data["json"] = tools.CommonBatchLogicDelete(ids, new(models.Menu))
	m.ServeJSON()
}
