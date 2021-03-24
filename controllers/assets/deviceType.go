package assets

import (
	"dominos_cmdb/models"
	"dominos_cmdb/tools"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type DTypeController struct {
	beego.Controller
}

func (d *DTypeController) List() {
	uName := d.GetSession("username").(string)
	if !tools.HasPermission(uName, "view_device_type") {
		d.Redirect(beego.URLFor("ErrorController.Error403"), 302)
	}
	o := orm.NewOrm()
	qs := o.QueryTable(new(models.DeviceType)).Filter("is_delete", 0)
	var deviceTypes []models.DeviceType
	var count int64
	limitNum := 10
	currentPage, err := d.GetInt("page")
	if err != nil {
		// 说明没有获取到当前页
		currentPage = 1
	}
	offsetNum := limitNum * (currentPage - 1)
	searchKey := d.GetString("ss")
	if searchKey == "" {
		// 搜索无值
		count, _ = qs.Count()
		_, _ = qs.Limit(limitNum).Offset(offsetNum).All(&deviceTypes)
	} else {
		count, _ = qs.Filter("DeviceType__contains", searchKey).Count()
		_, _ = qs.Filter("DeviceType__contains", searchKey).Limit(limitNum).Offset(offsetNum).All(&deviceTypes)
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
	d.Data["count"] = count
	d.Data["deviceTypes"] = deviceTypes
	d.Data["previousPage"] = previousPage
	d.Data["nextPage"] = nextPage
	d.Data["firstPageNum"] = firstPageNum
	d.Data["leftPages"] = leftPages
	d.Data["rightPages"] = rightPages
	d.Data["leftHasMore"] = leftHasMore
	d.Data["rightHasMore"] = rightHasMore
	d.Data["hasPrev"] = hasPrev
	d.Data["hasNext"] = hasNext
	d.Data["currentPage"] = currentPage
	d.Data["endPageNum"] = pageNum
	d.Data["ss"] = searchKey
	d.TplName = "assets/device_type_list.html"
}

func (d *DTypeController) Json() {
	o := orm.NewOrm()
	var deviceTypes []models.DeviceType
	_, err := o.QueryTable(new(models.DeviceType)).All(&deviceTypes)
	if err == nil {
		d.Data["json"] = deviceTypes
	} else {
		d.Data["json"] = map[string]interface{}{"msg": err.Error()}
	}
	d.ServeJSON()
}

func (d *DTypeController) Add() {
	uName := d.GetSession("username").(string)
	if !tools.HasPermission(uName, "add_device_type") {
		d.Redirect(beego.URLFor("ErrorController.Error403"), 302)
	}
	if d.Ctx.Input.IsGet() {
		d.Data["action"] = "Add"
		d.TplName = "assets/device_type_form.html"
	} else if d.Ctx.Input.IsPost() {
		deviceType := d.GetString("deviceType")
		desc := d.GetString("desc")
		isActive, _ := d.GetInt("is_active")
		o := orm.NewOrm()
		deviceTypeObj := models.DeviceType{DeviceType: deviceType, Desc: desc, IsActive: isActive}
		_, err := o.Insert(&deviceTypeObj)
		resMsg := map[string]interface{}{}
		if err == nil {
			resMsg["code"] = 0
			resMsg["msg"] = "添加成功"
		} else {
			resMsg["code"] = 1001
			resMsg["msg"] = err.Error()
		}
		d.Data["json"] = resMsg
		d.ServeJSON()
	}
}

func (d *DTypeController) Edit() {
	uName := d.GetSession("username").(string)
	if !tools.HasPermission(uName, "change_device_type") {
		d.Redirect(beego.URLFor("ErrorController.Error403"), 302)
	}
	if d.Ctx.Input.IsGet() {
		id, _ := d.GetInt("id")
		o := orm.NewOrm()
		dType := models.DeviceType{}
		err := o.QueryTable(new(models.DeviceType)).Filter("id", id).One(&dType)
		if err == nil {
			d.Data["d_type"] = dType
		} else {
			fmt.Println(err.Error())
		}
		d.Data["action"] = "Edit"
		d.TplName = "assets/device_type_form.html"
	} else if d.Ctx.Input.IsPost() {
		id, _ := d.GetInt("id")
		deviceType := d.GetString("deviceType")
		desc := d.GetString("desc")
		isActive, _ := d.GetInt("is_active")
		o := orm.NewOrm()
		_, err := o.QueryTable(new(models.DeviceType)).Filter("id", id).Update(orm.Params{
			"DeviceType": deviceType,
			"Desc":       desc,
			"IsActive":   isActive,
		})
		resMsg := map[string]interface{}{}
		if err == nil {
			resMsg["code"] = 0
			resMsg["msg"] = "更新成功"
		} else {
			resMsg["code"] = 1001
			resMsg["msg"] = err.Error()
		}
		d.Data["json"] = resMsg
		d.ServeJSON()
	}
}

func (d *DTypeController) Delete() {
	uName := d.GetSession("username").(string)
	if !tools.HasPermission(uName, "change_device_type") {
		d.Redirect(beego.URLFor("ErrorController.Error403"), 302)
	}
	id, _ := d.GetInt("id")
	d.Data["json"] = tools.CommonLogicDelete(id, new(models.DeviceType))
	d.ServeJSON()
}

func (d *DTypeController) BatchDelete() {
	uName := d.GetSession("username").(string)
	if !tools.HasPermission(uName, "change_device_type") {
		d.Redirect(beego.URLFor("ErrorController.Error403"), 302)
	}
	ids := d.GetString("ids")
	d.Data["json"] = tools.CommonBatchLogicDelete(ids, new(models.DeviceType))
	d.ServeJSON()
}

func (d *DTypeController) IsActive() {
	uName := d.GetSession("username").(string)
	if !tools.HasPermission(uName, "change_device_type") {
		d.Redirect(beego.URLFor("ErrorController.Error403"), 302)
	}
	isActive, _ := d.GetInt("is_active")
	id, _ := d.GetInt("id")
	d.Data["json"] = tools.CommonUpdateIsActive(new(models.DeviceType), id, isActive)
	d.ServeJSON()
}
