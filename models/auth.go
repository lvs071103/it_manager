package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type ContentType struct {
	Id          int           `orm:"pk;auto"`
	AppLabel    string        `orm:"column(app_label);size(100)"`
	Model       string        `orm:"column(model);size(100);unique"`
	Permissions []*Permission `orm:"reverse(many)"`
}

func (c *ContentType) TableName() string {
	return "sys_content_type"
}

func init() {
	orm.RegisterModel(new(ContentType))
}

type Permission struct {
	Id          int          `orm:"pk;auto"`
	Name        string       `orm:"column(name);description(名称);size(64)"`
	ContentType *ContentType `orm:"rel(fk)"`
	CodeName    string       `orm:"column(code_name);description(code_name);unique"`
	Users       []*User      `orm:"reverse(many)"`
	Groups      []*Group     `orm:"reverse(many)"`
}

func (p *Permission) TableName() string {
	return "auth_permission"
}

func init() {
	orm.RegisterModel(new(Permission))
}

type Group struct {
	Id          int
	GroupName   string        `orm:"column(group_name);unique;size(64)"`
	IsActive    int           `orm:"column(is_active);description(1代表启用0代表停用);default(1)"`
	IsDelete    int           `orm:"column(is_delete);description(1删除0未删除);default(0)"`
	Users       []*User       `orm:"reverse(many)"`
	Desc        string        `orm:"column(desc);size(255);null"`
	Permissions []*Permission `orm:"rel(m2m)"`
	Computers   []*Computer   `orm:"reverse(many)"`
}

func (g *Group) TableName() string {
	return "groups"
}

func init() {
	orm.RegisterModel(new(Group))
}

type User struct {
	Id          int       `orm:"column(id);pk;auto"`
	Groups      *Group    `orm:"rel(fk);description(用户组id);null"`
	UserName    string    `orm:"column(username);unique;size(64);description(用户名)"`
	Password    string    `orm:"column(password);size(32);description(登陆密码)"`
	Email       string    `orm:"column(email);null;size(320);description(邮箱)"`
	Age         int       `orm:"column(age);null;description(年龄)"`
	Gender      int       `orm:"column(gender);null;description(性别:1代表男,2代表女)"`
	Phone       int64     `orm:"column(phone);null;description(电话号码)"`
	Address     string    `orm:"column(address);null;size(255);description(住址)"`
	IsSuperUser int       `orm:"column(is_superuser);description(是否是管理员);default(0)"`
	IsStaff     int       `orm:"column(is_staff);description(是否是员工);default(0)"`
	IsActive    int       `orm:"column(is_active);description(1启用0停用);default(1)"`
	IsDelete    int       `orm:"column(is_delete);description(1删除0未删除);default(0)"`
	LastLogin   time.Time `orm:"column(last_login);auto_now_add;type(datetime);description(最近登陆);null"`
	CreateTime  time.Time `orm:"column(create_time);auto_now;type(datetime);description(创建时间);null"`
	ChineseName string    `orm:"column(chinese_name);description(中文名);size(64);null"`
	Job         string    `orm:"column(job);description(职位);null"`
	//Computer    []*Computer     `orm:"reverse(many)"`
	Computers   []*Computer   `orm:"reverse(many)"`
	Messages    []*Message    `orm:"reverse(many)"`
	Permissions []*Permission `orm:"rel(m2m)"`
}

func (u *User) TableName() string {
	return "auth_user"
}

func init() {
	orm.RegisterModel(new(User))
}
