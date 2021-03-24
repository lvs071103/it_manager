package tickets

import (
	"dominos_cmdb/models"
	"dominos_cmdb/tools"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type TicketController struct {
	beego.Controller
}

func (t *TicketController) List() {
	var TcList []models.Ticket
	o := orm.NewOrm()
	qs := o.QueryTable(new(models.Ticket)).Filter("is_delete", 0)
	// 定义每页显示的记录数
	limitNum := 10
	// 获取当前页
	currentPage, err := t.GetInt("page")
	if err != nil {
		// 说明没有获取到当前页
		currentPage = 1
	}
	offsetNum := limitNum * (currentPage - 1)
	var count int64 = 0
	searchKey := t.GetString("ss")
	if searchKey == "" {
		// 未输入搜索字段
		count, _ = qs.Count()
		_, _ = qs.RelatedSel().OrderBy("-Policy").Limit(limitNum).Offset(offsetNum).All(&TcList)
	} else {
		count, _ = qs.Filter("Title__contains", searchKey).Count()
		_, _ = qs.Filter("Title__contains", searchKey).Limit(limitNum).Offset(offsetNum).RelatedSel().All(&TcList)
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
	t.Data["previousPage"] = previousPage
	t.Data["nextPage"] = nextPage
	t.Data["currentPage"] = currentPage
	t.Data["count"] = count
	t.Data["hasPrev"] = hasPrev
	t.Data["hasNext"] = hasNext
	t.Data["leftHasMore"] = leftHasMore
	t.Data["rightHasMore"] = rightHasMore
	t.Data["leftPages"] = leftPages
	t.Data["rightPages"] = rightPages
	t.Data["firstPage"] = firstPageNum
	t.Data["endPageNum"] = pageNum
	t.Data["ss"] = searchKey
	t.Data["t_list"] = TcList
	t.TplName = "tickets/ticket_list.html"
}

func (t *TicketController) MyTickets() {
	currentUser := t.GetSession("username").(string)
	var TcList []models.Ticket
	o := orm.NewOrm()
	qs := o.QueryTable(new(models.Ticket)).Filter("is_delete", 0)
	// 定义每页显示的记录数
	limitNum := 10
	// 获取当前页
	currentPage, err := t.GetInt("page")
	if err != nil {
		// 说明没有获取到当前页
		currentPage = 1
	}
	offsetNum := limitNum * (currentPage - 1)
	var count int64 = 0
	searchKey := t.GetString("ss")
	if searchKey == "" {
		// 未输入搜索字段
		count, _ = qs.Filter("SubmitUser", currentUser).Count()
		_, _ = qs.Filter("SubmitUser", currentUser).RelatedSel().OrderBy("-Policy").Limit(limitNum).
			Offset(offsetNum).All(&TcList)
	} else {
		count, _ = qs.Filter("SubmitUser", currentUser).Filter("Title__contains", searchKey).Count()
		_, _ = qs.Filter("SubmitUser", currentUser).Filter("Title__contains", searchKey).Limit(limitNum).
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
	t.Data["previousPage"] = previousPage
	t.Data["nextPage"] = nextPage
	t.Data["currentPage"] = currentPage
	t.Data["count"] = count
	t.Data["hasPrev"] = hasPrev
	t.Data["hasNext"] = hasNext
	t.Data["leftHasMore"] = leftHasMore
	t.Data["rightHasMore"] = rightHasMore
	t.Data["leftPages"] = leftPages
	t.Data["rightPages"] = rightPages
	t.Data["firstPage"] = firstPageNum
	t.Data["endPageNum"] = pageNum
	t.Data["ss"] = searchKey
	t.Data["t_list"] = TcList
	t.TplName = "tickets/my_tickets.html"
}

func (t *TicketController) Add() {
	if t.Ctx.Input.IsGet() {
		var categories []models.TicketCategory
		o := orm.NewOrm()
		_, _ = o.QueryTable(new(models.TicketCategory)).All(&categories)
		fmt.Println(categories)
		t.Data["categories"] = categories
		t.TplName = "tickets/ticket_form.html"
	} else if t.Ctx.Input.IsPost() {
		categoryId, _ := t.GetInt("category")
		title := t.GetString("title")
		content := t.GetString("content")
		policy, _ := t.GetInt("policy")
		o := orm.NewOrm()
		category := models.TicketCategory{Id: categoryId}
		ticket := models.Ticket{
			Title:      title,
			Content:    content,
			Policy:     policy,
			Category:   &category,
			SubmitUser: t.GetSession("username").(string),
			CreateTime: time.Now().Format("2006-01-01"),
			Status:     2,
		}
		_, err := o.Insert(&ticket)
		resMsg := map[string]interface{}{}
		if err == nil {
			resMsg["code"] = 0
			resMsg["msg"] = "添加成功"
			// 此时，应该发送一则消息给管理员，提醒管理员，提交了一条工单
			var users []models.User
			_, _ = o.QueryTable(new(models.User)).Filter("username", "admin").All(&users)
			// 查出所有职位为桌面工程师所有人
			for _, user := range users {
				notifiedUser := models.User{Id: user.Id}
				message := models.Message{
					Flag:    1,
					Title:   title,
					Content: content,
					User:    &notifiedUser,
					IsRead:  0,
				}
				_, _ = o.Insert(&message)
			}
			// 发送邮件通知相应的桌面工程师
		} else {
			resMsg["code"] = 1001
			resMsg["msg"] = err.Error()
		}
		t.Data["json"] = resMsg

		t.ServeJSON()
	}
}

func (t *TicketController) Edit() {
	if t.Ctx.Input.IsGet() {
		id, _ := t.GetInt("id")
		ticket := models.Ticket{}
		o := orm.NewOrm()
		_ = o.QueryTable(new(models.Ticket)).Filter("id", id).RelatedSel().One(&ticket)
		t.Data["ticket"] = ticket
		t.TplName = "tickets/ticket_edit.html"
	} else if t.Ctx.Input.IsPost() {
		id, _ := t.GetInt("id")
		processUser := t.GetSession("username")
		processResult := t.GetString("process_result")
		status, _ := t.GetInt("status")
		progress, _ := t.GetInt("progress")
		o := orm.NewOrm()
		qs := o.QueryTable(new(models.Ticket))
		_, err := qs.Filter("id", id).Update(orm.Params{
			"ProcessUser":   processUser,
			"ProcessResult": processResult,
			"Status":        status,
			"Progress":      progress,
			"UpdateTime":    time.Now().Format("2006-01-01"),
		})
		returnMp := map[string]interface{}{}
		if err == nil {
			returnMp["code"] = 0
			returnMp["msg"] = "更新成功"
		} else {
			returnMp["code"] = 1001
			returnMp["msg"] = err.Error()
		}
		t.Data["json"] = returnMp
		t.ServeJSON()
	}
}

func (t *TicketController) Del() {
	id, _ := t.GetInt("id")
	t.Data["json"] = tools.CommonLogicDelete(id, new(models.Ticket))
	t.ServeJSON()
}

func (t *TicketController) BatchDelete() {
	ids := t.GetString("ids")
	t.Data["json"] = tools.CommonBatchLogicDelete(ids, new(models.Ticket))
	t.ServeJSON()
}

func (t *TicketController) IsActive() {
	isActive, _ := t.GetInt("is_active")
	id, _ := t.GetInt("id")
	t.Data["json"] = tools.CommonUpdateIsActive(new(models.Ticket), id, isActive)
	t.ServeJSON()
}

func (t *TicketController) Detail() {
	id, _ := t.GetInt("id")
	ticket := models.Ticket{}
	o := orm.NewOrm()
	_ = o.QueryTable(new(models.Ticket)).Filter("id", id).RelatedSel().One(&ticket)
	t.Data["ticket"] = ticket
	t.TplName = "tickets/ticket_detail.html"
}

func (t *TicketController) Assign() {
	if t.Ctx.Input.IsGet() {
		id, _ := t.GetInt("id")
		ticket := models.Ticket{}
		var users []models.User
		o := orm.NewOrm()
		_ = o.QueryTable(new(models.Ticket)).Filter("id", id).RelatedSel().One(&ticket)
		//查询desktop组内用户
		r := o.Raw("select u.id, u.username from auth_user u JOIN auth_group g on "+
			"u.groups_id = g.id WHERE g.group_name = ?", "desktop")
		_, _ = r.QueryRows(&users)
		fmt.Println(users)
		t.Data["users"] = users
		t.Data["ticket"] = ticket
		t.TplName = "tickets/tickets_assign.html"
	} else if t.Ctx.Input.IsPost() {
		id, _ := t.GetInt("id")
		processUser := t.GetString("username")
		user := models.User{}
		o := orm.NewOrm()
		qs := o.QueryTable(new(models.Ticket))
		_, err := qs.Filter("id", id).Update(orm.Params{
			"ProcessUser": processUser,
			"Status":      3,
			"UpdateTime":  time.Now().Format("2006-01-01"),
		})
		_ = o.QueryTable(new(models.User)).Filter("username", processUser).One(&user)
		receiveUser := models.User{Id: user.Id}
		returnMp := map[string]interface{}{}
		if err == nil {
			returnMp["code"] = 0
			returnMp["msg"] = "分配成功"
			content := fmt.Sprintf("管理员给你分配一个新的工单,工单号: %d, 请及时查看", id)
			// 查出所有职位为桌面工程师所有人
			message := models.Message{
				Flag:    1,
				Title:   "管理员分配了新的工单",
				Content: content,
				User:    &receiveUser,
				IsRead:  0,
			}
			_, _ = o.Insert(&message)
		} else {
			returnMp["code"] = 1001
			returnMp["msg"] = err.Error()
		}
		t.Data["json"] = returnMp
		t.ServeJSON()
	}
}
