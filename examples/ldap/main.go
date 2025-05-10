package main

import (
	"flag"
	"log"
	"strings"
)

var addr, baseDN, bindDN, bindPassword, groupFilter, password, newPassword, userFilter, username, certFile, clientCert string
var list bool
var skipTLSVerify bool
var attrs []string

func main() {
	flag.Parse()

	client := &LDAPClient{
		Addr:               addr,
		CertFile:           certFile,
		BaseDN:             baseDN,
		BindDN:             bindDN,
		BindPassword:       bindPassword,
		UserFilter:         userFilter,
		GroupFilter:        groupFilter,
		Attributes:         attrs,
		InsecureSkipVerify: skipTLSVerify,
	}
	if err := client.Init(); err != nil {
		panic(err)
	}

	if list {
		entries, err := client.ListUsers()
		if err != nil {
			log.Fatalf("Error listing users")
		}
		for _, entry := range entries {
			log.Println("========")
			for _, attr := range entry.Attributes {
				log.Printf("Entry: %+v -  %+v\n", attr.Name, attr.Values)
			}
		}
		return
	}

	if newPassword != "" {
		err := client.SetPassword(username, password, newPassword)
		if err != nil {
			log.Fatalf("Error changing password for user %s: %+v", username, err)
		}
		return
	}

	ok, user, err := client.Authenticate(username, password)
	if err != nil {
		log.Fatalf("Error authenticating user %s: %+v", username, err)
	}
	if !ok {
		log.Fatalf("Authenticating failed for user %s", username)
	}
	log.Printf("User: %+v", user)

	groups, err := client.GetGroupsOfUser(username)
	if err != nil {
		log.Fatalf("Error getting groups for user %s: %+v", username, err)
	}
	log.Printf("Groups: %+v", groups)
}

func init() {
	flag.StringVar(&baseDN, "base-dn", "dc=demo,dc=dev", "BaseDN LDAP")
	flag.StringVar(&bindDN, "bind-dn", "uid=readonlysuer,ou=People,dc=demo,dc=dev", "Bind DN")
	flag.StringVar(&bindPassword, "bind-pwd", "readonlypassword", "Bind password")
	flag.StringVar(&groupFilter, "group-filter", "(memberUid=%s)", "Group filter")
	flag.StringVar(&addr, "addr", "ldap://demo.dev:389", "LDAP addr")
	flag.BoolVar(&list, "list", false, "List users")
	flag.StringVar(&userFilter, "user-filter", "(uid=%s)", "User filter")
	flag.StringVar(&username, "username", "", "Username")
	flag.StringVar(&password, "password", "", "Password")
	flag.StringVar(&newPassword, "new-pwd", "", "New password")
	flag.StringVar(&certFile, "cert-file", "", "root cert file")
	flag.StringVar(&clientCert, "client-cert", "", "client cert file")
	flag.BoolVar(&skipTLSVerify, "skip-tls-verify", false, "Skip TLS verify")
	flag.Func("attrs", "Comma-separated list of attributes", func(value string) error {
		if len(strings.TrimSpace(value)) == 0 {
			attrs = []string{"givenName", "sn", "mail", "uid", "accountExpires", "userPrincipalName"}
		} else {
			attrs = strings.Split(value, ",")
		}
		return nil
	})
}
