package main

import (
	"fmt"
	"github.com/go-ldap/ldap"
	"log"
)

func main() {
	l, err := ldap.Dial("tcp",
		//fmt.Sprintf("%s:%d", "ldap.example.com", 389))
		fmt.Sprintf("%s:%d", "ldap.corp.tele2.ru", 389))
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	//000004DC: LdapErr: DSID-0C090A71, comment: In order to perform this operation a successful bind must be completed on the connection., data 0, v3839
	//err = l.Bind("domain\\login", "password")
	err = l.Bind("t2ru\\denis.tirskikh", "Norm@ndy2$")
	if err != nil {
		log.Println("Ошибка аутентификации")
	}

	searchRequest := ldap.NewSearchRequest(
		//"dc=example,dc=com", // The base dn to search
		"BaseDN",
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		//"(&(objectClass=organizationalPerson))", // The filter to apply
		//"(&(ObjectClass=user))",
		"(&(sAMAccountName=denis.tirskikh))",
		//"(&(sAMAccountName=login)(MemberOf=CN=LdapRole,LdapRoleDN))",
		//[]string{"dn", "cn"}, // A list attributes to retrieve
		[]string{"DisplayName", "objectSid"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Fatal(err)
	}

	if len(sr.Entries) > 0 {
		if sr.Entries[0].Attributes[0].Values[0] != "" {
			fmt.Println(sr.Entries[0].Attributes[0].Values[0])
		} else {
			fmt.Println("Нулевой трибут пустой")
		}
		if sr.Entries[0].Attributes[1].Values[0] != "" {
			fmt.Println(sr.Entries[0].Attributes[1].Values[0])
		} else {
			fmt.Println("Первый трибут пустой")
		}
	} else {
		//return nil, fmt.Errorf("Couldn't fetch search entries")
		log.Println("Результат LDAP-поиска пустой. Нет доступного соответствия логина и роли")
		fmt.Errorf("Результат LDAP-поиска пустой")
	}

	for _, entry := range sr.Entries {
		fmt.Printf("%s: %v\n", entry.DN, entry.GetAttributeValue("cn"))
	}
}
