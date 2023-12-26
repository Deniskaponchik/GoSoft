package main

import (
	"fmt"
	"github.com/go-ldap/ldap"
	"log"
	"os"
	//"github.com/go-ldap/ldap/v3"
)

const (
// BindUsername = "domain\\login"
// BindPassword = "password"
// FQDN = "DC.example.com"
// BaseDN = "cn=Family Name,OU=PointOut,dc=example,dc=com"
// Filter = "(objectClass=*)"
// Filter = "(objectClass=user)"
// Filter = "(SamAccountName=*)"
// Filter = "(CN=Family Name)"
// Filter = "(&(objectClass=organizationalPerson))"
// Filter = "(&(objectClass=user)(cn=Family Name))"
// Filter = "(&(objectClass=user)(sAMAccountName=login))"
// Filter = "(&(sAMAccountName=login))"
)

func main() {
	//TLS Connection
	l, err := ConnectTLS()
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	//Normal Bind and Search
	result, err := BindAndSearch(l)
	if err != nil {
		log.Fatal(err)
	}

	/* Non-TLS Connection
	l, err := Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	// Anonymous Bind and Search
	result, err := AnonymousBindAndSearch(l)
	if err != nil {
		log.Fatal(err)
	}*/

	fmt.Println("Вывод результата:")
	result.Entries[0].Print()

}

// Anonymous Bind and Search
func AnonymousBindAndSearch(l *ldap.Conn) (*ldap.SearchResult, error) {
	//err := l.UnauthenticatedBind("") //НЕ РАБОТАЕТ без аутентификации
	//err := l.Bind(BindUsername, BindPassword)
	err := l.Bind(os.Args[1], os.Args[2])
	if err != nil {
		log.Println("Ошибка аутентификации")
		return nil, err
	}

	anonReq := ldap.NewSearchRequest(
		//"",
		//BaseDN,
		os.Args[4],
		//ldap.ScopeBaseObject, // you can also use ldap.ScopeWholeSubtree
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		//Filter,
		//fmt.Sprintf("(&(objectClass=organizationalPerson)(uid=%s))", "denis.tirskikh"),
		fmt.Sprintf("(&(sAMAccountName=%s))", os.Args[5]),

		//[]string{},
		[]string{"dn", "cn"}, //"SamAccountName"
		nil,
	)
	result, err := l.Search(anonReq)
	if err != nil {
		return nil, fmt.Errorf("Anonymous Bind Search Error: %s", err)
	}

	if len(result.Entries) > 0 {
		//result.Entries[0].Print()
		return result, nil
	} else {
		return nil, fmt.Errorf("Couldn't fetch anonymous bind search entries")
	}
}

// Normal Bind and Search
func BindAndSearch(l *ldap.Conn) (*ldap.SearchResult, error) {
	//err := l.Bind(BindUsername, BindPassword)
	err := l.Bind(os.Args[1], os.Args[2])
	if err != nil {
		log.Println("Ошибка аутентификации")
		return nil, err
	}

	searchReq := ldap.NewSearchRequest(
		//BaseDN,
		os.Args[4],
		//ldap.ScopeBaseObject, // you can also use ldap.ScopeWholeSubtree
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		//Filter,
		//fmt.Sprintf("(&(objectClass=organizationalPerson)(uid=%s))", "denis.tirskikh"),
		fmt.Sprintf("(&(sAMAccountName=%s))", os.Args[5]),

		[]string{"dn", "cn"},
		nil,
	)
	result, err := l.Search(searchReq)
	if err != nil {
		return nil, fmt.Errorf("Search Error: %s", err)
	}

	if len(result.Entries) > 0 {
		return result, nil
	} else {
		return nil, fmt.Errorf("Couldn't fetch search entries")
	}
}

// Ldap Connection with TLS
func ConnectTLS() (*ldap.Conn, error) {
	// You can also use IP instead of FQDN
	//l, err := ldap.DialURL(fmt.Sprintf("ldaps://%s:636", FQDN))
	l, err := ldap.DialURL(fmt.Sprintf("ldaps://%s:636", os.Args[3]))
	if err != nil {
		return nil, err
	}

	return l, nil
}

// Ldap Connection without TLS
func Connect() (*ldap.Conn, error) {
	// You can also use IP instead of FQDN
	//l, err := ldap.DialURL(fmt.Sprintf("ldap://%s:389", FQDN))
	l, err := ldap.DialURL(fmt.Sprintf("ldap://%s:389", os.Args[3]))
	//l, err := ldap.DialURL(fmt.Sprintf("ldap://%s:3268", FQDN))
	if err != nil {
		return nil, err
	}

	return l, nil
}
