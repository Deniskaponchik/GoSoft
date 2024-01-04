package fokusov

import (
	"github.com/deniskaponchik/GoSoft/internal/entity"
	"github.com/gin-gonic/gin"
	//"github.com/golang-jwt/jwt"
	"net/http"
	//"time"
)

func (fok *Fokusov) showAdminkaPage(c *gin.Context) {
	//based on showArticleCreationPage function
	// Call the render function with the name of the template to render

	msgSlice := adminkaPageMsg
	adminkaPageMsg = make([]string, 10)

	//sapcnArr := [10]string{"БиДВ_BRK", "БиДВ_BRK", "БиДВ_BRK", "БиДВ_BRK", "БиДВ_BRK", "БиДВ_BRK", "БиДВ_BRK", "БиДВ_BRK", "БиДВ_BRK", "БиДВ_BRK"}
	sapcnArr := fok.Urest.GetSapcnSortSliceForAdminkaPage()

	userGivenName := ""
	//var errCook error
	//var errGet bool
	//if userGivenName, errCook := c.Cookie("userGivenName"); errCook != nil || userGivenName == "" {
	cookieGivenName, errCook := c.Cookie("userGivenName")
	if errCook != nil {
		fok.Logger.Println("userGivenName не удалось получить из куки")

		storeGivenName, _ := c.Get("userGivenName")
		if storeGivenName == nil {
			fok.Logger.Println("userGivenName не удалось получить из store")
		} else {
			userGivenName = storeGivenName.(string)
		}
	} else {
		userGivenName = cookieGivenName
	}

	render(c, gin.H{
		"page_adminka": true,
		"sapcnArr":     sapcnArr,
		//"arr0":       adminkaPageMsg[0],
		"arr0":          msgSlice[0],
		"arr2":          msgSlice[2],
		"title":         "Adminka",
		"userGivenName": userGivenName},
		"adminka.html")
}

func showLoginPage(c *gin.Context) {
	// Call the render function with the name of the template to render
	render(c, gin.H{
		"page_login": true,
		"title":      "Login",
	}, "login.html")
}

func (fok *Fokusov) performLogin(c *gin.Context) {

	//username := c.PostForm("username")   password := c.PostForm("password")
	userUnifi := &entity.User{
		Login:    c.PostForm("username"),
		Password: c.PostForm("password"),
	}

	var sameSiteCookie http.SameSite

	//if isUserValid(username, password) {
	err := fok.Urest.CheckUser(userUnifi) //в user добавился GivenName
	//err := fok.Authentication.AuthSecur(userUnifi)
	if err == nil {
		//token := generateSessionToken()
		//token, errGenToken := fok.generateSessionToken(userUnifi)
		token, errGenToken := fok.Urest.GetToken(userUnifi)
		if errGenToken == nil {
			fok.Logger.Println(token)
			fok.Logger.Println(userUnifi.GivenName)
			fok.Logger.Println(userUnifi.Login)

			c.SetSameSite(sameSiteCookie)

			//c.SetCookie("token", token, 3600, "", "", sameSiteCookie, false, true) 	//original
			c.SetCookie("token", token, 3600, "", "", false, true)
			c.SetCookie("userGivenName", userUnifi.GivenName, 3600, "", "", false, true)
			//c.SetCookie("is_logged_in", true, 3600, "/", "", false, true)
			c.Set("userGivenName", userUnifi.GivenName)
			//c.Set("userLogin", userUnifi.Login)
			c.Set("is_logged_in", true)
			//c.Header("Authorization", token)

			//render(c, gin.H{"title": "Successful Login"}, "login-successful.html")
			//c.Redirect(http.StatusTemporaryRedirect, "/user/adminka")
			fok.showAdminkaPage(c)
		} else {

			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"ErrorTitle": "Token Failed",
				//"ErrorMessage": "Invalid credentials provided"})
				"ErrorMessage": err})
		}
	} else {
		// If the username/password combination is invalid,
		fok.Logger.Println(userUnifi.Login)

		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"ErrorTitle": "Login Failed",
			//"ErrorMessage": "Invalid credentials provided"})
			"ErrorMessage": err})
	}
}

func logout(c *gin.Context) {

	var sameSiteCookie http.SameSite
	c.SetSameSite(sameSiteCookie)

	//c.SetCookie("token", "", -1, "", "", sameSiteCookie, false, true)  //ORIGINAL
	c.SetCookie("token", "", -1, "", "", false, true)
	c.SetCookie("userGivenName", "", -1, "", "", false, true)
	//c.Set("userLogin", nil)
	//c.Set("userGivenName", nil)
	c.Set("is_logged_in", false)
	//c.Header("Authorization", "")

	//c.Redirect(http.StatusTemporaryRedirect, "/")
	c.Redirect(http.StatusTemporaryRedirect, "/user/login")
}

/*МОЁ
func (fok *Fokusov) generateSessionToken(user *entity.User) (string, error) {
	//https://ru.hexlet.io/courses/go-web-development/lessons/auth/theory_unit
	jwtSecretKey := []byte(fok.JwtKey)

	//https://www.iana.org/assignments/jwt/jwt.xhtml
	// Генерируем полезные данные, которые будут храниться в токене
	payload := jwt.MapClaims{
		//"sub": user.Email,
		"nickname": user.Login,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}

	// Создаем новый JWT-токен и подписываем его по алгоритму HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	t, err := token.SignedString(jwtSecretKey)
	if err != nil {
		//logrus.WithError(err).Error("JWT token signing")
		//return c.SendStatus(fiber.StatusInternalServerError)
		fok.Logger.Println("JWT token НЕ БЫЛ подписан")
		return "", err
	} else {
		return t, nil
	}
}*/
/*
//ORIGINAL
func generateSessionToken() string {
	// We're using a random 16 character string as the session token.
	//This is NOT a secure way of generating session tokens
	return strconv.FormatInt(rand.Int63(), 16)
}*/

//REGISTRATION
/*
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
}*/
