package auth

import (
	"dominos_cmdb/models"
	"dominos_cmdb/tools"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
	"time"
)

type UserController struct {
	beego.Controller
}

func (u *UserController) List() {
	username := u.GetSession("username").(string)
	if !tools.HasPermission(username, "view_user") {
		u.Redirect(beego.URLFor("ErrorController.Error403"), 302)
	}
	var users []models.User
	o := orm.NewOrm()
	qs := o.QueryTable(new(models.User))
	// 定义每页显示的记录数
	limitNum := 10
	// 获取前端传过来的参数
	currentPage, err := u.GetInt("page")
	if err != nil { // 说明没有获取到当前页
		currentPage = 1
	}
	offsetNum := limitNum * (currentPage - 1)
	searchKey := u.GetString("ss")
	var count int64 = 0
	if searchKey != "" { // 搜索有值
		count, _ = qs.Filter("is_delete", 0).Filter("username__contains", searchKey).Count()
		_, _ = qs.Filter("is_delete", 0).Filter("username__contains", searchKey).RelatedSel().
			Limit(limitNum).Offset(offsetNum).All(&users)
	} else {
		count, _ = qs.Filter("is_delete", 0).Count()
		_, _ = qs.Filter("is_delete", 0).RelatedSel().Limit(limitNum).Offset(offsetNum).All(&users)
	}
	// countPage := int(math.Ceil(float64(count)) / float64(limitNum))
	pageNum := tools.GetPageNum(count, limitNum)
	/*
		分面逻辑    当前页  offset limit   offset计算方法
					1       0      2       2 * (1 - 1)
					2       2      2       2 * (2 - 1)
					3       4      2       2 * (3 - 1)
				limitNum * (currentPage - 1)
	*/
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
	hasPrev, hasNext, prevPageNum, nextPageNum := tools.HasNext(currentPage, pageNum)

	u.Data["users"] = users
	u.Data["previousPage"] = previousPage
	u.Data["nextPage"] = nextPage
	u.Data["currentPage"] = currentPage
	u.Data["count"] = count
	u.Data["hasPrev"] = hasPrev
	u.Data["hasNext"] = hasNext
	u.Data["nextPageNum"] = nextPageNum
	u.Data["prevPageNum"] = prevPageNum
	u.Data["leftPages"] = leftPages
	u.Data["rightPages"] = rightPages
	u.Data["leftHasMore"] = leftHasMore
	u.Data["rightHasMore"] = rightHasMore
	u.Data["firstPage"] = firstPageNum
	u.Data["endPageNum"] = pageNum
	u.Data["ss"] = searchKey

	u.TplName = "auth/user_list.html"
}

func (u *UserController) Json() {
	var users []models.User
	o := orm.NewOrm()
	_, _ = o.QueryTable(new(models.User)).All(&users)
	u.Data["json"] = users
	u.ServeJSON()
}

func (u *UserController) Add() {
	if u.Ctx.Input.IsPost() {
		gid, _ := u.GetInt("gid")
		username := u.GetString("username")
		password := u.GetString("password")
		age, _ := u.GetInt("age")
		email := u.GetString("email")
		phone := u.GetString("phone")
		address := u.GetString("address")
		isActive, _ := u.GetInt("is_active")
		gender, _ := u.GetInt("gender")
		isSuperuser, _ := u.GetInt("is_superuser")
		isStaff, _ := u.GetInt("is_staff")
		job := u.GetString("job")
		phoneInt64, _ := strconv.ParseInt(phone, 10, 64)
		groups := models.Group{Id: gid}
		userData := models.User{
			UserName:    username,
			Password:    tools.GetMD5(password),
			Age:         age,
			Phone:       phoneInt64,
			Email:       email,
			Address:     address,
			IsActive:    isActive,
			Gender:      gender,
			Groups:      &groups,
			IsSuperUser: isSuperuser,
			IsStaff:     isStaff,
			Job:         job,
		}
		o := orm.NewOrm()
		_, err := o.Insert(&userData)
		resMsg := map[string]interface{}{}
		if err != nil {
			resMsg["code"] = 1000
			resMsg["msg"] = err.Error()
			u.Data["json"] = resMsg
		} else {
			resMsg["code"] = 0
			resMsg["msg"] = "添加成功"
			u.Data["json"] = resMsg
		}
		u.ServeJSON()
	} else {
		username := u.GetSession("username").(string)
		if !tools.HasPermission(username, "add_user") {
			u.Redirect(beego.URLFor("ErrorController.Error403"), 302)
		}
		var groups []models.Group
		o := orm.NewOrm()
		_, err := o.QueryTable(new(models.Group)).All(&groups)
		if err != nil {
			fmt.Println(err.Error())
		}
		u.Data["groups"] = groups
		u.Data["action"] = "Add"
		u.TplName = "auth/user_form.html"
	}
}

