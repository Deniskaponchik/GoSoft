// handlers.client.go
package fokusov

import (
	"github.com/deniskaponchik/GoSoft/internal/entity"
	"github.com/gin-gonic/gin"
	"strings"
)

func (fok *Fokusov) getClientTest(c *gin.Context) {
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
			"page_test":    true,
			"title":        client.Hostname,
			"hostname":     client.Hostname,
			"countanomaly": client.CountAnomaly,
			"redmarker":    redMarker,
			//"anomalies_struct": client.SliceAnomalies},
			"anomalies_struct": sliceAnomalies},
			"test.html")

	} else {
		fok.Logger.Println("клиент НЕ найден в мапе клиентов")
		errMessage := "Client not found: " + clientHostname
		// If the client is not found, abort with an error
		//c.AbortWithError(http.StatusNotFound, err)
		render(c, gin.H{
			"title":    "Client not found",
			"hostname": errMessage},
			//"anomalies_struct": client.SliceAnomalies},
			"test.html")
	}
}

func (fok *Fokusov) showClientTest(c *gin.Context) {
	// Call the render function with the name of the template to render
	render(c, gin.H{
		"page_test": true,
		"title":     "Client Request Page"},
		"test.html")
}
