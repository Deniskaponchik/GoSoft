package fokusov

import (
	"github.com/deniskaponchik/GoSoft/Unifi/internal/entity"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (fok *Fokusov) officeChange(c *gin.Context) {
	sapcn := c.PostForm("01_sapcn")
	action := c.PostForm("01_dropdown")
	newValue := c.PostForm("01_field")

	switch action {

	case "Sapcn Name":
		fok.Logger.Println("Выбрано изменение имени sapcn")
		err := fok.Urest.OfficeSapcnChange(sapcn, newValue)
		if err != nil {
			fok.Logger.Println("Изменить sapcn не удалось")
			fok.Logger.Println(err)
			adminkaPageMsg[1] = err.Error()
		} else {
			adminkaPageMsg[1] = "sapcn успешно изменён"
		}

	case "Login":
		fok.Logger.Println("Выбрано изменение логина")
		err := fok.Urest.OfficeLoginChange(sapcn, newValue)
		if err != nil {
			fok.Logger.Println("Изменить логин не удалось")
			fok.Logger.Println(err)
			adminkaPageMsg[1] = err.Error()
		} else {
			adminkaPageMsg[1] = "Логин успешно изменён"
		}

	case "Time Zone":
		fok.Logger.Println("Выбрано изменение Time Zone")
		err := fok.Urest.OfficeTimeZoneChange(sapcn, newValue)
		if err != nil {
			fok.Logger.Println("Изменить Time Zone не удалось")
			fok.Logger.Println(err)
			adminkaPageMsg[1] = err.Error()
		} else {
			adminkaPageMsg[1] = "Time Zone успешно изменён"
		}

	case "Exception":
		fok.Logger.Println("Выбрано изменение Exception")
		err := fok.Urest.OfficeExceptionChange(sapcn, newValue)
		if err != nil {
			fok.Logger.Println("Изменить Exception не удалось")
			fok.Logger.Println(err)
			adminkaPageMsg[1] = err.Error()
		} else {
			adminkaPageMsg[1] = "Exception успешно изменён"
		}
	}

	fok.showAdminkaPage(c)
}

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
			adminkaPageMsg[2] = err.Error() //"Создать новый офис логин не удалось"
		} else {
			adminkaPageMsg[2] = "Новый офис успешно создан"
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

	//err := fok.Urest.ChangeSapcnLogin(sapcn, newLogin)
	err := fok.Urest.OfficeLoginChange(sapcn, newLogin)
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

func functionWithErrorExample(c *gin.Context) {
	//old name performLogin
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