func (u *UserController) IsActive() {
	username := u.GetSession("username").(string)
	if !tools.HasPermission(username, "change_user") {
		u.Redirect(beego.URLFor("ErrorController.Error403"), 302)
	}
	isActive, _ := u.GetInt("is_active")
	id, _ := u.GetInt("id")
	u.Data["json"] = tools.CommonUpdateIsActive(new(models.User), id, isActive)
	u.ServeJSON()
}

func (u *UserController) Delete() {
	username := u.GetSession("username").(string)
	if !tools.HasPermission(username, "change_user") {
		u.Redirect(beego.URLFor("ErrorController.Error403"), 302)
	}
	id, _ := u.GetInt("id")
	u.Data["json"] = tools.CommonLogicDelete(id, new(models.User))
	u.ServeJSON()
}

func (u *UserController) ResetPass() {
	username := u.GetSession("id").(string)
	if !tools.HasPermission(username, "change_user") {
		u.Redirect(beego.URLFor("ErrorController.Error403"), 302)
	}
	id, _ := u.GetInt("id")
	password := u.GetString("password")
	confirmPassword := u.GetString("confirm_password")
	resMsg := map[string]interface{}{}
	if password != confirmPassword {
		resMsg["msg"] = "两次输入不一致"
		resMsg["code"] = 111
	} else if len(password) < 6 {
		resMsg["msg"] = "密码长度不足6位"
		resMsg["code"] = 112
	} else {
		md5Password := tools.GetMD5(password)
		o := orm.NewOrm()
		_, err := o.QueryTable(new(models.User)).Filter("id", id).Update(orm.Params{
			"password": md5Password,
		})
		if err != nil {
			resMsg["msg"] = err
			resMsg["code"] = 1001
		} else {
			resMsg["msg"] = "OK"
			resMsg["code"] = 0
		}
	}
	u.Data["json"] = resMsg
	u.ServeJSON()
}

func (u *UserController) Edit() {
	if u.Ctx.Input.IsPost() {
		uName := u.GetSession("username").(string)
		if !tools.HasPermission(uName, "change_user") {
			u.Redirect(beego.URLFor("ErrorController.Error403"), 302)
		}
		id, _ := u.GetInt("uid")
		gid, _ := u.GetInt("gid")
		username := u.GetString("username")
		age, _ := u.GetInt("age")
		email := u.GetString("email")
		phone := u.GetString("phone")
		address := u.GetString("address")
		isActive, _ := u.GetInt("is_active")
		gender, _ := u.GetInt("gender")
		isSuperuser, _ := u.GetInt("is_superuser")
		isStaff, _ := u.GetInt("is_staff")
		job := u.GetString("job")
		phoneInt64, _ := strconv.ParseInt(phone, 10, 64)
		o := orm.NewOrm()
		_, err := o.QueryTable(new(models.User)).Filter("id", id).Update(orm.Params{
			"groups_id":    gid,
			"username":     username,
			"age":          age,
			"email":        email,
			"address":      address,
			"is_active":    isActive,
			"gender":       gender,
			"phone":        phoneInt64,
			"is_superuser": isSuperuser,
			"is_staff":     isStaff,
			"create_time":  time.Now(),
			"job":          job,
		})
		resMsg := map[string]interface{}{}
		if err != nil {
			resMsg["code"] = 1062
			resMsg["msg"] = err.Error()
			u.Data["json"] = resMsg
		} else {
			resMsg["code"] = 0
			resMsg["msg"] = "编辑成功"
			u.Data["json"] = resMsg
		}
		u.ServeJSON()
	} else if u.Ctx.Input.IsGet() {
		username := u.GetSession("username").(string)
		if !tools.HasPermission(username, "change_user") {
			u.Redirect(beego.URLFor("ErrorController.Error403"), 302)
		}
		id, _ := u.GetInt("uid")
		user := models.User{}
		var groups []models.Group
		o := orm.NewOrm()
		_, _ = o.QueryTable(new(models.Group)).All(&groups)
		err := o.QueryTable(new(models.User)).Filter("id", id).One(&user)
		if err == nil {
			u.Data["user"] = user
			u.Data["successURI"] = u.Ctx.Request.Referer()
		}
		u.Data["action"] = "Edit"
		u.Data["groups"] = groups
		u.TplName = "auth/user_form.html"
	}
}

