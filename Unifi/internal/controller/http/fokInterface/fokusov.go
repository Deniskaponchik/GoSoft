package fokusov

import (
	"github.com/deniskaponchik/GoSoft/Unifi/internal/usecase"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	//"time"
)

var router *gin.Engine

type Fokusov struct {
	//Router    *gin.Engine
	Port string

	//UnifiUC *usecase.UnifiUseCase
	Urest       usecase.UnifiRest //interface. НЕ ИСПОЛЬЗОВАТЬ разыменовыватель *
	LogFileName string
	Logger      *log.Logger
}

func New(uuc *usecase.UnifiUseCase, port string, logFileName string) *Fokusov { //router *gin.Engine,rest *usecase.Rest
	return &Fokusov{
		Port: port,
		//Router: router,
		//Router:  *gin.Engine,

		//UnifiUC: uuc,
		Urest:       uuc, //использовать структуру, реализующие методы интерфейса usecase.UnifiRest
		LogFileName: logFileName,
	}
}

func (fok *Fokusov) Start() {

	//FileNameGin := "Unifi_Gin_" + time.Now().Format("2006-01-02_15.04.05") + ".log"
	fileLogGin, err := os.OpenFile(fok.LogFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	multiWriter := io.MultiWriter(os.Stdout, fileLogGin)

	gin.DefaultWriter = multiWriter
	gin.DefaultErrorWriter = multiWriter
	fok.Logger = log.New(multiWriter, "", 0)

	gin.SetMode(gin.ReleaseMode) // Set Gin to production mode

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

	//port := ":" + fok.Port
	//fmt.Println(fok.Port)
	fok.Logger.Println(fok.Port)
	//gin.Context{}.String()
	//router.Run(":8081")
	err = router.Run(":" + fok.Port)
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
