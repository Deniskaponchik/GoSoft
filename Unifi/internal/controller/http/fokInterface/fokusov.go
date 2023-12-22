package fokusov

import (
	"github.com/deniskaponchik/GoSoft/Unifi/internal/usecase"
	_ "github.com/evrone/go-clean-template/docs"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	//swaggerFiles "github.com/swaggo/files"
	//ginSwagger "github.com/swaggo/gin-swagger"
	// Swagger docs.
	//_ "github.com/deniskaponchik/GoSoft/Unifi/docs"
)

var router *gin.Engine

type Fokusov struct {
	//Router    *gin.Engine
	Port string

	//UnifiUC *usecase.UnifiUseCase
	Urest       usecase.UnifiRestIn //interface. НЕ ИСПОЛЬЗОВАТЬ разыменовыватель *
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

var adminkaPageMsg []string

// NewRouter -.
// Swagger spec:
// @title       Go Clean Template API
// @description Using a translation service as an example
// @version     1.0
// @host        localhost:8080
// @BasePath    /v1
func (fok *Fokusov) Start() {

	adminkaPageMsg = make([]string, 10)

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
	router = gin.Default() //fok.Router = gin.Default()

	// Process the templates at the start so that they don't have to be loaded from the disk again. This makes serving HTML pages very fast.
	router.LoadHTMLGlob("../../web/templates/*")

	// Initialize the routes
	fok.initializeRoutes()

	fok.Logger.Println("Port : " + fok.Port)
	//router.Run(":8081")
	err = router.Run(":" + fok.Port)
	if err != nil {
		log.Fatalf(err.Error())
	}

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
