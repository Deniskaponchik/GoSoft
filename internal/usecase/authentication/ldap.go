package authentication

type Ldap struct {
	DN       string
	Domain   string
	Login    string
	Password string
	RoleDN   string
	Server   string
}

func NewLdap(dn, dm, l, p, r, s string) *Ldap {
	return &Ldap{
		DN:       dn,
		Domain:   dm,
		Login:    l,
		Password: p,
		RoleDN:   r,
		Server:   s,
	}
}
