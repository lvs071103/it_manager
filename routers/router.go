package routers

import (
	"dominos_cmdb/controllers/assets"
	"dominos_cmdb/controllers/auth"
	"dominos_cmdb/controllers/authentication"
	"dominos_cmdb/controllers/dashboard"
	"dominos_cmdb/controllers/echarts"
	"dominos_cmdb/controllers/messages"
	"dominos_cmdb/controllers/myError"
	"dominos_cmdb/controllers/my_workspaces"
	"dominos_cmdb/controllers/settings"
	"dominos_cmdb/controllers/tickets"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/login", &authentication.LoginController{}, "get:Login;post:Login")
	beego.Router("/ldap/login", &authentication.LdapAuthentication{}, "get:Login;post:Login")
	beego.Router("/logout", &authentication.LoginController{}, "get:Logout")
	beego.Router("/refreshCaptcha", &authentication.LoginController{}, "get:RefreshCaptcha")
	beego.Router("/api/home", &dashboard.HomeController{})
	beego.Router("/api/errors/403", &myError.ErrorController{}, "get:Error403")
	beego.Router("/api/loginWelcome", &dashboard.HomeController{}, "get:LoginWelcome")
	// group
	beego.Router("/api/group/list", &auth.GroupController{}, "get:List")
	beego.Router("/api/group/add", &auth.GroupController{}, "get:Add;post:Add")
	beego.Router("/api/group/edit", &auth.GroupController{}, "get:Edit;post:Edit")
	beego.Router("/api/group/isActive", &auth.GroupController{}, "post:IsActive")
	beego.Router("/api/group/delete", &auth.GroupController{}, "post:Delete")
	beego.Router("/api/group/batchDelete", &auth.GroupController{}, "post:BatchDelete")
	// user
	beego.Router("/api/user/list", &auth.UserController{}, "get:List")
	beego.Router("/api/user/json", &auth.UserController{}, "get:Json")
	beego.Router("/api/user/add", &auth.UserController{}, "get:Add;post:Add")
	beego.Router("/api/user/edit", &auth.UserController{}, "get:Edit;post:Edit")
	beego.Router("/api/user/resetPass", &auth.UserController{}, "post:ResetPass")
	beego.Router("/api/user/isActive", &auth.UserController{}, "post:IsActive")
	beego.Router("/api/user/delete", &auth.UserController{}, "post:Delete")
	beego.Router("/api/user/batchDelete", &auth.UserController{}, "post:BatchDelete")
	beego.Router("/api/user/getPerJson", &auth.UserController{},
		"get:GetPerJson;post:GetPerJson")
	beego.Router("/api/user/bindPer", &auth.UserController{}, "get:BindPer;post:BindPer")
	// icon
	beego.Router("/api/settings/icon/list", &settings.IconController{}, "get:List")
	beego.Router("/api/settings/icon/add", &settings.IconController{}, "get:Add;post:Add")
	beego.Router("/api/settings/icon/edit", &settings.IconController{}, "get:Edit;post:Edit")
	beego.Router("/api/settings/icon/del", &settings.IconController{}, "post:Delete")
	beego.Router("/api/settings/icon/isActive", &settings.IconController{}, "post:IsActive")
	beego.Router("/api/settings/icon/batchDel", &settings.IconController{}, "post:BatchDelete")
	// menu
	beego.Router("/api/settings/menu/list", &settings.MenuController{}, "get:List")
	beego.Router("/api/settings/menu/add", &settings.MenuController{}, "get:Add;post:Add")
	beego.Router("/api/settings/menu/edit", &settings.MenuController{}, "get:Edit;post:Edit")
	beego.Router("/api/settings/menu/delete", &settings.MenuController{}, "post:Delete")
	beego.Router("/api/settings/menu/isActive", &settings.MenuController{}, "post:IsActive")
	beego.Router("/api/settings/menu/batchDel", &settings.MenuController{}, "post:BatchDelete")
	// brand
	beego.Router("/api/assets/brand/list", &assets.BrandController{}, "get:List")
	beego.Router("/api/assets/brand/json", &assets.BrandController{}, "get:Json")
	beego.Router("/api/assets/brand/add", &assets.BrandController{}, "get:Add;post:Add")
	beego.Router("/api/assets/brand/edit", &assets.BrandController{}, "get:Edit;post:Edit")
	beego.Router("/api/assets/brand/isActive", &assets.BrandController{}, "post:IsActive")
	beego.Router("/api/assets/brand/delete", &assets.BrandController{}, "post:Delete")
	beego.Router("/api/assets/brand/batchDelete", &assets.BrandController{}, "post:BatchDelete")
	// deviceType
	beego.Router("/api/assets/device_type/list", &assets.DTypeController{}, "get:List")
	beego.Router("/api/assets/device_type/json", &assets.DTypeController{}, "get:Json")
	beego.Router("/api/assets/device_type/add", &assets.DTypeController{}, "get:Add;post:Add")
	beego.Router("/api/assets/device_type/edit", &assets.DTypeController{}, "get:Edit;post:Edit")
	beego.Router("/api/assets/device_type/delete", &assets.DTypeController{}, "post:Delete")
	beego.Router("/api/assets/device_type/batchDel", &assets.DTypeController{}, "post:BatchDelete")
	beego.Router("/api/assets/device_type/isActive", &assets.DTypeController{}, "post:IsActive")
	// pc
	beego.Router("/api/assets/pc/list", &assets.PCController{}, "get:List")
	beego.Router("/api/assets/pc/json", &assets.PCController{}, "get:Json")
	beego.Router("/api/assets/pc/add", &assets.PCController{}, "get:Add;post:Add")
	beego.Router("/api/assets/pc/edit", &assets.PCController{}, "get:Edit;post:Edit")
	beego.Router("/api/assets/pc/delete", &assets.PCController{}, "post:Delete")
	beego.Router("/api/assets/pc/broken", &assets.PCController{}, "post:Broken")
	beego.Router("/api/assets/pc/batchDel", &assets.PCController{}, "post:BatchDelete")
	beego.Router("/api/assets/pc/isActive", &assets.PCController{}, "post:IsActive")
	beego.Router("/api/assets/pc/bindUser", &assets.PCController{}, "get:BindUser;post:BindUser")
	// tickets
	beego.Router("/api/tickets/category/list", &tickets.TCController{}, "get:List")
	beego.Router("/api/tickets/category/add", &tickets.TCController{}, "get:Add;post:Add")
	beego.Router("/api/tickets/category/edit", &tickets.TCController{}, "get:Edit;post:Edit")
	beego.Router("/api/tickets/category/isActive", &tickets.TCController{}, "post:IsActive")
	beego.Router("/api/tickets/category/delete", &tickets.TCController{}, "post:Del")
	beego.Router("/api/tickets/category/batch_delete", &tickets.TCController{}, "post:BatchDelete")
	beego.Router("/api/tickets/list", &tickets.TicketController{}, "get:List")
	beego.Router("/api/tickets/add", &tickets.TicketController{}, "get:Add;post:Add")
	beego.Router("/api/tickets/edit", &tickets.TicketController{}, "get:Edit;post:Edit")
	beego.Router("/api/tickets/isActive", &tickets.TicketController{}, "post:IsActive")
	beego.Router("/api/tickets/delete", &tickets.TicketController{}, "post:Del")
	beego.Router("/api/tickets/batch_delete", &tickets.TicketController{}, "post:BatchDelete")
	beego.Router("/api/tickets/assign", &tickets.TicketController{}, "get:Assign;post:Assign")
	// my profile
	beego.Router("/api/profile/my/edit", &my_workspaces.ProfileController{}, "get:Edit;post:Edit")
	beego.Router("/api/tickets/my/ticket_list", &tickets.TicketController{}, "get:MyTickets")
	beego.Router("/api/tickets/my/ticket_detail", &tickets.TicketController{}, "get:Detail")
	beego.Router("/api/tickets/my/preprocess", &my_workspaces.MyTCController{}, "get:Preprocess")
	beego.Router("/api/assets/my/computer", &assets.PCController{}, "get:MyComputer;post:MyComputer")
	// messages
	beego.Router("/api/messages/list", &messages.MessageController{}, "get:List")
	beego.Router("/api/messages/isRead", &messages.MessageController{}, "get:IsRead;post:IsRead")
	beego.Router("/api/messages/delete", &messages.MessageController{}, "post:Del")
	beego.Router("/api/messages/batchDelete", &messages.MessageController{}, "post:BatchDel")
	beego.Router("/api/messages/setRead", &messages.MessageController{}, "post:SetRead")
	beego.Router("/api/messages/my_list", &messages.MessageController{}, "get:MyList")
	// e_chart
	beego.Router("/api/tickets/e_chars", &echarts.EchController{}, "get:TicketsChart")
}
