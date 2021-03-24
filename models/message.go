package models

import "github.com/astaxie/beego/orm"

type Message struct {
	Id       int    `orm:"column(id);pk;auto"`
	User     *User  `orm:"column(uid);rel(fk)"`
	Flag     int    `orm:"column(flag);description(1代表所有通知2代表工单通知3代表添加新主机);default(1)"`
	Title    string `orm:"column(title);description(标题);size(64)"`
	Content  string `orm:"column(content);description(内容);type(text)"`
	IsRead   int    `orm:"column(is_read);description(1代表已读,0代表未读);default(0)"`
	IsDelete int    `orm:"column(is_delete);description(1代表删除,0代表未删除);default(0)"`
}

func (m *Message) TableName() string {
	return "message_notification"
}

func init() {
	orm.RegisterModel(new(Message))
}
