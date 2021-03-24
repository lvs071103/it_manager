package main

import (
	"fmt"
	"github.com/jtblin/go-ldap-client"
	"github.com/korylprince/go-ad-auth/v3"
	"log"
)

func LDAPClientAuthenticate() {
	client := &ldap.LDAPClient{
		Base:         "dc=dashbrands, dc=local",
		Host:         "172.19.1.2",
		Port:         389,
		UseSSL:       false,
		BindDN:       "ou=Office_Shanghai, ou=Dominos_China_Users, dc=dashbrands, dc=local",
		BindPassword: "Jack@07110398186",
		UserFilter:   "(uid=%s)",
		GroupFilter:  "(memberUid=%s)",
		Attributes:   []string{"givenName", "sn", "mail", "uid"},
	}
	defer client.Close()

	ok, user, err := client.Authenticate("furong.zhou", "Jack@07110398186")
	if err != nil {
		fmt.Println(err.Error())
		//fmt.Printf("Error authenticating user %s: %+v\n", "furong.zhou", err.Error())
	}
	if !ok {
		fmt.Println(ok)
		//fmt.Printf("Authenticating failed for user %s\n", "furong.zhou")
	}
	fmt.Printf("User: %+v\n", user)

}

// 显示如何检索用户组
func GetGroupsOfUser() {
	client := &ldap.LDAPClient{
		Base:        "dc=example,dc=com",
		Host:        "172.19.1.2",
		Port:        389,
		GroupFilter: "(memberUid=%s)",
	}
	defer client.Close()
	groups, err := client.GetGroupsOfUser("username")
	if err != nil {
		log.Fatalf("Error getting groups for user %s: %+v", "username", err)
	}
	log.Printf("Groups: %+v", groups)
}

func GoAdAuth() {
	config := &auth.Config{
		Server:   "172.19.1.2",
		Port:     389,
		BaseDN:   "OU=Dominos_China_Users,DC=dashbrands,DC=local",
		Security: auth.SecurityStartTLS,
	}

	username := "furong.zhou@dashbrands.local"
	password := "Jack@07110398186"

	status, err := auth.Authenticate(config, username, password)

	if err != nil {
		//handle err
		fmt.Println(err)
	}

	if !status {
		//handle failed authentication
		fmt.Println(status)
	}
}

func main() {
	GoAdAuth()
}
