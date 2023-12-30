package fokusov

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

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
			} else {

			}
		}
	}
}

/* Если пользователь не залогинен, то перенаправит на страничку залогинивания. Если нет, то НЕ ДЕЛАЕТ НИЧЕГО
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

// This middleware sets whether the user is logged in or not
func setUserStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		if token, err := c.Cookie("token"); err == nil || token != "" {
			c.Set("is_logged_in", true)
		} else {
			c.Set("is_logged_in", false)
		}
	}
}
