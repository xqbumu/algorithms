package main

import (
	"flag"
	"log"
)

var base, bindDN, bindPassword, groupFilter, host, password, newPassword, serverName, userFilter, username string
var port int
var useSSL bool
var skipTLS, skipTLSVerify bool

type server struct{}

func main() {
	flag.Parse()

	client := &LDAPClient{
		Base:               base,
		Host:               host,
		Port:               port,
		UseSSL:             useSSL,
		SkipTLS:            skipTLS,
		BindDN:             bindDN,
		BindPassword:       bindPassword,
		UserFilter:         userFilter,
		GroupFilter:        groupFilter,
		Attributes:         []string{"givenName", "sn", "mail", "uid"},
		ServerName:         serverName,
		InsecureSkipVerify: skipTLSVerify,
	}
	defer client.Close()

	if newPassword != "" {
		err := client.ChnagePassword(username, password, newPassword)
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
	flag.StringVar(&base, "base", "dc=demo,dc=dev", "Base LDAP")
	flag.StringVar(&bindDN, "bind-dn", "uid=readonlysuer,ou=People,dc=demo,dc=dev", "Bind DN")
	flag.StringVar(&bindPassword, "bind-pwd", "readonlypassword", "Bind password")
	flag.StringVar(&groupFilter, "group-filter", "(memberUid=%s)", "Group filter")
	flag.StringVar(&host, "host", "ldap.demo.dev", "LDAP host")
	flag.StringVar(&password, "password", "", "Password")
	flag.StringVar(&newPassword, "new-pwd", "", "Password")
	flag.IntVar(&port, "port", 389, "LDAP port")
	flag.StringVar(&userFilter, "user-filter", "(uid=%s)", "User filter")
	flag.StringVar(&username, "username", "", "Username")
	flag.StringVar(&serverName, "server-name", "", "Server name for SSL (if use-ssl is set)")
	flag.BoolVar(&useSSL, "use-ssl", false, "Use SSL")
	flag.BoolVar(&skipTLS, "skip-tls", false, "Skip TLS start")
	flag.BoolVar(&skipTLSVerify, "skip-tls-verify", false, "Skip TLS verify")
}
