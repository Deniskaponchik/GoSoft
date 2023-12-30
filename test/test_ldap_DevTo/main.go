package main

import (
	"fmt"
	"github.com/go-ldap/ldap"
	"log"
	"os"
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
	//result.Entries[0].Print()
	if result.Entries[0].Attributes[0].Values[0] != "" {
		fmt.Println(result.Entries[0].Attributes[0].Values[0])
		//result.Entries[0].Attributes[0].
	} else {
		fmt.Println("Атрибут DisplayName пустой")
	}
	if result.Entries[0].Attributes[1].Values[0] != "" {
		fmt.Println([]byte(result.Entries[0].Attributes[1].Values[0]))
		sid := string([]byte(result.Entries[0].Attributes[1].Values[0]))
		fmt.Println(sid)

		/*https://github.com/bwmarrin/go-objectsid
		//НЕ ЗАРАБОТАЛО
		// A objectSID in base64
		// This is just here for the purpose of an example program
		//b64sid := `AQUAAAAAAAUVAAAArC22DNydmGz4WUTnUAQAAA==`
		b64sid := result.Entries[0].Attributes[1].Values[0]

		// Convert the above into binary form. This is the value you would get from a LDAP query on AD.
		bsid, _ := base64.StdEncoding.DecodeString(b64sid)

		// Decode the binary objectsid into a SID object
		sid := objectsid.Decode(bsid)

		// Print out just one component of the SID
		fmt.Println(sid.Authority)

		// Print out the relative identifier
		fmt.Println(sid.RID())

		// Print the entire ObjectSID
		fmt.Println(sid)
		*/

	} else {
		fmt.Println("Атрибут SID пустой")
	}

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
		//[]string{"dn", "cn"}, //"SamAccountName"
		[]string{"DisplayName", "objectSid"},
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

		//[]string{"dn", "cn"},
		[]string{"DisplayName", "objectSid"},
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
