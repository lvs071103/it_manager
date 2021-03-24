package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type Icon struct {
	Id         int       `orm:"pk;auto"`
	Name       string    `orm:"column(name);description(Icon名称);unique;"`
	ClassName  string    `orm:"column(class_name);description(类名);unique;"`
	IsActive   int       `orm:"column(is_active);description(1启用0停用);default(1)"`
	IsDelete   int       `orm:"column(is_delete);description(1删除0未删除);default(0)"`
	CreateTime time.Time `orm:"column(create_time);auto_now_add;null"`
	Desc       string    `orm:"column(desc);description(描述);null"`
	Menus      []*Menu   `orm:"reverse(many)"`
}

func (i *Icon) TableName() string {
	return "sys_icon"
}

func init() {
	orm.RegisterModel(new(Icon))
}

type Menu struct {
	Id         int       `orm:"column(id);pk;auto"`
	Name       string    `orm:"column(name);unique;size(64);description(菜单名称)"`
	UrlFor     string    `orm:"column(url_for);size(255);description(url)"`
	Pid        int       `orm:"column(pid);description(父节点id)"`
	IsActive   int       `orm:"column(is_active);description(1启用0停用);default(1)"`
	IsDelete   int       `orm:"column(is_delete);description(1删除0未删除);default(0)"`
	CreateTime time.Time `orm:"column(create_time);description(创建时间);auto_now_add;type(datetime);null"`
	Weight     int       `orm:"column(weight);description(权重,数值越大权重越高)"`
	Desc       string    `orm:"column(desc);size(255);description(描述);null"`
	Icons      *Icon     `orm:"rel(fk)"`
}

func (p *Menu) TableName() string {
	return "sys_menu"
}

func init() {
	orm.RegisterModel(new(Menu))
}

type MenuTree struct {
	Id       int
	Name     string
	UrlFor   string
	Weight   int
	Icon     string
	Children []*MenuTree
}
