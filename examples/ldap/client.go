// Package ldap provides a simple ldap client to authenticate,
// retrieve basic information and groups for a user.
package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"slices"

	"github.com/go-ldap/ldap/v3"
	"golang.org/x/text/encoding/unicode"
)

type LDAPClient struct {
	Attributes         []string
	Base               string
	BindDN             string
	BindPassword       string
	GroupFilter        string // e.g. "(memberUid=%s)"
	Addr               string
	ServerName         string
	UserFilter         string // e.g. "(uid=%s)"
	Conn               *ldap.Conn
	InsecureSkipVerify bool
	ClientCertificates []tls.Certificate // Adding client certificates
}

// Connect connects to the ldap backend.
func (lc *LDAPClient) Connect() error {
	if lc.Conn == nil {
		var l *ldap.Conn
		var err error

		config := &tls.Config{
			InsecureSkipVerify: lc.InsecureSkipVerify,
			ServerName:         lc.ServerName,
		}
		if len(lc.ClientCertificates) > 0 {
			config.Certificates = lc.ClientCertificates
		}

		l, err = ldap.DialURL(lc.Addr, ldap.DialWithTLSConfig(config))
		if err != nil {
			return err
		}

		lc.Conn = l
	}
	return nil
}

// Close closes the ldap backend connection.
func (lc *LDAPClient) Close() {
	if lc.Conn != nil {
		lc.Conn.Close()
		lc.Conn = nil
	}
}

// Authenticate authenticates the user against the ldap backend.
func (lc *LDAPClient) Authenticate(username, password string) (bool, map[string]string, error) {
	err := lc.Connect()
	if err != nil {
		return false, nil, err
	}

	entry, err := lc.findUser()
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

	// Bind as the user to verify their password
	err = lc.Conn.Bind(userDN, password)
	if err != nil {
		return false, user, err
	}

	// Rebind as the read only user for any further queries
	if lc.BindDN != "" && lc.BindPassword != "" {
		err = lc.Conn.Bind(lc.BindDN, lc.BindPassword)
		if err != nil {
			return true, user, err
		}
	}

	return true, user, nil
}

// GetGroupsOfUser returns the group for a user.
func (lc *LDAPClient) GetGroupsOfUser(username string) ([]string, error) {
	err := lc.Connect()
	if err != nil {
		return nil, err
	}

	searchRequest := ldap.NewSearchRequest(
		lc.Base,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(lc.GroupFilter, username),
		[]string{"cn"}, // can it be something else than "cn"?
		nil,
	)
	sr, err := lc.Conn.Search(searchRequest)
	if err != nil {
		return nil, err
	}
	groups := []string{}
	for _, entry := range sr.Entries {
		groups = append(groups, entry.GetAttributeValue("cn"))
	}
	return groups, nil
}

func (lc *LDAPClient) ChnagePassword(username, oldPwd, newPwd string) error {
	err := lc.Connect()
	if err != nil {
		return err
	}
	// First bind with a read only user
	if lc.BindDN != "" && lc.BindPassword != "" {
		err := lc.Conn.Bind(lc.BindDN, lc.BindPassword)
		if err != nil {
			return err
		}
	}

	entry, err := lc.findUser()
	if err != nil {
		return err
	}

	// Bind as the user to verify their password
	err = lc.Conn.Bind(entry.DN, oldPwd)
	if err != nil {
		return err
	}
	defer func() {
		// Rebind as the read only user for any further queries
		if lc.BindDN != "" && lc.BindPassword != "" {
			err = lc.Conn.Bind(lc.BindDN, lc.BindPassword)
			if err != nil {
				panic(err)
			}
		}
	}()

	// 修改密码
	utf16 := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	// Remember to add double quotation marks to your password string !!!!!
	pwdEncoded, _ := utf16.NewEncoder().String(fmt.Sprintf(`"%s"`, newPwd))
	pwdModifyReq := ldap.NewModifyRequest(entry.DN, nil)
	pwdModifyReq.Replace("unicodePwd", []string{pwdEncoded})
	err = lc.Conn.Modify(pwdModifyReq)
	return err
}

func (lc *LDAPClient) findUser() (*ldap.Entry, error) {
	err := lc.Connect()
	if err != nil {
		return nil, err
	}

	// First bind with a read only user
	if lc.BindDN != "" && lc.BindPassword != "" {
		err := lc.Conn.Bind(lc.BindDN, lc.BindPassword)
		if err != nil {
			log.Printf("error binding as read only user: %+v", err)
			return nil, err
		}
	}

	attributes := make([]string, 0, len(lc.Attributes))
	attributes = append(attributes, lc.Attributes...)
	if !slices.Contains(attributes, "*") {
		attributes = append(attributes, "dn")
	}
	// Search for the given username
	searchRequest := ldap.NewSearchRequest(
		lc.Base,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(lc.UserFilter, username),
		attributes,
		nil,
	)

	sr, err := lc.Conn.Search(searchRequest)
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
