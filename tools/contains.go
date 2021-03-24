package tools

import (
	"dominos_cmdb/models"
	"github.com/astaxie/beego/orm"
)

func HasPermission(username string, perStr string) bool {
	user := models.User{UserName: username}
	o := orm.NewOrm()
	// _ = o.Read(&user)
	_ = o.QueryTable(new(models.User)).Filter("UserName", username).RelatedSel().One(&user)
	// fmt.Println(user)
	// 当前用户为超级管理员则返回true
	if user.IsSuperUser == 1 {
		return true
	}

	// 查询当前用户拥有的权限列表
	_, _ = o.LoadRelated(&user, "Permissions")

	var codeNames []string

	for _, per := range user.Permissions {
		codeNames = append(codeNames, per.CodeName)
	}

	for _, v := range codeNames {
		if v == perStr {
			return true
		}
	}

	return false
}
