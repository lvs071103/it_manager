package tools

import (
	"github.com/astaxie/beego/orm"
)

func CommonLogicDelete(sid int, tableName interface{}) interface{} {
	var messages map[string]interface{}
	o := orm.NewOrm()
	_, err := o.QueryTable(tableName).Filter("id", sid).Update(orm.Params{"is_delete": 1})
	if err == nil {
		messages = map[string]interface{}{"code": 0, "msg": "删除成功"}
	} else {
		messages = map[string]interface{}{"code": 1001, "msg": err.Error()}
	}

	return messages
}
