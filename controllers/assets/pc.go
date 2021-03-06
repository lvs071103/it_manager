package assets

import (
	"dominos_cmdb/models"
	"dominos_cmdb/tools"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/xuri/excelize/v2"
	_ "github.com/xuri/excelize/v2"
	"strconv"
	_ "strconv"
	"strings"
)

type PCController struct {
	beego.Controller
}

func (p *PCController) List() {
	uName := p.GetSession("username").(string)
	if !tools.HasPermission(uName, "view_computer") {
		p.Redirect(beego.URLFor("ErrorController.Error403"), 302)
	}
	o := orm.NewOrm()
	qs := o.QueryTable(new(models.Computer)).Filter("is_delete", 0)
	var computers []models.Computer
	var count int64
	limitNum := 10
	currentPage, err := p.GetInt("page")
	if err != nil { // 说明没有获取到当前页
		currentPage = 1
	}
	offsetNum := limitNum * (currentPage - 1)
	searchKey := p.GetString("ss")

	if searchKey == "" {
		// 搜索无值
		count, _ = qs.Count()
		_, _ = qs.Limit(limitNum).Offset(offsetNum).RelatedSel().All(&computers)
	} else {
		count, _ = qs.RelatedSel().Filter("Groups__GroupName__contains", searchKey).Count()
		_, _ = qs.RelatedSel().Filter("Groups__GroupName__contains", searchKey).Limit(limitNum).
			Offset(offsetNum).All(&computers)
		// fmt.Print(computers)
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

	p.Data["count"] = count
	p.Data["computers"] = computers
	p.Data["previousPage"] = previousPage
	p.Data["nextPage"] = nextPage
	p.Data["firstPageNum"] = firstPageNum
	p.Data["leftPages"] = leftPages
	p.Data["rightPages"] = rightPages
	p.Data["leftHasMore"] = leftHasMore
	p.Data["rightHasMore"] = rightHasMore
	p.Data["hasPrev"] = hasPrev
	p.Data["hasNext"] = hasNext
	p.Data["currentPage"] = currentPage
	p.Data["endPageNum"] = pageNum
	p.Data["ss"] = searchKey
	p.TplName = "assets/computer_list.html"
}

func (p *PCController) Json() {
	o := orm.NewOrm()
	var pc []models.Computer
	_, err := o.QueryTable(new(models.Computer)).All(&pc)
	if err == nil {
		p.Data["json"] = pc
	} else {
		p.Data["json"] = map[string]interface{}{"msg": err.Error()}
	}
	p.ServeJSON()
}

func (p *PCController) Add() {
	if p.Ctx.Input.IsGet() {
		var users []models.User
		var groups []models.Group
		var deviceTypes []models.DeviceType
		var brands []models.Brand
		o := orm.NewOrm()
		_, _ = o.QueryTable(new(models.User)).Filter("is_delete", 0).All(&users)
		_, _ = o.QueryTable(new(models.Group)).Filter("is_delete", 0).All(&groups)
		_, _ = o.QueryTable(new(models.DeviceType)).Filter("is_delete", 0).All(&deviceTypes)
		_, _ = o.QueryTable(new(models.Brand)).Filter("is_delete", 0).All(&brands)
		p.Data["users"] = users
		p.Data["groups"] = groups
		p.Data["deviceTypes"] = deviceTypes
		p.Data["brands"] = brands
		p.Data["action"] = "Add"
		p.TplName = "assets/computer_form.html"
	} else if p.Ctx.Input.IsPost() {
		location := p.GetString("location")
		name := p.GetString("name")
		deviceTypeId, _ := p.GetInt("deviceType")
		brandId, _ := p.GetInt("brand")
		model := p.GetString("model")
		username := p.GetString("username")
		deptId, _ := p.GetInt("dept")
		assetNumber := p.GetString("name")
		purchaseDate := p.GetString("purchase_date")
		sn := p.GetString("sn")
		quickServiceCode := p.GetString("quick_service_code")
		warranty := p.GetString("warranty")
		remark := p.GetString("remark")
		ipAddress := p.GetString("ip_address")
		isActive, _ := p.GetInt("is_active")
		isBreakdown, _ := p.GetInt("is_breakdown")
		inRepository, _ := p.GetInt("in_repository")
		cpu := p.GetString("cpu")
		memory := p.GetString("memory")
		disk := p.GetString("disk")
		systemV := p.GetString("os")
		fmt.Println(cpu, memory, disk)
		deptObj := models.Group{Id: deptId}
		deviceTypes := models.DeviceType{Id: deviceTypeId}
		binds := models.Brand{Id: brandId}
		computer := models.Computer{Location: location, Name: name, DeviceTypes: &deviceTypes, Brands: &binds,
			Model: model, UserName: username, Groups: &deptObj, AssetNo: assetNumber, PurchaseDate: purchaseDate,
			SN: sn, QuickServiceCode: quickServiceCode, Warranty: warranty, Remark: remark, IsActive: isActive,
			IsBreakdown: isBreakdown, InRepository: inRepository, IpAddress: ipAddress, Cpu: cpu, Memory: memory,
			Disk: disk, OS: systemV,
		}
		o := orm.NewOrm()
		_, err := o.Insert(&computer)
		if err == nil {
			p.Data["json"] = map[string]interface{}{"code": 0, "msg": "添加成功"}
		} else {
			p.Data["json"] = map[string]interface{}{"code": 1001, "msg": err.Error()}
		}
		p.ServeJSON()
	}
}

func (p *PCController) BindUser() {
	if p.Ctx.Input.IsGet() {
		id, _ := p.GetInt("id")
		o := orm.NewOrm()
		computer := models.Computer{}
		var userList []models.User
		// var users []models.User
		_ = o.QueryTable(new(models.Computer)).Filter("id", id).RelatedSel().One(&computer)
		_, _ = o.LoadRelated(&computer, "PreUsers")
		// 判断之前没有绑定，如果有绑定PreUsers > 0
		if len(computer.PreUsers) > 0 {
			_, _ = o.QueryTable(new(models.User)).
				Filter("is_delete", 0).Filter("is_active", 1).Exclude("id__in", computer.PreUsers).All(&userList)
		} else {
			_, _ = o.QueryTable(new(models.User)).
				Filter("is_delete", 0).Filter("is_active", 1).All(&userList)
		}
		//_, _ = o.QueryTable(new(models.User)).All(&users)
		p.Data["computer"] = computer
		p.Data["userList"] = userList
		// p.Data["users"] = users
		p.TplName = "assets/pc_user_bind.html"
	} else if p.Ctx.Input.IsPost() {
		id, _ := p.GetInt("id")
		username := p.GetString("username")
		preUserStr := p.GetString("previous_users")
		preUserArr := strings.Split(preUserStr, "&")
		o := orm.NewOrm()
		computer := models.Computer{Id: id}
		// 更新当前用户
		_, _ = o.QueryTable(new(models.Computer)).Filter("id", id).Update(orm.Params{"username": username})
		// 查询已绑定的数据并清理
		m2m := o.QueryM2M(&computer, "PreUsers")
		_, _ = m2m.Clear()
		// 如果先前用户列表大于0, 并且userId不为空串
		if len(preUserArr) > 0 {
			for _, userId := range preUserArr {
				if userId != "" {
					userObj := models.User{Id: tools.StrToInt(userId)}
					m2m := o.QueryM2M(&computer, "PreUsers")
					_, _ = m2m.Add(userObj)
				}
			}
		}
		p.Data["json"] = map[string]interface{}{"code": 0, "msg": "绑定成功"}
		p.ServeJSON()
	}
}

func (p *PCController) Broken() {
	id, _ := p.GetInt("id")
	computer := models.Computer{Id: id}
	o := orm.NewOrm()
	_, err := o.QueryTable(new(models.Computer)).Filter("id", id).Update(orm.Params{
		"is_active":     0,
		"is_breakdown":  1,
		"in_repository": 1,
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	m2m := o.QueryM2M(&computer, "PreUsers")
	_, _ = m2m.Clear()
	p.Data["json"] = map[string]interface{}{"code": 0, "msg": "处理完成"}
	p.ServeJSON()
}

func (p *PCController) Delete() {
	uName := p.GetSession("username").(string)
	if !tools.HasPermission(uName, "change_computer") {
		p.Redirect(beego.URLFor("ErrorController.Error403"), 302)
	}
	id, _ := p.GetInt("id")
	p.Data["json"] = tools.CommonLogicDelete(id, new(models.Computer))
	p.ServeJSON()
}

func (p *PCController) BatchDelete() {
	uName := p.GetSession("username").(string)
	if !tools.HasPermission(uName, "change_computer") {
		p.Redirect(beego.URLFor("ErrorController.Error403"), 302)
	}
	ids := p.GetString("ids")
	p.Data["json"] = tools.CommonBatchLogicDelete(ids, new(models.Computer))
	p.ServeJSON()
}

func (p *PCController) IsActive() {
	uName := p.GetSession("username").(string)
	if !tools.HasPermission(uName, "change_computer") {
		p.Redirect(beego.URLFor("ErrorController.Error403"), 302)
	}
	isActive, _ := p.GetInt("is_active")
	id, _ := p.GetInt("id")
	o := orm.NewOrm()
	qs := o.QueryTable(new(models.Computer)).Filter("id", id)
	resMsg := map[string]interface{}{}
	if isActive == 0 {
		_, err := qs.Update(orm.Params{"is_active": 1, "in_repository": 0})
		if err == nil {
			resMsg["msg"] = "启用成功"
		} else {
			resMsg["msg"] = err.Error()
		}
	} else if isActive == 1 {
		_, err := qs.Update(orm.Params{"is_active": 0, "in_repository": 1})
		if err == nil {
			resMsg["msg"] = "停用成功"
		} else {
			resMsg["msg"] = err.Error()
		}
	}
	p.Data["json"] = resMsg
	p.ServeJSON()
}

func (p *PCController) Edit() {
	if p.Ctx.Input.IsGet() {
		id, _ := p.GetInt("id")
		var users []models.User
		var groups []models.Group
		var deviceTypes []models.DeviceType
		var brands []models.Brand
		var computer models.Computer
		o := orm.NewOrm()
		_, _ = o.QueryTable(new(models.User)).Filter("is_delete", 0).All(&users)
		_, _ = o.QueryTable(new(models.Group)).Filter("is_delete", 0).All(&groups)
		_, _ = o.QueryTable(new(models.DeviceType)).Filter("is_delete", 0).All(&deviceTypes)
		_, _ = o.QueryTable(new(models.Brand)).Filter("is_delete", 0).All(&brands)
		_ = o.QueryTable(new(models.Computer)).Filter("id", id).Filter("is_delete", 0).RelatedSel().One(&computer)

		p.Data["users"] = users
		p.Data["groups"] = groups
		p.Data["deviceTypes"] = deviceTypes
		p.Data["brands"] = brands
		p.Data["computer"] = computer
		p.Data["action"] = "Edit"
		p.TplName = "assets/computer_form.html"
	} else if p.Ctx.Input.IsPost() {
		id, _ := p.GetInt("id")
		location := p.GetString("location")
		deviceTypeId, _ := p.GetInt("deviceType")
		name := p.GetString("name")
		brandId, _ := p.GetInt("brand")
		model := p.GetString("model")
		username := p.GetString("username")
		deptId, _ := p.GetInt("dept")
		assetNumber := p.GetString("name")
		purchaseDate := p.GetString("purchase_date")
		sn := p.GetString("sn")
		quickServiceCode := p.GetString("quick_service_code")
		warranty := p.GetString("warranty")
		remark := p.GetString("remark")
		isActive, _ := p.GetInt("is_active")
		isBreakdown, _ := p.GetInt("is_breakdown")
		inRepository, _ := p.GetInt("in_repository")
		cpu := p.GetString("cpu")
		memory := p.GetString("memory")
		disk := p.GetString("disk")
		systemV := p.GetString("os")
		o := orm.NewOrm()
		_, err := o.QueryTable(new(models.Computer)).Filter("id", id).Update(orm.Params{"location": location,
			"device_types_id": deviceTypeId, "brands_id": brandId, "model": model, "username": username,
			"groups_id": deptId, "asset_no": assetNumber, "purchase_date": purchaseDate, "sn": sn,
			"quick_service_code": quickServiceCode, "warranty": warranty, "is_active": isActive, "remark": remark,
			"is_breakdown": isBreakdown, "in_repository": inRepository, "name": name, "cpu": cpu, "memory": memory,
			"disk": disk, "os": systemV,
		})
		if err == nil {
			p.Data["json"] = map[string]interface{}{"code": 0, "msg": "更新成功"}
		} else {
			p.Data["json"] = map[string]interface{}{"code": 1001, "msg": err.Error()}
		}
		p.ServeJSON()
	}
}

func (p *PCController) MyComputer() {
	if p.Ctx.Input.IsGet() {
		username := p.GetSession("username").(string)
		// currentUser := models.User{}
		// var users []models.User
		var groups []models.Group
		var deviceTypes []models.DeviceType
		var brands []models.Brand
		var computer models.Computer
		o := orm.NewOrm()
		// _ = o.QueryTable(new(models.User)).Filter("username", username).One(&currentUser)
		_, _ = o.QueryTable(new(models.Group)).Filter("is_delete", 0).All(&groups)
		_, _ = o.QueryTable(new(models.DeviceType)).Filter("is_delete", 0).All(&deviceTypes)
		_, _ = o.QueryTable(new(models.Brand)).Filter("is_delete", 0).All(&brands)
		err := o.QueryTable(new(models.Computer)).Filter("username", username).Filter("is_delete", 0).
			RelatedSel().One(&computer)
		if err == nil {
			// p.Data["users"] = users
			p.Data["groups"] = groups
			p.Data["deviceTypes"] = deviceTypes
			p.Data["brands"] = brands
			p.Data["computer"] = computer
			p.TplName = "assets/my_computer.html"
		} else {
			retMsg := map[string]interface{}{"err": "当前用户没有查询到主机记录"}
			p.Data["json"] = retMsg
			p.ServeJSON()
		}
	} else if p.Ctx.Input.IsPost() {
		id, _ := p.GetInt("id")
		location := p.GetString("location")
		deviceTypeId, _ := p.GetInt("deviceType")
		name := p.GetString("name")
		brandId, _ := p.GetInt("brand")
		model := p.GetString("model")
		deptId, _ := p.GetInt("dept")
		assetNumber := p.GetString("name")
		purchaseDate := p.GetString("purchase_date")
		sn := p.GetString("sn")
		quickServiceCode := p.GetString("quick_service_code")
		warranty := p.GetString("warranty")
		o := orm.NewOrm()
		_, err := o.QueryTable(new(models.Computer)).Filter("id", id).Update(orm.Params{"location": location,
			"device_types_id": deviceTypeId, "brands_id": brandId, "model": model, "groups_id": deptId,
			"asset_no": assetNumber, "purchase_date": purchaseDate, "sn": sn, "quick_service_code": quickServiceCode,
			"warranty": warranty, "name": name,
		})

		if err == nil {
			p.Data["json"] = map[string]interface{}{"code": 0, "msg": "更新成功"}
		} else {
			p.Data["json"] = map[string]interface{}{"code": 1001, "msg": err.Error()}
		}
		p.ServeJSON()
	}
}

func (p *PCController) Export2Excel() {
	uName := p.GetSession("username").(string)
	if !tools.HasPermission(uName, "view_computer") {
		p.Redirect(beego.URLFor("ErrorController.Error403"), 302)
	}
	o := orm.NewOrm()
	qs := o.QueryTable(new(models.Computer)).Filter("is_delete", 0)
	var computers []models.Computer
	_, _ = qs.RelatedSel().All(&computers)

	titles := []string{
		"位置", "主机名", "设备类型", "品牌", "型号", "当前使用者", "资产分组", "资产编号", "购买时间", "序列号", "快速服务代码",
		"质保时间", "逻辑删除1删除0未删除", "1启用0停用", "1报废0正常", "1仓库中0仓库外", "IP地址", "备注", "主机cpu", "主机内存",
		"主机磁盘", "操作系统",
	}
	cells := []string{
		"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U",
		"V",
	}
	// 新建工作簿
	f := excelize.NewFile()
	index := f.NewSheet("Sheet1")

	// 列标题赋值
	for i, v := range titles {
		_ = f.SetCellValue("Sheet1", cells[i]+"1", v)
	}

	var computerList [][]string

	for _, element := range computers {
		tmpArr := []string{element.Location, element.Name, element.DeviceTypes.DeviceType, element.Brands.Name,
			element.Model, element.UserName, element.Groups.GroupName, element.AssetNo, element.PurchaseDate,
			element.SN, element.QuickServiceCode, element.Warranty, strconv.Itoa(element.IsDelete),
			strconv.Itoa(element.IsActive), strconv.Itoa(element.IsBreakdown), strconv.Itoa(element.InRepository),
			element.IpAddress, element.Remark, element.Cpu, element.Memory, element.Disk, element.OS,
		}
		computerList = append(computerList, tmpArr)
	}

	for i, v := range computerList {
		// i 行数
		// v 行内容
		columnNum := 0 // 列数
		for _, e := range v {
			pos := cells[columnNum] + strconv.Itoa(i+2)
			fmt.Println(pos)
			fmt.Println(e)
			_ = f.SetCellValue("Sheet1", pos, e)
			columnNum++
		}
	}

	f.SetActiveSheet(index)
	f.Path = "static/excels/computers.xlsx"
	if err := f.Save(); err != nil {
		fmt.Println(err)
	}
	p.Ctx.Output.Download(f.Path, "computers.xlsx")
}
