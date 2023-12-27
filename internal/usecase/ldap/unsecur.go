package ldap

/*Отключаю за неиспользование. Если понадобится, привести по образу и подобию с версией TLS
func (ldap *Ldap) AuthUnSecur(userLogin, userPassword string) (err error) {
	//Non-TLS Connection
	l, err := Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	// Anonymous Bind and Search
	result, err := AnonymousBindAndSearch(l)
	if err != nil {
		log.Fatal(err)
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
}*/
