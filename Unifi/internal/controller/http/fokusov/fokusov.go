package fokusov

import (
	"github.com/deniskaponchik/GoSoft/Unifi/internal/entity"
	"github.com/deniskaponchik/GoSoft/Unifi/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

var router *gin.Engine

type Fokusov struct {
	Router      *gin.Engine
	UnifiUC     *usecase.UnifiUseCase
	UnifiClient *entity.Client
}

func New(uuc *usecase.UnifiUseCase) *Fokusov { //router *gin.Engine,
	return &Fokusov{
		Router: router,
		//Router:  *gin.Engine,
		UnifiUC: uuc,
	}
}

func (fok *Fokusov) Start() {
	// Set Gin to production mode
	gin.SetMode(gin.ReleaseMode)

	// Set the router as the default one provided by Gin
	//router = gin.Default()
	fok.Router = gin.Default()

	// Process the templates at the start so that they don't have to be loaded from the disk again. This makes serving HTML pages very fast.
	//router.LoadHTMLGlob("templates/*")
	fok.Router.LoadHTMLGlob("templates/*")
	//router.LoadHTMLGlob("web/templates/*")
	//fok.Router.LoadHTMLGlob("web/templates/*")

	// Initialize the routes
	//initializeRoutes()
	fok.Router.GET("/", showIndexPage)
	clientRoutes := fok.Router.Group("/client")
	{
		// Handle GET requests at /article/view/some_article_id
		clientRoutes.GET("/view/:client_hostname", fok.getClient)
	}

	// Start serving the application
	//router.Run()
	fok.Router.Run()
}

func (fok *Fokusov) getClient(c *gin.Context) {
	// Check if the client hostname is valid
	//if articleID, err := strconv.Atoi(c.Param("article_id")); err == nil {
	clientHostname := c.Param("client_hostname")

	// Check if the client exists
	//if client, err := getArticleByID(articleID); err == nil {
	//client := fok.UnifiUC.GetClientForRest(clientHostname)
	fok.UnifiClient = fok.UnifiUC.GetClientForRest(clientHostname)
	if fok.UnifiClient != nil {

		// Call the render function with the title, article and the name of the
		// template
		render(c, gin.H{
			"title":            fok.UnifiClient.Hostname,
			"hostname":         fok.UnifiClient.Hostname,
			"anomalies_struct": fok.UnifiClient.SliceAnomalies},
			"article.html")

	} else {
		// If the client is not found, abort with an error
		//c.AbortWithError(http.StatusNotFound, err)
		render(c, gin.H{
			"title":    "Client did not found",
			"hostname": "Client did not found"},
			//"anomalies_struct": client.SliceAnomalies},
			"article.html")
	}
}

func showIndexPage(c *gin.Context) {
	//articles := getAllArticles()

	// Call the render function with the name of the template to render
	render(c, gin.H{
		"title": "Home Page"},
		"index.html")
	//"payload": articles}, "index.html")
}

// Render one of HTML, JSON or CSV based on the 'Accept' header of the request
// If the header doesn't specify this, HTML is rendered, provided that the template name is present
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
