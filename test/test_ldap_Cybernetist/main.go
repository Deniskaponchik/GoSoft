package main

//https://cybernetist.com/2020/05/18/getting-started-with-go-ldap/
//КОД НЕ ДОВЁЛ ДО РАБОЧЕГО СОСТОЯНИЯ

import (
	"fmt"
	"github.com/go-ldap/ldap"
	"log"
)

func main() {

	ldapURL := "ldap://example.com:389" //
	l, err := ldap.DialURL(ldapURL)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	//err = l.UnauthenticatedBind("cn=read-only-admin,dc=example,dc=com")
	err = l.UnauthenticatedBind("")
	if err != nil {
		log.Fatal(err)
	}

	//
	user := ""                         //""
	baseDN := "DC=corp,DC=tele2,DC=ru" //"DC=example,DC=com"
	filter := fmt.Sprintf("(CN=%s)", ldap.EscapeFilter(user))

	// Filters must start and finish with ()!
	searchReq := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree, 0, 0, 0, false,
		filter,                     //поиск по CN
		[]string{"sAMAccountName"}, //SamAccountName  sAMAccountName  //вывод
		[]ldap.Control{})

	result, err := l.Search(searchReq)
	if err != nil {
		//return fmt.Errorf("failed to query LDAP: %w", err)
		fmt.Errorf("failed to query LDAP: %w", err)
	}

	if len(result.Entries) > 0 {
		log.Println("result is empty")
	} else {
		log.Println("Got", len(result.Entries), "search results")
	}

}
