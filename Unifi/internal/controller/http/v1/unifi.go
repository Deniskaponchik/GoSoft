package v1

import (
	"github.com/deniskaponchik/GoSoft/Unifi/internal/entity"
	"github.com/deniskaponchik/GoSoft/Unifi/internal/usecase"
	"github.com/deniskaponchik/GoSoft/Unifi/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type unifiRoutes struct {
	//t usecase.Translation
	u *usecase.UnifiUseCase
	l logger.Interface
}

func newUnifiRoutes(handler *gin.RouterGroup, u *usecase.UnifiUseCase, l logger.Interface) {
	r := &unifiRoutes{u, l}

	clientRoutes := handler.Group("/client")
	{
		//h.GET("/history", r.history)
		//h.POST("/do-translate", r.doTranslate)
		clientRoutes.GET("/view/:client_hostname", r.getClient)
	}
	/*
		apsRoutes := handler.Group("/ap")
		{
			apsRoutes.GET("/region", r.getAps)
		}*/
}

type clientResponse struct {
	Client entity.Client `json:"history"`
}

func (r *unifiRoutes) getClient(c *gin.Context) {
	//translations, err := r.t.History(c.Request.Context())
	client := r.u.GetClientForRest(c.Param("client_hostname"))
	if client != nil {
		render(c, gin.H{
			//"title":   article.Title,
			"title":    client.Hostname, //используется в header.html
			"hostname": client.Hostname,
			//"payload": article},
			"anomalies_struct": client.SliceAnomalies},
			//"anomalies_string": client.SliceAnomalies},
			"client.html")
	} else {

	}
}

// Render one of HTML, JSON or CSV based on the 'Accept' header of the request
// If the header doesn't specify this, HTML is rendered, provided that
// the template name is present
func render(c *gin.Context, data gin.H, templateName string) {
	loggedInInterface, _ := c.Get("is_logged_in")
	data["is_logged_in"] = loggedInInterface.(bool)

	switch c.Request.Header.Get("Accept") {
	case "application/json":
		// Respond with JSON
		c.JSON(http.StatusOK, data["payload"])
	case "application/xml":
		// Respond with XML
		c.XML(http.StatusOK, data["payload"])
	default:
		// Respond with HTML
		c.HTML(http.StatusOK, templateName, data)
	}
}

/*https://github.com/evrone/go-clean-template/blob/master/internal/controller/http/v1/translation.go
type historyResponse struct {
	History []entity.Translation `json:"history"`
}
// @Summary     Show history
// @Description Show all translation history
// @ID          history
// @Tags  	    translation
// @Accept      json
// @Produce     json
// @Success     200 {object} historyResponse
// @Failure     500 {object} response
// @Router      /translation/history [get]
func (r *translationRoutes) history(c *gin.Context) {
	translations, err := r.t.History(c.Request.Context())
	if err != nil {
		r.l.Error(err, "http - v1 - history")
		errorResponse(c, http.StatusInternalServerError, "database problems")
		return
	}

	c.JSON(http.StatusOK, historyResponse{translations})
}

type doTranslateRequest struct {
	Source      string `json:"source"       binding:"required"  example:"auto"`
	Destination string `json:"destination"  binding:"required"  example:"en"`
	Original    string `json:"original"     binding:"required"  example:"текст для перевода"`
}

// @Summary     Translate
// @Description Translate a text
// @ID          do-translate
// @Tags  	    translation
// @Accept      json
// @Produce     json
// @Param       request body doTranslateRequest true "Set up translation"
// @Success     200 {object} entity.Translation
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /translation/do-translate [post]
func (r *translationRoutes) doTranslate(c *gin.Context) {
	var request doTranslateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - v1 - doTranslate")
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	translation, err := r.t.Translate(
		c.Request.Context(),
		entity.Translation{
			Source:      request.Source,
			Destination: request.Destination,
			Original:    request.Original,
		},
	)
	if err != nil {
		r.l.Error(err, "http - v1 - doTranslate")
		errorResponse(c, http.StatusInternalServerError, "translation service problems")
		return
	}

	c.JSON(http.StatusOK, translation)
}*/
