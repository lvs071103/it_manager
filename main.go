package main

import (
	"dominos_cmdb/controllers/dashboard"
	"dominos_cmdb/controllers/myError"
	"dominos_cmdb/models"
	_ "dominos_cmdb/models"
	_ "dominos_cmdb/routers"
	"dominos_cmdb/tools"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/astaxie/beego/session/mysql"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

func AutoInsertContentType(label, model string) {
	contentType := models.ContentType{AppLabel: label, Model: model}
	o := orm.NewOrm()
	if !o.QueryTable(new(models.ContentType)).Filter("model", model).Exist() {
		_, err := o.Insert(&contentType)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func AutoInsertPermission(modelName string) {
	listMp := []map[string]string{
		{"name": "Can Add " + strings.Title(modelName), "code_name": "add_" + modelName},
		{"name": "Can Change " + strings.Title(modelName), "code_name": "change_" + modelName},
		{"name": "Can Delete " + strings.Title(modelName), "code_name": "delete_" + modelName},
		{"name": "Can View " + strings.Title(modelName), "code_name": "view_" + modelName},
	}
	o := orm.NewOrm()
	contentType := models.ContentType{Model: modelName}
	_ = o.QueryTable(new(models.ContentType)).Filter("model", modelName).One(&contentType)
	for item := range listMp {
		permission := models.Permission{
			Name:        listMp[item]["name"],
			ContentType: &contentType,
			CodeName:    listMp[item]["code_name"],
		}
		if !o.QueryTable(new(models.Permission)).Filter("code_name", listMp[item]["code_name"]).Exist() {
			_, err := o.Insert(&permission)
			if err != nil {
				fmt.Println(err.Error())
			}
		}

	}
}

func init() {
	username := beego.AppConfig.String("username")
	password := beego.AppConfig.String("password")
	host := beego.AppConfig.String("host")
	port := beego.AppConfig.String("port")
	db := beego.AppConfig.String("db")
	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.BConfig.WebConfig.Session.SessionProvider = "mysql"
	beego.BConfig.WebConfig.Session.SessionProviderConfig =
		username + ":" + password + "@tcp(" + host + ":" + port + ")/" + db + "?charset=utf8&loc=Local"
	beego.BConfig.WebConfig.Session.SessionGCMaxLifetime = 36000
	beego.Router("/", &dashboard.HomeController{})
}

func DataInit() {
	label := "auth"
	modelStr := beego.AppConfig.String("models")
	modelsArr := strings.Split(modelStr, " ")
	for item := range modelsArr {
		AutoInsertContentType(label, modelsArr[item])
		AutoInsertPermission(modelsArr[item])
	}

}

func init() {
	username := beego.AppConfig.String("username")
	password := beego.AppConfig.String("password")
	host := beego.AppConfig.String("host")
	port := beego.AppConfig.String("port")
	db := beego.AppConfig.String("db")
	dataSource := username + ":" + password + "@tcp(" + host + ":" + port + ")/" + db + "?charset=utf8&loc=Local"
	_ = orm.RegisterDriver("mysql", orm.DRMySQL)
	_ = orm.RegisterDataBase("default", "mysql", dataSource, 30)
}

func main() {
	orm.RunCommand()
	beego.InsertFilter("/api/*", beego.BeforeRouter, tools.LoginFilter)
	beego.ErrorController(&myError.ErrorController{})
	//orm.Debug = true
	DataInit()
	beego.Run()
}
