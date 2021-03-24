package dashboard

import (
	"dominos_cmdb/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type HomeController struct {
	beego.Controller
}

func (h *HomeController) Get() {
	// 后端主页
	o := orm.NewOrm()
	// 根据session获得uid
	userName := h.GetSession("username")
	if userName == nil {
		h.Ctx.Redirect(302, "/login")
	}
	loginType := h.GetSession("loginType").(string)
	if loginType == "ldap" {
		user := models.User{UserName: userName.(string)}
		o := orm.NewOrm()
		_ = o.QueryTable(new(models.User)).Filter("username", userName).RelatedSel().One(&user)
		messageCount, _ := o.QueryTable(new(models.Message)).Filter("uid", user.Id).Filter("is_read", 0).Count()

		qs := o.QueryTable(new(models.Menu))
		var menuList []models.Menu
		// 查询当前用户的可查看的一级菜单列表（Filter(id__in)）降序
		_, _ = qs.Filter("pid", 0).Filter("is_delete", 0).OrderBy("-weight").All(&menuList)

		var trees []models.MenuTree
		for _, perMenu := range menuList {
			pid := perMenu.Id
			icons := models.Icon{}
			_ = o.QueryTable(new(models.Icon)).Filter("id", perMenu.Icons.Id).One(&icons)
			menuTreeData := models.MenuTree{
				Id: perMenu.Id, Name: perMenu.Name, Icon: icons.ClassName, UrlFor: perMenu.UrlFor, Weight: perMenu.Weight,
				Children: []*models.MenuTree{},
			}
			GetChildrenNode(pid, &menuTreeData)
			trees = append(trees, menuTreeData)
		}

		h.Data["menu_trees"] = trees
		h.Data["user"] = user
		h.Data["messageCount"] = messageCount
		h.Data["currentUser"] = userName
	} else {
		user := models.User{Id: h.GetSession("id").(int)}
		_ = o.Read(&user)
		_, _ = o.LoadRelated(&user, "Permissions")
		messageCount, _ := o.QueryTable(new(models.Message)).Filter("uid", h.GetSession("id").(int)).
			Filter("is_read", 0).Count()
		qs := o.QueryTable(new(models.Menu))
		var menuList []models.Menu
		// 查询当前用户的可查看的一级菜单列表（Filter(id__in)）降序
		_, _ = qs.Filter("pid", 0).Filter("is_delete", 0).OrderBy("-weight").All(&menuList)

		var trees []models.MenuTree
		for _, perMenu := range menuList {
			pid := perMenu.Id
			icons := models.Icon{}
			_ = o.QueryTable(new(models.Icon)).Filter("id", perMenu.Icons.Id).One(&icons)
			menuTreeData := models.MenuTree{
				Id: perMenu.Id, Name: perMenu.Name, Icon: icons.ClassName, UrlFor: perMenu.UrlFor, Weight: perMenu.Weight,
				Children: []*models.MenuTree{},
			}
			GetChildrenNode(pid, &menuTreeData)
			trees = append(trees, menuTreeData)
		}

		h.Data["menu_trees"] = trees
		h.Data["user"] = user
		h.Data["messageCount"] = messageCount
		h.Data["currentUser"] = userName
	}

	h.TplName = "index.html"
}

func GetChildrenNode(pid int, treeNode *models.MenuTree) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(models.Menu))
	var menuList []models.Menu
	_, err := qs.Filter("pid", pid).Filter("is_delete", 0).OrderBy("-weight").All(&menuList)
	if err != nil { // 有错误返回
		return
	}
	// 查询三级及以上数据
	for i := 0; i < len(menuList); i++ {
		pid := menuList[i].Id // 根据pid获取子节点
		icons := models.Icon{}
		_ = o.QueryTable(new(models.Icon)).Filter("id", menuList[i].Icons.Id).One(&icons)
		menuTreeData := models.MenuTree{Id: menuList[i].Id, Name: menuList[i].Name, Icon: icons.ClassName,
			UrlFor: menuList[i].UrlFor, Weight: menuList[i].Weight, Children: []*models.MenuTree{}}
		treeNode.Children = append(treeNode.Children, &menuTreeData)
		GetChildrenNode(pid, &menuTreeData)
	}
	return
}

func (h *HomeController) LoginWelcome() {
	h.TplName = "dashboard/welcome.html"
}

func (h *HomeController) NullMenu() string {
	return "#"
}
