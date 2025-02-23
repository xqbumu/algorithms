// Package ldap provides a simple ldap client to authenticate,
// retrieve basic information and groups for a user.
package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"slices"

	"github.com/go-ldap/ldap/v3"
)

type LDAPClient struct {
	Addr               string
	BaseDN             string
	BindDN             string
	BindPassword       string
	UserFilter         string // e.g. "(uid=%s)"
	GroupFilter        string // e.g. "(memberUid=%s)"
	Attributes         []string
	InsecureSkipVerify bool
	CertPool           *x509.CertPool
	ClientCertificates []tls.Certificate // Adding client certificates
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
	}

	return l, nil
}

// Authenticate authenticates the user against the ldap backend.
func (lc *LDAPClient) Authenticate(username, password string) (bool, map[string]string, error) {
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
func (lc *LDAPClient) ChnagePassword(username, oldPwd, newPwd string) error {
	entry, err := lc.findUser()
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

	err = ChangeUserPassword(l, entry.DN, oldPwd, newPwd)

	return err
}

func (lc *LDAPClient) findUser() (*ldap.Entry, error) {
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

func encodePassword(pwd string) string {
	pwdQuoted := fmt.Sprintf("\"%s\"", pwd) // Surround with quotes

	utf16 := []uint16{}
	for _, r := range pwdQuoted {
		utf16 = append(utf16, uint16(r))
	}
	buf := make([]byte, len(utf16)*2)
	for i, v := range utf16 {
		binary.LittleEndian.PutUint16(buf[i*2:], v)
	}

	return string(buf)
	// return base64.StdEncoding.EncodeToString(buf) // 将 encodePwd 进行 Base64 编码
}

func ChangeUserPassword(ldapConn *ldap.Conn, userDN, oldPassword, newPassword string) error {
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

func SetUserPassword(ldapConn *ldap.Conn, userDN, password string) error {
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
