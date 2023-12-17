package fokusov

import (
	"github.com/deniskaponchik/GoSoft/Unifi/internal/entity"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (fok *Fokusov) officeNew(c *gin.Context) {
	newSapcn := c.PostForm("02_sapcn")
	login := c.PostForm("02_login")
	timeZoneStr := c.PostForm("02_zone")
	timeZone, err := strconv.Atoi(timeZoneStr)
	if err == nil {
		office := &entity.Office{
			Site_ApCutName: newSapcn,
			UserLogin:      login,
			TimeZone:       timeZone,
			TimeZoneStr:    timeZoneStr,
		}

		fok.Logger.Println(office.Site_ApCutName)
		fok.Logger.Println(office.UserLogin)
		fok.Logger.Println(office.TimeZone)

		err = fok.Urest.OfficeNew(office) //newSapcn, login, timeZone)
		if err != nil {
			fok.Logger.Println("Создать новый офис не удалось")
			fok.Logger.Println(err)
			adminkaPageMsg[0] = err.Error() //"Создать новый офис логин не удалось"
		} else {
			adminkaPageMsg[0] = "логин успешно изменён"
		}

		fok.showAdminkaPage(c)
		//c.Redirect(http.StatusTemporaryRedirect, "/user/adminka")
		//c.Request.Method = "GET"
		//c.HTML(http.StatusSeeOther, "adminka.html", nil)
	}
}

func (fok *Fokusov) officeLoginChange(c *gin.Context) {
	sapcn := c.PostForm("00_sapcn")
	newLogin := c.PostForm("00_login")

	fok.Logger.Println(sapcn)
	fok.Logger.Println(newLogin)

	err := fok.Urest.ChangeSapcnLogin(sapcn, newLogin)
	if err != nil {
		fok.Logger.Println("Изменить логин не удалось")
		fok.Logger.Println(err)
		adminkaPageMsg[0] = err.Error() //"Изменить логин не удалось"
	} else {
		adminkaPageMsg[0] = "логин успешно изменён"
	}

	fok.showAdminkaPage(c)
	//c.Redirect(http.StatusTemporaryRedirect, "/user/adminka")
	//c.Request.Method = "GET"
	//c.HTML(http.StatusSeeOther, "adminka.html", nil)
}

func functionWithErrorExample(c *gin.Context) { //old name performLogin
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
