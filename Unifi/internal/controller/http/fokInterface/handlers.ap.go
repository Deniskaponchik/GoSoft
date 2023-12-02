// handlers.client.go
package fokusov

import (
	"fmt"
	"github.com/deniskaponchik/GoSoft/Unifi/internal/entity"
	"github.com/gin-gonic/gin"
	"log"
	"strings"
)

func (fok *Fokusov) getAP(c *gin.Context) {
	log.Println("")
	// Check if the client hostname is valid
	var apHostname string
	apHostname = c.PostForm("ap_hostname")
	if apHostname == "" {
		apHostname = c.Param("ap_hostname")
		fmt.Println("Точка взята из метода GET")
	}
	apHostname = strings.ToUpper(apHostname)
	fmt.Println(apHostname)

	//if client, err := getArticleByID(articleID); err == nil {
	//ap := fok.UnifiUC.GetApForRest(apHostname)
	ap := fok.Urest.GetApForRest(apHostname)

	if ap != nil {
		fmt.Println("точка найдена в мапе")
		fmt.Println(ap.Name)
		sliceAnomalies := []*entity.Anomaly{}

		var date string
		anomalyTempMap := make(map[string]string)
		//пересобираем массив в обратную сторону
		for i := len(ap.SliceAnomalies) - 1; i > -1; i-- {
			sliceAnomalies = append(sliceAnomalies, ap.SliceAnomalies[i])

			date = strings.Split(ap.SliceAnomalies[i].DateHour, " ")[0]
			anomalyTempMap[date] = date
		}
		ap.CountAnomaly = len(anomalyTempMap)
		redMarker := false
		if ap.CountAnomaly > 9 {
			redMarker = true
		}

		// Call the render function with the title, article and the name of the
		// template
		render(c, gin.H{
			"title":        ap.Name,
			"hostname":     ap.Name,
			"countanomaly": ap.CountAnomaly,
			"redmarker":    redMarker,
			//"anomalies_struct": client.SliceAnomalies},
			"anomalies_struct": sliceAnomalies},
			"ap.html")

	} else {
		fmt.Println("Точка НЕ найдена в мапе")
		errMessage := "Access point not found: " + apHostname
		// If the client is not found, abort with an error
		//c.AbortWithError(http.StatusNotFound, err)
		render(c, gin.H{
			"title":    "Access point not found",
			"hostname": errMessage},
			//"anomalies_struct": client.SliceAnomalies},
			"ap.html")
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

func (fok *Fokusov) showApRequestPage(c *gin.Context) {
	// Call the render function with the name of the template to render
	render(c, gin.H{
		"title": "Ap Request Page"}, "ap.html")
}
