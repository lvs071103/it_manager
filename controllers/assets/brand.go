package assets

import (
	"dominos_cmdb/models"
	"dominos_cmdb/tools"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type BrandController struct {
	beego.Controller
}

func (b *BrandController) List() {
	uName := b.GetSession("username").(string)
	if !tools.HasPermission(uName, "view_brand") {
		b.Redirect(beego.URLFor("ErrorController.Error403"), 302)
	}
	o := orm.NewOrm()
	qs := o.QueryTable(new(models.Brand)).Filter("is_delete", 0)
	var brands []models.Brand
	var count int64
	limitNum := 10
	currentPage, err := b.GetInt("page")
	if err != nil {
		// 说明没有获取到当前页
		currentPage = 1
	}
	offsetNum := limitNum * (currentPage - 1)
	searchKey := b.GetString("ss")
	if searchKey == "" {
		// 搜索有值
		count, _ = qs.Count()
		_, _ = qs.Limit(limitNum).Offset(offsetNum).All(&brands)
	} else {
		count, _ = qs.Filter("Name__contains", searchKey).Count()
		_, _ = qs.Filter("Name__contains", searchKey).Limit(limitNum).Offset(offsetNum).All(&brands)
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

	b.Data["count"] = count
	b.Data["brands"] = brands
	b.Data["previousPage"] = previousPage
	b.Data["nextPage"] = nextPage
	b.Data["firstPageNum"] = firstPageNum
	b.Data["leftPages"] = leftPages
	b.Data["rightPages"] = rightPages
	b.Data["leftHasMore"] = leftHasMore
	b.Data["rightHasMore"] = rightHasMore
	b.Data["hasPrev"] = hasPrev
	b.Data["hasNext"] = hasNext
	b.Data["currentPage"] = currentPage
	b.Data["endPageNum"] = currentPage
	b.Data["ss"] = searchKey
	b.TplName = "assets/brand_list.html"
}

func (b *BrandController) Json() {
	o := orm.NewOrm()
	var bands []models.Brand
	_, err := o.QueryTable(new(models.Brand)).All(&bands)
	if err == nil {
		b.Data["json"] = bands
	} else {
		b.Data["json"] = map[string]interface{}{"msg": err.Error()}
	}
	b.ServeJSON()
}

func (b *BrandController) Add() {
	//uid := b.GetSession("id").(int)
	//if !tools.HasPermission(uid, "add_brand") {
	//	b.Redirect(beego.URLFor("ErrorController.Error403"), 302)
	//}
	if b.Ctx.Input.IsGet() {
		b.Data["action"] = "Add"
		b.TplName = "assets/brand_form.html"
	} else if b.Ctx.Input.IsPost() {
		name := b.GetString("brand_name")
		isActive, _ := b.GetInt("is_active")
		desc := b.GetString("desc")
		o := orm.NewOrm()
		brands := models.Brand{Name: name, Desc: desc, IsActive: isActive}
		_, err := o.Insert(&brands)
		resMsg := map[string]interface{}{}
		if err == nil {
			resMsg["code"] = 0
			resMsg["msg"] = "添加成功"
		} else {
			resMsg["code"] = 1001
			resMsg["msg"] = err.Error()
		}
		b.Data["json"] = resMsg
		b.ServeJSON()
	}
}

func (b *BrandController) Edit() {
	uName := b.GetSession("username").(string)
	if !tools.HasPermission(uName, "change_brand") {
		b.Redirect(beego.URLFor("ErrorController.Error403"), 302)
	}
	if b.Ctx.Input.IsGet() {
		id, _ := b.GetInt("id")
		o := orm.NewOrm()
		brand := models.Brand{}
		err := o.QueryTable(new(models.Brand)).Filter("id", id).One(&brand)
		if err == nil {
			b.Data["brand"] = brand
		} else {
			fmt.Println(err.Error())
		}
		b.Data["action"] = "Edit"
		b.TplName = "assets/brand_form.html"
	} else if b.Ctx.Input.IsPost() {
		id, _ := b.GetInt("id")
		brandName := b.GetString("brand_name")
		desc := b.GetString("desc")
		o := orm.NewOrm()
		_, err := o.QueryTable(new(models.Brand)).Filter("id", id).Update(orm.Params{
			"Name": brandName,
			"Desc": desc,
		})
		resMsg := map[string]interface{}{}
		if err == nil {
			resMsg["code"] = 0
			resMsg["msg"] = "更新成功"
		} else {
			resMsg["code"] = 1001
			resMsg["msg"] = err.Error()
		}
		b.Data["json"] = resMsg
		b.ServeJSON()
	}
}

func (b *BrandController) Delete() {
	uName := b.GetSession("username").(string)
	if !tools.HasPermission(uName, "change_brand") {
		b.Redirect(beego.URLFor("ErrorController.Error403"), 302)
	}
	id, _ := b.GetInt("id")
	b.Data["json"] = tools.CommonLogicDelete(id, new(models.Brand))
	b.ServeJSON()
}

func (b *BrandController) BatchDelete() {
	uName := b.GetSession("username").(string)
	if !tools.HasPermission(uName, "change_brand") {
		b.Redirect(beego.URLFor("ErrorController.Error403"), 302)
	}
	ids := b.GetString("ids")
	b.Data["json"] = tools.CommonBatchLogicDelete(ids, new(models.Brand))
	b.ServeJSON()
}

func (b *BrandController) IsActive() {
	uName := b.GetSession("username").(string)
	if !tools.HasPermission(uName, "change_brand") {
		b.Redirect(beego.URLFor("ErrorController.Error403"), 302)
	}
	isActive, _ := b.GetInt("is_active")
	id, _ := b.GetInt("id")
	b.Data["json"] = tools.CommonUpdateIsActive(new(models.Brand), id, isActive)
	b.ServeJSON()
}
