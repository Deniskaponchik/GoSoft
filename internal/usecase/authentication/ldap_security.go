package authentication

import (
	"fmt"
	"github.com/deniskaponchik/GoSoft/internal/entity"
	"github.com/go-ldap/ldap"
	"log"
	"strings"
)

func (ldp *Ldap) AuthSecur(user *entity.User) (err error) { //userLogin, userPassword string
	//TLS Connection
	//l, err := ConnectTLS()
	connection, err := ldp.ConnectTLS()
	if err != nil {
		//log.Fatal(err)
		return err
	}
	//defer l.Close()
	defer connection.Close()

	//Normal Bind and Search
	//result, err := BindAndSearch(l)
	err = ldp.BindAndSearch(connection, user) //userLogin, userPassword
	if err != nil {
		//log.Fatal(err)
		return err
	} else {
		return nil
	}
}

// Normal Bind and Search
func (ldp *Ldap) BindAndSearch(l *ldap.Conn, user *entity.User) (err error) { //(*ldap.SearchResult, error) {  //userLogin, userPassword string

	//полученными о ПОЛЬЗОВАТЕЛЯ логином и паролем пытаемся авторизоваться в LDAP
	//err := l.Bind(BindUsername, BindPassword)
	//err := l.Bind(os.Args[1], os.Args[2])
	//bindUsername := user.Login + "@domain"
	bindUsername := user.Login + "@" + ldp.Domain //change after reboot PC
	err = l.Bind(bindUsername, user.Password)
	if err != nil {
		log.Println("Ошибка аутентификации на LDAP")
		log.Println(err)
		//return nil, err
		return err
	}

	filter := fmt.Sprintf("(&(sAMAccountName=%s)(MemberOf=%s))", user.Login, ldp.RoleDN)
	searchReq := ldap.NewSearchRequest(
		ldp.DN, //BaseDN,	//os.Args[4],
		//ldap.ScopeBaseObject, // you can also use ldap.ScopeWholeSubtree
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		//Filter,
		//fmt.Sprintf("(&(objectClass=organizationalPerson)(uid=%s))", "login"),
		//fmt.Sprintf("(&(sAMAccountName=%s))", os.Args[5]),
		filter,
		[]string{"DisplayName"}, //, "objectSid"
		nil,
	)

	result, err := l.Search(searchReq)
	if err != nil {
		//return nil, fmt.Errorf("Search Error: %s", err)
		//return fmt.Errorf("Search Error: %s", err)
		log.Println("Ошибка поискового запроса:")
		log.Println(err)
		return err
	}

	if len(result.Entries) > 0 {
		displayName := result.Entries[0].Attributes[0].Values[0]
		if displayName != "" {
			fioSlice := strings.Split(displayName, " ")
			user.FIO = displayName
			user.GivenName = fioSlice[1]
			user.SurName = fioSlice[0]
			user.MiddleName = fioSlice[2]
		}
		return nil
	} else {
		//return nil, fmt.Errorf("Couldn't fetch search entries")
		log.Println("Результат LDAP-поиска пустой. Нет доступного соответствия логина и роли")
		return fmt.Errorf("Результат LDAP-поиска пустой")
	}
}

// Ldap Connection with TLS
func (ldp *Ldap) ConnectTLS() (*ldap.Conn, error) {
	// You can also use IP instead of FQDN
	//l, err := ldap.DialURL(fmt.Sprintf("ldaps://%s:636", FQDN))
	//l, err := ldap.DialURL(fmt.Sprintf("ldaps://%s:636", os.Args[3]))
	l, err := ldap.DialURL(fmt.Sprintf("ldaps://%s:636", ldp.Server))
	if err != nil {
		log.Println("Ошибка создания и проверки подключения к LDAP")
		log.Println(err)
		return nil, err
	}
	return l, nil
}
