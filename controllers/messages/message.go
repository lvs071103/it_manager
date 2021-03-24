package messages

import (
	"dominos_cmdb/models"
	"dominos_cmdb/tools"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strings"
)

type MessageController struct {
	beego.Controller
}

func (m *MessageController) List() {
	var Messages []models.Message
	o := orm.NewOrm()
	qs := o.QueryTable(new(models.Message)).Filter("is_delete", 0)
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
		count, _ = qs.Count()
		_, _ = qs.RelatedSel().Limit(limitNum).Offset(offsetNum).All(&Messages)
	} else {
		count, _ = qs.Filter("Title__contains", searchKey).Count()
		_, _ = qs.Filter("Title__contains", searchKey).Limit(limitNum).Offset(offsetNum).RelatedSel().All(&Messages)
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
	m.Data["messages"] = Messages
	m.TplName = "messages/message_list.html"
}

func (m *MessageController) MyList() {
	username := m.GetSession("username")
	fmt.Println(username)
	currentUser := models.User{}
	_ = orm.NewOrm().QueryTable(new(models.User)).Filter("username", username).One(&currentUser)
	var Messages []models.Message
	o := orm.NewOrm()
	qs := o.QueryTable(new(models.Message)).Filter("is_delete", 0).Filter("uid", currentUser.Id)
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
		count, _ = qs.Count()
		_, _ = qs.RelatedSel().Limit(limitNum).Offset(offsetNum).All(&Messages)
	} else {
		count, _ = qs.Filter("Title__contains", searchKey).Count()
		_, _ = qs.Filter("Title__contains", searchKey).Limit(limitNum).Offset(offsetNum).RelatedSel().All(&Messages)
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
	m.Data["messages"] = Messages
	m.TplName = "messages/my_messages.html"
}

func (m *MessageController) IsRead() {
	id, _ := m.GetInt("id")
	message := models.Message{}
	o := orm.NewOrm()
	qs := o.QueryTable(new(models.Message)).Filter("id", id)
	_, _ = qs.Update(orm.Params{
		"is_read": 1,
	})
	_ = qs.RelatedSel().One(&message)
	fmt.Println(m.Ctx.Request.Referer())
	m.Data["message"] = message
	m.Data["referer_url"] = m.Ctx.Request.Referer()
	m.TplName = "messages/message_form.html"
}

func (m *MessageController) SetRead() {
	ids := m.GetString("ids")
	idArr := strings.Split(ids, "&")
	o := orm.NewOrm()
	qs := o.QueryTable(new(models.Message)).Filter("is_delete", 0)
	var errs []string
	var retMsg map[string]interface{}
	for _, id := range idArr {
		_, err := qs.Filter("id", id).Update(orm.Params{"is_read": 1})
		if err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		retMsg = map[string]interface{}{"code": 1001, "msg": "操作失败"}
	} else {
		retMsg = map[string]interface{}{"code": 0, "msg": "操作成功"}
	}
	m.Data["json"] = retMsg
	m.ServeJSON()

}

func (m *MessageController) Del() {
	id, _ := m.GetInt("id")
	fmt.Println(id)
	m.Data["json"] = tools.CommonLogicDelete(id, new(models.Message))
	m.ServeJSON()
}

func (m *MessageController) BatchDel() {
	ids := m.GetString("ids")
	m.Data["json"] = tools.CommonBatchLogicDelete(ids, new(models.Message))
	m.ServeJSON()
}
