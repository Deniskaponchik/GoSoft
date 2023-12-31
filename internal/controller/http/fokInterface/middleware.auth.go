package fokusov

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// https://github.com/zhashkevych/todo-app/blob/master/pkg/handler/middleware.go
func (fok *Fokusov) setUserStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			//newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
			c.Set("userLogin", nil)
			fok.Logger.Println("Empty Authorization key in Header")
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			//newErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
			c.Set("userLogin", nil)
			fok.Logger.Println("invalid Authorization header")
			return
		}

		if len(headerParts[1]) == 0 {
			//newErrorResponse(c, http.StatusUnauthorized, "token is empty")
			c.Set("userLogin", nil)
			fok.Logger.Println("token is empty")
			return
		}

		userId, err := h.services.Authorization.ParseToken(headerParts[1])
		userLogin, err := fok.Urest.
		if err != nil {
			newErrorResponse(c, http.StatusUnauthorized, err.Error())
			return
		}

		c.Set(userCtx, userId)
	}
}

/* This middleware sets whether the user is logged in or not
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
			fok.Logger.Println("Ключ is_logged_in в запросе отсутствует")
			if token, err := c.Cookie("token"); err != nil || token == "" {
				fok.Logger.Println("Token пустой или отсутствует")
				c.AbortWithStatus(http.StatusUnauthorized)
				c.Redirect(http.StatusTemporaryRedirect, "/user/login")
			}
		}
	}
}

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
}
