// handlers.client.go
package fokusov

import (
	"fmt"
	"github.com/deniskaponchik/GoSoft/Unifi/internal/entity"
	"github.com/gin-gonic/gin"
)

func (fok *Fokusov) getAP(c *gin.Context) {
	// Check if the client hostname is valid
	var apHostname string
	apHostname = c.PostForm("hostname")
	if apHostname == "" {
		apHostname = c.Param("client_hostname")
		fmt.Println("Client взят из метода GET")
	}
	fmt.Println(apHostname)

	// Check if the client exists
	//if client, err := getArticleByID(articleID); err == nil {
	ap := fok.UnifiUC.GetClientForRest(apHostname)
	//fok.UnifiClient = fok.UnifiUC.GetClientForRest(clientHostname)

	if ap != nil {
		fmt.Println("точка найдена в мапе")
		fmt.Println(ap.Hostname)
		sliceAnomalies := []*entity.Anomaly{}
		sliceAnomalies = ap.SliceAnomalies

		// Call the render function with the title, article and the name of the
		// template
		render(c, gin.H{
			"title":    ap.Hostname,
			"hostname": ap.Hostname,
			//"anomalies_struct": client.SliceAnomalies},
			"anomalies_struct": sliceAnomalies},
			"client.html")

	} else {
		fmt.Println("клиент НЕ найден в мапе клиентов")
		errMessage := "Client not found: " + apHostname
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
