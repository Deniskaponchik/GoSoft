package fokusov

/* CAN DELETE MY CODE
func (fok *Fokusov) isUserValid(userUnifi *entity.User) error {
	err := fok.Urest.LdapCheckUser(userUnifi)
	if err == nil{
		return nil
	} else{
		return err
	}
}*/

/* Original version with slice user
type user struct {
	Username string `json:"username"`
	Password string `json:"-"`
}

func (fok *Fokusov) isUserValid(username, password string) bool {
	for slice user
	for _, u := range userList {
		if u.Username == username && u.Password == password {
			return true
		}
	}
	return false
}

// we're storing the user list in memory. We also have some users predefined.
// In a real application, this list will most likely be fetched from a database.
// Moreover, in production settings, you should store passwords securely
// by salting and hashing them instead of using them as we're doing in this demo
var userList = []user{
	user{Username: "user1", Password: "pass1"},
	user{Username: "user2", Password: "pass2"},
	user{Username: "user3", Password: "pass3"},
}

// Register a new user with the given username and password
func registerNewUser(username, password string) (*user, error) {
	if strings.TrimSpace(password) == "" {
		return nil, errors.New("The password can't be empty")
	} else if !isUsernameAvailable(username) {
		return nil, errors.New("The username isn't available")
	}

	u := user{Username: username, Password: password}

	userList = append(userList, u)

	return &u, nil
}

// Check if the supplied username is available
func isUsernameAvailable(username string) bool {
	for _, u := range userList {
		if u.Username == username {
			return false
		}
	}
	return true
}*/
