// handlers.client.go
package fokusov

import (
	"github.com/deniskaponchik/GoSoft/internal/entity"
	"github.com/gin-gonic/gin"
	"strings"
)

/*
func (fok *Fokusov) postClient(c *gin.Context){
	// Obtain the POSTed hostname
	hostname := c.PostForm("hostname")

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

func (fok *Fokusov) getClient(c *gin.Context) {
	fok.Logger.Println("")
	// Check if the client hostname is valid
	var clientHostname string
	clientHostname = c.PostForm("cl_hostname")
	if clientHostname == "" {
		clientHostname = c.Param("client_hostname")
		fok.Logger.Println("Client взят из метода GET")
	}
	clientHostname = strings.ToUpper(clientHostname)
	fok.Logger.Println(clientHostname)

	//if client, err := getArticleByID(articleID); err == nil {
	//client := fok.UnifiUC.GetClientForRest(clientHostname)
	client := fok.Urest.GetClientForRest(clientHostname)

	if client != nil {
		fok.Logger.Println("клиент найден в мапе клиентов")
		//fmt.Println(client.Hostname)
		//sliceAnomalies := []*entity.Anomaly{}

		j := 0
		var date string
		anomalyTempMap := make(map[string]string)

		//пересобираем массив в обратную сторону
		lenClSliceAnom := len(client.SliceAnomalies)
		sliceAnomalies := make([]*entity.Anomaly, int(lenClSliceAnom))
		//for i := len(client.SliceAnomalies) - 1; i > -1; i-- {
		for i := lenClSliceAnom - 1; i > -1; i-- {
			//sliceAnomalies = append(sliceAnomalies, client.SliceAnomalies[i])
			sliceAnomalies[j] = client.SliceAnomalies[i]
			j++

			date = strings.Split(client.SliceAnomalies[i].DateHour, " ")[0]
			anomalyTempMap[date] = date
		}
		client.CountAnomaly = len(anomalyTempMap)
		redMarker := false
		if client.CountAnomaly > 9 {
			redMarker = true
		}

		// Call the render function with the title, article and the name of the
		// template
		render(c, gin.H{
			"page_client":  true,
			"title":        client.Hostname,
			"hostname":     client.Hostname,
			"countanomaly": client.CountAnomaly,
			"redmarker":    redMarker,
			//"anomalies_struct": client.SliceAnomalies},
			"anomalies_struct": sliceAnomalies},
			"client.html")

	} else {
		fok.Logger.Println("клиент НЕ найден в мапе клиентов")
		errMessage := "Client not found: " + clientHostname
		// If the client is not found, abort with an error
		//c.AbortWithError(http.StatusNotFound, err)
		render(c, gin.H{
			"title":    "Client not found",
			"hostname": errMessage},
			//"anomalies_struct": client.SliceAnomalies},
			"client.html")
	}
}

/*
func (fok *Fokusov) getClientFok(c *gin.Context) {
	// Check if the client hostname is valid
	//if articleID, err := strconv.Atoi(c.Param("article_id")); err == nil {
	clientHostname := c.Param("client_hostname")
	fmt.Println(clientHostname)

	// Check if the client exists
	//if client, err := getArticleByID(articleID); err == nil {
	client := fok.UnifiUC.GetClientForRest(clientHostname)
	//fok.UnifiClient = fok.UnifiUC.GetClientForRest(clientHostname)
	if fok.UnifiClient != nil {
		fmt.Println("клиент найден в мапе клиентов")
		fmt.Println(fok.UnifiClient.Hostname)
		// Call the render function with the title, article and the name of the
		// template
		fok.render(c, gin.H{
			"title":            fok.UnifiClient.Hostname,
			"hostname":         fok.UnifiClient.Hostname,
			"anomalies_struct": fok.UnifiClient.SliceAnomalies},
			"client.html")

	} else {
		fmt.Println("клиент НЕ найден в мапе клиентов")
		// If the client is not found, abort with an error
		//c.AbortWithError(http.StatusNotFound, err)
		fok.render(c, gin.H{
			"title":    "Client did not found",
			"hostname": "Client did not found"},
			//"anomalies_struct": client.SliceAnomalies},
			"client.html")
	}
}*/

func (fok *Fokusov) showClientRequestPage(c *gin.Context) {
	// Call the render function with the name of the template to render
	render(c, gin.H{
		"page_client": true,
		"title":       "Client Request Page"},
		"client.html")
}
