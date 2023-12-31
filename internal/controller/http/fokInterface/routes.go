package fokusov

import (
	_ "github.com/deniskaponchik/GoSoft/docs"  // Swagger docs.
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

// NewRouter -.
// Swagger spec:
// @title       IT Support App
// @description Handling Unifi Wi-Fi controller
// @version     1.0
// @host        localhost:8081
// @BasePath    /
func (fok *Fokusov) initializeRoutes() {
	// Use the setUserStatus middleware for every route to set a flag indicating
	// whether the request was from an authenticated user or not
	router.Use(fok.setUserStatus()) //router.Use(setUserStatus())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	siteRoutes := router.Group("/office")
	{
		//adminka page
		siteRoutes.POST("/login_change", fok.officeLoginChange)
		siteRoutes.POST("/change", fok.officeChange)
		siteRoutes.POST("/sapcn_add", fok.officeNew)

		//siteRoutes.POST("/exception_add", fok.addException)
		//siteRoutes.POST("/exception_del", fok.delException)
	}

	userRoutes := router.Group("/user")
	{
		userRoutes.GET("/login", ensureNotLoggedIn(), showLoginPage)
		userRoutes.POST("/login", ensureNotLoggedIn(), fok.performLogin)

		userRoutes.GET("/logout", fok.ensureLoggedIn(), logout)
		userRoutes.GET("/adminka", fok.ensureLoggedIn(), fok.showAdminkaPage)
		//redirect from POST performLogin
		userRoutes.POST("/adminka", fok.ensureLoggedIn(), fok.showAdminkaPage)
	}

	testRoutes := router.Group("/test")
	{
		testRoutes.GET("/request", fok.showClientTest)
		testRoutes.POST("/request", fok.getClientTest)
	}

	clientRoutes := router.Group("/client")
	{
		clientRoutes.GET("/request", fok.showClientRequestPage)
		clientRoutes.POST("/request", fok.getClient)
		clientRoutes.GET("/view/:client_hostname", fok.getClient)
	}

	apRoutes := router.Group("/ap")
	{
		apRoutes.GET("/request", fok.showApRequestPage)
		apRoutes.POST("/request", fok.getAP)
		apRoutes.GET("/view/:ap_hostname", fok.getAP)

		//apRoutes.POST("/exception_add", fok.addException)
		//apRoutes.POST("/exception_del", fok.delException)
	}

	/*
		// Use the setUserStatus middleware for every route to set a flag
		// indicating whether the request was from an authenticated user or not
		//router.Use(setUserStatus())

		// Handle the index route
		router.GET("/", showIndexPage)

		// Group article related routes together
		//articleRoutes := router.Group("/article")
		clientRoutes := router.Group("/client")
		{
			// Handle GET requests at /article/view/some_article_id
			clientRoutes.GET("/view/:client_hostname", getClient)

			// Handle the GET requests at /article/create
			// Show the article creation page
			// Ensure that the user is logged in by using the middleware
			//articleRoutes.GET("/create", ensureLoggedIn(), showArticleCreationPage)

			// Handle POST requests at /article/create
			// Ensure that the user is logged in by using the middleware
			//articleRoutes.POST("/create", ensureLoggedIn(), createArticle)
		}

		// Group user related routes together
		userRoutes := router.Group("/u")
		{
			// Handle the GET requests at /u/login
			// Show the login page
			// Ensure that the user is not logged in by using the middleware
			userRoutes.GET("/login", ensureNotLoggedIn(), showLoginPage)

			// Handle POST requests at /u/login
			// Ensure that the user is not logged in by using the middleware
			userRoutes.POST("/login", ensureNotLoggedIn(), performLogin)

			// Handle GET requests at /u/logout
			// Ensure that the user is logged in by using the middleware
			userRoutes.GET("/logout", ensureLoggedIn(), logout)

			// Handle the GET requests at /u/register
			// Show the registration page
			// Ensure that the user is not logged in by using the middleware
			userRoutes.GET("/register", ensureNotLoggedIn(), showRegistrationPage)

			// Handle POST requests at /u/register
			// Ensure that the user is not logged in by using the middleware
			userRoutes.POST("/register", ensureNotLoggedIn(), register)
		}*/
}
