// Package ldap provides a simple ldap client to authenticate,
// retrieve basic information and groups for a user.
package main

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"log"
	"slices"

	"github.com/go-ldap/ldap/v3"
	"golang.org/x/text/encoding/unicode"
)

type LDAPClient struct {
	Addr               string            `yaml:"addr" json:"addr,omitempty"`
	BaseDN             string            `yaml:"baseDn" json:"base_dn,omitempty"`
	BindDN             string            `yaml:"bindDn" json:"bind_dn,omitempty"`
	BindPassword       string            `yaml:"bindPassword" json:"bind_password,omitempty"`
	UserFilter         string            `yaml:"userFilter" json:"user_filter,omitempty"`   // e.g. "(uid=%s)"
	GroupFilter        string            `yaml:"groupFilter" json:"group_filter,omitempty"` // e.g. "(memberUid=%s)"
	Attributes         []string          `yaml:"attributes" json:"attributes,omitempty"`
	InsecureSkipVerify bool              `yaml:"insecureSkipVerify" json:"insecure_skip_verify,omitempty"`
	CertPool           *x509.CertPool    `yaml:"certPool" json:"cert_pool,omitempty"`
	ClientCertificates []tls.Certificate `yaml:"clientCertificates" json:"client_certificates,omitempty"` // Adding client certificates
}

// Connect connects to the ldap backend.
func (lc *LDAPClient) Connect(bind bool) (*ldap.Conn, error) {
	var l *ldap.Conn
	var err error

	config := &tls.Config{
		InsecureSkipVerify: lc.InsecureSkipVerify,
	}
	if lc.CertPool != nil {
		config.RootCAs = lc.CertPool
	}
	if len(lc.ClientCertificates) > 0 {
		config.Certificates = lc.ClientCertificates
	}

	l, err = ldap.DialURL(lc.Addr, ldap.DialWithTLSConfig(config))
	if err != nil {
		return nil, err
	}

	if bind {
		// First bind with a read only user
		if lc.BindDN != "" && lc.BindPassword != "" {
			err := l.Bind(lc.BindDN, lc.BindPassword)
			if err != nil {
				log.Printf("error binding as read only user: %+v", err)
				return nil, err
			}
		}
		// check server type
		searchRequest := ldap.NewSearchRequest(
			lc.BindDN,
			ldap.ScopeBaseObject, ldap.NeverDerefAliases, 0, 0, false,
			"(objectClass=*)",
			[]string{"*"},
			nil,
		)
		result, err := l.Search(searchRequest)
		log.Println(result, err)
	}

	return l, nil
}

// Authenticate authenticates the user against the ldap backend.
func (lc *LDAPClient) Authenticate(username, password string) (bool, map[string]string, error) {
	entry, err := lc.findUser(username)
	if err != nil {
		return false, nil, err
	}

	userDN := entry.DN
	user := map[string]string{}
	if slices.Contains(lc.Attributes, "*") {
		for _, attr := range entry.Attributes {
			user[attr.Name] = entry.GetAttributeValue(attr.Name)
		}
	} else {
		for _, attr := range lc.Attributes {
			user[attr] = entry.GetAttributeValue(attr)
		}
	}

	l, err := lc.Connect(false)
	if err != nil {
		return false, nil, err
	}
	defer l.Close()

	// Bind as the user to verify their password
	err = l.Bind(userDN, password)
	if err != nil {
		return false, user, err
	}

	return true, user, nil
}

// GetGroupsOfUser returns the group for a user.
func (lc *LDAPClient) GetGroupsOfUser(username string) ([]string, error) {
	l, err := lc.Connect(true)
	if err != nil {
		return nil, err
	}
	defer l.Close()

	searchRequest := ldap.NewSearchRequest(
		lc.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(lc.GroupFilter, username),
		[]string{"cn"}, // can it be something else than "cn"?
		nil,
	)
	sr, err := l.Search(searchRequest)
	if err != nil {
		return nil, err
	}
	groups := []string{}
	for _, entry := range sr.Entries {
		groups = append(groups, entry.GetAttributeValue("cn"))
	}
	return groups, nil
}

// 修改密码
func (lc *LDAPClient) ChangePassword(username, oldPwd, newPwd string) error {
	entry, err := lc.findUser(username)
	if err != nil {
		return err
	}

	l, err := lc.Connect(false)
	if err != nil {
		return err
	}
	defer l.Close()

	// Bind as the user to verify their password
	err = l.Bind(entry.DN, oldPwd)
	if err != nil {
		return err
	}

	err = ldapChangeUserPassword(l, entry.DN, oldPwd, newPwd)

	return err
}

