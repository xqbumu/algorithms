package main

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/go-ldap/ldap/v3"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"golang.org/x/text/encoding/unicode"
)

// Server: https://github.com/glauth/glauth

// Example_userAuthentication shows how a typical application can verify a login attempt
// Refer to https://github.com/go-ldap/ldap/issues/93 for issues revolving around unauthenticated binds, with zero length passwords
func Test_userAuthentication(t *testing.T) {
	bindusername := "cn=serviceuser,ou=svcaccts,dc=glauth,dc=com"
	bindpassword := "mysecret"

	// The username and password we want to check
	username := "johndoe"
	password := "dogood"

	l, err := ldap.DialURL("ldap://ubuntu.orb.local:3893")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	// First bind with a read only user
	// Reconnect with TLS
	// err = l.StartTLS(&tls.Config{InsecureSkipVerify: true})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	err = l.Bind(bindusername, bindpassword)
	if err != nil {
		log.Fatal(err)
	}

	user := checkBind(l, username, password, "")
	// 修改密码
	utf16 := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	pwdEncoded, _ := utf16.NewEncoder().String("\"testpassword\"") //Remember to add double quotation marks to your password string !!!!!
	log.Println("password is :", pwdEncoded)
	passwordModify := ldap.NewModifyRequest(user.DN, nil)
	passwordModify.Replace("unicodePwd", []string{pwdEncoded})
	err = l.Modify(passwordModify)
	if err != nil {
		log.Printf("dn: %s, err: %v", user.DN, err)
	}

	// Rebind as the read only user for any further queries
	err = l.Bind(bindusername, bindpassword)
	if err != nil {
		log.Fatal(err)
	}
	// otpusername := "otpuser"
	// otppassword := "mysecret"
	// otpsecret := "3hnvnk4ycv44glzigd6s25j4dougs3rk"
	// checkBind(l, otpusername, otppassword, generatePassCode(otpsecret))
}

func checkBind(l *ldap.Conn, username, password, otp string) *ldap.Entry {
	// Search for the given username
	// ldapsearch -LLL \
	//   -H ldap://localhost:3893 \
	//   -D "cn=serviceuser,ou=svcaccts,dc=glauth,dc=com" \
	//   -w "mysecret" \
	//   -x -bdc=glauth,dc=com \
	//   "cn=johndoe"
	searchRequest := ldap.NewSearchRequest(
		"dc=glauth,dc=com",
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0,
		false,
		fmt.Sprintf("(uid=%s)", ldap.EscapeFilter(username)),
		// fmt.Sprintf("(&(objectClass=organizationalPerson)(uid=%s))", ldap.EscapeFilter(username)),
		[]string{"*"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Fatal(err)
	}

	if len(sr.Entries) != 1 {
		log.Fatal("User does not exist or too many entries returned")
	}

	user := sr.Entries[0]

	// Bind as the user to verify their password
	err = l.Bind(user.DN, password+otp)
	if err != nil {
		log.Fatalf("dc: %s, err: %s", username, err)
	}
	log.Printf("dc: %s, login done", username)

	return user
}

// Demo function, not used in main
// Generates Passcode using a UTF-8 (not base32) secret and custom parameters
func generatePassCode(secret string) string {
	// secret := base32.StdEncoding.EncodeToString([]byte(secret))
	passcode, err := totp.GenerateCodeCustom(secret, time.Now(), totp.ValidateOpts{
		Period:    30,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})
	if err != nil {
		panic(err)
	}
	return passcode
}
