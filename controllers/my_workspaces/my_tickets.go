package my_workspaces

import (
	"dominos_cmdb/models"
	"dominos_cmdb/tools"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type MyTCController struct {
	beego.Controller
}

func (m *MyTCController) Preprocess() {
	currentUser := m.GetSession("username").(string)
	var TcList []models.Ticket
	o := orm.NewOrm()
	qs := o.QueryTable(new(models.Ticket)).Filter("is_delete", 0)
	// 定义每页显示的记录数
	limitNum := 10
	// 获取当前页
	currentPage, err := m.GetInt("page")
	if err != nil {
		// 说明没有获取到当前页
		currentPage = 1
	}
	offsetNum := limitNum * (currentPage - 1)
	var count int64 = 0
	searchKey := m.GetString("ss")
	if searchKey == "" {
		// 未输入搜索字段
		count, _ = qs.Filter("ProcessUser", currentUser).Count()
		_, _ = qs.Filter("ProcessUser", currentUser).RelatedSel().OrderBy("-Policy").Limit(limitNum).
			Offset(offsetNum).All(&TcList)
	} else {
		count, _ = qs.Filter("ProcessUser", currentUser).Filter("Title__contains", searchKey).Count()
		_, _ = qs.Filter("ProcessUser", currentUser).Filter("Title__contains", searchKey).Limit(limitNum).
			Offset(offsetNum).RelatedSel().All(&TcList)
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
	m.Data["previousPage"] = previousPage
	m.Data["nextPage"] = nextPage
	m.Data["currentPage"] = currentPage
	m.Data["count"] = count
	m.Data["hasPrev"] = hasPrev
	m.Data["hasNext"] = hasNext
	m.Data["leftHasMore"] = leftHasMore
	m.Data["rightHasMore"] = rightHasMore
	m.Data["leftPages"] = leftPages
	m.Data["rightPages"] = rightPages
	m.Data["firstPage"] = firstPageNum
	m.Data["endPageNum"] = pageNum
	m.Data["ss"] = searchKey
	m.Data["t_list"] = TcList
	m.TplName = "tickets/preprocess_tickets.html"
}
