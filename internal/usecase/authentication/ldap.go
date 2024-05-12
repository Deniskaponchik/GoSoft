package authentication

type Ldap struct {
	DN       string
	Domain   string
	Login    string
	Password string
	RoleDN   string
	Server   string
}

func NewLdap(dn, dm, login, password, r, s string) *Ldap {
	return &Ldap{
		DN:       dn,
		Domain:   dm,
		Login:    login,
		Password: password,
		RoleDN:   r,
		Server:   s,
	}
}
