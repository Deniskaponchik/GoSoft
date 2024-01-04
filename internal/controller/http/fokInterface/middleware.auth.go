package fokusov

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (fok *Fokusov) setUserStatus() gin.HandlerFunc {
	fok.Logger.Println("")
	fok.Logger.Println("")

	return func(c *gin.Context) {
		if token, err := c.Cookie("token"); err == nil || token != "" {
			fok.Logger.Println(token)

			//userGivenName, errCheckToken := fok.Urest.CheckToken(token)
			userLogin, errCheckToken := fok.Urest.CheckToken(token)
			if errCheckToken != nil {
				fok.Logger.Println(errCheckToken.Error())
				c.SetCookie("userGivenName", "", -1, "", "", false, true)
				//c.Set("userLogin", nil)
				//c.Set("userGivenName", nil)
				c.Set("is_logged_in", false)
			} else {
				//fok.Logger.Println(userGivenName)
				fok.Logger.Println(userLogin)
				//c.Set("userLogin", userLogin)  		 //в куках НЕТ. только в токене зашит
				//c.Set("userGivenName", userGivenName)  //пришёл в куках из предыдущих запросов
				c.Set("is_logged_in", true) //это не куки, а store, действующий в рамках одного контекста
			}
		} else {
			fok.Logger.Println("токен в запросе пустой")
			c.SetCookie("userGivenName", "", -1, "", "", false, true)
			//c.Set("userLogin", nil)
			//c.Set("userGivenName", nil)
			c.Set("is_logged_in", false)
		}
	}
}

/* https://github.com/zhashkevych/todo-app/blob/master/pkg/handler/middleware.go
func (fok *Fokusov) setUserStatus() gin.HandlerFunc { //  c *gin.Context
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			//newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
			//c.Set("userLogin", nil)
			c.Set("userGivenName", nil)
			c.Set("is_logged_in", false)
			fok.Logger.Println("Empty Authorization key in Header")
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			//newErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
			//c.Set("userLogin", nil)
			c.Set("userGivenName", nil)
			c.Set("is_logged_in", false)
			fok.Logger.Println("invalid Authorization header")
			return
		}

		if len(headerParts[1]) == 0 {
			//newErrorResponse(c, http.StatusUnauthorized, "token is empty")
			//c.Set("userLogin", nil)
			c.Set("userGivenName", nil)
			c.Set("is_logged_in", false)
			fok.Logger.Println("token is empty")
			return
		}

		//userId, err := h.services.Authorization.ParseToken(headerParts[1])
		//userLogin, err := fok.Authorization.ParseToken(headerParts[1])
		userGivenName, err := fok.Urest.CheckToken(headerParts[1])
		if err != nil {
			//newErrorResponse(c, http.StatusUnauthorized, err.Error())
			//c.Set("userLogin", nil)
			c.Set("userGivenName", nil)
			c.Set("is_logged_in", false)
			fok.Logger.Println(err.Error())
			return
		}

		//c.Set(userCtx, userId)
		//c.Set("userLogin", userLogin)
		c.Set("userGivenName", userGivenName)
		c.Set("is_logged_in", true)
	}
}*/

/* ORIGINAL. Функция запускается при ВСЕХ входящих запросах
func setUserStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		if token, err := c.Cookie("token"); err == nil || token != "" {
			c.Set("is_logged_in", true)
		} else {
			c.Set("is_logged_in", false)
		}
	}
}*/

func (fok *Fokusov) ensureLoggedIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		// If there's an error or if the token is empty the user is not logged in
		loggedInInterface, _ := c.Get("is_logged_in")
		loggedIn := loggedInInterface.(bool)
		if !loggedIn {
			//if token, err := c.Cookie("token"); err != nil || token == "" {
			//c.AbortWithStatus(http.StatusUnauthorized)
			c.Redirect(http.StatusTemporaryRedirect, "/user/login")
		}
	}
}

/*
func (fok *Fokusov) ensureLoggedIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		//loggedInInterface, _ := c.Get("userLogin")
		loggedInInterface, _ := c.Get("userGivenName")
		if loggedInInterface == nil {
			fok.Logger.Println("Ключ userGivenName в запросе отсутствует")
			c.AbortWithStatus(http.StatusUnauthorized)
			c.Redirect(http.StatusTemporaryRedirect, "/user/login")
		}
	}
}*/
/*ORIGINAL
//Если пользователь не залогинен, то перенаправит на страничку залогинивания. Если нет, то НЕ ДЕЛАЕТ НИЧЕГО
func ensureLoggedIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		// If there's an error or if the token is empty the user is not logged in
		loggedInInterface, _ := c.Get("is_logged_in")
		loggedIn := loggedInInterface.(bool)
		if !loggedIn {
			if token, err := c.Cookie("token"); err != nil || token == "" {
				c.AbortWithStatus(http.StatusUnauthorized)
				c.Redirect(http.StatusTemporaryRedirect, "/user/login")
			}
		}
	}
}*/

func (fok *Fokusov) ensureNotLoggedIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		loggedInInterface, _ := c.Get("is_logged_in")
		loggedIn := loggedInInterface.(bool)
		if loggedIn {
			//if token, err := c.Cookie("token"); err == nil || token != "" {
			//c.AbortWithStatus(http.StatusUnauthorized)
			c.Redirect(http.StatusTemporaryRedirect, "/user/adminka")
		}
	}
}

/* Моя неудачная попытка
func (fok *Fokusov) ensureNotLoggedIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		// If there's no error or if the token is not empty the user is already logged in
		//loggedInInterface, _ := c.Get("userLogin")
		loggedInInterface, _ := c.Get("userGivenName")
		if loggedInInterface != nil {
			fok.Logger.Println("Ключ userGivenName в запросе присутствует")
			c.AbortWithStatus(http.StatusUnauthorized)
			c.Redirect(http.StatusTemporaryRedirect, "/user/adminka")
		}
	}
}*/
/*ORIGINAL
// Если пользователь залогинен, то перенаправит на страничку админки. Если нет, то НЕ ДЕЛАЕТ НИЧЕГО
func ensureNotLoggedIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		// If there's no error or if the token is not empty the user is already logged in
		loggedInInterface, _ := c.Get("is_logged_in")
		loggedIn := loggedInInterface.(bool)
		if loggedIn {
			if token, err := c.Cookie("token"); err == nil || token != "" {
				c.AbortWithStatus(http.StatusUnauthorized)
				c.Redirect(http.StatusTemporaryRedirect, "/user/adminka")
			}
		}
	}
}*/

//https://fokusov.com/posts/razrabotka-web-prilozhenij-i-mikroservisov-na-golang-s-gin/
//https://github.com/demo-apps/go-gin-app/blob/master/middleware.auth.go
