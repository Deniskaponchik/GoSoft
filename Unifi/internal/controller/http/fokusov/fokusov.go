package fokusov

import (
	"fmt"
	"github.com/deniskaponchik/GoSoft/Unifi/internal/usecase"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var router *gin.Engine

type Fokusov struct {
	//Router    *gin.Engine
	UnifiUC *usecase.UnifiUseCase
	//Rest 		*usecase.Rest  //interface
	//UnifiClient *entity.Client
	Port string
}

func New(uuc *usecase.UnifiUseCase, port string) *Fokusov { //router *gin.Engine,rest *usecase.Rest
	return &Fokusov{
		//Router: router,
		//Router:  *gin.Engine,
		UnifiUC: uuc,
		//Rest:    rest,
		Port: port,
	}
}

func (fok *Fokusov) Start() {
	// Set Gin to production mode
	gin.SetMode(gin.ReleaseMode)

	// Set the router as the default one provided by Gin
	router = gin.Default()
	//fok.Router = gin.Default()

	// Process the templates at the start so that they don't have to be loaded from the disk again. This makes serving HTML pages very fast.
	//router.LoadHTMLGlob("templates/*")
	router.LoadHTMLGlob("../../web/templates/*")
	//fok.Router.LoadHTMLGlob("templates/*")
	//fok.Router.LoadHTMLGlob("./internal/controller/http/fokusov/templates/*") //НЕ РАБОТАЕТ
	//fok.Router.LoadHTMLGlob("../../web/templates/*")

	// Initialize the routes
	//initializeRoutes()
	// Use the setUserStatus middleware for every route to set a flag
	// indicating whether the request was from an authenticated user or not
	//router.Use(setUserStatus())

	//router.GET("/", showIndexPage)
	//router.GET("/", fok.showIndexPage)
	//fok.Router.GET("/", fok.showIndexPage)

	clientRoutes := router.Group("/client")
	//clientRoutes := fok.Router.Group("/client")
	{
		// Handle GET requests at
		clientRoutes.GET("/request", fok.showClientRequestPage)
		// Handle POST requests at
		clientRoutes.POST("/request", fok.getClient)
		// Handle GET requests at /article/view/some_article_id
		clientRoutes.GET("/view/:client_hostname", fok.getClient)
	}
	apRoutes := router.Group("/ap")
	//clientRoutes := fok.Router.Group("/client")
	{
		// Handle GET requests at
		apRoutes.GET("/request", fok.showApRequestPage)
		// Handle POST requests at
		apRoutes.POST("/request", fok.getAP)
		// Handle GET requests at /article/view/some_article_id
		apRoutes.GET("/view/:ap_hostname", fok.getAP)
	}

	// Start serving the application
	//port := ":" + fok.Port
	fmt.Println(fok.Port)
	//router.Run(":8081")
	err := router.Run(":" + fok.Port)
	//fok.Router.Run()
	if err != nil {
		log.Fatalf(err.Error())
	}

}

// Render one of HTML, JSON or CSV based on the 'Accept' header of the request
// If the header doesn't specify this, HTML is rendered, provided that the template name is present
func render(c *gin.Context, data gin.H, templateName string) {
	//loggedInInterface, _ := c.Get("is_logged_in")
	//data["is_logged_in"] = loggedInInterface.(bool)

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
