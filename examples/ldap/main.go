package main

import (
	"crypto/x509"
	"flag"
	"log"
	"os"
)

var addr, baseDN, bindDN, bindPassword, groupFilter, password, newPassword, userFilter, username, certFile, clientCert string
var skipTLSVerify bool

func main() {
	flag.Parse()

	// 创建证书池
	certPool := x509.NewCertPool()

	// 读取证书
	if len(certFile) > 0 {
		certPEM, err := os.ReadFile(certFile)
		if err != nil {
			log.Fatal(err)
		}
		// 将证书添加到池中
		if ok := certPool.AppendCertsFromPEM(certPEM); !ok {
			log.Fatal("Failed to append certificate")
		}
	}

	client := &LDAPClient{
		Addr:               addr,
		BaseDN:             baseDN,
		BindDN:             bindDN,
		BindPassword:       bindPassword,
		UserFilter:         userFilter,
		GroupFilter:        groupFilter,
		Attributes:         []string{"givenName", "sn", "mail", "uid"},
		InsecureSkipVerify: skipTLSVerify,
		CertPool:           certPool,
		ClientCertificates: nil,
	}

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
	flag.StringVar(&baseDN, "base-dn", "dc=demo,dc=dev", "BaseDN LDAP")
	flag.StringVar(&bindDN, "bind-dn", "uid=readonlysuer,ou=People,dc=demo,dc=dev", "Bind DN")
	flag.StringVar(&bindPassword, "bind-pwd", "readonlypassword", "Bind password")
	flag.StringVar(&groupFilter, "group-filter", "(memberUid=%s)", "Group filter")
	flag.StringVar(&addr, "addr", "ldap://demo.dev:389", "LDAP addr")
	flag.StringVar(&userFilter, "user-filter", "(uid=%s)", "User filter")
	flag.StringVar(&username, "username", "", "Username")
	flag.StringVar(&password, "password", "", "Password")
	flag.StringVar(&newPassword, "new-pwd", "", "New password")
	flag.StringVar(&certFile, "cert-file", "", "root cert file")
	flag.StringVar(&clientCert, "client-cert", "", "client cert file")
	flag.BoolVar(&skipTLSVerify, "skip-tls-verify", false, "Skip TLS verify")
}
