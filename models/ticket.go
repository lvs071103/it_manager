package models

import (
	"github.com/astaxie/beego/orm"
)

type TicketCategory struct {
	Id       int       `orm:"column(id);pk;auto"`
	Name     string    `orm:"column(name);size(64);description(分类名称);unique"`
	Desc     string    `orm:"column(desc);size(255);description(备注);null"`
	IsActive int       `orm:"column(is_active);description(1代表启用0代表停用);default(1)"`
	IsDelete int       `orm:"column(is_delete);description(1代表删除0代表未删除);default(0)"`
	Tickets  []*Ticket `orm:"reverse(many)"`
}

type Ticket struct {
	Id            int             `orm:"column(id);pk;auto"`
	Title         string          `orm:"column(title);size(64);description(工单标题)"`
	Content       string          `orm:"column(content);type(text);"`
	Category      *TicketCategory `orm:"rel(fk)"`
	Policy        int             `orm:"column(policy);description(100代表重要,99代表普通);default(99)"`
	SubmitUser    string          `orm:"column(submit_user);description(提交用户);size(64);null"`
	Progress      int             `orm:"column(progress);description(进度);null"`
	ProcessUser   string          `orm:"column(process_user);description(处理用户);size(64);null"`
	ProcessResult string          `orm:"column(process_result);description(处理结果);size(64);type(text);null"`
	CreateTime    string          `orm:"column(create_time);description(提交时间);type(datetime);auto_now_add;null"`
	UpdateTime    string          `orm:"column(update_time);description(更新时间);type(datetime);auto_now;null"`
	Status        int             `orm:"column(status);description(状态1代表已完成,2代表未分配,3代表已分配,4代表处理中);default(0)"`
	IsDelete      int             `orm:"column(is_delete);description(是否删除1代表删除0代表未删除);default(0);"`
}

func (c *TicketCategory) TableName() string {
	return "tickets_category"
}

func (t *Ticket) TableName() string {
	return "tickets"
}

func init() {
	orm.RegisterModel(new(TicketCategory), new(Ticket))
}
