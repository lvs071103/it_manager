package authentication

import (
	"crypto/tls"
	"dominos_cmdb/models"
	"dominos_cmdb/tools"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"gopkg.in/ldap.v2"
	"regexp"
	"strings"
	"time"
)

type LdapAuthentication struct {
	beego.Controller
}

func (l *LdapAuthentication) Login() {
	if l.Ctx.Input.IsGet() {
		l.TplName = "login/login.html"
	} else if l.Ctx.Input.IsPost() {
		// 用来认证的用户名和密码
		username := l.GetString("username")
		password := l.GetString("password")
		// 用来获取查询权限的 bind 用户。如果 ldap 禁止了匿名查询，那我们就需要先用这个帐户 bind 以下才能开始查询
		// bind 的账号通常要使用完整的 DN 信息。例如 cn=manager,dc=example,dc=org
		// 在 AD 上，则可以用诸如 mananger@example.org 的方式来 bind
		bindUsername := beego.AppConfig.String("ldap_user")
		bindPassword := beego.AppConfig.String("ldap_pass")
		resMsg := map[string]interface{}{}

		lc, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", "172.19.1.2", 389))
		if err != nil {
			//log.Fatal(err)
			fmt.Println(err)
		}
		defer lc.Close()

		// Reconnect with TLS
		err = lc.StartTLS(&tls.Config{InsecureSkipVerify: true})
		if err != nil {
			//log.Fatal(err)
			fmt.Println(err)
		}
		// First bind with a read only user
		// 先用我们的 bind 账号给 bind 上去
		err = lc.Bind(bindUsername, bindPassword)
		if err != nil {
			fmt.Println(err)
		}

		chars := []string{"]", "^", "\\\\", ".", "[", "(", ")", "-"}
		r := strings.Join(chars, "")
		re := regexp.MustCompile("[" + r + "]+")
		queryUsername := re.ReplaceAllString(username, "*")
		fmt.Println(queryUsername)

		// queryUserName := fmt.Sprintf("%s@dashbrands.local", username)

		fmt.Println(fmt.Sprintf("(&(objectClass=Person)(cn=%s*))", queryUsername))
		// Search for the given username
		searchRequest := ldap.NewSearchRequest("ou=Dominos_China_Users,dc=dashbrands,dc=local",
			ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
			fmt.Sprintf("(&(objectClass=organizationalPerson)(cn=%s*))", queryUsername),
			[]string{"dn"},
			nil,
		)

		sr, err := lc.Search(searchRequest)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(sr.Entries)

		if len(sr.Entries) != 1 {
			resMsg["code"] = 123
			resMsg["msg"] = "User does not exist or too many entries returned"
			l.Data["json"] = resMsg
		} else {
			// 如果没有意外，那么我们就可以获取用户的实际 DN 了
			userDN := sr.Entries[0].DN
			// Bind as the user to verify their password
			// 拿这个 dn 和他的密码去做 bind 验证
			err = lc.Bind(userDN, password)
			if err != nil {
				resMsg["code"] = 125
				resMsg["msg"] = err.Error()
				l.Data["json"] = resMsg
			} else {
				l.SetSession("username", username)
				l.SetSession("loginType", "ldap")
				CreateOrUpdate(username, password)
				UpdateLastLogin(username)
				resMsg["code"] = 0
				resMsg["msg"] = "登陆成功"
				l.Data["json"] = resMsg
			}
		}

		l.ServeJSON()
	}
}

func CreateOrUpdate(username, password string) {
	user := models.User{
		UserName:    username,
		Email:       username + "@dominos.com.cn",
		Password:    tools.GetMD5(password),
		IsActive:    1,
		IsStaff:     1,
		IsDelete:    0,
		IsSuperUser: 0,
	}
	o := orm.NewOrm()
	exists := o.QueryTable(new(models.User)).Filter("username", username).Exist()
	if !exists {
		_, _ = o.Insert(&user)
	} else {
		_, _ = o.QueryTable(new(models.User)).Filter("username", username).Update(orm.Params{
			"password": tools.GetMD5(password),
			"is_staff": 1,
		})
	}
}

func UpdateLastLogin(username string) bool {
	o := orm.NewOrm()
	_, err := o.QueryTable(new(models.User)).Filter("username", username).Update(orm.Params{
		"last_login": time.Now(),
	})

	if err != nil {
		return false
	} else {
		return true
	}
}
