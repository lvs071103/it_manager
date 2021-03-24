package echarts

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type EchController struct {
	beego.Controller
}

func (e *EchController) TicketsChart() {
	var processDate orm.ParamsList
	var userList orm.ParamsList
	//var seriesData []map[string]interface{}
	o := orm.NewOrm()
	_, _ = o.Raw("select distinct update_time from tickets limit 15").ValuesFlat(&processDate)
	_, _ = o.Raw("select u.username from auth_user u JOIN auth_group g "+
		"on u.groups_id = g.id WHERE g.group_name = ?", "desktop").ValuesFlat(&userList)

	//var userTickets map[string]interface{}
	var series []map[string]interface{}

	for i := 0; i < len(userList); i++ {
		seriesMap := map[string]interface{}{}
		seriesMap["name"] = userList[i]
		seriesMap["type"] = "line"
		seriesMap["stack"] = "总量"
		var countList []int
		for d := 0; d < len(processDate); d++ {
			var count int
			_ = o.Raw("select count(*) from tickets "+
				"where process_user = ? and update_time = ?", userList[i], processDate[d]).QueryRow(&count)
			countList = append(countList, count)
		}
		seriesMap["data"] = countList
		series = append(series, seriesMap)
	}

	retMsg := map[string]interface{}{
		"processDate": processDate,
		"userList":    userList,
		"series":      series,
	}

	e.Data["json"] = retMsg
	e.ServeJSON()

}
