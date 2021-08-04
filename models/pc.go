package models

import "github.com/astaxie/beego/orm"

type Brand struct {
	Id        int         `orm:"pk;auto"`
	Name      string      `orm:"column(name);size(64);description(品牌名称);unique"`
	IsDelete  int         `orm:"column(is_delete);description(1删除0未删除);default(0)"`
	IsActive  int         `orm:"column(is_active);description(1启用0停用);default(1)"`
	Desc      string      `orm:"column(desc);size(255);description(描述)"`
	Computers []*Computer `orm:"reverse(many)"`
}

type DeviceType struct {
	Id         int         `orm:"pk;auto"`
	DeviceType string      `orm:"column(device_type);size(64);description(设备类型);unique"`
	IsActive   int         `orm:"column(is_active);description(1启用0停用);default(1)"`
	IsDelete   int         `orm:"column(is_delete);description(1删除0未删除);default(0)"`
	Desc       string      `orm:"column(desc);size(255);null;description(描述)"`
	Computers  []*Computer `orm:"reverse(many)"`
}

type Computer struct {
	Id               int         `orm:"pk;auto"`
	Location         string      `orm:"column(location);description(位置,比如上海办公室)"`
	Name             string      `orm:"column(name);size(64);description(主机名);unique"`
	DeviceTypes      *DeviceType `orm:"rel(fk);description(设备类型)"`
	Brands           *Brand      `orm:"rel(fk);description(品牌);"`
	Model            string      `orm:"column(model);description(型号)"`
	CurrentUser      *User       `orm:"rel(one);null"`
	Department       *Group      `orm:"rel(fk);null"`
	AssetNo          string      `orm:"column(asset_no);size(12);description(资产编号);null"`
	PurchaseDate     string      `orm:"column(purchase_date);type(datetime);description(购买时间);null"`
	SN               string      `orm:"column(sn);size(20);description(序列号);null"`
	QuickServiceCode string      `orm:"column(quick_service_code);description(快速服务代码);null"`
	Warranty         string      `orm:"column(warranty);type(datetime);description(质保时间);null"`
	IsDelete         int         `orm:"column(is_delete);description(1删除0未删除);default(0)"`
	IsActive         int         `orm:"column(is_active);description(1启用0停用);default(1)"`
	IsBreakdown      int         `orm:"column(is_breakdown);description(1报废0正常);default(0)"`
	InRepository     int         `orm:"column(in_repository);description(1仓库中0仓库外);default(1)"`
	IpAddress        string      `orm:"column(ip_address);description(最近一次使用的IP地址);null"`
	Remark           string      `orm:"column(remark);description(备注);null"`
	PreUsers         []*User     `orm:"rel(m2m)"`
	Cpu				 string      `orm:"column(cpu);size(64);description(主机cpu);null"`
	Memory           string		 `orm:"column(memory);size(64);description(主机内存);null"`
	Disk             string      `orm:"column(disk);size(255);description(主机磁盘);null"`
}

func (c *Computer) TableName() string {
	return "asset_pc"
}

func (d *DeviceType) TableName() string {
	return "computer_device_type"
}

func (b *Brand) TableName() string {
	return "computer_brand"
}

func init() {
	orm.RegisterModel(new(Brand), new(DeviceType), new(Computer))
}
