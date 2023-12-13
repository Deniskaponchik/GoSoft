package fokusov

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (fok *Fokusov) officeLoginChange(c *gin.Context) {
	sapcn := c.PostForm("00_sapcn")
	newLogin := c.PostForm("00_login")

	//check if correct

	fok.Logger.Println(sapcn)
	fok.Logger.Println(newLogin)

	fok.AdminkaArr[0] = "логин успешно изменён"

	fok.showAdminkaPage(c)
}

func performLog(c *gin.Context) {
	// Obtain the POSTed username and password values
	username := c.PostForm("username")
	password := c.PostForm("password")

	var sameSiteCookie http.SameSite

	if isUserValid(username, password) {
		token := generateSessionToken()
		c.SetSameSite(sameSiteCookie)
		//c.SetCookie("token", token, 3600, "", "", sameSiteCookie, false, true)
		c.SetCookie("token", token, 3600, "", "", false, true) //моё
		c.Set("is_logged_in", true)

		//render(c, gin.H{"title": "Successful Login"}, "login-successful.html")
		c.Redirect(http.StatusTemporaryRedirect, "/user/adminka")

	} else {
		// If the username/password combination is invalid,
		// show the error message on the login page
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"ErrorTitle":   "Login Failed",
			"ErrorMessage": "Invalid credentials provided"})
	}
}

func reg(c *gin.Context) {
	// Obtain the POSTed username and password values
	username := c.PostForm("username")
	password := c.PostForm("password")

	var sameSiteCookie http.SameSite

	if _, err := registerNewUser(username, password); err == nil {
		// If the user is created, set the token in a cookie and log the user in
		token := generateSessionToken()
		c.SetSameSite(sameSiteCookie)
		//c.SetCookie("token", token, 3600, "", "", sameSiteCookie, false, true)
		c.SetCookie("token", token, 3600, "", "", false, true)
		c.Set("is_logged_in", true)

		render(c, gin.H{
			"title": "Successful registration & Login"}, "login-successful.html")

	} else {
		// If the username/password combination is invalid,
		// show the error message on the login page
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"ErrorTitle":   "Registration Failed",
			"ErrorMessage": err.Error()})

	}
}
