package fokusov

import (
	"github.com/deniskaponchik/GoSoft/internal/usecase"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
)

var router *gin.Engine
var adminkaPageMsg []string

type Fokusov struct {
	//Router    *gin.Engine
	Port string
	//JwtKey      string
	Urest usecase.UnifiRestIn //interface. НЕ ИСПОЛЬЗОВАТЬ разыменовыватель *
	//Authentication usecase.Authentication //interface. НЕ ИСПОЛЬЗОВАТЬ разыменовыватель *
	//Authorization  usecase.Authorization  //interface. НЕ ИСПОЛЬЗОВАТЬ разыменовыватель *
	LogFileName string
	Logger      *log.Logger
	CookieTTL   int
}

func New(uuc *usecase.UnifiUseCase, port string, logFileName string, ct int) *Fokusov { // jwtKey string,
	//router *gin.Engine,rest *usecase.Rest
	return &Fokusov{
		Port:      port,
		CookieTTL: ct * 60,
		//Router:  *gin.Engine,
		//UnifiUC: uuc,
		Urest: uuc, //использовать структуру, реализующую методы интерфейса usecase.UnifiRest
		//Authentication: uuc,
		LogFileName: logFileName,
	}
}

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

	// Process the templates at the start so that they don't have to be loaded from the disk again.
	//This makes serving HTML pages very fast.
	//router.LoadHTMLGlob("../../web/templates/*")
	//router.LoadHTMLGlob("templates/*")
	router.LoadHTMLGlob("web/templates/*")

	fok.initializeRoutes() // Initialize the routes

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
		c.JSON(http.StatusOK, data["payload"]) // Respond with JSON
	case "application/xml":
		c.XML(http.StatusOK, data["payload"]) // Respond with XML
	default:
		c.HTML(http.StatusOK, templateName, data) // Respond with HTML
	}
}

func (fok *Fokusov) Stop() {

}