func (u *UserController) BatchDelete() {
	username := u.GetSession("username").(string)
	if !tools.HasPermission(username, "change_user") {
		u.Redirect(beego.URLFor("ErrorController.Error403"), 302)
	}
	ids := u.GetString("ids")
	u.Data["json"] = tools.CommonBatchLogicDelete(ids, new(models.User))
	u.ServeJSON()
}

func (u *UserController) GetPerJson() {
	userId, _ := u.GetInt("uid")
	o := orm.NewOrm()
	// 已绑定的权限
	currentUser := models.User{Id: userId}
	_, err := o.LoadRelated(&currentUser, "Permissions")
	if err != nil {
		fmt.Println(err)
	}
	var selectIds []int
	for _, per := range currentUser.Permissions {
		selectIds = append(selectIds, per.Id)
	}
	// 所有权限
	var permissions []models.Permission
	_, _ = o.QueryTable(new(models.Permission)).All(&permissions)
	var arrayMap []map[string]interface{}

	for _, item := range permissions {
		perMap := map[string]interface{}{"value": item.Id, "title": item.Name, "align": "center"}
		arrayMap = append(arrayMap, perMap)
	}

	resMap := map[string]interface{}{"selectedIds": selectIds, "permissions": arrayMap}
	u.Data["json"] = resMap
	u.ServeJSON()
}

func (u *UserController) BindPer() {
	if u.Ctx.Input.IsGet() {
		uName := u.GetSession("username").(string)
		if !tools.HasPermission(uName, "change_user") {
			u.Redirect(beego.URLFor("ErrorController.Error403"), 302)
		}
		userId, _ := u.GetInt("uid")
		o := orm.NewOrm()
		qs := o.QueryTable(new(models.User))
		user := models.User{}
		_ = qs.Filter("id", userId).One(&user)
		u.Data["user"] = user
		u.TplName = "auth/user_per_bind.html"
	} else if u.Ctx.Input.IsPost() {
		uName := u.GetSession("username").(string)
		if !tools.HasPermission(uName, "change_user") {
			u.Redirect(beego.URLFor("ErrorController.Error403"), 302)
		}
		userId, _ := u.GetInt("user_id")
		idsStr := u.GetString("ids")
		o := orm.NewOrm()
		user := models.User{Id: userId}
		m2m := o.QueryM2M(&user, "Permissions")
		if idsStr == "" {
			// 如果检测为空字符串，直接清理permission表
			_, _ = m2m.Clear()
		} else {
			idsArr := strings.Split(idsStr, "&")
			_, _ = m2m.Clear()
			for _, perId := range idsArr {
				perId := tools.StrToInt(perId)
				permission := models.Permission{Id: perId}
				m2m := o.QueryM2M(&user, "Permissions")
				_, _ = m2m.Add(permission)
			}
		}

		u.Data["json"] = map[string]interface{}{"code": 0, "msg": "绑定完成"}
		u.ServeJSON()
	}
}
