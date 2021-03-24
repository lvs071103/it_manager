package tools

import (
	"github.com/astaxie/beego/orm"
	"strings"
)

func CommonBatchLogicDelete(ids string, tableName interface{}) interface{} {
	idArr := strings.Split(ids, "&")
	o := orm.NewOrm()
	qs := o.QueryTable(tableName)
	var errs []string
	var messages map[string]interface{}
	for _, id := range idArr {
		_, err := qs.Filter("id", id).Update(orm.Params{"is_delete": 1})
		if err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		messages = map[string]interface{}{"code": 1001, "msg": "批量删除失败"}
	} else {
		messages = map[string]interface{}{"code": 0, "msg": "批量删除成功"}
	}
	return messages
}
