package tickets

import (
	"dominos_cmdb/models"
	"dominos_cmdb/tools"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type TCController struct {
	beego.Controller
}

func (t *TCController) List() {
	var TcList []models.TicketCategory
	o := orm.NewOrm()
	qs := o.QueryTable(new(models.TicketCategory)).Filter("is_delete", 0)
	// 定义每页显示的记录数
	limitNum := 10
	// 获取当前页
	currentPage, err := t.GetInt("page")
	if err != nil {
		// 说明没有获取到当前页
		currentPage = 1
	}
	offsetNum := limitNum * (currentPage - 1)
	var count int64 = 0
	searchKey := t.GetString("ss")
	if searchKey == "" {
		// 未输入搜索字段
		count, _ = qs.Count()
		_, _ = qs.Limit(limitNum).Offset(offsetNum).All(&TcList)
	} else {
		count, _ = qs.Filter("name__contains", searchKey).Count()
		_, _ = qs.Filter("name__contains", searchKey).Limit(limitNum).Offset(offsetNum).All(&TcList)
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
	t.Data["previousPage"] = previousPage
	t.Data["nextPage"] = nextPage
	t.Data["currentPage"] = currentPage
	t.Data["count"] = count
	t.Data["hasPrev"] = hasPrev
	t.Data["hasNext"] = hasNext
	t.Data["leftHasMore"] = leftHasMore
	t.Data["rightHasMore"] = rightHasMore
	t.Data["leftPages"] = leftPages
	t.Data["rightPages"] = rightPages
	t.Data["firstPage"] = firstPageNum
	t.Data["endPageNum"] = pageNum
	t.Data["ss"] = searchKey
	t.Data["tc_list"] = TcList
	t.TplName = "tickets/category_list.html"
}

func (t *TCController) Add() {
	if t.Ctx.Input.IsGet() {
		t.Data["action"] = "Add"
		t.TplName = "tickets/category_form.html"
	} else if t.Ctx.Input.IsPost() {
		name := t.GetString("name")
		desc := t.GetString("desc")
		isActive, _ := t.GetInt("is_active")
		fmt.Println(name, desc, isActive)
		o := orm.NewOrm()
		category := models.TicketCategory{Name: name, Desc: desc, IsActive: isActive}
		_, err := o.Insert(&category)
		resMsg := map[string]interface{}{}
		if err == nil {
			resMsg["code"] = 0
			resMsg["msg"] = "添加成功"
		} else {
			resMsg["code"] = 1001
			resMsg["msg"] = err.Error()
		}
		t.Data["json"] = resMsg

		t.ServeJSON()
	}
}

func (t *TCController) Edit() {
	if t.Ctx.Input.IsGet() {
		id, _ := t.GetInt("id")
		o := orm.NewOrm()
		category := models.TicketCategory{}
		qs := o.QueryTable(new(models.TicketCategory))
		_ = qs.Filter("id", id).One(&category)
		t.Data["category"] = category
		t.Data["action"] = "Edit"
		t.TplName = "tickets/category_form.html"
	} else if t.Ctx.Input.IsPost() {
		id, _ := t.GetInt("id")
		name := t.GetString("name")
		desc := t.GetString("desc")
		isActive, _ := t.GetInt("is_active")
		o := orm.NewOrm()
		qs := o.QueryTable(new(models.TicketCategory))
		_, err := qs.Filter("id", id).Update(orm.Params{
			"name":      name,
			"desc":      desc,
			"is_active": isActive,
		})
		returnMp := map[string]interface{}{}
		if err == nil {
			returnMp["code"] = 0
			returnMp["msg"] = "更新成功"
		} else {
			returnMp["code"] = 1001
			returnMp["msg"] = err.Error()
		}
		t.Data["json"] = returnMp
		t.ServeJSON()
	}
}

func (t *TCController) Del() {
	id, _ := t.GetInt("id")
	t.Data["json"] = tools.CommonLogicDelete(id, new(models.TicketCategory))
	t.ServeJSON()
}

func (t *TCController) BatchDelete() {
	ids := t.GetString("ids")
	t.Data["json"] = tools.CommonBatchLogicDelete(ids, new(models.TicketCategory))
	t.ServeJSON()
}

func (t *TCController) IsActive() {
	isActive, _ := t.GetInt("is_active")
	id, _ := t.GetInt("id")
	t.Data["json"] = tools.CommonUpdateIsActive(new(models.TicketCategory), id, isActive)
	t.ServeJSON()
}
