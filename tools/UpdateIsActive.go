package tools

import (
	"github.com/astaxie/beego/orm"
)

func CommonUpdateIsActive(tableName interface{}, id, isActive int) interface{} {
	o := orm.NewOrm()
	qs := o.QueryTable(tableName).Filter("id", id)
	resMsg := map[string]interface{}{}
	if isActive == 0 {
		_, err := qs.Update(orm.Params{"is_active": 1})
		if err == nil {
			resMsg["msg"] = "启用成功"
		} else {
			resMsg["msg"] = err.Error()
		}
	} else if isActive == 1 {
		_, err := qs.Update(orm.Params{"is_active": 0})
		if err == nil {
			resMsg["msg"] = "停用成功"
		} else {
			resMsg["msg"] = err.Error()
		}
	}

	return resMsg
}
