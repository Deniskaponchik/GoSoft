package fokusov

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"strconv"
)

func showLoginPage(c *gin.Context) {
	// Call the render function with the name of the template to render
	render(c, gin.H{
		"page_login": true,
		"title":      "Login",
	}, "login.html")
}

func performLogin(c *gin.Context) {
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

func generateSessionToken() string {
	// We're using a random 16 character string as the session token. This is NOT a secure way of generating session tokens
	// DO NOT USE THIS IN PRODUCTION
	return strconv.FormatInt(rand.Int63(), 16)
}

func logout(c *gin.Context) {

	var sameSiteCookie http.SameSite

	// Clear the cookie
	c.SetSameSite(sameSiteCookie)
	//c.SetCookie("token", "", -1, "", "", sameSiteCookie, false, true)
	c.SetCookie("token", "", -1, "", "", false, true) //моё

	//c.Redirect(http.StatusTemporaryRedirect, "/")
	c.Redirect(http.StatusTemporaryRedirect, "/user/login")
}

func showAdminkaPage(c *gin.Context) {
	//based on showArticleCreationPage function
	// Call the render function with the name of the template to render
	render(c, gin.H{
		"page_adminka": true,
		"title":        "Adminka"},
		"adminka.html")
}

func showRegistrationPage(c *gin.Context) {
	// Call the render function with the name of the template to render
	render(c, gin.H{
		//"title": "Register"}, "register.html")
		"title": "Register"},
		"register.html")
}

func register(c *gin.Context) {
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