// 设置密码
func (lc *LDAPClient) SetPassword(username, oldPwd, newPwd string) error {
	// 查找用户
	entry, err := lc.findUser(username)
	if err != nil {
		return err
	}

	// 创建连接
	l, err := lc.Connect(false)
	if err != nil {
		return err
	}
	defer l.Close()

	// 检查密码
	err = l.Bind(entry.DN, oldPwd)
	if err != nil {
		return err
	}

	// bind with a read only user
	if lc.BindDN != "" && lc.BindPassword != "" {
		err := l.Bind(lc.BindDN, lc.BindPassword)
		if err != nil {
			log.Printf("error binding as read only user: %+v", err)
			return err
		}
	}

	err = ldapSetUserPassword(l, entry.DN, newPwd)
	if err != nil {
		return err
	}

	return err
}

// 列出全部用户
func (lc *LDAPClient) ListUsers(filterOptions ...LDAPSearchRequestOption) ([]*ldap.Entry, error) {
	l, err := lc.Connect(true)
	if err != nil {
		return nil, err
	}
	defer l.Close()

	attributes := make([]string, 0, len(lc.Attributes))
	attributes = append(attributes, lc.Attributes...)
	if !slices.Contains(attributes, "*") {
		attributes = append(attributes, "dn")
	}
	// Search for the given username
	searchRequest := ldap.NewSearchRequest(
		lc.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(lc.UserFilter, "*"),
		attributes,
		nil,
	)
	if len(filterOptions) > 0 {
		for _, option := range filterOptions {
			option(searchRequest)
		}
	}

	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Printf("error listing users: %+v", err)
		return nil, err
	}

	return sr.Entries, nil
}

type LDAPSearchRequestOption func(searchRequest *ldap.SearchRequest)

func (lc *LDAPClient) findUser(username string, filterOptions ...LDAPSearchRequestOption) (*ldap.Entry, error) {
	l, err := lc.Connect(true)
	if err != nil {
		return nil, err
	}
	defer l.Close()

	attributes := make([]string, 0, len(lc.Attributes))
	attributes = append(attributes, lc.Attributes...)
	if !slices.Contains(attributes, "*") {
		attributes = append(attributes, "dn")
	}
	// Search for the given username
	searchRequest := ldap.NewSearchRequest(
		lc.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(lc.UserFilter, username),
		attributes,
		nil,
	)
	if len(filterOptions) > 0 {
		for _, option := range filterOptions {
			option(searchRequest)
		}
	}

	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Printf("error searching for user %s: %+v", username, err)
		return nil, err
	}

	if len(sr.Entries) < 1 {
		return nil, errors.New("user does not exist")
	}

	if len(sr.Entries) > 1 {
		return nil, errors.New("too many entries returned")
	}

	return sr.Entries[0], nil
}

// EncodePwd string to unicode string
func encodePassword(pwd string) string {
	utf16 := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM) // 使用小端编码
	pwdEncoded, err := utf16.NewEncoder().String("\"" + pwd + "\"")
	if err != nil {
		return ""
	}
	return pwdEncoded
}

func ldapChangeUserPassword(ldapConn *ldap.Conn, userDN, oldPassword, newPassword string) error {
	oldPasswordEncoded := encodePassword(oldPassword)
	newPasswordEncoded := encodePassword(newPassword)

	modifyRequest := ldap.NewModifyRequest(userDN, nil)
	modifyRequest.Delete("unicodePwd", []string{oldPasswordEncoded})
	modifyRequest.Add("unicodePwd", []string{newPasswordEncoded})

	log.Printf("oldPasswordEncoded: %s, newPasswordEncoded: %s", oldPasswordEncoded, newPasswordEncoded)

	if err := ldapConn.Modify(modifyRequest); err != nil {
		return fmt.Errorf("password change failed: %v", err)
	}

	fmt.Println("Password successfully changed!")
	return nil
}

func ldapSetUserPassword(ldapConn *ldap.Conn, userDN, password string) error {
	passwordEncoded := encodePassword(password)

	// Create modify request
	modifyRequest := ldap.NewModifyRequest(userDN, nil)
	modifyRequest.Replace("unicodePwd", []string{string(passwordEncoded)})

	if err := ldapConn.Modify(modifyRequest); err != nil {
		return fmt.Errorf("password set failed: %v", err)
	}

	fmt.Println("Password successfully set!")
	return nil
}
